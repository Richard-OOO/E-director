package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net/mail"
	"strings"
	"time"

	"github.com/Richard-OOO/E-director/apps/gateway/internal/domain"
	mysqlmodels "github.com/Richard-OOO/E-director/apps/gateway/internal/models/mysql"
	redismodels "github.com/Richard-OOO/E-director/apps/gateway/internal/models/redis"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	Create(ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByID(ctx context.Context, id string) (domain.User, error)
}

type SessionStore interface {
	Save(ctx context.Context, token string, session redismodels.Session) error
	Get(ctx context.Context, token string) (redismodels.Session, error)
	Delete(ctx context.Context, token string) error
}

type VerificationCodeStore interface {
	Save(ctx context.Context, email, codeHash string) error
	Get(ctx context.Context, email string) (string, error)
	SaveCode(ctx context.Context, email, codeHash string) error
	GetCodeHash(ctx context.Context, email string) (string, error)
	DeleteCode(ctx context.Context, email string) error
	HasCooldown(ctx context.Context, email string) (bool, error)
	SetCooldown(ctx context.Context, email string) error
	IncrementAttempts(ctx context.Context, email string) (int, error)
	GetAttempts(ctx context.Context, email string) (int, error)
	ClearAttempts(ctx context.Context, email string) error
	TTL() time.Duration
	Cooldown() time.Duration
}

type AuthService struct {
	users           UserStore
	sessions        SessionStore
	codes           VerificationCodeStore
	mailer          Mailer
	maxCodeAttempts int
}

type AuthInput struct {
	Email            string
	Password         string
	VerificationCode string
	DisplayName      string
}

type SendCodeResult struct {
	CooldownSeconds  int `json:"cooldown_seconds"`
	ExpiresInSeconds int `json:"expires_in_seconds"`
}

func NewAuthService(users UserStore, sessions SessionStore, codes VerificationCodeStore, options ...AuthOption) *AuthService {
	service := &AuthService{users: users, sessions: sessions, codes: codes, maxCodeAttempts: 5}
	for _, option := range options {
		option(service)
	}
	return service
}

type AuthOption func(*AuthService)

func WithMailer(mailer Mailer) AuthOption {
	return func(s *AuthService) {
		s.mailer = mailer
	}
}

func WithMaxCodeAttempts(maxAttempts int) AuthOption {
	return func(s *AuthService) {
		if maxAttempts > 0 {
			s.maxCodeAttempts = maxAttempts
		}
	}
}

func (s *AuthService) SendCode(ctx context.Context, email string) (SendCodeResult, error) {
	email = normalizeEmail(email)
	if !isValidEmail(email) {
		return SendCodeResult{}, ErrInvalidInput
	}
	if s.codes == nil {
		return SendCodeResult{}, ErrStorageUnavailable
	}
	coolingDown, err := s.codes.HasCooldown(ctx, email)
	if err != nil {
		return SendCodeResult{}, err
	}
	if coolingDown {
		return SendCodeResult{}, ErrCodeCooldown
	}
	code, err := newVerificationCode()
	if err != nil {
		return SendCodeResult{}, err
	}
	if err := s.codes.SaveCode(ctx, email, hashVerificationCode(email, code)); err != nil {
		return SendCodeResult{}, err
	}
	if s.mailer == nil {
		_ = s.codes.DeleteCode(ctx, email)
		return SendCodeResult{}, ErrMailUnavailable
	}
	if err := s.mailer.SendLoginCode(ctx, email, code, s.codes.TTL()); err != nil {
		_ = s.codes.DeleteCode(ctx, email)
		return SendCodeResult{}, ErrMailUnavailable
	}
	if err := s.codes.SetCooldown(ctx, email); err != nil {
		return SendCodeResult{}, err
	}
	return SendCodeResult{CooldownSeconds: int(s.codes.Cooldown().Seconds()), ExpiresInSeconds: int(s.codes.TTL().Seconds())}, nil
}

func (s *AuthService) Register(ctx context.Context, input AuthInput) (domain.AuthResult, error) {
	input = normalizeAuthInput(input)
	if !isValidEmail(input.Email) || input.Password == "" || input.VerificationCode == "" {
		return domain.AuthResult{}, ErrInvalidInput
	}
	if s.users == nil || s.codes == nil {
		return domain.AuthResult{}, ErrStorageUnavailable
	}
	if err := s.verifyRegistrationCode(ctx, input.Email, input.VerificationCode); err != nil {
		return domain.AuthResult{}, err
	}
	if _, err := s.users.FindByEmail(ctx, input.Email); err == nil {
		return domain.AuthResult{}, ErrEmailRegistered
	} else if !errors.Is(err, mysqlmodels.ErrNotFound) {
		return domain.AuthResult{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.AuthResult{}, err
	}
	now := time.Now().UTC()
	user := domain.User{
		ID:           newUserID(input.Email, now),
		Email:        input.Email,
		DisplayName:  defaultDisplayName(input),
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.users.Create(ctx, user); err != nil {
		return domain.AuthResult{}, err
	}
	_ = s.codes.DeleteCode(ctx, input.Email)
	_ = s.codes.ClearAttempts(ctx, input.Email)
	return s.createAuthResult(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, input AuthInput) (domain.AuthResult, error) {
	input = normalizeAuthInput(input)
	if !isValidEmail(input.Email) || input.Password == "" {
		return domain.AuthResult{}, ErrInvalidInput
	}
	if s.users == nil {
		return domain.AuthResult{}, ErrStorageUnavailable
	}
	user, err := s.users.FindByEmail(ctx, input.Email)
	if errors.Is(err, mysqlmodels.ErrNotFound) {
		return domain.AuthResult{}, ErrInvalidCredentials
	}
	if err != nil {
		return domain.AuthResult{}, err
	}
	if user.PasswordHash == "" || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)) != nil {
		return domain.AuthResult{}, ErrInvalidCredentials
	}
	return s.createAuthResult(ctx, user)
}

func (s *AuthService) Me(ctx context.Context, token string) (domain.User, error) {
	token = strings.TrimSpace(token)
	if token == "" || s.sessions == nil || s.users == nil {
		return domain.User{}, ErrUnauthorized
	}
	session, err := s.sessions.Get(ctx, token)
	if err != nil {
		return domain.User{}, ErrUnauthorized
	}
	return s.users.FindByID(ctx, session.UserID)
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	token = strings.TrimSpace(token)
	if token == "" || s.sessions == nil {
		return ErrUnauthorized
	}
	return s.sessions.Delete(ctx, token)
}

func (s *AuthService) verifyRegistrationCode(ctx context.Context, email, code string) error {
	attempts, err := s.codes.GetAttempts(ctx, email)
	if err != nil {
		return err
	}
	if attempts >= s.maxCodeAttempts {
		return ErrTooManyAttempts
	}
	storedHash, err := s.codes.GetCodeHash(ctx, email)
	if errors.Is(err, redis.Nil) {
		return ErrCodeExpired
	}
	if err != nil {
		return err
	}
	if !verifyCodeHash(email, code, storedHash) {
		attempts, err := s.codes.IncrementAttempts(ctx, email)
		if err != nil {
			return err
		}
		if attempts >= s.maxCodeAttempts {
			return ErrTooManyAttempts
		}
		return ErrInvalidCredentials
	}
	return nil
}

func (s *AuthService) createAuthResult(ctx context.Context, user domain.User) (domain.AuthResult, error) {
	token := newToken(user.Email)
	if s.sessions != nil {
		if err := s.sessions.Save(ctx, token, redismodels.Session{UserID: user.ID, Email: user.Email}); err != nil {
			return domain.AuthResult{}, err
		}
	}
	return domain.AuthResult{SessionToken: token, User: user}, nil
}

func normalizeAuthInput(input AuthInput) AuthInput {
	input.Email = normalizeEmail(input.Email)
	input.Password = strings.TrimSpace(input.Password)
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	input.VerificationCode = strings.TrimSpace(input.VerificationCode)
	return input
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func isValidEmail(email string) bool {
	address, err := mail.ParseAddress(email)
	return err == nil && address.Address == email
}

func defaultDisplayName(input AuthInput) string {
	if input.DisplayName != "" {
		return input.DisplayName
	}
	if before, _, ok := strings.Cut(input.Email, "@"); ok && before != "" {
		return before
	}
	return input.Email
}

func newUserID(seed string, now time.Time) string {
	sum := sha256.Sum256([]byte(seed + now.Format(time.RFC3339Nano)))
	return "user_" + hex.EncodeToString(sum[:8])
}

func newToken(seed string) string {
	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
	sum := sha256.Sum256(append(buf, []byte(seed)...))
	return hex.EncodeToString(sum[:])
}

func newVerificationCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func hashVerificationCode(email, code string) string {
	sum := sha256.Sum256([]byte(normalizeEmail(email) + ":" + strings.TrimSpace(code)))
	return hex.EncodeToString(sum[:])
}

func verifyCodeHash(email, code, storedHash string) bool {
	expected := hashVerificationCode(email, code)
	return subtle.ConstantTimeCompare([]byte(expected), []byte(storedHash)) == 1
}
