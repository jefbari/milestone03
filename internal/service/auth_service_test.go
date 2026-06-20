package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"letter-square-api/internal/apperror"
	"letter-square-api/internal/dto"
	"letter-square-api/internal/entity"
	"letter-square-api/internal/helper"
	"letter-square-api/internal/service"
)

// Mock

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) Create(u *entity.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *mockUserRepo) FindByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockUserRepo) FindByID(id int64) (*entity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockUserRepo) FindByUsername(username string) (*entity.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

// Tests
func TestRegister_Success(t *testing.T) {
	repo := new(mockUserRepo)
	svc := service.NewAuthService(repo, "secret", 24)

	req := &dto.RegisterRequest{
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "password123",
	}

	repo.On("FindByEmail", req.Email).Return(nil, nil)
	repo.On("FindByUsername", req.Username).Return(nil, nil)
	repo.On("Create", mock.AnythingOfType("*entity.User")).Return(nil)

	user, token, err := svc.Register(req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, req.Username, user.Username)
	assert.Equal(t, req.Email, user.Email)
	repo.AssertExpectations(t)
}

func TestRegister_EmailTaken(t *testing.T) {
	repo := new(mockUserRepo)
	svc := service.NewAuthService(repo, "secret", 24)

	req := &dto.RegisterRequest{
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "password123",
	}

	existingUser := &entity.User{ID: 1, Email: req.Email}
	repo.On("FindByEmail", req.Email).Return(existingUser, nil)

	user, token, err := svc.Register(req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	assert.Equal(t, apperror.ErrEmailTaken, err)
	repo.AssertExpectations(t)
}

func TestRegister_UsernameTaken(t *testing.T) {
	repo := new(mockUserRepo)
	svc := service.NewAuthService(repo, "secret", 24)

	req := &dto.RegisterRequest{
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "password123",
	}

	existingUser := &entity.User{ID: 2, Username: req.Username}
	repo.On("FindByEmail", req.Email).Return(nil, nil)
	repo.On("FindByUsername", req.Username).Return(existingUser, nil)

	user, token, err := svc.Register(req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	assert.Equal(t, apperror.ErrUsernameTaken, err)
	repo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	repo := new(mockUserRepo)
	svc := service.NewAuthService(repo, "secret", 24)

	password := "password123"
	req := &dto.LoginRequest{
		Email:    "john@example.com",
		Password: password,
	}

	user := &entity.User{
		ID:       1,
		Username: "johndoe",
		Email:    req.Email,
		Password: helper.HashPassword(password),
	}
	repo.On("FindByEmail", req.Email).Return(user, nil)

	loggedUser, token, err := svc.Login(req)

	assert.NoError(t, err)
	assert.NotNil(t, loggedUser)
	assert.NotEmpty(t, token)
	assert.Equal(t, user.ID, loggedUser.ID)
	repo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	repo := new(mockUserRepo)
	svc := service.NewAuthService(repo, "secret", 24)

	req := &dto.LoginRequest{
		Email:    "john@example.com",
		Password: "wrongpassword",
	}

	user := &entity.User{
		ID:       1,
		Username: "johndoe",
		Email:    req.Email,
		Password: helper.HashPassword("correctpassword"),
	}
	repo.On("FindByEmail", req.Email).Return(user, nil)

	loggedUser, token, err := svc.Login(req)

	assert.Error(t, err)
	assert.Nil(t, loggedUser)
	assert.Empty(t, token)
	assert.Equal(t, apperror.ErrInvalidCredential, err)
	repo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := new(mockUserRepo)
	svc := service.NewAuthService(repo, "secret", 24)

	req := &dto.LoginRequest{
		Email:    "notexist@example.com",
		Password: "password123",
	}

	repo.On("FindByEmail", req.Email).Return(nil, nil)

	loggedUser, token, err := svc.Login(req)

	assert.Error(t, err)
	assert.Nil(t, loggedUser)
	assert.Empty(t, token)
	assert.Equal(t, apperror.ErrInvalidCredential, err)
	repo.AssertExpectations(t)
}
