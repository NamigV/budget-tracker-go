package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailTaken         = errors.New("email already taken")
	ErrUserNotFound       = errors.New("user not found")
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
