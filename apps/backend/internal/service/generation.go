package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Richard-OOO/E-director/apps/backend/internal/domain"
	mysqlmodels "github.com/Richard-OOO/E-director/apps/backend/internal/models/mysql"
	"gorm.io/gorm"
)

type GenerationStore interface {
	CreateProject(ctx context.Context, project domain.GenerationProject) error
	CreateJob(ctx context.Context, job domain.GenerationJob) error
	UpdateProjectCurrentJob(ctx context.Context, projectID, jobID string) error
	CreateChapters(ctx context.Context, chapters []domain.GenerationChapter) error
	UpdateJobStatus(ctx context.Context, jobID string, status domain.GenerationJobStatus) error
	UpdateChapterResult(ctx context.Context, chapter domain.GenerationChapter) error
	CreateScenes(ctx context.Context, scenes []domain.GenerationScene) error
	UpdateSceneYAML(ctx context.Context, userID, projectID, sceneID, yaml string) error
	DeleteProject(ctx context.Context, userID, projectID string) error
	GetProjectSnapshot(ctx context.Context, userID, projectID string) (mysqlmodels.ProjectSnapshot, error)
	ListProjects(ctx context.Context, userID string) ([]mysqlmodels.ProjectListItem, error)
}

type GenerationService struct {
	store         GenerationStore
	splitter      ChapterSplitter
	promptBuilder ChapterPromptBuilder
	generator     ChapterGenerator
	events        *GenerationEventBus
}

type CreateProjectInput struct {
	UserID     string
	Title      string
	Language   string
	SourceType string
	Content    string
}

type CreateProjectResult struct {
	ProjectID string
	JobID     string
	Status    domain.GenerationJobStatus
}

type UpdateSceneYAMLInput struct {
	UserID    string
	ProjectID string
	SceneID   string
	YAML      string
}

func NewGenerationService(store GenerationStore, splitter ChapterSplitter, promptBuilder ChapterPromptBuilder, generator ChapterGenerator, events *GenerationEventBus) *GenerationService {
	return &GenerationService{store: store, splitter: splitter, promptBuilder: promptBuilder, generator: generator, events: events}
}

func (s *GenerationService) CreateProjectAndStart(ctx context.Context, input CreateProjectInput) (CreateProjectResult, error) {
	if s.store == nil || s.generator == nil {
		return CreateProjectResult{}, ErrStorageUnavailable
	}
	if strings.TrimSpace(input.UserID) == "" || strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.Content) == "" {
		return CreateProjectResult{}, ErrInvalidInput
	}
	now := time.Now().UTC()
	projectID := newID("project", input.UserID, input.Title, input.Content, now.Format(time.RFC3339Nano))
	jobID := newID("job", projectID, input.Title, now.Format(time.RFC3339Nano))
	project := domain.GenerationProject{
		ID:             projectID,
		UserID:         input.UserID,
		Title:          strings.TrimSpace(input.Title),
		Language:       defaultString(input.Language, "zh-CN"),
		SourceType:     defaultString(input.SourceType, "plain_text"),
		SourceTextHash: hashText(input.Content),
		SourceText:     input.Content,
		ChapterCount:   0,
		CurrentJobID:   jobID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	job := domain.GenerationJob{
		ID:                     jobID,
		ProjectID:              projectID,
		UserID:                 input.UserID,
		Status:                 domain.GenerationJobSplitting,
		SchemaVersion:          "1.0",
		TargetFormat:           "yaml",
		SceneGranularity:       "medium",
		EnableDynamicSchema:    false,
		EnableDesignReasons:    false,
		EnableCameraDirections: true,
		EnableDialogues:        true,
		EnableEmotionTags:      true,
		MaxConcurrentChapters:  1,
		MaxConcurrentScenes:    1,
		StreamMode:             "scene_delta",
		CreatedAt:              now,
		UpdatedAt:              now,
	}
	if err := s.store.CreateProject(ctx, project); err != nil {
		return CreateProjectResult{}, err
	}
	if err := s.store.CreateJob(ctx, job); err != nil {
		return CreateProjectResult{}, err
	}
	if err := s.store.UpdateProjectCurrentJob(ctx, projectID, jobID); err != nil {
		return CreateProjectResult{}, err
	}

	chapters, err := s.splitter.Split(ChapterSplitInput{ProjectID: projectID, JobID: jobID, Title: project.Title, Language: project.Language, Content: input.Content})
	if err != nil {
		return CreateProjectResult{}, err
	}
	if err := s.store.CreateChapters(ctx, chapters); err != nil {
		return CreateProjectResult{}, err
	}
	go s.runProjectGeneration(context.Background(), project, jobID, chapters)
	return CreateProjectResult{ProjectID: projectID, JobID: jobID, Status: domain.GenerationJobChapterProcessing}, nil
}

func (s *GenerationService) runProjectGeneration(ctx context.Context, project domain.GenerationProject, jobID string, chapters []domain.GenerationChapter) {
	projectID := project.ID
	_ = s.store.UpdateJobStatus(ctx, jobID, domain.GenerationJobChapterProcessing)
	s.publishGenerationEvent(domain.GenerationEventStarted, projectID, jobID, 0, map[string]any{
		"project_id":    projectID,
		"job_id":        jobID,
		"chapter_count": len(chapters),
	})

	previousContext := ChapterCarryContext{}
	for i, chapter := range chapters {
		chapterProgress := domain.GenerationProgressPayload{
			CompletedChapters: i,
			TotalChapters:     len(chapters),
			OverallProgress:   progressPercent(i, len(chapters)),
		}
		s.publishChapterEvent(domain.GenerationEventChapterStarted, projectID, jobID, i+1, chapter, chapterProgress)

		prompt := s.promptBuilder.BuildChapterPrompt(ChapterPromptInput{
			NovelTitle:      project.Title,
			Language:        project.Language,
			ChapterID:       chapter.ChapterID,
			ChapterTitle:    chapter.Title,
			ChapterIndex:    chapter.ChapterIndex,
			ChapterContent:  chapter.Content,
			PreviousContext: previousContext,
			Config: ChapterPromptConfig{
				TargetFormat:           "yaml",
				SceneGranularity:       "medium",
				EnableDesignReasons:    false,
				EnableCameraDirections: true,
				EnableDialogues:        true,
				EnableEmotionTags:      true,
			},
		})
		result, err := s.generator.GenerateChapter(ctx, prompt)
		if err != nil {
			fmt.Printf("[e-director:generation] GenerateChapter failed project_id=%s job_id=%s chapter_id=%s err=%v\n", projectID, jobID, chapter.ChapterID, err)
			chapter.Status = domain.GenerationChapterFailed
			chapter.ErrorMessage = err.Error()
			chapter.UpdatedAt = time.Now().UTC()
			_ = s.store.UpdateChapterResult(ctx, chapter)
			_ = s.store.UpdateJobStatus(ctx, jobID, domain.GenerationJobFailed)
			s.publishFailureEvent(projectID, jobID, chapter, err, i, len(chapters))
			s.closeEventStream(projectID, jobID)
			return
		}
		chapter.Status = domain.GenerationChapterCompleted
		chapter.Summary = result.ChapterSummary
		chapter.DesignNoteYAML = renderDesignNoteYAML(result.ChapterSchemaDesignNote)
		chapter.SceneCount = len(result.Scenes)
		chapter.CompletedSceneCount = len(result.Scenes)
		chapter.FinishedAt = ptrTime(time.Now().UTC())
		chapter.UpdatedAt = time.Now().UTC()
		if err := s.store.UpdateChapterResult(ctx, chapter); err != nil {
			_ = s.store.UpdateJobStatus(ctx, jobID, domain.GenerationJobFailed)
			s.publishFailureEvent(projectID, jobID, chapter, err, i, len(chapters))
			s.closeEventStream(projectID, jobID)
			return
		}
		now := time.Now().UTC()
		scenes := make([]domain.GenerationScene, 0, len(result.Scenes))
		for _, scene := range result.Scenes {
			scenes = append(scenes, domain.GenerationScene{
				ID:               newID("scene", jobID, chapter.ChapterID, scene.SceneID),
				JobID:            jobID,
				ProjectID:        projectID,
				ChapterRecordID:  chapter.ID,
				ChapterID:        chapter.ChapterID,
				SceneID:          scene.SceneID,
				SceneIndex:       scene.SceneIndex,
				Title:            scene.Title,
				Summary:          scene.Summary,
				Status:           domain.GenerationSceneCompleted,
				GeneratedYAML:    scene.YAMLContent,
				EditableYAML:     scene.YAMLContent,
				DesignReasonYAML: "",
				YAMLHash:         hashText(scene.YAMLContent),
				CreatedAt:        now,
				UpdatedAt:        now,
			})
		}
		if err := s.store.CreateScenes(ctx, scenes); err != nil {
			_ = s.store.UpdateJobStatus(ctx, jobID, domain.GenerationJobFailed)
			s.publishFailureEvent(projectID, jobID, chapter, err, i, len(chapters))
			s.closeEventStream(projectID, jobID)
			return
		}
		previousContext = result.CarryContext
		s.publishSchemaSummaryEvent(projectID, jobID, chapter, result, i+1)
		s.publishChapterCompletedEvent(projectID, jobID, chapter, result, i+1, len(chapters))
	}
	_ = s.store.UpdateJobStatus(ctx, jobID, domain.GenerationJobCompleted)
	s.publishGenerationCompletedEvent(projectID, jobID, len(chapters))
	s.closeEventStream(projectID, jobID)
}

func (s *GenerationService) GetProject(ctx context.Context, userID, projectID string) (mysqlmodels.ProjectSnapshot, error) {
	if s.store == nil {
		return mysqlmodels.ProjectSnapshot{}, ErrStorageUnavailable
	}
	return s.store.GetProjectSnapshot(ctx, userID, projectID)
}

func (s *GenerationService) UpdateSceneYAML(ctx context.Context, input UpdateSceneYAMLInput) error {
	if s.store == nil {
		return ErrStorageUnavailable
	}
	if strings.TrimSpace(input.UserID) == "" || strings.TrimSpace(input.ProjectID) == "" || strings.TrimSpace(input.SceneID) == "" || strings.TrimSpace(input.YAML) == "" {
		return ErrInvalidInput
	}
	if err := s.store.UpdateSceneYAML(ctx, input.UserID, input.ProjectID, input.SceneID, input.YAML); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *GenerationService) DeleteProject(ctx context.Context, userID, projectID string) error {
	if s.store == nil {
		return ErrStorageUnavailable
	}
	if strings.TrimSpace(userID) == "" || strings.TrimSpace(projectID) == "" {
		return ErrInvalidInput
	}
	if err := s.store.DeleteProject(ctx, userID, projectID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *GenerationService) ListProjects(ctx context.Context, userID string) ([]mysqlmodels.ProjectListItem, error) {
	if s.store == nil {
		return nil, ErrStorageUnavailable
	}
	return s.store.ListProjects(ctx, userID)
}

func (s *GenerationService) publishGenerationEvent(eventType domain.GenerationEventType, projectID, jobID string, sequence int, payload any) {
	if s.events == nil {
		return
	}
	s.events.Publish(domain.GenerationEvent{EventID: newID("event", projectID, jobID, fmt.Sprintf("%d", sequence), string(eventType)), EventType: eventType, ProjectID: projectID, JobID: jobID, Sequence: sequence, CreatedAt: time.Now().UTC(), Payload: payload})
}

func (s *GenerationService) publishChapterEvent(eventType domain.GenerationEventType, projectID, jobID string, sequence int, chapter domain.GenerationChapter, progress domain.GenerationProgressPayload) {
	if s.events == nil {
		return
	}
	payload := domain.ChapterStatusPayload{ChapterID: chapter.ChapterID, ChapterTitle: chapter.Title, ChapterIndex: chapter.ChapterIndex, Status: chapter.Status, Progress: progress}
	s.events.Publish(domain.GenerationEvent{EventID: newID("event", projectID, jobID, chapter.ChapterID, string(eventType)), EventType: eventType, ProjectID: projectID, JobID: jobID, Sequence: sequence, CreatedAt: time.Now().UTC(), Payload: payload})
}

func (s *GenerationService) publishChapterCompletedEvent(projectID, jobID string, chapter domain.GenerationChapter, result ChapterGenerationResult, sequence, total int) {
	if s.events == nil {
		return
	}
	progress := domain.GenerationProgressPayload{
		CompletedChapters: sequence,
		TotalChapters:     total,
		CompletedScenes:   len(result.Scenes),
		TotalScenes:       len(result.Scenes),
		OverallProgress:   progressPercent(sequence, total),
	}
	payload := domain.ChapterCompletedPayload{
		ChapterStatusPayload: domain.ChapterStatusPayload{
			ChapterID:    chapter.ChapterID,
			ChapterTitle: chapter.Title,
			ChapterIndex: chapter.ChapterIndex,
			Status:       domain.GenerationChapterCompleted,
			Progress:     progress,
		},
		ChapterSummary:          result.ChapterSummary,
		ChapterSchemaDesignNote: schemaDesignNotePayload(result.ChapterSchemaDesignNote),
		CarryContextSummary:     result.CarryContext.PreviousChapterSummary,
		Scenes:                  sceneYAMLPayloads(result.Scenes),
	}
	s.events.Publish(domain.GenerationEvent{EventID: newID("event", projectID, jobID, chapter.ChapterID, "completed"), EventType: domain.GenerationEventChapterCompleted, ProjectID: projectID, JobID: jobID, Sequence: sequence, CreatedAt: time.Now().UTC(), Payload: payload})
}

func (s *GenerationService) publishSchemaSummaryEvent(projectID, jobID string, chapter domain.GenerationChapter, result ChapterGenerationResult, sequence int) {
	if s.events == nil {
		return
	}
	payload := map[string]any{
		"chapter_id":                 chapter.ChapterID,
		"chapter_title":              chapter.Title,
		"chapter_index":              chapter.ChapterIndex,
		"chapter_schema_design_note": schemaDesignNotePayload(result.ChapterSchemaDesignNote),
		"carry_context_summary":      result.CarryContext.PreviousChapterSummary,
	}
	s.events.Publish(domain.GenerationEvent{EventID: newID("event", projectID, jobID, chapter.ChapterID, "schema-summary"), EventType: domain.GenerationEventSchemaSummary, ProjectID: projectID, JobID: jobID, Sequence: sequence, CreatedAt: time.Now().UTC(), Payload: payload})
}

func schemaDesignNotePayload(note ChapterSchemaDesignNote) domain.SchemaDesignNotePayload {
	reasons := make([]domain.SchemaKeyReasonPayload, 0, len(note.KeyReasons))
	for _, reason := range note.KeyReasons {
		reasons = append(reasons, domain.SchemaKeyReasonPayload{FieldName: reason.FieldName, Reason: reason.Reason})
	}
	return domain.SchemaDesignNotePayload{Summary: note.Summary, KeyReasons: reasons}
}

func sceneYAMLPayloads(scenes []GeneratedScene) []domain.ChapterSceneYAMLPayload {
	payloads := make([]domain.ChapterSceneYAMLPayload, 0, len(scenes))
	for _, scene := range scenes {
		payloads = append(payloads, domain.ChapterSceneYAMLPayload{
			SceneID:     scene.SceneID,
			SceneIndex:  scene.SceneIndex,
			Title:       scene.Title,
			Summary:     scene.Summary,
			YAMLContent: scene.YAMLContent,
		})
	}
	return payloads
}

func (s *GenerationService) publishFailureEvent(projectID, jobID string, chapter domain.GenerationChapter, err error, sequence, total int) {
	if s.events == nil {
		return
	}
	payload := domain.ChapterStatusPayload{
		ChapterID:    chapter.ChapterID,
		ChapterTitle: chapter.Title,
		ChapterIndex: chapter.ChapterIndex,
		Status:       domain.GenerationChapterFailed,
		Progress: domain.GenerationProgressPayload{
			CompletedChapters: sequence,
			TotalChapters:     total,
			OverallProgress:   progressPercent(sequence, total),
		},
		Error: &domain.GenerationError{Message: err.Error(), Retryable: false},
	}
	s.events.Publish(domain.GenerationEvent{EventID: newID("event", projectID, jobID, chapter.ChapterID, "failed"), EventType: domain.GenerationEventFailed, ProjectID: projectID, JobID: jobID, Sequence: sequence, CreatedAt: time.Now().UTC(), Payload: payload})
}

func (s *GenerationService) publishGenerationCompletedEvent(projectID, jobID string, total int) {
	if s.events == nil {
		return
	}
	payload := domain.GenerationProgressPayload{CompletedChapters: total, TotalChapters: total, OverallProgress: 100}
	s.events.Publish(domain.GenerationEvent{EventID: newID("event", projectID, jobID, "completed"), EventType: domain.GenerationEventCompleted, ProjectID: projectID, JobID: jobID, Sequence: total, CreatedAt: time.Now().UTC(), Payload: payload})
}

func (s *GenerationService) closeEventStream(projectID, jobID string) {
	if s.events == nil {
		return
	}
	s.events.Close(projectID, jobID)
}

func progressPercent(completed, total int) float64 {
	if total <= 0 {
		return 0
	}
	return float64(completed) / float64(total) * 100
}

func renderDesignNoteYAML(note ChapterSchemaDesignNote) string {
	if strings.TrimSpace(note.Summary) == "" && len(note.KeyReasons) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("summary: \"")
	b.WriteString(strings.ReplaceAll(note.Summary, "\"", "'"))
	b.WriteString("\"\nkey_reasons:\n")
	for _, reason := range note.KeyReasons {
		b.WriteString("  - field_name: \"")
		b.WriteString(strings.ReplaceAll(reason.FieldName, "\"", "'"))
		b.WriteString("\"\n    reason: \"")
		b.WriteString(strings.ReplaceAll(reason.Reason, "\"", "'"))
		b.WriteString("\"\n")
	}
	return strings.TrimSpace(b.String())
}

func renderDesignReasonsYAML(reasons []domain.DesignReason) string {
	if len(reasons) == 0 {
		return "[]"
	}
	var b strings.Builder
	b.WriteString("- reasons:\n")
	for _, reason := range reasons {
		b.WriteString(fmt.Sprintf("  - reason_id: %s\n", reason.ReasonID))
		b.WriteString(fmt.Sprintf("    target_path: %s\n", reason.TargetPath))
		b.WriteString(fmt.Sprintf("    field_name: %s\n", reason.FieldName))
		b.WriteString(fmt.Sprintf("    reason_type: %s\n", reason.ReasonType))
		b.WriteString(fmt.Sprintf("    title: %s\n", reason.Title))
		b.WriteString(fmt.Sprintf("    description: %s\n", reason.Description))
		b.WriteString(fmt.Sprintf("    hover_text: %s\n", reason.HoverText))
	}
	return strings.TrimSpace(b.String())
}

func newID(parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, ":")))
	return parts[0] + "_" + hex.EncodeToString(sum[:8])
}

func ptrTime(t time.Time) *time.Time { return &t }
