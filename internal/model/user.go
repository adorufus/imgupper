package model

import (
	"errors"
	"regexp"
	"time"
)

// User represents a user entity
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never expose password in JSON responses
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate validates user data
func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}

	if len(u.Name) < 3 {
		return errors.New("name must be at least 3 characters")
	}

	if u.Email == "" {
		return errors.New("email is required")
	}

	// Simple email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}

	return nil
}
