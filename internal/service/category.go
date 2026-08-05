package service

import (
	"context"
	"fmt"

	"github.com/NamigV/budget-tracker-go/internal/domain"
)

func (s *Service) CreateCategory(ctx context.Context, userID int64, name, categoryType string) (domain.Category, error) {
	if name == "" {
		return domain.Category{}, fmt.Errorf("%w: name is empty", domain.ErrInvalidCategory)
	}
	if categoryType != "income" && categoryType != "expense" {
		return domain.Category{}, fmt.Errorf("%w: type must be 'income' or 'expense'", domain.ErrInvalidCategory)
	}

	c := &domain.Category{
		UserID: userID,
		Name:   name,
		Type:   categoryType,
	}
	if err := s.repo.CreateCategory(ctx, c); err != nil {
		return domain.Category{}, fmt.Errorf("create category: %w", err)
	}
	return *c, nil
}

func (s *Service) ListCategories(ctx context.Context, userID int64) ([]domain.Category, error) {
	categories, err := s.repo.ListCategories(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	return categories, nil
}
