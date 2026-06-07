package mysql

import "time"

type GenerationScene struct {
	ID                   string     `gorm:"primaryKey;size:64" json:"id"`
	JobID                string     `gorm:"index:idx_generation_scenes_job_status;uniqueIndex:idx_generation_scenes_job_scene;uniqueIndex:idx_generation_scenes_job_chapter_scene_index;size:64;not null" json:"job_id"`
	ProjectID            string     `gorm:"index;size:64;not null" json:"project_id"`
	ChapterRecordID      string     `gorm:"index:idx_generation_scenes_chapter_record_index;size:64;not null" json:"chapter_record_id"`
	ChapterID            string     `gorm:"uniqueIndex:idx_generation_scenes_job_chapter_scene_index;size:64;not null" json:"chapter_id"`
	SceneID              string     `gorm:"uniqueIndex:idx_generation_scenes_job_scene;size:96;not null" json:"scene_id"`
	SceneIndex           int        `gorm:"index:idx_generation_scenes_chapter_record_index;uniqueIndex:idx_generation_scenes_job_chapter_scene_index;not null" json:"scene_index"`
	Title                string     `gorm:"size:255" json:"title"`
	Summary              string     `gorm:"type:text" json:"summary"`
	Status               string     `gorm:"index:idx_generation_scenes_job_status;size:32;not null" json:"status"`
	GeneratedYAML        string     `gorm:"type:longtext" json:"generated_yaml"`
	EditableYAML         string     `gorm:"type:longtext" json:"editable_yaml"`
	DesignReasonYAML     string     `gorm:"type:longtext" json:"design_reason_yaml"`
	YAMLHash             string     `gorm:"size:64" json:"yaml_hash"`
	IsEdited             bool       `json:"is_edited"`
	EditedByUserID       string     `gorm:"index;size:64" json:"edited_by_user_id"`
	EditedAt             *time.Time `json:"edited_at"`
	GenerationStartedAt  *time.Time `json:"generation_started_at"`
	GenerationFinishedAt *time.Time `json:"generation_finished_at"`
	ErrorCode            string     `gorm:"size:128" json:"error_code"`
	ErrorMessage         string     `gorm:"type:text" json:"error_message"`
	Retryable            bool       `json:"retryable"`
	RetryCount           int        `json:"retry_count"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}
