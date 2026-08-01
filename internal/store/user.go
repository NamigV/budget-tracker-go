package store

import (
	"context"
	"fmt"

	"github.com/NamigV/budget-tracker-go/internal/domain"
)

func (s *Store) CreateUser(ctx context.Context, u *domain.User) error {
	if err := s.db.WithContext(ctx).Create(u).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}
