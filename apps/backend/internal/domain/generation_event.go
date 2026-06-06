package domain

import "time"

type GenerationEventType string

const (
	GenerationEventStarted          GenerationEventType = "generation_started"
	GenerationEventChapterStarted   GenerationEventType = "chapter_started"
	GenerationEventSceneStarted     GenerationEventType = "scene_started"
	GenerationEventSceneDelta       GenerationEventType = "scene_delta"
	GenerationEventSceneCompleted   GenerationEventType = "scene_completed"
	GenerationEventChapterCompleted GenerationEventType = "chapter_completed"
	GenerationEventSchemaSummary    GenerationEventType = "schema_summary"
	GenerationEventCompleted        GenerationEventType = "generation_completed"
	GenerationEventFailed           GenerationEventType = "generation_failed"
)

type GenerationEvent struct {
	EventID   string
	EventType GenerationEventType
	ProjectID string
	JobID     string
	SessionID string
	Sequence  int
	CreatedAt time.Time
	Payload   any
}

type GenerationProgressPayload struct {
	CompletedChapters int     `json:"completed_chapters"`
	TotalChapters     int     `json:"total_chapters"`
	CompletedScenes   int     `json:"completed_scenes"`
	TotalScenes       int     `json:"total_scenes"`
	OverallProgress   float64 `json:"overall_progress"`
}

type ChapterStatusPayload struct {
	ChapterID    string                    `json:"chapter_id"`
	ChapterTitle string                    `json:"chapter_title"`
	ChapterIndex int                       `json:"chapter_index"`
	Status       GenerationChapterStatus   `json:"status"`
	Progress     GenerationProgressPayload `json:"progress"`
	Error        *GenerationError          `json:"error,omitempty"`
}

type SceneStatusPayload struct {
	ChapterID    string                    `json:"chapter_id"`
	ChapterTitle string                    `json:"chapter_title"`
	SceneID      string                    `json:"scene_id"`
	SceneIndex   int                       `json:"scene_index"`
	Status       GenerationSceneStatus     `json:"status"`
	Progress     GenerationProgressPayload `json:"progress"`
	Error        *GenerationError          `json:"error,omitempty"`
}

type SceneYAMLDeltaPayload struct {
	ChapterID     string                    `json:"chapter_id"`
	ChapterTitle  string                    `json:"chapter_title"`
	SceneID       string                    `json:"scene_id"`
	SceneIndex    int                       `json:"scene_index"`
	DeltaIndex    int                       `json:"delta_index"`
	YAMLContent   string                    `json:"yaml_content"`
	DesignReasons []DesignReason            `json:"design_reasons"`
	Progress      GenerationProgressPayload `json:"progress"`
}

type DesignReason struct {
	ReasonID    string  `json:"reason_id"`
	TargetPath  string  `json:"target_path"`
	FieldName   string  `json:"field_name"`
	ReasonType  string  `json:"reason_type"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	HoverText   string  `json:"hover_text"`
	Confidence  float64 `json:"confidence,omitempty"`
}
