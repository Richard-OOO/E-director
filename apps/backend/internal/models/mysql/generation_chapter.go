package mysql

import "time"

type GenerationChapter struct {
	ID                  string `gorm:"primaryKey;size:64"`
	JobID               string `gorm:"index:idx_generation_chapters_job_status;uniqueIndex:idx_generation_chapters_job_chapter;uniqueIndex:idx_generation_chapters_job_index;size:64;not null"`
	ProjectID           string `gorm:"index;size:64;not null"`
	ChapterID           string `gorm:"uniqueIndex:idx_generation_chapters_job_chapter;size:64;not null"`
	ChapterIndex        int    `gorm:"uniqueIndex:idx_generation_chapters_job_index;not null"`
	Title               string `gorm:"size:255"`
	Content             string `gorm:"type:longtext"`
	ContentHash         string `gorm:"size:64"`
	Summary             string `gorm:"type:text"`
	Status              string `gorm:"index:idx_generation_chapters_job_status;size:32;not null"`
	SceneCount          int
	CompletedSceneCount int
	FailedSceneCount    int
	DesignNoteYAML      string `gorm:"type:longtext"`
	ErrorCode           string `gorm:"size:128"`
	ErrorMessage        string `gorm:"type:text"`
	Retryable           bool
	StartedAt           *time.Time
	FinishedAt          *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
