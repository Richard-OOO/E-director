package mysql

import "time"

type GenerationJob struct {
	ID                     string     `gorm:"primaryKey;size:64" json:"id"`
	ProjectID              string     `gorm:"index;size:64;not null" json:"project_id"`
	UserID                 string     `gorm:"index;size:64;not null" json:"user_id"`
	Status                 string     `gorm:"index;size:32;not null" json:"status"`
	SchemaVersion          string     `gorm:"size:32" json:"schema_version"`
	TargetFormat           string     `gorm:"size:32" json:"target_format"`
	SceneGranularity       string     `gorm:"size:32" json:"scene_granularity"`
	EnableDynamicSchema    bool       `json:"enable_dynamic_schema"`
	EnableDesignReasons    bool       `json:"enable_design_reasons"`
	EnableCameraDirections bool       `json:"enable_camera_directions"`
	EnableDialogues        bool       `json:"enable_dialogues"`
	EnableEmotionTags      bool       `json:"enable_emotion_tags"`
	MaxConcurrentChapters  int        `json:"max_concurrent_chapters"`
	MaxConcurrentScenes    int        `json:"max_concurrent_scenes"`
	StreamMode             string     `gorm:"size:32" json:"stream_mode"`
	TotalChapters          int        `json:"total_chapters"`
	CompletedChapters      int        `json:"completed_chapters"`
	FailedChapters         int        `json:"failed_chapters"`
	TotalScenes            int        `json:"total_scenes"`
	CompletedScenes        int        `json:"completed_scenes"`
	FailedScenes           int        `json:"failed_scenes"`
	ProgressPercent        int        `json:"progress_percent"`
	ErrorCode              string     `gorm:"size:128" json:"error_code"`
	ErrorMessage           string     `gorm:"type:text" json:"error_message"`
	Retryable              bool       `json:"retryable"`
	StartedAt              *time.Time `json:"started_at"`
	FinishedAt             *time.Time `json:"finished_at"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}
