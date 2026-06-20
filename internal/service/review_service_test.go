package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"letter-square-api/internal/apperror"
	"letter-square-api/internal/dto"
	"letter-square-api/internal/entity"
	"letter-square-api/internal/service"
)

// ─── Mock Review Repo ────────────────────────────────────────────────────────

type mockReviewRepo struct {
	mock.Mock
}

func (m *mockReviewRepo) Create(r *entity.Review) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *mockReviewRepo) FindByID(id int64) (*entity.Review, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Review), args.Error(1)
}

func (m *mockReviewRepo) FindByMovieID(movieID int64, limit, offset int) ([]*entity.Review, error) {
	args := m.Called(movieID, limit, offset)
	return args.Get(0).([]*entity.Review), args.Error(1)
}

func (m *mockReviewRepo) FindByUserID(userID int64, limit, offset int) ([]*entity.Review, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]*entity.Review), args.Error(1)
}

func (m *mockReviewRepo) Update(r *entity.Review) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *mockReviewRepo) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockReviewRepo) ExistsByUserAndMovie(userID, movieID int64) (bool, error) {
	args := m.Called(userID, movieID)
	return args.Bool(0), args.Error(1)
}

// ─── Mock Movie Repo ─────────────────────────────────────────────────────────

type mockMovieRepo struct {
	mock.Mock
}

func (m *mockMovieRepo) Create(mv *entity.Movie) error {
	args := m.Called(mv)
	return args.Error(0)
}

func (m *mockMovieRepo) FindByID(id int64) (*entity.Movie, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Movie), args.Error(1)
}

func (m *mockMovieRepo) FindAll(search, genre string, limit, offset int) ([]*entity.Movie, error) {
	args := m.Called(search, genre, limit, offset)
	return args.Get(0).([]*entity.Movie), args.Error(1)
}

func (m *mockMovieRepo) Update(mv *entity.Movie) error {
	args := m.Called(mv)
	return args.Error(0)
}

func (m *mockMovieRepo) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockMovieRepo) FindTopRated(limit int) ([]*entity.Movie, error) {
	args := m.Called(limit)
	return args.Get(0).([]*entity.Movie), args.Error(1)
}

// ─── Tests ───────────────────────────────────────────────────────────────────

func TestCreateReview_Success(t *testing.T) {
	reviewRepo := new(mockReviewRepo)
	movieRepo := new(mockMovieRepo)
	svc := service.NewReviewService(reviewRepo, movieRepo)

	movie := &entity.Movie{ID: 1, Title: "Inception"}
	review := &entity.Review{ID: 1, UserID: 10, MovieID: 1, Rating: 4.5, Body: "Great film!", MovieTitle: "Inception"}

	movieRepo.On("FindByID", int64(1)).Return(movie, nil)
	reviewRepo.On("ExistsByUserAndMovie", int64(10), int64(1)).Return(false, nil)
	reviewRepo.On("Create", mock.AnythingOfType("*entity.Review")).Return(nil)
	reviewRepo.On("FindByID", int64(0)).Return(review, nil) // after create, ID still 0 until DB sets it

	req := &dto.CreateReviewRequest{Rating: 4.5, Body: "Great film!"}
	result, err := svc.Create(10, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	reviewRepo.AssertExpectations(t)
	movieRepo.AssertExpectations(t)
}

func TestCreateReview_MovieNotFound(t *testing.T) {
	reviewRepo := new(mockReviewRepo)
	movieRepo := new(mockMovieRepo)
	svc := service.NewReviewService(reviewRepo, movieRepo)

	movieRepo.On("FindByID", int64(99)).Return(nil, nil)

	req := &dto.CreateReviewRequest{Rating: 4.0, Body: "Good movie"}
	result, err := svc.Create(10, 99, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, apperror.ErrMovieNotFound, err)
	movieRepo.AssertExpectations(t)
}

func TestCreateReview_Duplicate(t *testing.T) {
	reviewRepo := new(mockReviewRepo)
	movieRepo := new(mockMovieRepo)
	svc := service.NewReviewService(reviewRepo, movieRepo)

	movie := &entity.Movie{ID: 1, Title: "Inception"}
	movieRepo.On("FindByID", int64(1)).Return(movie, nil)
	reviewRepo.On("ExistsByUserAndMovie", int64(10), int64(1)).Return(true, nil)

	req := &dto.CreateReviewRequest{Rating: 3.0, Body: "Seen it"}
	result, err := svc.Create(10, 1, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, apperror.ErrDuplicateReview, err)
	reviewRepo.AssertExpectations(t)
	movieRepo.AssertExpectations(t)
}

func TestDeleteReview_Forbidden(t *testing.T) {
	reviewRepo := new(mockReviewRepo)
	movieRepo := new(mockMovieRepo)
	svc := service.NewReviewService(reviewRepo, movieRepo)

	review := &entity.Review{ID: 5, UserID: 99} // belongs to user 99
	reviewRepo.On("FindByID", int64(5)).Return(review, nil)

	err := svc.Delete(5, 10) // user 10 tries to delete

	assert.Error(t, err)
	assert.Equal(t, apperror.ErrForbidden, err)
	reviewRepo.AssertExpectations(t)
}

func TestUpdateReview_Success(t *testing.T) {
	reviewRepo := new(mockReviewRepo)
	movieRepo := new(mockMovieRepo)
	svc := service.NewReviewService(reviewRepo, movieRepo)

	review := &entity.Review{ID: 5, UserID: 10, Rating: 3.0, Body: "OK"}
	reviewRepo.On("FindByID", int64(5)).Return(review, nil)
	reviewRepo.On("Update", mock.AnythingOfType("*entity.Review")).Return(nil)

	newRating := 5.0
	newBody := "Actually amazing!"
	req := &dto.UpdateReviewRequest{Rating: &newRating, Body: &newBody}

	result, err := svc.Update(5, 10, req)

	assert.NoError(t, err)
	assert.Equal(t, 5.0, result.Rating)
	assert.Equal(t, "Actually amazing!", result.Body)
	reviewRepo.AssertExpectations(t)
}
