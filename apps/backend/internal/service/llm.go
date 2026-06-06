package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Richard-OOO/E-director/apps/backend/internal/domain"
)

type ChapterGenerator interface {
	GenerateChapter(ctx context.Context, prompt ChapterPrompt) (ChapterGenerationResult, error)
}

type MockChapterGenerator struct{}

func NewMockChapterGenerator() *MockChapterGenerator {
	return &MockChapterGenerator{}
}

func (g *MockChapterGenerator) GenerateChapter(ctx context.Context, prompt ChapterPrompt) (ChapterGenerationResult, error) {
	_ = ctx
	chapterID := extractChapterID(prompt.UserPrompt)
	if chapterID == "" {
		chapterID = "chapter_001"
	}
	chapterIndex := extractChapterIndex(chapterID)
	sceneCount := 2
	if chapterIndex == 1 {
		sceneCount = 2
	}
	scenes := make([]GeneratedScene, 0, sceneCount)
	for i := 1; i <= sceneCount; i++ {
		sceneID := fmt.Sprintf("%s_scene_%03d", chapterID, i)
		scenes = append(scenes, GeneratedScene{
			SceneID:    sceneID,
			SceneIndex: i,
			Title:      fmt.Sprintf("Scene %d of %s", i, chapterID),
			Summary:    fmt.Sprintf("Mock summary for %s", sceneID),
			YAMLContent: strings.TrimSpace(fmt.Sprintf(`chapter_id: %s
scene_id: %s
scene_index: %d
title: %s
summary: %s
camera: "static close-up"
emotion_tags:
  - tension
visual_elements:
  - rain
  - neon`, chapterID, sceneID, i, fmt.Sprintf("Scene %d", i), fmt.Sprintf("Mock summary for %s", sceneID))),
			DesignReasons: []domain.DesignReason{{
				ReasonID:    sceneID + "_reason_1",
				TargetPath:  "camera",
				FieldName:   "camera",
				ReasonType:  "cinematography",
				Title:       "Maintain visual focus",
				Description: "Use a compact camera setup to keep the scene readable.",
				HoverText:   "Keeps the visual language consistent with the previous chapter.",
				Confidence:  0.88,
			}},
		})
	}
	chapterTitle := chapterID
	if title := extractChapterTitle(prompt.UserPrompt); title != "" {
		chapterTitle = title
	}
	return ChapterGenerationResult{
		ChapterID:      chapterID,
		ChapterTitle:   chapterTitle,
		ChapterSummary: fmt.Sprintf("Mock summary for %s", chapterID),
		Scenes:         scenes,
		ChapterSchemaDesignNote: ChapterSchemaDesignNote{
			Summary: fmt.Sprintf("This chapter uses a compact mock schema for %s.", chapterID),
			KeyReasons: []SchemaKeyReason{{
				FieldName: "camera",
				Reason:    "Needed to preserve visual continuity in the next chapter.",
			}},
		},
		CarryContext: ChapterCarryContext{
			PreviousChapterID:      chapterID,
			PreviousChapterTitle:   chapterTitle,
			PreviousChapterSummary: fmt.Sprintf("Mock summary for %s", chapterID),
			CharacterTimeline:      []CharacterState{{Name: "Protagonist", Role: "lead", LastKnownState: "moving forward", Motivation: "continue"}},
			TimeAndPlaceNotes:      []string{"night", "neon city"},
			ContinuityNotes:        []string{"keep rain and neon imagery"},
			SchemaDesignSummary:    "retain camera, emotion_tags, and visual_elements",
			ReusableYAMLFields:     []string{"camera", "emotion_tags", "visual_elements"},
		},
	}, nil
}

func extractChapterID(text string) string {
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Current chapter:") {
			parts := strings.Split(line, "/")
			if len(parts) > 0 {
				piece := strings.TrimSpace(strings.TrimPrefix(parts[0], "Current chapter:"))
				if piece != "" {
					return piece
				}
			}
		}
	}
	return ""
}

func extractChapterTitle(text string) string {
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Current chapter:") {
			parts := strings.Split(line, "/")
			if len(parts) >= 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}

func extractChapterIndex(chapterID string) int {
	var n int
	if _, err := fmt.Sscanf(chapterID, "chapter_%d", &n); err == nil && n > 0 {
		return n
	}
	return 1
}
