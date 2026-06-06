package domain

import "time"

type ProjectListItem struct {
	ProjectID    string              `json:"project_id"`
	Title        string              `json:"title"`
	Language     string              `json:"language"`
	Status       GenerationJobStatus `json:"status"`
	ChapterCount int                 `json:"chapter_count"`
	SceneCount   int                 `json:"scene_count"`
	CurrentJobID string              `json:"current_job_id"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

type ProjectSnapshot struct {
	Project  GenerationProject `json:"project"`
	Job      GenerationJob     `json:"job"`
	Chapters []ChapterSnapshot `json:"chapters"`
}

type ChapterSnapshot struct {
	Chapter GenerationChapter `json:"chapter"`
	Scenes  []GenerationScene `json:"scenes"`
}
