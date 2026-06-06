package domain

import "time"

type GenerationJobStatus string

const (
	GenerationJobPending           GenerationJobStatus = "pending"
	GenerationJobSplitting         GenerationJobStatus = "splitting"
	GenerationJobChapterProcessing GenerationJobStatus = "chapter_processing"
	GenerationJobAssembling        GenerationJobStatus = "assembling"
	GenerationJobCompleted         GenerationJobStatus = "completed"
	GenerationJobFailed            GenerationJobStatus = "failed"
	GenerationJobCancelled         GenerationJobStatus = "cancelled"
)

type GenerationChapterStatus string

const (
	GenerationChapterPending    GenerationChapterStatus = "pending"
	GenerationChapterProcessing GenerationChapterStatus = "processing"
	GenerationChapterCompleted  GenerationChapterStatus = "completed"
	GenerationChapterFailed     GenerationChapterStatus = "failed"
	GenerationChapterSkipped    GenerationChapterStatus = "skipped"
)

type GenerationSceneStatus string

const (
	GenerationScenePending    GenerationSceneStatus = "pending"
	GenerationSceneProcessing GenerationSceneStatus = "processing"
	GenerationSceneCompleted  GenerationSceneStatus = "completed"
	GenerationSceneFailed     GenerationSceneStatus = "failed"
)

type GenerationProject struct {
	ID             string
	UserID         string
	Title          string
	Language       string
	SourceType     string
	SourceTextHash string
	SourceText     string
	ChapterCount   int
	CurrentJobID   string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type GenerationJob struct {
	ID                     string
	ProjectID              string
	UserID                 string
	Status                 GenerationJobStatus
	SchemaVersion          string
	TargetFormat           string
	SceneGranularity       string
	EnableDynamicSchema    bool
	EnableDesignReasons    bool
	EnableCameraDirections bool
	EnableDialogues        bool
	EnableEmotionTags      bool
	MaxConcurrentChapters  int
	MaxConcurrentScenes    int
	StreamMode             string
	TotalChapters          int
	CompletedChapters      int
	FailedChapters         int
	TotalScenes            int
	CompletedScenes        int
	FailedScenes           int
	ProgressPercent        int
	ErrorCode              string
	ErrorMessage           string
	Retryable              bool
	StartedAt              *time.Time
	FinishedAt             *time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type GenerationChapter struct {
	ID                  string
	JobID               string
	ProjectID           string
	ChapterID           string
	ChapterIndex        int
	Title               string
	Content             string
	ContentHash         string
	Summary             string
	Status              GenerationChapterStatus
	SceneCount          int
	CompletedSceneCount int
	FailedSceneCount    int
	DesignNoteYAML      string
	ErrorCode           string
	ErrorMessage        string
	Retryable           bool
	StartedAt           *time.Time
	FinishedAt          *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type GenerationScene struct {
	ID                   string
	JobID                string
	ProjectID            string
	ChapterRecordID      string
	ChapterID            string
	SceneID              string
	SceneIndex           int
	Title                string
	Summary              string
	Status               GenerationSceneStatus
	GeneratedYAML        string
	EditableYAML         string
	DesignReasonYAML     string
	YAMLHash             string
	IsEdited             bool
	EditedByUserID       string
	EditedAt             *time.Time
	GenerationStartedAt  *time.Time
	GenerationFinishedAt *time.Time
	ErrorCode            string
	ErrorMessage         string
	Retryable            bool
	RetryCount           int
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type GenerationError struct {
	Code      string
	Message   string
	Retryable bool
	ChapterID string
	SceneID   string
}

func IsTerminalJobStatus(status GenerationJobStatus) bool {
	switch status {
	case GenerationJobCompleted, GenerationJobFailed, GenerationJobCancelled:
		return true
	default:
		return false
	}
}

func IsTerminalChapterStatus(status GenerationChapterStatus) bool {
	switch status {
	case GenerationChapterCompleted, GenerationChapterFailed, GenerationChapterSkipped:
		return true
	default:
		return false
	}
}

func IsTerminalSceneStatus(status GenerationSceneStatus) bool {
	switch status {
	case GenerationSceneCompleted, GenerationSceneFailed:
		return true
	default:
		return false
	}
}
