package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Richard-OOO/E-director/apps/backend/internal/domain"
	"gorm.io/gorm"
)

const (
	DefaultSystemPromptKey = "chapter_system"
	StylePromptKey         = "style_export"
)

const defaultSystemPromptContent = `You are E-Director, a professional film director and AI video script formatter.
Convert one novel chapter into scene-level YAML scripts using the fixed E-director YAML Schema.
Return only JSON that matches the requested outer response shape. Do not return markdown.
Do not invent new YAML schema fields; only fill the fixed YAML schema fields with content from the chapter.`

const defaultStylePromptContent = `你是 E-Director 的剧本生成提示词配置。请按照用户选择的导出风格调整场景 YAML 的表达方式，但不要改动固定 YAML Schema 字段。

风格配置示例：
🔘 剧本节奏：[极快（多动作） / 舒缓（保留环境描写）]
🔘 对话风格：[偏向口语化 / 保留原著文学性]
🔘 额外要求输入框：[“请重点刻画主角的心理挣扎，多加一些微表情的描写”]

用户可以自行编辑以上风格要求。生成时必须保持章节、场景、镜头、动作、台词、情绪和 ai_video_prompt 的结构完整。`

type PromptStore interface {
	ListPrompts(ctx context.Context, userID string) ([]domain.Prompt, error)
	GetPrompt(ctx context.Context, userID, promptID string) (domain.Prompt, error)
	GetPromptByKey(ctx context.Context, userID, key string) (domain.Prompt, error)
	CreatePrompt(ctx context.Context, prompt domain.Prompt) error
	UpdatePrompt(ctx context.Context, userID, promptID, name, content string) (domain.Prompt, error)
	DeletePrompt(ctx context.Context, userID, promptID string) error
}

type PromptService struct {
	store PromptStore
}

type CreatePromptInput struct {
	UserID  string
	Name    string
	Content string
}

type UpdatePromptInput struct {
	UserID   string
	PromptID string
	Name     string
	Content  string
}

func NewPromptService(store PromptStore) *PromptService {
	return &PromptService{store: store}
}

func (s *PromptService) ListPrompts(ctx context.Context, userID string) ([]domain.Prompt, error) {
	if err := s.ensureReady(userID); err != nil {
		return nil, err
	}
	if err := s.ensureDefaultPrompts(ctx, userID); err != nil {
		return nil, err
	}
	return s.store.ListPrompts(ctx, userID)
}

func (s *PromptService) GetPrompt(ctx context.Context, userID, promptID string) (domain.Prompt, error) {
	if err := s.ensureReady(userID); err != nil {
		return domain.Prompt{}, err
	}
	if strings.TrimSpace(promptID) == "" {
		return domain.Prompt{}, ErrInvalidInput
	}
	if err := s.ensureDefaultPrompts(ctx, userID); err != nil {
		return domain.Prompt{}, err
	}
	prompt, err := s.store.GetPrompt(ctx, userID, promptID)
	if err != nil {
		return domain.Prompt{}, mapNotFound(err)
	}
	return prompt, nil
}

func (s *PromptService) CreatePrompt(ctx context.Context, input CreatePromptInput) (domain.Prompt, error) {
	if err := s.ensureReady(input.UserID); err != nil {
		return domain.Prompt{}, err
	}
	name := strings.TrimSpace(input.Name)
	content := strings.TrimSpace(input.Content)
	if name == "" || content == "" {
		return domain.Prompt{}, ErrInvalidInput
	}
	now := time.Now().UTC()
	prompt := domain.Prompt{
		ID:             newID("prompt", input.UserID, name, now.Format(time.RFC3339Nano)),
		UserID:         input.UserID,
		Key:            newID("custom", input.UserID, name, now.Format(time.RFC3339Nano)),
		Name:           name,
		Content:        content,
		DefaultContent: "",
		IsDefault:      false,
		IsEditable:     true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := s.store.CreatePrompt(ctx, prompt); err != nil {
		return domain.Prompt{}, err
	}
	return prompt, nil
}

func (s *PromptService) UpdatePrompt(ctx context.Context, input UpdatePromptInput) (domain.Prompt, error) {
	if err := s.ensureReady(input.UserID); err != nil {
		return domain.Prompt{}, err
	}
	name := strings.TrimSpace(input.Name)
	content := strings.TrimSpace(input.Content)
	if strings.TrimSpace(input.PromptID) == "" || name == "" || content == "" {
		return domain.Prompt{}, ErrInvalidInput
	}
	prompt, err := s.store.GetPrompt(ctx, input.UserID, input.PromptID)
	if err != nil {
		return domain.Prompt{}, mapNotFound(err)
	}
	if !prompt.IsEditable {
		return domain.Prompt{}, ErrInvalidInput
	}
	updated, err := s.store.UpdatePrompt(ctx, input.UserID, input.PromptID, name, content)
	if err != nil {
		return domain.Prompt{}, mapNotFound(err)
	}
	return updated, nil
}

func (s *PromptService) DeletePrompt(ctx context.Context, userID, promptID string) error {
	if err := s.ensureReady(userID); err != nil {
		return err
	}
	if strings.TrimSpace(promptID) == "" {
		return ErrInvalidInput
	}
	prompt, err := s.store.GetPrompt(ctx, userID, promptID)
	if err != nil {
		return mapNotFound(err)
	}
	if prompt.IsDefault {
		return ErrInvalidInput
	}
	return mapNotFound(s.store.DeletePrompt(ctx, userID, promptID))
}

func (s *PromptService) ResetPrompt(ctx context.Context, userID, promptID string) (domain.Prompt, error) {
	if err := s.ensureReady(userID); err != nil {
		return domain.Prompt{}, err
	}
	prompt, err := s.store.GetPrompt(ctx, userID, promptID)
	if err != nil {
		return domain.Prompt{}, mapNotFound(err)
	}
	if strings.TrimSpace(prompt.DefaultContent) == "" {
		return domain.Prompt{}, ErrInvalidInput
	}
	updated, err := s.store.UpdatePrompt(ctx, userID, promptID, prompt.Name, prompt.DefaultContent)
	if err != nil {
		return domain.Prompt{}, mapNotFound(err)
	}
	return updated, nil
}

func (s *PromptService) GenerationPrompts(ctx context.Context, userID string) (systemPrompt string, stylePrompt string, err error) {
	if err := s.ensureReady(userID); err != nil {
		return "", "", err
	}
	if err := s.ensureDefaultPrompts(ctx, userID); err != nil {
		return "", "", err
	}
	system, err := s.store.GetPromptByKey(ctx, userID, DefaultSystemPromptKey)
	if err != nil {
		return "", "", mapNotFound(err)
	}
	style, err := s.store.GetPromptByKey(ctx, userID, StylePromptKey)
	if err != nil {
		return "", "", mapNotFound(err)
	}
	return system.Content, style.Content, nil
}

func (s *PromptService) ensureReady(userID string) error {
	if s == nil || s.store == nil {
		return ErrStorageUnavailable
	}
	if strings.TrimSpace(userID) == "" {
		return ErrInvalidInput
	}
	return nil
}

func (s *PromptService) ensureDefaultPrompts(ctx context.Context, userID string) error {
	defaults := []domain.Prompt{
		defaultPrompt(userID, DefaultSystemPromptKey, "默认剧本生成提示词", defaultSystemPromptContent, false),
		defaultPrompt(userID, StylePromptKey, "导出风格提示词", defaultStylePromptContent, true),
	}
	for _, prompt := range defaults {
		if _, err := s.store.GetPromptByKey(ctx, userID, prompt.Key); err == nil {
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err := s.store.CreatePrompt(ctx, prompt); err != nil {
			return err
		}
	}
	return nil
}

func defaultPrompt(userID, key, name, content string, editable bool) domain.Prompt {
	now := time.Now().UTC()
	return domain.Prompt{
		ID:             newID("prompt", userID, key),
		UserID:         userID,
		Key:            key,
		Name:           name,
		Content:        content,
		DefaultContent: content,
		IsDefault:      true,
		IsEditable:     editable,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func mapNotFound(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}
