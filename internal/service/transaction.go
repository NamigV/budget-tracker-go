package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/NamigV/budget-tracker-go/internal/domain"
)

func (s *Service) CreateTransaction(ctx context.Context, userID, categoryID, amountCents int64, currency, comment string, occurredOn time.Time) (domain.Transaction, error) {
	if amountCents <= 0 {
		return domain.Transaction{}, fmt.Errorf("%w: amount_cents must be positive", domain.ErrInvalidTransaction)
	}

	if currency == "" {
		currency = "RUB"
	}
	currency = strings.ToUpper(currency)
	if len(currency) != 3 {
		return domain.Transaction{}, fmt.Errorf("%w: currency must be a 3-letter code", domain.ErrInvalidTransaction)
	}

	if occurredOn.IsZero() {
		occurredOn = time.Now()
	}

	if _, err := s.repo.GetCategory(ctx, userID, categoryID); err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			return domain.Transaction{}, fmt.Errorf("%w: category not found", domain.ErrInvalidTransaction)
		}
		return domain.Transaction{}, fmt.Errorf("verify category: %w", err)
	}

	var commentPtr *string
	if comment != "" {
		commentPtr = &comment
	}

	t := &domain.Transaction{
		UserID:      userID,
		CategoryID:  categoryID,
		AmountCents: amountCents,
		Currency:    currency,
		Comment:     commentPtr,
		OccurredOn:  occurredOn,
	}
	if err := s.repo.CreateTransaction(ctx, t); err != nil {
		return domain.Transaction{}, fmt.Errorf("create transaction: %w", err)
	}
	return *t, nil
}

func (s *Service) ListTransactions(ctx context.Context, userID int64) ([]domain.Transaction, error) {
	transactions, err := s.repo.ListTransactions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list transactions: %w", err)
	}
	return transactions, nil
}
