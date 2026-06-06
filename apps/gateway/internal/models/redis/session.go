package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Session struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type SessionStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewSessionStore(client *redis.Client, ttl time.Duration) *SessionStore {
	return &SessionStore{client: client, ttl: ttl}
}

func (s *SessionStore) Save(ctx context.Context, token string, session Session) error {
	payload, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, s.key(token), payload, s.ttl).Err()
}

func (s *SessionStore) Get(ctx context.Context, token string) (Session, error) {
	value, err := s.client.Get(ctx, s.key(token)).Result()
	if err != nil {
		return Session{}, err
	}
	var session Session
	if err := json.Unmarshal([]byte(value), &session); err != nil {
		return Session{}, err
	}
	return session, nil
}

func (s *SessionStore) Delete(ctx context.Context, token string) error {
	return s.client.Del(ctx, s.key(token)).Err()
}

func (s *SessionStore) key(token string) string {
	return "auth:session:" + token
}
