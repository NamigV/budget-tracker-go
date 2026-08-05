package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/NamigV/budget-tracker-go/internal/domain"
)

const sessionTTL = 7 * 24 * time.Hour

func (s *Service) Register(ctx context.Context, email, password string) (domain.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || !strings.Contains(email, "@") {
		return domain.User{}, fmt.Errorf("%w: invalid email", domain.ErrInvalidCredentials)
	}
	if len(password) < 8 {
		return domain.User{}, fmt.Errorf("%w: password must be at least 8 characters", domain.ErrInvalidCredentials)
	}

	switch _, err := s.repo.GetUserByEmail(ctx, email); {
	case err == nil:
		return domain.User{}, domain.ErrEmailTaken
	case errors.Is(err, domain.ErrUserNotFound):
	default:
		return domain.User{}, fmt.Errorf("check email: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, fmt.Errorf("hash password: %w", err)
	}

	u := &domain.User{Email: email, PasswordHash: string(hash)}
	if err := s.repo.CreateUser(ctx, u); err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}
	return *u, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (domain.Session, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	u, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.Session{}, domain.ErrInvalidCredentials
		}
		return domain.Session{}, fmt.Errorf("get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return domain.Session{}, domain.ErrInvalidCredentials
	}

	token, err := randomToken()
	if err != nil {
		return domain.Session{}, fmt.Errorf("generate token: %w", err)
	}

	if err := s.sessions.Save(ctx, token, u.ID, sessionTTL); err != nil {
		return domain.Session{}, fmt.Errorf("save session: %w", err)
	}

	return domain.Session{
		Token:     token,
		UserID:    u.ID,
		ExpiresAt: time.Now().Add(sessionTTL),
	}, nil
}

func (s *Service) Authenticate(ctx context.Context, token string) (int64, error) {
	return s.sessions.UserID(ctx, token)
}

func (s *Service) Logout(ctx context.Context, token string) error {
	return s.sessions.Delete(ctx, token)
}

func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
