package mysql

import (
	"context"
	"errors"

	"github.com/Richard-OOO/E-director/apps/gateway/internal/domain"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(ctx context.Context, user domain.User) error {
	return s.db.WithContext(ctx).Create(&User{
		ID:           user.ID,
		Email:        user.Email,
		DisplayName:  user.DisplayName,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}).Error
}

func (s *UserStore) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var user User
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, ErrNotFound
		}
		return domain.User{}, err
	}
	return toDomainUser(user), nil
}

func (s *UserStore) FindByID(ctx context.Context, id string) (domain.User, error) {
	var user User
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, ErrNotFound
		}
		return domain.User{}, err
	}
	return toDomainUser(user), nil
}

func toDomainUser(user User) domain.User {
	return domain.User{
		ID:           user.ID,
		Email:        user.Email,
		DisplayName:  user.DisplayName,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
