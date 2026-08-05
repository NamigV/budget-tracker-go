package domain

import (
	"errors"
	"time"
)

var ErrSessionNotFound = errors.New("session not found")

type Session struct {
	Token     string
	UserID    int64
	ExpiresAt time.Time
}
