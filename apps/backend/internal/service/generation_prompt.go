package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Richard-OOO/E-director/apps/backend/internal/domain"
)

type ChapterPromptBuilder struct {
	BaseSystemPrompt string
	SchemaVersion    string
}

type ChapterPrompt struct {
	SystemPrompt string
	UserPrompt   string
	OutputSchema string
}

type ChapterPromptInput struct {
	NovelTitle      string
	Language        string
	ChapterID       string
	ChapterTitle    string
	ChapterIndex    int
	ChapterContent  string
	PreviousContext ChapterCarryContext
	Config          ChapterPromptConfig
}

type ChapterPromptConfig struct {
	TargetFormat           string
	SceneGranularity       string
	EnableDesignReasons    bool
	EnableCameraDirections bool
	EnableDialogues        bool
	EnableEmotionTags      bool
}

type ChapterCarryContext struct {
	PreviousChapterID      string           `json:"previous_chapter_id"`
	PreviousChapterTitle   string           `json:"previous_chapter_title"`
	PreviousChapterSummary string           `json:"previous_chapter_summary"`
	CharacterTimeline      []CharacterState `json:"character_timeline"`
	TimeAndPlaceNotes      []string         `json:"time_and_place_notes"`
	ContinuityNotes        []string         `json:"continuity_notes"`
	SchemaDesignSummary    string           `json:"schema_design_summary"`
	ReusableYAMLFields     []string         `json:"reusable_yaml_fields"`
}

type CharacterState struct {
	Name           string   `json:"name"`
	Role           string   `json:"role"`
	LastKnownState string   `json:"last_known_state"`
	Motivation     string   `json:"motivation"`
	Relationships  []string `json:"relationships"`
}

type ChapterGenerationResult struct {
	ChapterID               string                  `json:"chapter_id"`
	ChapterTitle            string                  `json:"chapter_title"`
	ChapterSummary          string                  `json:"chapter_summary"`
	Scenes                  []GeneratedScene        `json:"scenes"`
	ChapterSchemaDesignNote ChapterSchemaDesignNote `json:"chapter_schema_design_note"`
	CarryContext            ChapterCarryContext     `json:"carry_context"`
}

type GeneratedScene struct {
	SceneID       string                `json:"scene_id"`
	SceneIndex    int                   `json:"scene_index"`
	Title         string                `json:"title"`
	Summary       string                `json:"summary"`
	YAMLContent   string                `json:"yaml_content"`
	DesignReasons []domain.DesignReason `json:"design_reasons"`
}

type ChapterSchemaDesignNote struct {
	Summary    string            `json:"summary"`
	KeyReasons []SchemaKeyReason `json:"key_reasons"`
}

type SchemaKeyReason struct {
	FieldName string `json:"field_name"`
	Reason    string `json:"reason"`
}

func NewChapterPromptBuilder() ChapterPromptBuilder {
	return ChapterPromptBuilder{
		SchemaVersion: "1.0",
		BaseSystemPrompt: strings.TrimSpace(`You are E-Director, a professional film director and structured YAML designer.
Convert one novel chapter into scene-level structured output for AI video production.
Return only JSON that matches the requested schema. Do not return markdown.`),
	}
}

func (b ChapterPromptBuilder) BuildChapterPrompt(input ChapterPromptInput) ChapterPrompt {
	return ChapterPrompt{
		SystemPrompt: defaultString(b.BaseSystemPrompt, NewChapterPromptBuilder().BaseSystemPrompt),
		UserPrompt: strings.Join([]string{
			fmt.Sprintf("Novel title: %s", input.NovelTitle),
			fmt.Sprintf("Language: %s", input.Language),
			fmt.Sprintf("Current chapter: %s / %s / index %d", input.ChapterID, input.ChapterTitle, input.ChapterIndex),
			b.BuildContinuityBlock(input.PreviousContext),
			b.BuildRequirementsBlock(input.Config),
			"Current chapter content:",
			input.ChapterContent,
		}, "\n\n"),
		OutputSchema: b.BuildOutputSchema(),
	}
}

func (b ChapterPromptBuilder) BuildContinuityBlock(context ChapterCarryContext) string {
	if context.PreviousChapterID == "" && context.PreviousChapterSummary == "" && len(context.CharacterTimeline) == 0 {
		return "Previous compact context: none. This is the first chapter."
	}
	payload, _ := json.MarshalIndent(context, "", "  ")
	return "Previous compact context from earlier chapters. Use it only for continuity; do not copy prior chapter text:\n" + string(payload)
}

func (b ChapterPromptBuilder) BuildRequirementsBlock(config ChapterPromptConfig) string {
	return fmt.Sprintf(`Output requirements:
- Target format for scene content: %s.
- Scene granularity: %s.
- Include stable chapter_id and scene_id on every scene.
- Include chapter_schema_design_note explaining why this chapter needs its YAML shape.
- Include carry_context for the next chapter: concise summary, character timeline, time/place notes, continuity notes, reusable YAML fields.
- Design reasons enabled: %t.
- Camera directions enabled: %t.
- Dialogues enabled: %t.
- Emotion tags enabled: %t.`, defaultString(config.TargetFormat, "yaml"), defaultString(config.SceneGranularity, "medium"), config.EnableDesignReasons, config.EnableCameraDirections, config.EnableDialogues, config.EnableEmotionTags)
}

func (b ChapterPromptBuilder) BuildOutputSchema() string {
	return strings.TrimSpace(`{
  "chapter_id": "chapter_001",
  "chapter_title": "Chapter title",
  "chapter_summary": "Compact chapter summary for continuity",
  "scenes": [
    {
      "scene_id": "chapter_001_scene_001",
      "scene_index": 1,
      "title": "Scene title",
      "summary": "Scene summary",
      "yaml_content": "chapter_id: chapter_001\nscene_id: chapter_001_scene_001\n...",
      "design_reasons": []
    }
  ],
  "chapter_schema_design_note": {
    "summary": "Why this chapter uses this YAML shape",
    "key_reasons": [{"field_name":"camera","reason":"..."}]
  },
  "carry_context": {
    "previous_chapter_id": "chapter_001",
    "previous_chapter_title": "Chapter title",
    "previous_chapter_summary": "Short summary",
    "character_timeline": [],
    "time_and_place_notes": [],
    "continuity_notes": [],
    "schema_design_summary": "Reusable schema decisions",
    "reusable_yaml_fields": ["camera", "emotion_tags"]
  }
}`)
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
