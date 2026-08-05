package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/NamigV/budget-tracker-go/internal/domain"
)

type createCategoryRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type categoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func toCategoryResponse(c domain.Category) categoryResponse {
	return categoryResponse{ID: c.ID, Name: c.Name, Type: c.Type}
}

func (h *Handler) createCategory(w http.ResponseWriter, r *http.Request) {
	var req createCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	c, err := h.categories.CreateCategory(r.Context(), userID(r), req.Name, req.Type)

	switch {
	case err == nil:
		writeJSON(w, http.StatusCreated, toCategoryResponse(c))
	case errors.Is(err, domain.ErrInvalidCategory):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		log.Printf("create category: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) listCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categories.ListCategories(r.Context(), userID(r))
	if err != nil {
		log.Printf("list categories: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := make([]categoryResponse, 0, len(categories))
	for _, c := range categories {
		resp = append(resp, toCategoryResponse(c))
	}
	writeJSON(w, http.StatusOK, resp)
}
