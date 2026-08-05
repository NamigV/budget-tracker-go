package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/NamigV/budget-tracker-go/internal/domain"
)

type createTransactionRequest struct {
	CategoryID  int64  `json:"category_id"`
	AmountCents int64  `json:"amount_cents"`
	Currency    string `json:"currency"`
	Comment     string `json:"comment"`
	OccurredOn  string `json:"occurred_on"`
}

type transactionResponse struct {
	ID          int64  `json:"id"`
	CategoryID  int64  `json:"category_id"`
	AmountCents int64  `json:"amount_cents"`
	Currency    string `json:"currency"`
	Comment     string `json:"comment"`
	OccurredOn  string `json:"occurred_on"`
}

func toTransactionResponse(t domain.Transaction) transactionResponse {
	comment := ""
	if t.Comment != nil {
		comment = *t.Comment
	}
	return transactionResponse{
		ID:          t.ID,
		CategoryID:  t.CategoryID,
		AmountCents: t.AmountCents,
		Currency:    t.Currency,
		Comment:     comment,
		OccurredOn:  t.OccurredOn.Format("2006-01-02"),
	}
}

func (h *Handler) createTransaction(w http.ResponseWriter, r *http.Request) {
	var req createTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	var occurredOn time.Time
	if req.OccurredOn != "" {
		parsed, err := time.Parse("2006-01-02", req.OccurredOn)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid occurred_on, use YYYY-MM-DD")
			return
		}
		occurredOn = parsed
	}

	t, err := h.transactions.CreateTransaction(r.Context(), userID(r), req.CategoryID, req.AmountCents, req.Currency, req.Comment, occurredOn)

	switch {
	case err == nil:
		writeJSON(w, http.StatusCreated, toTransactionResponse(t))
	case errors.Is(err, domain.ErrInvalidTransaction):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		log.Printf("create transaction: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) listTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.transactions.ListTransactions(r.Context(), userID(r))
	if err != nil {
		log.Printf("list transactions: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := make([]transactionResponse, 0, len(transactions))
	for _, t := range transactions {
		resp = append(resp, toTransactionResponse(t))
	}
	writeJSON(w, http.StatusOK, resp)
}
