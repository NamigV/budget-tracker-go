package http

import (
	"context"
	"net/http"
	"time"

	"github.com/NamigV/budget-tracker-go/internal/domain"
)

type categoryService interface {
	CreateCategory(ctx context.Context, userID int64, name, categoryType string) (domain.Category, error)
	ListCategories(ctx context.Context, userID int64) ([]domain.Category, error)
}

type transactionService interface {
	CreateTransaction(ctx context.Context, userID, categoryID, amountCents int64, currency, comment string, occurredOn time.Time) (domain.Transaction, error)
	ListTransactions(ctx context.Context, userID int64) ([]domain.Transaction, error)
}

type service interface {
	categoryService
	transactionService
	authService
}

type Handler struct {
	categories   categoryService
	transactions transactionService
	auth         authService
}

func New(svc service) *Handler {
	return &Handler{categories: svc, transactions: svc, auth: svc}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)

	mux.HandleFunc("POST /api/v1/register", h.register)
	mux.HandleFunc("POST /api/v1/login", h.login)
	mux.HandleFunc("POST /api/v1/logout", h.logout)

	mux.HandleFunc("POST /api/v1/categories", h.requireAuth(h.createCategory))
	mux.HandleFunc("GET /api/v1/categories", h.requireAuth(h.listCategories))
	mux.HandleFunc("POST /api/v1/transactions", h.requireAuth(h.createTransaction))
	mux.HandleFunc("GET /api/v1/transactions", h.requireAuth(h.listTransactions))

	return mux
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}
