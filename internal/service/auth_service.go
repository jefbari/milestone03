package service

import (
	"letter-square-api/internal/apperror"
	"letter-square-api/internal/dto"
	"letter-square-api/internal/entity"
	"letter-square-api/internal/helper"
	"letter-square-api/internal/repository"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) (*entity.User, string, error)
	Login(req *dto.LoginRequest) (*entity.User, string, error)
}

type authService struct {
	userRepo      repository.UserRepository
	jwtSecret     string
	jwtExpiryHour int
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtExpiryHour int) AuthService {
	return &authService{
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		jwtExpiryHour: jwtExpiryHour,
	}
}

func (s *authService) Register(req *dto.RegisterRequest) (*entity.User, string, error) {
	// Check email uniqueness
	existing, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, "", apperror.ErrInternal
	}
	if existing != nil {
		return nil, "", apperror.ErrEmailTaken
	}

	// Check username uniqueness
	existingByUsername, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, "", apperror.ErrInternal
	}
	if existingByUsername != nil {
		return nil, "", apperror.ErrUsernameTaken
	}

	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: helper.HashPassword(req.Password),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", apperror.ErrInternal
	}

	token, err := helper.GenerateToken(user.ID, user.Username, s.jwtSecret, s.jwtExpiryHour)
	if err != nil {
		return nil, "", apperror.ErrInternal
	}

	return user, token, nil
}

func (s *authService) Login(req *dto.LoginRequest) (*entity.User, string, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, "", apperror.ErrInternal
	}
	if user == nil || !helper.CheckPassword(req.Password, user.Password) {
		return nil, "", apperror.ErrInvalidCredential
	}

	token, err := helper.GenerateToken(user.ID, user.Username, s.jwtSecret, s.jwtExpiryHour)
	if err != nil {
		return nil, "", apperror.ErrInternal
	}

	return user, token, nil
}
