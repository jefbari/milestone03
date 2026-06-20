package handler

import (
	"net/http"

	"letter-square-api/internal/dto"
	"letter-square-api/internal/helper"
	"letter-square-api/internal/service"
)

type AuthHandler struct {
	authSvc service.AuthService
}

func NewAuthHandler(authSvc service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// POST /api/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := helper.DecodeJSON(r, &req); err != nil {
		helper.WriteError(w, &dto.ValidationError{Msg: "invalid request body"})
		return
	}
	if err := req.Validate(); err != nil {
		helper.WriteError(w, err)
		return
	}

	user, token, err := h.authSvc.Register(&req)
	if err != nil {
		helper.WriteError(w, err)
		return
	}

	helper.WriteSuccess(w, http.StatusCreated, "register successful", dto.AuthResponse{
		Token: token,
		User:  user,
	})
}

// POST /api/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := helper.DecodeJSON(r, &req); err != nil {
		helper.WriteError(w, &dto.ValidationError{Msg: "invalid request body"})
		return
	}
	if err := req.Validate(); err != nil {
		helper.WriteError(w, err)
		return
	}

	user, token, err := h.authSvc.Login(&req)
	if err != nil {
		helper.WriteError(w, err)
		return
	}

	helper.WriteSuccess(w, http.StatusOK, "login successful", dto.AuthResponse{
		Token: token,
		User:  user,
	})
}
