package mysql

import (
	"context"
	"errors"
	"time"

	"github.com/Richard-OOO/E-director/apps/backend/internal/domain"
	"gorm.io/gorm"
)

type GenerationStore struct {
	db *gorm.DB
}

func NewGenerationStore(db *gorm.DB) *GenerationStore {
	return &GenerationStore{db: db}
}

type ProjectSnapshot struct {
	Project  GenerationProject
	Job      GenerationJob
	Chapters []GenerationChapter
	Scenes   []GenerationScene
}

type ProjectListItem struct {
	ProjectID    string
	Title        string
	Language     string
	Status       string
	ChapterCount int
	SceneCount   int
	CurrentJobID string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (s *GenerationStore) CreateProject(ctx context.Context, project domain.GenerationProject) error {
	return s.db.WithContext(ctx).Create(&GenerationProject{
		ID:             project.ID,
		UserID:         project.UserID,
		Title:          project.Title,
		Language:       project.Language,
		SourceType:     project.SourceType,
		SourceTextHash: project.SourceTextHash,
		SourceText:     project.SourceText,
		ChapterCount:   project.ChapterCount,
		CurrentJobID:   project.CurrentJobID,
		CreatedAt:      project.CreatedAt,
		UpdatedAt:      project.UpdatedAt,
	}).Error
}

func (s *GenerationStore) CreateJob(ctx context.Context, job domain.GenerationJob) error {
	return s.db.WithContext(ctx).Create(&GenerationJob{
		ID:                     job.ID,
		ProjectID:              job.ProjectID,
		UserID:                 job.UserID,
		Status:                 string(job.Status),
		SchemaVersion:          job.SchemaVersion,
		TargetFormat:           job.TargetFormat,
		SceneGranularity:       job.SceneGranularity,
		EnableDynamicSchema:    job.EnableDynamicSchema,
		EnableDesignReasons:    job.EnableDesignReasons,
		EnableCameraDirections: job.EnableCameraDirections,
		EnableDialogues:        job.EnableDialogues,
		EnableEmotionTags:      job.EnableEmotionTags,
		MaxConcurrentChapters:  job.MaxConcurrentChapters,
		MaxConcurrentScenes:    job.MaxConcurrentScenes,
		StreamMode:             job.StreamMode,
		TotalChapters:          job.TotalChapters,
		CompletedChapters:      job.CompletedChapters,
		FailedChapters:         job.FailedChapters,
		TotalScenes:            job.TotalScenes,
		CompletedScenes:        job.CompletedScenes,
		FailedScenes:           job.FailedScenes,
		ProgressPercent:        job.ProgressPercent,
		ErrorCode:              job.ErrorCode,
		ErrorMessage:           job.ErrorMessage,
		Retryable:              job.Retryable,
		StartedAt:              job.StartedAt,
		FinishedAt:             job.FinishedAt,
		CreatedAt:              job.CreatedAt,
		UpdatedAt:              job.UpdatedAt,
	}).Error
}

func (s *GenerationStore) UpdateProjectCurrentJob(ctx context.Context, projectID, jobID string) error {
	return s.db.WithContext(ctx).Model(&GenerationProject{}).Where("id = ?", projectID).Update("current_job_id", jobID).Error
}

func (s *GenerationStore) CreateChapters(ctx context.Context, chapters []domain.GenerationChapter) error {
	rows := make([]GenerationChapter, 0, len(chapters))
	for _, chapter := range chapters {
		rows = append(rows, GenerationChapter{
			ID:                  chapter.ID,
			JobID:               chapter.JobID,
			ProjectID:           chapter.ProjectID,
			ChapterID:           chapter.ChapterID,
			ChapterIndex:        chapter.ChapterIndex,
			Title:               chapter.Title,
			Content:             chapter.Content,
			ContentHash:         chapter.ContentHash,
			Summary:             chapter.Summary,
			Status:              string(chapter.Status),
			SceneCount:          chapter.SceneCount,
			CompletedSceneCount: chapter.CompletedSceneCount,
			FailedSceneCount:    chapter.FailedSceneCount,
			DesignNoteYAML:      chapter.DesignNoteYAML,
			ErrorCode:           chapter.ErrorCode,
			ErrorMessage:        chapter.ErrorMessage,
			Retryable:           chapter.Retryable,
			StartedAt:           chapter.StartedAt,
			FinishedAt:          chapter.FinishedAt,
			CreatedAt:           chapter.CreatedAt,
			UpdatedAt:           chapter.UpdatedAt,
		})
	}
	if len(rows) == 0 {
		return nil
	}
	return s.db.WithContext(ctx).Create(&rows).Error
}

func (s *GenerationStore) UpdateJobStatus(ctx context.Context, jobID string, status domain.GenerationJobStatus) error {
	return s.db.WithContext(ctx).Model(&GenerationJob{}).Where("id = ?", jobID).Update("status", string(status)).Error
}

func (s *GenerationStore) UpdateChapterResult(ctx context.Context, chapter domain.GenerationChapter) error {
	updates := map[string]any{
		"summary":               chapter.Summary,
		"status":                string(chapter.Status),
		"scene_count":           chapter.SceneCount,
		"completed_scene_count": chapter.CompletedSceneCount,
		"failed_scene_count":    chapter.FailedSceneCount,
		"design_note_yaml":      chapter.DesignNoteYAML,
		"error_code":            chapter.ErrorCode,
		"error_message":         chapter.ErrorMessage,
		"retryable":             chapter.Retryable,
		"started_at":            chapter.StartedAt,
		"finished_at":           chapter.FinishedAt,
		"updated_at":            time.Now().UTC(),
	}
	return s.db.WithContext(ctx).Model(&GenerationChapter{}).Where("job_id = ? AND chapter_id = ?", chapter.JobID, chapter.ChapterID).Updates(updates).Error
}

func (s *GenerationStore) CreateScenes(ctx context.Context, scenes []domain.GenerationScene) error {
	rows := make([]GenerationScene, 0, len(scenes))
	for _, scene := range scenes {
		rows = append(rows, GenerationScene{
			ID:                   scene.ID,
			JobID:                scene.JobID,
			ProjectID:            scene.ProjectID,
			ChapterRecordID:      scene.ChapterRecordID,
			ChapterID:            scene.ChapterID,
			SceneID:              scene.SceneID,
			SceneIndex:           scene.SceneIndex,
			Title:                scene.Title,
			Summary:              scene.Summary,
			Status:               string(scene.Status),
			GeneratedYAML:        scene.GeneratedYAML,
			EditableYAML:         scene.EditableYAML,
			DesignReasonYAML:     scene.DesignReasonYAML,
			YAMLHash:             scene.YAMLHash,
			IsEdited:             scene.IsEdited,
			EditedByUserID:       scene.EditedByUserID,
			EditedAt:             scene.EditedAt,
			GenerationStartedAt:  scene.GenerationStartedAt,
			GenerationFinishedAt: scene.GenerationFinishedAt,
			ErrorCode:            scene.ErrorCode,
			ErrorMessage:         scene.ErrorMessage,
			Retryable:            scene.Retryable,
			RetryCount:           scene.RetryCount,
			CreatedAt:            scene.CreatedAt,
			UpdatedAt:            scene.UpdatedAt,
		})
	}
	if len(rows) == 0 {
		return nil
	}
	return s.db.WithContext(ctx).Create(&rows).Error
}

func (s *GenerationStore) GetProjectSnapshot(ctx context.Context, userID, projectID string) (ProjectSnapshot, error) {
	var project GenerationProject
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", projectID, userID).First(&project).Error; err != nil {
		return ProjectSnapshot{}, err
	}
	var job GenerationJob
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", project.CurrentJobID, userID).First(&job).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return ProjectSnapshot{}, err
		}
	}
	var chapters []GenerationChapter
	if err := s.db.WithContext(ctx).Where("project_id = ?", projectID).Order("chapter_index ASC").Find(&chapters).Error; err != nil {
		return ProjectSnapshot{}, err
	}
	var scenes []GenerationScene
	if err := s.db.WithContext(ctx).Where("project_id = ?", projectID).Order("chapter_id ASC, scene_index ASC").Find(&scenes).Error; err != nil {
		return ProjectSnapshot{}, err
	}
	return ProjectSnapshot{Project: project, Job: job, Chapters: chapters, Scenes: scenes}, nil
}

func (s *GenerationStore) ListProjects(ctx context.Context, userID string) ([]ProjectListItem, error) {
	var rows []struct {
		ProjectID    string
		Title        string
		Language     string
		Status       string
		ChapterCount int
		SceneCount   int
		CurrentJobID string
		CreatedAt    time.Time
		UpdatedAt    time.Time
	}
	if err := s.db.WithContext(ctx).
		Table("generation_projects AS p").
		Select("p.id as project_id, p.title, p.language, COALESCE(j.status, '') as status, p.chapter_count, COUNT(s.id) as scene_count, p.current_job_id, p.created_at, p.updated_at").
		Joins("LEFT JOIN generation_jobs j ON j.id = p.current_job_id").
		Joins("LEFT JOIN generation_scenes s ON s.project_id = p.id").
		Where("p.user_id = ?", userID).
		Group("p.id, p.title, p.language, j.status, p.chapter_count, p.current_job_id, p.created_at, p.updated_at").
		Order("p.updated_at DESC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]ProjectListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, ProjectListItem(row))
	}
	return items, nil
}
