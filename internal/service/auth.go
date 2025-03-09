package service

import (
	"context"
	"errors"
	"time"

	"github.com/adorufus/imgupper/internal/model"
	"github.com/adorufus/imgupper/pkg/middleware"
)

type AuthService interface {
	Register(ctx context.Context, req model.RegisterRequest) (model.AuthResponse, error)
	Login(ctx context.Context, req model.LoginRequest) (model.AuthResponse, error)
}

type authService struct {
	deps          Deps
	jwtConfig     middleware.JWTConfig
	tokenDuration time.Duration
}

func NewAuthService(deps Deps, jwtSecret string, tokenDuration time.Duration) AuthService {
	return &authService{
		deps: deps,
		jwtConfig: middleware.JWTConfig{
			Secret:         jwtSecret,
			ExpirationTime: tokenDuration,
		},
		tokenDuration: tokenDuration,
	}
}

// Login implements AuthService.
func (s *authService) Login(ctx context.Context, req model.LoginRequest) (model.AuthResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return model.AuthResponse{}, err
	}

	// Get user by email
	user, err := s.deps.Repos.User.GetByEmail(ctx, req.Email)
	if err != nil {
		s.deps.Logger.Error("Failed to get user by email", "error", err, "email", req.Email)
		return model.AuthResponse{}, errors.New("invalid email or password")
	}

	// Check password
	if !model.CheckPassword(req.Password, user.Password) {
		return model.AuthResponse{}, errors.New("invalid email or password")
	}

	// Generate token
	token, err := middleware.GenerateToken(user.ID, user.Email, s.jwtConfig)
	if err != nil {
		s.deps.Logger.Error("Failed to generate token", "error", err)
		return model.AuthResponse{}, errors.New("failed to generate auth token")
	}

	expiresAt := time.Now().Add(s.tokenDuration)

	return model.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user,
	}, nil
}

// Register implements AuthService.
func (s *authService) Register(ctx context.Context, req model.RegisterRequest) (model.AuthResponse, error) {
	if err := req.Validate(); err != nil {
		return model.AuthResponse{}, err
	}

	exists, err := s.deps.Repos.User.ExistsByEmail(ctx, req.Email)

	if err != nil {
		s.deps.Logger.Error("Failed to check user existence", "error", err)
		return model.AuthResponse{}, errors.New("internal error")
	}

	if exists {
		return model.AuthResponse{}, errors.New("user with this email already exists")
	}

	hashedPassword, err := model.HashPassword(req.Password)
	if err != nil {
		s.deps.Logger.Error("Failed to hash password", "error", err)
		return model.AuthResponse{}, errors.New("internal error")
	}

	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	createdUser, err := s.deps.Repos.User.Create(ctx, user)
	if err != nil {
		s.deps.Logger.Error("Failed to create user", "error", err)
		return model.AuthResponse{}, errors.New("Failed to create user")
	}

	token, err := middleware.GenerateToken(createdUser.ID, createdUser.Email, s.jwtConfig)
	if err != nil {
		s.deps.Logger.Error("Failed to generate token", "error", err)
		return model.AuthResponse{}, errors.New("failed to generate auth token")
	}

	expiresAt := time.Now().Add(s.tokenDuration)

	return model.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      createdUser,
	}, nil
}
