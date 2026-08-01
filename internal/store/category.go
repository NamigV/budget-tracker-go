package store

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/NamigV/budget-tracker-go/internal/domain"
)

func (s *Store) CreateCategory(ctx context.Context, c *domain.Category) error {
	if err := s.db.WithContext(ctx).Create(c).Error; err != nil {
		return fmt.Errorf("create category: %w", err)
	}
	return nil
}

func (s *Store) ListCategories(ctx context.Context, userID int64) ([]domain.Category, error) {
	var categories []domain.Category
	if err := s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("name").
		Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	return categories, nil
}

func (s *Store) GetCategory(ctx context.Context, userID, categoryID int64) (domain.Category, error) {
	var c domain.Category
	err := s.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", categoryID, userID).
		First(&c).Error

	switch {
	case err == nil:
		return c, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return domain.Category{}, domain.ErrCategoryNotFound
	default:
		return domain.Category{}, fmt.Errorf("get category: %w", err)
	}
}
