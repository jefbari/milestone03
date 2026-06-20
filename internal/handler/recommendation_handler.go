package handler

import (
	"net/http"

	"letter-square-api/internal/dto"
	"letter-square-api/internal/helper"
	"letter-square-api/internal/middleware"
	"letter-square-api/internal/service"
)

type RecommendationHandler struct {
	recSvc service.RecommendationService
}

func NewRecommendationHandler(recSvc service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{recSvc: recSvc}
}

// POST /api/recommendations/start
// Returns the opening line + the first question. This is the "landing page"
// of the recommendation flow, just expressed as JSON since this is an API.
func (h *RecommendationHandler) Start(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	result, err := h.recSvc.StartSession(userID)
	if err != nil {
		helper.WriteError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusCreated, "let's find your next favorite film", result)
}

// POST /api/recommendations/sessions/{id}/answer
// Submit an answer; returns either the next question, or - on the final
// answer - the Gemini-generated recommendation based on the watchlist.
func (h *RecommendationHandler) Answer(w http.ResponseWriter, r *http.Request) {
	sessionID, err := parseID(r, "id")
	if err != nil {
		helper.WriteError(w, &dto.ValidationError{Msg: "invalid session id"})
		return
	}
	userID := middleware.GetUserID(r)

	var req dto.AnswerRecommendationRequest
	if err := helper.DecodeJSON(r, &req); err != nil {
		helper.WriteError(w, &dto.ValidationError{Msg: "invalid request body"})
		return
	}
	if err := req.Validate(); err != nil {
		helper.WriteError(w, err)
		return
	}

	result, err := h.recSvc.AnswerQuestion(sessionID, userID, req.Answer)
	if err != nil {
		helper.WriteError(w, err)
		return
	}

	msg := "next question"
	if result.Status == "completed" {
		msg = "here's what we found for you"
	}
	helper.WriteSuccess(w, http.StatusOK, msg, result)
}

// GET /api/recommendations/sessions/{id}
func (h *RecommendationHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	sessionID, err := parseID(r, "id")
	if err != nil {
		helper.WriteError(w, &dto.ValidationError{Msg: "invalid session id"})
		return
	}
	userID := middleware.GetUserID(r)

	session, err := h.recSvc.GetSession(sessionID, userID)
	if err != nil {
		helper.WriteError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "session retrieved", session)
}
