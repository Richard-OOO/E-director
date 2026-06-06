package mysql

import "time"

type GenerationScene struct {
	ID                   string `gorm:"primaryKey;size:64"`
	JobID                string `gorm:"index:idx_generation_scenes_job_status;uniqueIndex:idx_generation_scenes_job_scene;uniqueIndex:idx_generation_scenes_job_chapter_scene_index;size:64;not null"`
	ProjectID            string `gorm:"index;size:64;not null"`
	ChapterRecordID      string `gorm:"index:idx_generation_scenes_chapter_record_index;size:64;not null"`
	ChapterID            string `gorm:"uniqueIndex:idx_generation_scenes_job_chapter_scene_index;size:64;not null"`
	SceneID              string `gorm:"uniqueIndex:idx_generation_scenes_job_scene;size:96;not null"`
	SceneIndex           int    `gorm:"index:idx_generation_scenes_chapter_record_index;uniqueIndex:idx_generation_scenes_job_chapter_scene_index;not null"`
	Title                string `gorm:"size:255"`
	Summary              string `gorm:"type:text"`
	Status               string `gorm:"index:idx_generation_scenes_job_status;size:32;not null"`
	GeneratedYAML        string `gorm:"type:longtext"`
	EditableYAML         string `gorm:"type:longtext"`
	DesignReasonYAML     string `gorm:"type:longtext"`
	YAMLHash             string `gorm:"size:64"`
	IsEdited             bool
	EditedByUserID       string `gorm:"index;size:64"`
	EditedAt             *time.Time
	GenerationStartedAt  *time.Time
	GenerationFinishedAt *time.Time
	ErrorCode            string `gorm:"size:128"`
	ErrorMessage         string `gorm:"type:text"`
	Retryable            bool
	RetryCount           int
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
