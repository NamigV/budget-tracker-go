package service

import (
	"context"
	"time"

	"github.com/NamigV/budget-tracker-go/internal/domain"
)

type repository interface {
	CreateCategory(ctx context.Context, c *domain.Category) error
	ListCategories(ctx context.Context, userID int64) ([]domain.Category, error)
	GetCategory(ctx context.Context, userID, categoryID int64) (domain.Category, error)
	CreateTransaction(ctx context.Context, t *domain.Transaction) error
	ListTransactions(ctx context.Context, userID int64) ([]domain.Transaction, error)
	CreateUser(ctx context.Context, u *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}

type sessionStore interface {
	Save(ctx context.Context, token string, userID int64, ttl time.Duration) error
	UserID(ctx context.Context, token string) (int64, error)
	Delete(ctx context.Context, token string) error
}

type Service struct {
	repo     repository
	sessions sessionStore
}

func New(repo repository, sessions sessionStore) *Service {
	return &Service{repo: repo, sessions: sessions}
}
