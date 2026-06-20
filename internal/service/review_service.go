package service

import (
	"letter-square-api/internal/apperror"
	"letter-square-api/internal/dto"
	"letter-square-api/internal/entity"
	"letter-square-api/internal/repository"
)

type ReviewService interface {
	Create(userID, movieID int64, req *dto.CreateReviewRequest) (*entity.Review, error)
	GetByMovie(movieID int64, page, limit int) ([]*entity.Review, error)
	GetByUser(userID int64, page, limit int) ([]*entity.Review, error)
	Update(reviewID, userID int64, req *dto.UpdateReviewRequest) (*entity.Review, error)
	Delete(reviewID, userID int64) error
}

type reviewService struct {
	reviewRepo repository.ReviewRepository
	movieRepo  repository.MovieRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository, movieRepo repository.MovieRepository) ReviewService {
	return &reviewService{reviewRepo: reviewRepo, movieRepo: movieRepo}
}

func (s *reviewService) Create(userID, movieID int64, req *dto.CreateReviewRequest) (*entity.Review, error) {
	// Movie must exist
	m, err := s.movieRepo.FindByID(movieID)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if m == nil {
		return nil, apperror.ErrMovieNotFound
	}

	// One review per user per movie
	exists, err := s.reviewRepo.ExistsByUserAndMovie(userID, movieID)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if exists {
		return nil, apperror.ErrDuplicateReview
	}

	r := &entity.Review{
		UserID:  userID,
		MovieID: movieID,
		Rating:  req.Rating,
		Body:    req.Body,
	}
	if err := s.reviewRepo.Create(r); err != nil {
		return nil, apperror.ErrInternal
	}

	// Return with joins
	created, err := s.reviewRepo.FindByID(r.ID)
	if err != nil || created == nil {
		return r, nil
	}
	return created, nil
}

func (s *reviewService) GetByMovie(movieID int64, page, limit int) ([]*entity.Review, error) {
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	reviews, err := s.reviewRepo.FindByMovieID(movieID, limit, offset)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	return reviews, nil
}

func (s *reviewService) GetByUser(userID int64, page, limit int) ([]*entity.Review, error) {
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	reviews, err := s.reviewRepo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	return reviews, nil
}

func (s *reviewService) Update(reviewID, userID int64, req *dto.UpdateReviewRequest) (*entity.Review, error) {
	rv, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if rv == nil {
		return nil, apperror.ErrReviewNotFound
	}
	if rv.UserID != userID {
		return nil, apperror.ErrForbidden
	}

	if req.Rating != nil {
		rv.Rating = *req.Rating
	}
	if req.Body != nil {
		rv.Body = *req.Body
	}

	if err := s.reviewRepo.Update(rv); err != nil {
		return nil, apperror.ErrInternal
	}
	return rv, nil
}

func (s *reviewService) Delete(reviewID, userID int64) error {
	rv, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		return apperror.ErrInternal
	}
	if rv == nil {
		return apperror.ErrReviewNotFound
	}
	if rv.UserID != userID {
		return apperror.ErrForbidden
	}
	if err := s.reviewRepo.Delete(reviewID); err != nil {
		return apperror.ErrInternal
	}
	return nil
}
