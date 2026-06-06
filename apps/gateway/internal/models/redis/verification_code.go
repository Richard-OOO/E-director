package redis

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type VerificationCodeStore struct {
	client     *redis.Client
	ttl        time.Duration
	cooldown   time.Duration
	attemptTTL time.Duration
}

func NewVerificationCodeStore(client *redis.Client, ttl time.Duration) *VerificationCodeStore {
	return NewVerificationCodeStoreWithLimits(client, ttl, time.Minute, 15*time.Minute)
}

func NewVerificationCodeStoreWithLimits(client *redis.Client, ttl, cooldown, attemptTTL time.Duration) *VerificationCodeStore {
	return &VerificationCodeStore{client: client, ttl: ttl, cooldown: cooldown, attemptTTL: attemptTTL}
}

func (s *VerificationCodeStore) Save(ctx context.Context, email, codeHash string) error {
	return s.SaveCode(ctx, email, codeHash)
}

func (s *VerificationCodeStore) Get(ctx context.Context, email string) (string, error) {
	return s.GetCodeHash(ctx, email)
}

func (s *VerificationCodeStore) SaveCode(ctx context.Context, email, codeHash string) error {
	return s.client.Set(ctx, s.codeKey(email), codeHash, s.ttl).Err()
}

func (s *VerificationCodeStore) GetCodeHash(ctx context.Context, email string) (string, error) {
	return s.client.Get(ctx, s.codeKey(email)).Result()
}

func (s *VerificationCodeStore) DeleteCode(ctx context.Context, email string) error {
	return s.client.Del(ctx, s.codeKey(email)).Err()
}

func (s *VerificationCodeStore) HasCooldown(ctx context.Context, email string) (bool, error) {
	err := s.client.Get(ctx, s.cooldownKey(email)).Err()
	if err == nil {
		return true, nil
	}
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	return false, err
}

func (s *VerificationCodeStore) SetCooldown(ctx context.Context, email string) error {
	return s.client.Set(ctx, s.cooldownKey(email), "1", s.cooldown).Err()
}

func (s *VerificationCodeStore) IncrementAttempts(ctx context.Context, email string) (int, error) {
	key := s.attemptsKey(email)
	count, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if count == 1 {
		if err := s.client.Expire(ctx, key, s.attemptTTL).Err(); err != nil {
			return 0, err
		}
	}
	return int(count), nil
}

func (s *VerificationCodeStore) GetAttempts(ctx context.Context, email string) (int, error) {
	value, err := s.client.Get(ctx, s.attemptsKey(email)).Result()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	count, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *VerificationCodeStore) ClearAttempts(ctx context.Context, email string) error {
	return s.client.Del(ctx, s.attemptsKey(email)).Err()
}

func (s *VerificationCodeStore) TTL() time.Duration {
	return s.ttl
}

func (s *VerificationCodeStore) Cooldown() time.Duration {
	return s.cooldown
}

func (s *VerificationCodeStore) codeKey(email string) string {
	return "auth:code:" + email
}

func (s *VerificationCodeStore) cooldownKey(email string) string {
	return "auth:code:cooldown:" + email
}

func (s *VerificationCodeStore) attemptsKey(email string) string {
	return "auth:code:attempts:" + email
}
