package mysql

import "time"

type GenerationProject struct {
	ID             string    `gorm:"primaryKey;size:64" json:"id"`
	UserID         string    `gorm:"index;size:64;not null" json:"user_id"`
	Title          string    `gorm:"size:255;not null" json:"title"`
	Language       string    `gorm:"size:32;not null" json:"language"`
	SourceType     string    `gorm:"size:32;not null" json:"source_type"`
	SourceTextHash string    `gorm:"index;size:64" json:"source_text_hash"`
	SourceText     string    `gorm:"type:longtext" json:"source_text"`
	ChapterCount   int       `json:"chapter_count"`
	CurrentJobID   string    `gorm:"index;size:64" json:"current_job_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
