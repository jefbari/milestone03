package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"letter-square-api/internal/apperror"
	"letter-square-api/internal/dto"
	"letter-square-api/internal/entity"
	"letter-square-api/internal/repository"
)

const recommendationIntro = "Your next favorite film is probably one you've never heard of.\nLet's find it."

// recommendationQuestions are asked in order, one per AnswerQuestion call,
// before Gemini is asked to generate the final recommendation.
var recommendationQuestions = []string{
	"What's your mood today — calm and curious, or craving something thrilling?",
	"How much time have you got? A quick watch under 90 minutes, or no limit?",
	"Want something close to the vibe of your watchlist, or a total curveball?",
}

// GeminiClient is the minimal interface this service needs from the Gemini
// third-party client. Defined here (consumer side) so it can be mocked in tests.
type GeminiClient interface {
	GenerateText(prompt string) (string, error)
}

type RecommendationService interface {
	StartSession(userID int64) (*dto.StartRecommendationResponse, error)
	AnswerQuestion(sessionID, userID int64, answer string) (*dto.AnswerRecommendationResponse, error)
	GetSession(sessionID, userID int64) (*entity.RecommendationSession, error)
}

type recommendationService struct {
	sessionRepo   repository.RecommendationSessionRepository
	watchlistRepo repository.WatchlistRepository
	gemini        GeminiClient
}

func NewRecommendationService(
	sessionRepo repository.RecommendationSessionRepository,
	watchlistRepo repository.WatchlistRepository,
	gemini GeminiClient,
) RecommendationService {
	return &recommendationService{
		sessionRepo:   sessionRepo,
		watchlistRepo: watchlistRepo,
		gemini:        gemini,
	}
}

func (s *recommendationService) StartSession(userID int64) (*dto.StartRecommendationResponse, error) {
	session := &entity.RecommendationSession{
		UserID:  userID,
		Status:  "in_progress",
		Step:    0,
		Answers: "[]",
	}
	if err := s.sessionRepo.Create(session); err != nil {
		return nil, apperror.ErrInternal
	}

	return &dto.StartRecommendationResponse{
		SessionID:  session.ID,
		Intro:      recommendationIntro,
		Question:   recommendationQuestions[0],
		Step:       1,
		TotalSteps: len(recommendationQuestions),
	}, nil
}

func (s *recommendationService) AnswerQuestion(sessionID, userID int64, answer string) (*dto.AnswerRecommendationResponse, error) {
	if strings.TrimSpace(answer) == "" {
		return nil, apperror.New(400, "answer is required")
	}

	session, err := s.sessionRepo.FindByID(sessionID)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if session == nil {
		return nil, apperror.ErrNotFound
	}
	if session.UserID != userID {
		return nil, apperror.ErrForbidden
	}
	if session.Status == "completed" {
		return nil, apperror.New(409, "this recommendation session is already completed")
	}

	var answers []string
	_ = json.Unmarshal([]byte(session.Answers), &answers)
	answers = append(answers, answer)

	nextIndex := session.Step + 1

	// Still more questions to ask.
	if nextIndex < len(recommendationQuestions) {
		b, _ := json.Marshal(answers)
		session.Answers = string(b)
		session.Step = nextIndex
		if err := s.sessionRepo.Update(session); err != nil {
			return nil, apperror.ErrInternal
		}
		return &dto.AnswerRecommendationResponse{
			SessionID:  session.ID,
			Status:     "in_progress",
			Question:   recommendationQuestions[nextIndex],
			Step:       nextIndex + 1,
			TotalSteps: len(recommendationQuestions),
		}, nil
	}

	// Last answer received — generate the recommendation using the user's
	// watchlist as context, plus everything they answered.
	watchlistItems, err := s.watchlistRepo.FindByUserID(userID, 50, 0)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	var titles []string
	for _, w := range watchlistItems {
		if w.Movie != nil {
			titles = append(titles, w.Movie.Title)
		}
	}
	watchlistText := "empty — they haven't added anything yet"
	if len(titles) > 0 {
		watchlistText = strings.Join(titles, ", ")
	}

	prompt := buildRecommendationPrompt(watchlistText, answers)
	result, err := s.gemini.GenerateText(prompt)
	if err != nil {
		return nil, apperror.New(503, "recommendation service unavailable: "+err.Error())
	}

	b, _ := json.Marshal(answers)
	session.Answers = string(b)
	session.Step = nextIndex
	session.Status = "completed"
	session.Result = result
	if err := s.sessionRepo.Update(session); err != nil {
		return nil, apperror.ErrInternal
	}

	return &dto.AnswerRecommendationResponse{
		SessionID:      session.ID,
		Status:         "completed",
		Recommendation: result,
	}, nil
}

func (s *recommendationService) GetSession(sessionID, userID int64) (*entity.RecommendationSession, error) {
	session, err := s.sessionRepo.FindByID(sessionID)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if session == nil {
		return nil, apperror.ErrNotFound
	}
	if session.UserID != userID {
		return nil, apperror.ErrForbidden
	}
	return session, nil
}

func buildRecommendationPrompt(watchlist string, answers []string) string {
	mood := safeGet(answers, 0)
	timeAvailable := safeGet(answers, 1)
	direction := safeGet(answers, 2)

	return fmt.Sprintf(`You are a thoughtful movie recommendation assistant for a platform called LetterSquare.
A user is looking for their next favorite film — something they probably haven't heard of yet.

Their current watchlist: %s

Their mood today: %s
Time available: %s
Direction they want: %s

Based on this, recommend 3-5 movies. For each, briefly explain why it fits their mood and
how it connects to (or deliberately breaks from) their watchlist. Keep the tone warm and
personal, like a friend who really knows film. Respond in plain text with a numbered list.`,
		watchlist, mood, timeAvailable, direction)
}

func safeGet(arr []string, idx int) string {
	if idx < len(arr) {
		return arr[idx]
	}
	return "not specified"
}
