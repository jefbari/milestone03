package apperror

import "net/http"

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e *AppError) Error() string { return e.Message }

func New(code int, msg string) *AppError { return &AppError{Code: code, Message: msg} }

var (
	ErrBadRequest       = New(http.StatusBadRequest, "bad request")
	ErrUnauthorized     = New(http.StatusUnauthorized, "unauthorized")
	ErrForbidden        = New(http.StatusForbidden, "forbidden")
	ErrNotFound         = New(http.StatusNotFound, "not found")
	ErrConflict         = New(http.StatusConflict, "conflict")
	ErrInternal         = New(http.StatusInternalServerError, "internal server error")
	ErrEmailTaken       = New(http.StatusConflict, "email already registered")
	ErrUsernameTaken    = New(http.StatusConflict, "username already taken")
	ErrInvalidCredential = New(http.StatusUnauthorized, "invalid email or password")
	ErrMovieNotFound    = New(http.StatusNotFound, "movie not found")
	ErrReviewNotFound   = New(http.StatusNotFound, "review not found")
	ErrDuplicateReview  = New(http.StatusConflict, "you already reviewed this movie")
	ErrDuplicateWatchlist = New(http.StatusConflict, "movie already in watchlist")
	ErrRatingRange      = New(http.StatusBadRequest, "rating must be between 1.0 and 5.0")
)
