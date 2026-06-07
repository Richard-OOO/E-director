package mysql

import "time"

type GenerationChapter struct {
	ID                  string     `gorm:"primaryKey;size:64" json:"id"`
	JobID               string     `gorm:"index:idx_generation_chapters_job_status;uniqueIndex:idx_generation_chapters_job_chapter;uniqueIndex:idx_generation_chapters_job_index;size:64;not null" json:"job_id"`
	ProjectID           string     `gorm:"index;size:64;not null" json:"project_id"`
	ChapterID           string     `gorm:"uniqueIndex:idx_generation_chapters_job_chapter;size:64;not null" json:"chapter_id"`
	ChapterIndex        int        `gorm:"uniqueIndex:idx_generation_chapters_job_index;not null" json:"chapter_index"`
	Title               string     `gorm:"size:255" json:"title"`
	Content             string     `gorm:"type:longtext" json:"content"`
	ContentHash         string     `gorm:"size:64" json:"content_hash"`
	Summary             string     `gorm:"type:text" json:"summary"`
	Status              string     `gorm:"index:idx_generation_chapters_job_status;size:32;not null" json:"status"`
	SceneCount          int        `json:"scene_count"`
	CompletedSceneCount int        `json:"completed_scene_count"`
	FailedSceneCount    int        `json:"failed_scene_count"`
	DesignNoteYAML      string     `gorm:"type:longtext" json:"design_note_yaml"`
	ErrorCode           string     `gorm:"size:128" json:"error_code"`
	ErrorMessage        string     `gorm:"type:text" json:"error_message"`
	Retryable           bool       `json:"retryable"`
	StartedAt           *time.Time `json:"started_at"`
	FinishedAt          *time.Time `json:"finished_at"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}
