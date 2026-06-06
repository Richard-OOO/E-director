package service

import "errors"

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid email or verification code")
	ErrEmailRegistered    = errors.New("email already registered")
	ErrStorageUnavailable = errors.New("storage unavailable")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrCodeCooldown       = errors.New("verification code sent too frequently")
	ErrCodeExpired        = errors.New("verification code expired")
	ErrTooManyAttempts    = errors.New("too many verification attempts")
	ErrMailUnavailable    = errors.New("mail service unavailable")
)
