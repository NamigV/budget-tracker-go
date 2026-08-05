package store

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/NamigV/budget-tracker-go/internal/domain"
)

func (s *Store) CreateUser(ctx context.Context, u *domain.User) error {
	if err := s.db.WithContext(ctx).Create(u).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var u domain.User
	err := s.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	switch {
	case err == nil:
		return u, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return domain.User{}, domain.ErrUserNotFound
	default:
		return domain.User{}, fmt.Errorf("get user by email: %w", err)
	}
}
