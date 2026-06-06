package mysql

import "time"

type GenerationProject struct {
	ID             string `gorm:"primaryKey;size:64"`
	UserID         string `gorm:"index;size:64;not null"`
	Title          string `gorm:"size:255;not null"`
	Language       string `gorm:"size:32;not null"`
	SourceType     string `gorm:"size:32;not null"`
	SourceTextHash string `gorm:"index;size:64"`
	SourceText     string `gorm:"type:longtext"`
	ChapterCount   int
	CurrentJobID   string `gorm:"index;size:64"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
