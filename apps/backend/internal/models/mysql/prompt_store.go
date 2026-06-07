package mysql

import (
	"context"
	"time"

	"github.com/Richard-OOO/E-director/apps/backend/internal/domain"
	"gorm.io/gorm"
)

type PromptStore struct {
	db *gorm.DB
}

func NewPromptStore(db *gorm.DB) *PromptStore {
	return &PromptStore{db: db}
}

func (s *PromptStore) ListPrompts(ctx context.Context, userID string) ([]domain.Prompt, error) {
	var rows []Prompt
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Order("is_default DESC, created_at ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	prompts := make([]domain.Prompt, 0, len(rows))
	for _, row := range rows {
		prompts = append(prompts, promptToDomain(row))
	}
	return prompts, nil
}

func (s *PromptStore) GetPrompt(ctx context.Context, userID, promptID string) (domain.Prompt, error) {
	var row Prompt
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", promptID, userID).First(&row).Error; err != nil {
		return domain.Prompt{}, err
	}
	return promptToDomain(row), nil
}

func (s *PromptStore) GetPromptByKey(ctx context.Context, userID, key string) (domain.Prompt, error) {
	var row Prompt
	if err := s.db.WithContext(ctx).Where("user_id = ? AND `key` = ?", userID, key).First(&row).Error; err != nil {
		return domain.Prompt{}, err
	}
	return promptToDomain(row), nil
}

func (s *PromptStore) CreatePrompt(ctx context.Context, prompt domain.Prompt) error {
	return s.db.WithContext(ctx).Create(&Prompt{
		ID:             prompt.ID,
		UserID:         prompt.UserID,
		Key:            prompt.Key,
		Name:           prompt.Name,
		Content:        prompt.Content,
		DefaultContent: prompt.DefaultContent,
		IsDefault:      prompt.IsDefault,
		IsEditable:     prompt.IsEditable,
		CreatedAt:      prompt.CreatedAt,
		UpdatedAt:      prompt.UpdatedAt,
	}).Error
}

func (s *PromptStore) UpdatePrompt(ctx context.Context, userID, promptID, name, content string) (domain.Prompt, error) {
	now := time.Now().UTC()
	result := s.db.WithContext(ctx).Model(&Prompt{}).Where("id = ? AND user_id = ?", promptID, userID).Updates(map[string]any{
		"name":       name,
		"content":    content,
		"updated_at": now,
	})
	if result.Error != nil {
		return domain.Prompt{}, result.Error
	}
	if result.RowsAffected == 0 {
		return domain.Prompt{}, gorm.ErrRecordNotFound
	}
	return s.GetPrompt(ctx, userID, promptID)
}

func (s *PromptStore) DeletePrompt(ctx context.Context, userID, promptID string) error {
	result := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", promptID, userID).Delete(&Prompt{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func promptToDomain(row Prompt) domain.Prompt {
	return domain.Prompt{
		ID:             row.ID,
		UserID:         row.UserID,
		Key:            row.Key,
		Name:           row.Name,
		Content:        row.Content,
		DefaultContent: row.DefaultContent,
		IsDefault:      row.IsDefault,
		IsEditable:     row.IsEditable,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
	}
}
