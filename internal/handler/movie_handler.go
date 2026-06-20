package handler

import (
	"net/http"
	"strconv"

	"letter-square-api/internal/dto"
	"letter-square-api/internal/helper"
	"letter-square-api/internal/service"
)

type MovieHandler struct {
	movieSvc service.MovieService
}

func NewMovieHandler(movieSvc service.MovieService) *MovieHandler {
	return &MovieHandler{movieSvc: movieSvc}
}

// POST /api/movies
func (h *MovieHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateMovieRequest
	if err := helper.DecodeJSON(r, &req); err != nil {
		helper.WriteError(w, &dto.ValidationError{Msg: "invalid request body"})
		return
	}
	if err := req.Validate(); err != nil {
		helper.WriteError(w, err)
		return
	}
	movie, err := h.movieSvc.Create(&req)
	if err != nil {
		helper.WriteError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusCreated, "movie created", movie)
}

// GET /api/movies
func (h *MovieHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	search := q.Get("search")
	genre := q.Get("genre")
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	movies, err := h.movieSvc.GetAll(search, genre, page, limit)
	if err != nil {
		helper.WriteError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "movies retrieved", movies)
}

// GET /api/movies/:id
func (h *MovieHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		helper.WriteError(w, &dto.ValidationError{Msg: "invalid movie id"})
		return
	}
	movie, err := h.movieSvc.GetByID(id)
	if err != nil {
		helper.WriteError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "movie retrieved", movie)
}

// PUT /api/movies/:id
func (h *MovieHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		helper.WriteError(w, &dto.ValidationError{Msg: "invalid movie id"})
		return
	}
	var req dto.UpdateMovieRequest
	if err := helper.DecodeJSON(r, &req); err != nil {
		helper.WriteError(w, &dto.ValidationError{Msg: "invalid request body"})
		return
	}
	movie, err := h.movieSvc.Update(id, &req)
	if err != nil {
		helper.WriteError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "movie updated", movie)
}

// DELETE /api/movies/:id
func (h *MovieHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		helper.WriteError(w, &dto.ValidationError{Msg: "invalid movie id"})
		return
	}
	if err := h.movieSvc.Delete(id); err != nil {
		helper.WriteError(w, err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "movie deleted", nil)
}
