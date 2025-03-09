// internal/service/user.go
package service

import (
	"context"

	"github.com/adorufus/imgupper/internal/model"
)

// UserService defines the user service interface
type UserService interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	GetByID(ctx context.Context, id int64) (model.User, error)
	GetAll(ctx context.Context) ([]model.User, error)
	Update(ctx context.Context, user model.User) (model.User, error)
	Delete(ctx context.Context, id int64) error
}

// userService implements UserService
type userService struct {
	deps Deps
}

// NewUserService creates a new UserService
func NewUserService(deps Deps) UserService {
	return &userService{
		deps: deps,
	}
}

// Create creates a new user
func (s *userService) Create(ctx context.Context, user model.User) (model.User, error) {
	// Validate user data
	if err := user.Validate(); err != nil {
		return model.User{}, err
	}

	// Create user in repository
	return s.deps.Repos.User.Create(ctx, user)
}

// GetByID gets a user by ID
func (s *userService) GetByID(ctx context.Context, id int64) (model.User, error) {
	return s.deps.Repos.User.GetByID(ctx, id)
}

// GetAll gets all users
func (s *userService) GetAll(ctx context.Context) ([]model.User, error) {
	return s.deps.Repos.User.GetAll(ctx)
}

// Update updates a user
func (s *userService) Update(ctx context.Context, user model.User) (model.User, error) {
	// Validate user data
	if err := user.Validate(); err != nil {
		return model.User{}, err
	}

	// Update user in repository
	return s.deps.Repos.User.Update(ctx, user)
}

// Delete deletes a user
func (s *userService) Delete(ctx context.Context, id int64) error {
	return s.deps.Repos.User.Delete(ctx, id)
}
