package handler

import (
	"net/http"
	"strconv"

	"letter-square-api/internal/dto"
	"letter-square-api/internal/helper"
	"letter-square-api/internal/middleware"
	"letter-square-api/internal/service"
)

type WatchlistHandler struct {
	watchlistSvc service.WatchlistService
}

func NewWatchlistHandler(watchlistSvc service.WatchlistService) *WatchlistHandler {
	return &WatchlistHandler{watchlistSvc: watchlistSvc}
}

// POST /api/watchlist/:movieId
func (h *WatchlistHandler) Add(w http.ResponseWriter, r *http.Request) {
	movieID, err := parseID(r, "movieId")
	if err != nil {
		helper.WriteError(w, &dto.ValidationError{Msg: "invalid movie id"})
		return
	}
	userID := middleware.GetUserID(r)
	if err := h.watchlistSvc.Add(userID, movieID); err != nil {
		helper.WriteError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusCreated, "movie added to watchlist", nil)
}

// DELETE /api/watchlist/:movieId
func (h *WatchlistHandler) Remove(w http.ResponseWriter, r *http.Request) {
	movieID, err := parseID(r, "movieId")
	if err != nil {
		helper.WriteError(w, &dto.ValidationError{Msg: "invalid movie id"})
		return
	}
	userID := middleware.GetUserID(r)
	if err := h.watchlistSvc.Remove(userID, movieID); err != nil {
		helper.WriteError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "movie removed from watchlist", nil)
}

// GET /api/watchlist
func (h *WatchlistHandler) GetMyWatchlist(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	items, err := h.watchlistSvc.GetByUser(userID, page, limit)
	if err != nil {
		helper.WriteError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "watchlist retrieved", items)
}
