package helper

import (
	"encoding/json"
	"net/http"

	"letter-square-api/internal/apperror"
	"letter-square-api/internal/dto"
)

func WriteJSON(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(body)
}

func WriteSuccess(w http.ResponseWriter, code int, msg string, data interface{}) {
	WriteJSON(w, code, dto.Success(msg, data))
}

func WriteError(w http.ResponseWriter, err error) {
	var appErr *apperror.AppError
	var code int
	var msg string

	switch e := err.(type) {
	case *apperror.AppError:
		appErr = e
		code = appErr.Code
		msg = appErr.Message
	case *dto.ValidationError:
		code = http.StatusBadRequest
		msg = e.Error()
	default:
		code = http.StatusInternalServerError
		msg = "internal server error"
	}

	WriteJSON(w, code, dto.Fail(msg))
}

func DecodeJSON(r *http.Request, dst interface{}) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
