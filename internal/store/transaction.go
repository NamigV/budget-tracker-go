package store

import (
	"context"
	"fmt"

	"github.com/NamigV/budget-tracker-go/internal/domain"
)

func (s *Store) CreateTransaction(ctx context.Context, t *domain.Transaction) error {
	if err := s.db.WithContext(ctx).Create(t).Error; err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}
	return nil
}

func (s *Store) ListTransactions(ctx context.Context, userID int64) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	if err := s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("occurred_on DESC, id DESC").
		Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("list transactions: %w", err)
	}
	return transactions, nil
}
