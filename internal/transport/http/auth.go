package http

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/NamigV/budget-tracker-go/internal/domain"
)

type ctxKey int

const userIDKey ctxKey = 0

const sessionCookie = "session"

type authService interface {
	Register(ctx context.Context, email, password string) (domain.User, error)
	Login(ctx context.Context, email, password string) (domain.Session, error)
	Authenticate(ctx context.Context, token string) (int64, error)
	Logout(ctx context.Context, token string) error
}

type credentialsRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req credentialsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	u, err := h.auth.Register(r.Context(), req.Email, req.Password)
	switch {
	case err == nil:
		writeJSON(w, http.StatusCreated, map[string]any{"id": u.ID, "email": u.Email})
	case errors.Is(err, domain.ErrEmailTaken):
		writeError(w, http.StatusConflict, "email already taken")
	case errors.Is(err, domain.ErrInvalidCredentials):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		log.Printf("register: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req credentialsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	sess, err := h.auth.Login(r.Context(), req.Email, req.Password)
	switch {
	case err == nil:
		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookie,
			Value:    sess.Token,
			Path:     "/",
			Expires:  sess.ExpiresAt,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	case errors.Is(err, domain.ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, "invalid credentials")
	default:
		log.Printf("login: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(sessionCookie); err == nil {
		_ = h.auth.Logout(r.Context(), c.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(sessionCookie)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}
		uid, err := h.auth.Authenticate(r.Context(), c.Value)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}
		next(w, r.WithContext(context.WithValue(r.Context(), userIDKey, uid)))
	}
}
