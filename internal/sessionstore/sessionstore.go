package sessionstore

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/NamigV/budget-tracker-go/internal/config"
	"github.com/NamigV/budget-tracker-go/internal/domain"
)

type Store struct {
	client *redis.Client
}

func Connect(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	return client, nil
}

func New(client *redis.Client) *Store {
	return &Store{client: client}
}

func (s *Store) Save(ctx context.Context, token string, userID int64, ttl time.Duration) error {
	if err := s.client.Set(ctx, sessionKey(token), userID, ttl).Err(); err != nil {
		return fmt.Errorf("save session: %w", err)
	}
	return nil
}

func (s *Store) UserID(ctx context.Context, token string) (int64, error) {
	value, err := s.client.Get(ctx, sessionKey(token)).Result()
	switch {
	case err == nil:
		id, parseErr := strconv.ParseInt(value, 10, 64)
		if parseErr != nil {
			return 0, fmt.Errorf("parse session user id: %w", parseErr)
		}
		return id, nil
	case errors.Is(err, redis.Nil):
		return 0, domain.ErrSessionNotFound
	default:
		return 0, fmt.Errorf("get session: %w", err)
	}
}

func (s *Store) Delete(ctx context.Context, token string) error {
	if err := s.client.Del(ctx, sessionKey(token)).Err(); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func sessionKey(token string) string {
	return "session:" + token
}
