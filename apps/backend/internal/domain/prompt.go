package domain

import "time"

type Prompt struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Key            string    `json:"key"`
	Name           string    `json:"name"`
	Content        string    `json:"content"`
	DefaultContent string    `json:"default_content"`
	IsDefault      bool      `json:"is_default"`
	IsEditable     bool      `json:"is_editable"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
