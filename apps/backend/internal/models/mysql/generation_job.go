package mysql

import "time"

type GenerationJob struct {
	ID                     string `gorm:"primaryKey;size:64"`
	ProjectID              string `gorm:"index;size:64;not null"`
	UserID                 string `gorm:"index;size:64;not null"`
	Status                 string `gorm:"index;size:32;not null"`
	SchemaVersion          string `gorm:"size:32"`
	TargetFormat           string `gorm:"size:32"`
	SceneGranularity       string `gorm:"size:32"`
	EnableDynamicSchema    bool
	EnableDesignReasons    bool
	EnableCameraDirections bool
	EnableDialogues        bool
	EnableEmotionTags      bool
	MaxConcurrentChapters  int
	MaxConcurrentScenes    int
	StreamMode             string `gorm:"size:32"`
	TotalChapters          int
	CompletedChapters      int
	FailedChapters         int
	TotalScenes            int
	CompletedScenes        int
	FailedScenes           int
	ProgressPercent        int
	ErrorCode              string `gorm:"size:128"`
	ErrorMessage           string `gorm:"type:text"`
	Retryable              bool
	StartedAt              *time.Time
	FinishedAt             *time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
}
