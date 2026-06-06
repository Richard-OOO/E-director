package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/Richard-OOO/E-director/apps/backend/internal/domain"
	mysqlmodels "github.com/Richard-OOO/E-director/apps/backend/internal/models/mysql"
)

type GenerationStore interface {
	CreateProject(ctx context.Context, project domain.GenerationProject) error
	CreateJob(ctx context.Context, job domain.GenerationJob) error
	UpdateProjectCurrentJob(ctx context.Context, projectID, jobID string) error
	CreateChapters(ctx context.Context, chapters []domain.GenerationChapter) error
	UpdateJobStatus(ctx context.Context, jobID string, status domain.GenerationJobStatus) error
	UpdateChapterResult(ctx context.Context, chapter domain.GenerationChapter) error
	CreateScenes(ctx context.Context, scenes []domain.GenerationScene) error
	GetProjectSnapshot(ctx context.Context, userID, projectID string) (mysqlmodels.ProjectSnapshot, error)
	ListProjects(ctx context.Context, userID string) ([]mysqlmodels.ProjectListItem, error)
}

type GenerationService struct {
	store         GenerationStore
	splitter      ChapterSplitter
	promptBuilder ChapterPromptBuilder
	generator     ChapterGenerator
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

func NewGenerationService(store GenerationStore, splitter ChapterSplitter, promptBuilder ChapterPromptBuilder, generator ChapterGenerator) *GenerationService {
	return &GenerationService{store: store, splitter: splitter, promptBuilder: promptBuilder, generator: generator}
}

func (s *GenerationService) CreateProjectAndStart(ctx context.Context, input CreateProjectInput) (CreateProjectResult, error) {
	if s.store == nil || s.generator == nil {
		return CreateProjectResult{}, ErrStorageUnavailable
	}
	if strings.TrimSpace(input.UserID) == "" || strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.Content) == "" {
		return CreateProjectResult{}, ErrInvalidInput
	}
	now := time.Now().UTC()
	projectID := newID("project", input.UserID, input.Title, input.Content)
	jobID := newID("job", projectID, input.Title)
	project := domain.GenerationProject{
		ID:           projectID,
		UserID:       input.UserID,
		Title:        strings.TrimSpace(input.Title),
		Language:     defaultString(input.Language, "zh-CN"),
		SourceType:   defaultString(input.SourceType, "plain_text"),
		SourceText:   input.Content,
		ChapterCount: 0,
		CurrentJobID: jobID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	job := domain.GenerationJob{
		ID:                     jobID,
		ProjectID:              projectID,
		UserID:                 input.UserID,
		Status:                 domain.GenerationJobSplitting,
		SchemaVersion:          "1.0",
		TargetFormat:           "yaml",
		SceneGranularity:       "medium",
		EnableDynamicSchema:    true,
		EnableDesignReasons:    true,
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
	_ = s.store.UpdateJobStatus(ctx, jobID, domain.GenerationJobChapterProcessing)

	previousContext := ChapterCarryContext{}
	for i, chapter := range chapters {
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
				EnableDesignReasons:    true,
				EnableCameraDirections: true,
				EnableDialogues:        true,
				EnableEmotionTags:      true,
			},
		})
		result, err := s.generator.GenerateChapter(ctx, prompt)
		if err != nil {
			chapter.Status = domain.GenerationChapterFailed
			chapter.ErrorMessage = err.Error()
			chapter.UpdatedAt = time.Now().UTC()
			_ = s.store.UpdateChapterResult(ctx, chapter)
			_ = s.store.UpdateJobStatus(ctx, jobID, domain.GenerationJobFailed)
			return CreateProjectResult{}, err
		}
		chapter.Status = domain.GenerationChapterCompleted
		chapter.Summary = result.ChapterSummary
		chapter.DesignNoteYAML = renderDesignNoteYAML(result.ChapterSchemaDesignNote)
		chapter.SceneCount = len(result.Scenes)
		chapter.CompletedSceneCount = len(result.Scenes)
		chapter.FinishedAt = ptrTime(time.Now().UTC())
		chapter.UpdatedAt = time.Now().UTC()
		if err := s.store.UpdateChapterResult(ctx, chapter); err != nil {
			return CreateProjectResult{}, err
		}
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
				DesignReasonYAML: renderDesignReasonsYAML(scene.DesignReasons),
				YAMLHash:         hashText(scene.YAMLContent),
				CreatedAt:        now,
				UpdatedAt:        now,
			})
		}
		if err := s.store.CreateScenes(ctx, scenes); err != nil {
			return CreateProjectResult{}, err
		}
		previousContext = result.CarryContext
		_ = i
	}
	_ = s.store.UpdateJobStatus(ctx, jobID, domain.GenerationJobCompleted)
	return CreateProjectResult{ProjectID: projectID, JobID: jobID, Status: domain.GenerationJobCompleted}, nil
}

func (s *GenerationService) GetProject(ctx context.Context, userID, projectID string) (mysqlmodels.ProjectSnapshot, error) {
	if s.store == nil {
		return mysqlmodels.ProjectSnapshot{}, ErrStorageUnavailable
	}
	return s.store.GetProjectSnapshot(ctx, userID, projectID)
}

func (s *GenerationService) ListProjects(ctx context.Context, userID string) ([]mysqlmodels.ProjectListItem, error) {
	if s.store == nil {
		return nil, ErrStorageUnavailable
	}
	return s.store.ListProjects(ctx, userID)
}

func renderDesignNoteYAML(note ChapterSchemaDesignNote) string {
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
