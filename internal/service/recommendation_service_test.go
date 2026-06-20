package service_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"letter-square-api/internal/apperror"
	"letter-square-api/internal/entity"
	"letter-square-api/internal/service"
)

// Mock RecommendationSessionRepository

type mockSessionRepo struct {
	mock.Mock
}

func (m *mockSessionRepo) Create(s *entity.RecommendationSession) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m *mockSessionRepo) FindByID(id int64) (*entity.RecommendationSession, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.RecommendationSession), args.Error(1)
}

func (m *mockSessionRepo) Update(s *entity.RecommendationSession) error {
	args := m.Called(s)
	return args.Error(0)
}

// Mock WatchlistRepository

type mockWatchlistRepo struct {
	mock.Mock
}

func (m *mockWatchlistRepo) Add(userID, movieID int64) error {
	args := m.Called(userID, movieID)
	return args.Error(0)
}

func (m *mockWatchlistRepo) Remove(userID, movieID int64) error {
	args := m.Called(userID, movieID)
	return args.Error(0)
}

func (m *mockWatchlistRepo) FindByUserID(userID int64, limit, offset int) ([]*entity.Watchlist, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]*entity.Watchlist), args.Error(1)
}

func (m *mockWatchlistRepo) Exists(userID, movieID int64) (bool, error) {
	args := m.Called(userID, movieID)
	return args.Bool(0), args.Error(1)
}

// ─── Mock GeminiClient ────────────────────────────────────────────────────────

type mockGeminiClient struct {
	mock.Mock
}

func (m *mockGeminiClient) GenerateText(prompt string) (string, error) {
	args := m.Called(prompt)
	return args.String(0), args.Error(1)
}

// ─── Tests ────────────────────────────────────────────────────────────────────

func TestStartSession_Success(t *testing.T) {
	sessionRepo := new(mockSessionRepo)
	wlRepo := new(mockWatchlistRepo)
	gemini := new(mockGeminiClient)
	svc := service.NewRecommendationService(sessionRepo, wlRepo, gemini)

	sessionRepo.On("Create", mock.AnythingOfType("*entity.RecommendationSession")).
		Run(func(args mock.Arguments) {
			s := args.Get(0).(*entity.RecommendationSession)
			s.ID = 1
		}).Return(nil)

	result, err := svc.StartSession(42)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.SessionID)
	assert.Contains(t, result.Intro, "Your next favorite film")
	assert.Equal(t, 1, result.Step)
	assert.Equal(t, 3, result.TotalSteps)
	assert.NotEmpty(t, result.Question)
	sessionRepo.AssertExpectations(t)
}

func TestAnswerQuestion_ProgressesToNextQuestion(t *testing.T) {
	sessionRepo := new(mockSessionRepo)
	wlRepo := new(mockWatchlistRepo)
	gemini := new(mockGeminiClient)
	svc := service.NewRecommendationService(sessionRepo, wlRepo, gemini)

	session := &entity.RecommendationSession{
		ID:      1,
		UserID:  42,
		Status:  "in_progress",
		Step:    0,
		Answers: "[]",
	}
	sessionRepo.On("FindByID", int64(1)).Return(session, nil)
	sessionRepo.On("Update", mock.AnythingOfType("*entity.RecommendationSession")).Return(nil)

	result, err := svc.AnswerQuestion(1, 42, "feeling adventurous")

	assert.NoError(t, err)
	assert.Equal(t, "in_progress", result.Status)
	assert.Equal(t, 2, result.Step) // moved to question #2
	assert.NotEmpty(t, result.Question)
	assert.Empty(t, result.Recommendation)
	sessionRepo.AssertExpectations(t)
	wlRepo.AssertNotCalled(t, "FindByUserID", mock.Anything, mock.Anything, mock.Anything)
	gemini.AssertNotCalled(t, "GenerateText", mock.Anything)
}

func TestAnswerQuestion_FinalAnswerGeneratesRecommendation(t *testing.T) {
	sessionRepo := new(mockSessionRepo)
	wlRepo := new(mockWatchlistRepo)
	gemini := new(mockGeminiClient)
	svc := service.NewRecommendationService(sessionRepo, wlRepo, gemini)

	existingAnswers, _ := json.Marshal([]string{"adventurous", "under 90 minutes"})
	session := &entity.RecommendationSession{
		ID:      1,
		UserID:  42,
		Status:  "in_progress",
		Step:    2, // already answered question 0 and 1; this answer is for the last question (index 2)
		Answers: string(existingAnswers),
	}
	sessionRepo.On("FindByID", int64(1)).Return(session, nil)
	sessionRepo.On("Update", mock.AnythingOfType("*entity.RecommendationSession")).Return(nil)

	watchlist := []*entity.Watchlist{
		{Movie: &entity.Movie{Title: "Paprika"}},
		{Movie: &entity.Movie{Title: "Oldboy"}},
	}
	wlRepo.On("FindByUserID", int64(42), 50, 0).Return(watchlist, nil)
	gemini.On("GenerateText", mock.AnythingOfType("string")).
		Return("1. Perfect Blue - because you liked Paprika...", nil)

	result, err := svc.AnswerQuestion(1, 42, "surprise me")

	assert.NoError(t, err)
	assert.Equal(t, "completed", result.Status)
	assert.Contains(t, result.Recommendation, "Perfect Blue")
	assert.Empty(t, result.Question)
	sessionRepo.AssertExpectations(t)
	wlRepo.AssertExpectations(t)
	gemini.AssertExpectations(t)
}

func TestAnswerQuestion_Forbidden(t *testing.T) {
	sessionRepo := new(mockSessionRepo)
	wlRepo := new(mockWatchlistRepo)
	gemini := new(mockGeminiClient)
	svc := service.NewRecommendationService(sessionRepo, wlRepo, gemini)

	session := &entity.RecommendationSession{ID: 1, UserID: 99, Status: "in_progress"}
	sessionRepo.On("FindByID", int64(1)).Return(session, nil)

	result, err := svc.AnswerQuestion(1, 42, "some answer") // user 42 != owner 99

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, apperror.ErrForbidden, err)
	sessionRepo.AssertExpectations(t)
}

func TestAnswerQuestion_AlreadyCompleted(t *testing.T) {
	sessionRepo := new(mockSessionRepo)
	wlRepo := new(mockWatchlistRepo)
	gemini := new(mockGeminiClient)
	svc := service.NewRecommendationService(sessionRepo, wlRepo, gemini)

	session := &entity.RecommendationSession{ID: 1, UserID: 42, Status: "completed"}
	sessionRepo.On("FindByID", int64(1)).Return(session, nil)

	result, err := svc.AnswerQuestion(1, 42, "another answer")

	assert.Nil(t, result)
	assert.Error(t, err)
	sessionRepo.AssertExpectations(t)
}

func TestAnswerQuestion_EmptyAnswerRejected(t *testing.T) {
	sessionRepo := new(mockSessionRepo)
	wlRepo := new(mockWatchlistRepo)
	gemini := new(mockGeminiClient)
	svc := service.NewRecommendationService(sessionRepo, wlRepo, gemini)

	result, err := svc.AnswerQuestion(1, 42, "   ")

	assert.Nil(t, result)
	assert.Error(t, err)
	sessionRepo.AssertNotCalled(t, "FindByID", mock.Anything)
}
