package domain

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrInvalidCategory  = errors.New("invalid category")
)

type Category struct {
	ID        int64
	UserID    int64
	Name      string
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
