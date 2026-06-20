package handler

import (
	"net/http"
	"strconv"

	"letter-square-api/internal/dto"
	"letter-square-api/internal/helper"
	"letter-square-api/internal/middleware"
	"letter-square-api/internal/service"
)

type ReviewHandler struct {
	reviewSvc service.ReviewService
}

func NewReviewHandler(reviewSvc service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewSvc: reviewSvc}
}

// POST /api/movies/:id/reviews
func (h *ReviewHandler) Create(w http.ResponseWriter, r *http.Request) {
	movieID, err := parseID(r, "id")
	if err != nil {
		helper.WriteError(w, &dto.ValidationError{Msg: "invalid movie id"})
		return
	}
	userID := middleware.GetUserID(r)

	var req dto.CreateReviewRequest
	if err := helper.DecodeJSON(r, &req); err != nil {
		helper.WriteError(w, &dto.ValidationError{Msg: "invalid request body"})
		return
	}
	if err := req.Validate(); err != nil {
		helper.WriteError(w, err)
		return
	}

	review, err := h.reviewSvc.Create(userID, movieID, &req)
	if err != nil {
		helper.WriteError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusCreated, "review created", review)
}

// GET /api/movies/:id/reviews
func (h *ReviewHandler) GetByMovie(w http.ResponseWriter, r *http.Request) {
	movieID, err := parseID(r, "id")
	if err != nil {
		helper.WriteError(w, &dto.ValidationError{Msg: "invalid movie id"})
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	reviews, err := h.reviewSvc.GetByMovie(movieID, page, limit)
	if err != nil {
		helper.WriteError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "reviews retrieved", reviews)
}

// GET /api/users/me/reviews
func (h *ReviewHandler) GetMyReviews(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	reviews, err := h.reviewSvc.GetByUser(userID, page, limit)
	if err != nil {
		helper.WriteError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "reviews retrieved", reviews)
}

// PUT /api/reviews/:reviewId
func (h *ReviewHandler) Update(w http.ResponseWriter, r *http.Request) {
	reviewID, err := parseID(r, "reviewId")
	if err != nil {
		helper.WriteError(w, &dto.ValidationError{Msg: "invalid review id"})
		return
	}
	userID := middleware.GetUserID(r)

	var req dto.UpdateReviewRequest
	if err := helper.DecodeJSON(r, &req); err != nil {
		helper.WriteError(w, &dto.ValidationError{Msg: "invalid request body"})
		return
	}
	if err := req.Validate(); err != nil {
		helper.WriteError(w, err)
		return
	}

	review, err := h.reviewSvc.Update(reviewID, userID, &req)
	if err != nil {
		helper.WriteError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "review updated", review)
}

// DELETE /api/reviews/:reviewId
func (h *ReviewHandler) Delete(w http.ResponseWriter, r *http.Request) {
	reviewID, err := parseID(r, "reviewId")
	if err != nil {
		helper.WriteError(w, &dto.ValidationError{Msg: "invalid review id"})
		return
	}
	userID := middleware.GetUserID(r)

	if err := h.reviewSvc.Delete(reviewID, userID); err != nil {
		helper.WriteError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "review deleted", nil)
}
