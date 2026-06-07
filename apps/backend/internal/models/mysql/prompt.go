package mysql

import "time"

type Prompt struct {
	ID             string    `gorm:"primaryKey;size:64" json:"id"`
	UserID         string    `gorm:"uniqueIndex:idx_prompts_user_key;index;size:64;not null" json:"user_id"`
	Key            string    `gorm:"uniqueIndex:idx_prompts_user_key;size:96;not null" json:"key"`
	Name           string    `gorm:"size:255;not null" json:"name"`
	Content        string    `gorm:"type:longtext" json:"content"`
	DefaultContent string    `gorm:"type:longtext" json:"default_content"`
	IsDefault      bool      `gorm:"index" json:"is_default"`
	IsEditable     bool      `json:"is_editable"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
