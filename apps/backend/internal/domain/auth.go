package domain

import "time"

type User struct {
	ID           string
	Email        string
	DisplayName  string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type VerificationCode struct {
	Email     string
	Code      string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type AuthResult struct {
	SessionToken string
	User         User
}
