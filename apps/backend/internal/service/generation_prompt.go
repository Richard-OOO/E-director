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

func (c *ChapterCarryContext) UnmarshalJSON(data []byte) error {
	type chapterCarryContextAlias struct {
		PreviousChapterID      string          `json:"previous_chapter_id"`
		PreviousChapterTitle   string          `json:"previous_chapter_title"`
		PreviousChapterSummary string          `json:"previous_chapter_summary"`
		CharacterTimeline      json.RawMessage `json:"character_timeline"`
		TimeAndPlaceNotes      []string        `json:"time_and_place_notes"`
		ContinuityNotes        []string        `json:"continuity_notes"`
		SchemaDesignSummary    string          `json:"schema_design_summary"`
		ReusableYAMLFields     []string        `json:"reusable_yaml_fields"`
	}
	var raw chapterCarryContextAlias
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	c.PreviousChapterID = raw.PreviousChapterID
	c.PreviousChapterTitle = raw.PreviousChapterTitle
	c.PreviousChapterSummary = raw.PreviousChapterSummary
	c.TimeAndPlaceNotes = raw.TimeAndPlaceNotes
	c.ContinuityNotes = raw.ContinuityNotes
	c.SchemaDesignSummary = raw.SchemaDesignSummary
	c.ReusableYAMLFields = raw.ReusableYAMLFields

	states, err := parseCharacterTimeline(raw.CharacterTimeline)
	if err != nil {
		return err
	}
	c.CharacterTimeline = states
	return nil
}

func parseCharacterTimeline(data json.RawMessage) ([]CharacterState, error) {
	if len(data) == 0 || string(data) == "null" {
		return nil, nil
	}
	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("decode carry_context.character_timeline: %w", err)
	}
	states := make([]CharacterState, 0, len(items))
	for _, item := range items {
		var state CharacterState
		if err := json.Unmarshal(item, &state); err == nil {
			states = append(states, state)
			continue
		}
		var note string
		if err := json.Unmarshal(item, &note); err == nil {
			note = strings.TrimSpace(note)
			if note != "" {
				states = append(states, CharacterState{Name: note, LastKnownState: note})
			}
			continue
		}
		return nil, fmt.Errorf("decode carry_context.character_timeline item: expected object or string")
	}
	return states, nil
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
		BaseSystemPrompt: strings.TrimSpace(`You are E-Director, a professional film director and AI video script formatter.
Convert one novel chapter into scene-level YAML scripts using the fixed E-director YAML Schema.
Return only JSON that matches the requested outer response shape. Do not return markdown.
Do not invent new YAML schema fields; only fill the fixed YAML schema fields with content from the chapter.`),
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
- yaml_content must be a YAML string that follows the fixed E-director YAML Schema exactly.
- Do not design a new YAML schema per chapter. Do not add, rename, or remove YAML schema fields.
- Fill empty lists as [] and empty text fields as short empty-safe descriptions; do not output null values.
- Include carry_context for the next chapter: concise summary, character timeline, time/place notes, continuity notes, reusable YAML fields.
- Camera directions enabled: %t.
- Dialogues enabled: %t.
- Emotion tags enabled: %t.`, defaultString(config.TargetFormat, "yaml"), defaultString(config.SceneGranularity, "medium"), config.EnableCameraDirections, config.EnableDialogues, config.EnableEmotionTags)
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
      "yaml_content": "schema_version: \"1.0\"\nchapter:\n  chapter_id: \"chapter_001\"\n  chapter_title: \"Chapter title\"\n  chapter_index: 1\nscene:\n  scene_id: \"chapter_001_scene_001\"\n  scene_index: 1\n  title: \"Scene title\"\n  summary: \"Scene summary\"\n  location:\n    name: \"Location name\"\n    type: \"interior\"\n    time_of_day: \"night\"\n    atmosphere: \"Scene atmosphere\"\n  characters:\n    - name: \"Character name\"\n      role: \"protagonist\"\n      appearance: \"Visible appearance\"\n      emotional_state: \"Current emotion\"\n      goal: \"Scene goal\"\n  action:\n    - actor: \"Character name\"\n      description: \"Filmable action\"\n      purpose: \"Narrative purpose\"\n  dialogues: []\n  camera:\n    - shot_type: \"wide_shot\"\n      movement: \"static\"\n      subject: \"Shot subject\"\n      description: \"Camera description\"\n  visual_elements:\n    - name: \"Visual element\"\n      category: \"prop\"\n      description: \"Visual description\"\n  audio:\n    music: \"Music direction\"\n    sound_effects: []\n  emotion_tags: []\n  ai_video_prompt:\n    positive: \"Positive AI video prompt\"\n    negative: \"Negative AI video prompt\"\n  transition:\n    type: \"cut\"\n    description: \"Transition description\""
    }
  ],
  "carry_context": {
    "previous_chapter_id": "chapter_001",
    "previous_chapter_title": "Chapter title",
    "previous_chapter_summary": "Short summary",
    "character_timeline": [
      {
        "name": "Character name",
        "role": "protagonist",
        "last_known_state": "Where the character ends this chapter",
        "motivation": "Current motivation",
        "relationships": []
      }
    ],
    "time_and_place_notes": [],
    "continuity_notes": [],
    "schema_design_summary": "Fixed E-director YAML Schema v1.0 is used for every scene.",
    "reusable_yaml_fields": ["location", "characters", "action", "dialogues", "camera", "visual_elements", "audio", "emotion_tags", "ai_video_prompt", "transition"]
  }
}`)
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
