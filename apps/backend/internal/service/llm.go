package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Richard-OOO/E-director/apps/backend/internal/domain"
)

type ChapterGenerator interface {
	GenerateChapter(ctx context.Context, prompt ChapterPrompt) (ChapterGenerationResult, error)
}

type OpenAIChapterGeneratorConfig struct {
	BaseURL   string
	APIKey    string
	Model     string
	MaxTokens int
}

type OpenAIChapterGenerator struct {
	client *http.Client
	config OpenAIChapterGeneratorConfig
}

type openAIChatCompletionRequest struct {
	Model       string               `json:"model"`
	Messages    []openAIMessage      `json:"messages"`
	MaxTokens   int                  `json:"max_tokens,omitempty"`
	Temperature float64              `json:"temperature,omitempty"`
	ResponseFmt openAIResponseFormat `json:"response_format,omitempty"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponseFormat struct {
	Type string `json:"type"`
}

type openAIChatCompletionResponse struct {
	ID      string                       `json:"id,omitempty"`
	Model   string                       `json:"model,omitempty"`
	Choices []openAIChatCompletionChoice `json:"choices"`
}

type openAIChatCompletionChoice struct {
	Text    string `json:"text,omitempty"`
	Message struct {
		Content          openAIMessageContent `json:"content"`
		ReasoningContent string               `json:"reasoning_content,omitempty"`
	} `json:"message"`
}

type openAIMessageContent struct {
	Text  string
	Parts []openAIContentPart
}

type openAIContentPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (c *openAIMessageContent) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		c.Text = s
		return nil
	}

	var parts []openAIContentPart
	if err := json.Unmarshal(data, &parts); err == nil {
		c.Parts = parts
		return nil
	}

	if string(data) == "null" {
		return nil
	}

	return fmt.Errorf("unsupported message.content shape")
}

const maxDiagnosticBodyBytes = 8192
const maxDiagnosticPreviewChars = 600

func extractAssistantText(choice openAIChatCompletionChoice) (content string, reasoning string) {
	reasoning = strings.TrimSpace(choice.Message.ReasoningContent)

	if text := strings.TrimSpace(choice.Message.Content.Text); text != "" {
		return text, reasoning
	}

	var parts []string
	for _, part := range choice.Message.Content.Parts {
		if strings.EqualFold(part.Type, "text") && strings.TrimSpace(part.Text) != "" {
			parts = append(parts, strings.TrimSpace(part.Text))
		}
	}
	if len(parts) > 0 {
		return strings.Join(parts, "\n"), reasoning
	}

	if text := strings.TrimSpace(choice.Text); text != "" {
		return text, reasoning
	}

	return "", reasoning
}

func diagnosticPreview(value string) string {
	value = strings.TrimSpace(value)
	if len(value) <= maxDiagnosticPreviewChars {
		return value
	}
	return value[:maxDiagnosticPreviewChars] + "...[truncated]"
}

func diagnosticBody(value string) string {
	value = strings.TrimSpace(value)
	if len(value) <= maxDiagnosticBodyBytes {
		return value
	}
	return value[:maxDiagnosticBodyBytes] + "...[truncated]"
}

func NewOpenAIChapterGenerator(cfg OpenAIChapterGeneratorConfig) *OpenAIChapterGenerator {
	if strings.TrimSpace(cfg.BaseURL) == "" {
		cfg.BaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	}
	if strings.TrimSpace(cfg.Model) == "" {
		cfg.Model = "qwen-flash-2025-07-28"
	}
	if cfg.MaxTokens <= 0 {
		cfg.MaxTokens = 10000
	}
	return &OpenAIChapterGenerator{
		client: &http.Client{
			Timeout: 90 * time.Second,
			Transport: &http.Transport{
				Proxy: nil,
			},
		},
		config: cfg,
	}
}

func (g *OpenAIChapterGenerator) GenerateChapter(ctx context.Context, prompt ChapterPrompt) (ChapterGenerationResult, error) {
	if g == nil {
		return ChapterGenerationResult{}, errors.New("generator unavailable")
	}
	if strings.TrimSpace(g.config.APIKey) == "" {
		return ChapterGenerationResult{}, errors.New("OPENAI_API_KEY is required")
	}
	payload := openAIChatCompletionRequest{
		Model: g.config.Model,
		Messages: []openAIMessage{
			{Role: "system", Content: prompt.SystemPrompt},
			{Role: "user", Content: prompt.UserPrompt + "\n\nReturn JSON only that matches this schema:\n" + prompt.OutputSchema},
		},
		MaxTokens:   g.config.MaxTokens,
		Temperature: 0.2,
		ResponseFmt: openAIResponseFormat{Type: "json_object"},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return ChapterGenerationResult{}, err
	}
	endpoint := strings.TrimRight(g.config.BaseURL, "/") + "/chat/completions"
	log.Printf(
		"[e-director:llm] request provider=openai-compatible endpoint=%s model=%s max_tokens=%d message_count=%d response_format=%s",
		endpoint,
		g.config.Model,
		g.config.MaxTokens,
		len(payload.Messages),
		payload.ResponseFmt.Type,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return ChapterGenerationResult{}, err
	}
	req.Header.Set("Authorization", "Bearer "+g.config.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := g.client.Do(req)
	if err != nil {
		return ChapterGenerationResult{}, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChapterGenerationResult{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ChapterGenerationResult{}, fmt.Errorf(
			"openai-compatible api returned status %s model=%s body=%s",
			resp.Status,
			g.config.Model,
			diagnosticBody(string(respBody)),
		)
	}
	var completion openAIChatCompletionResponse
	if err := json.Unmarshal(respBody, &completion); err != nil {
		return ChapterGenerationResult{}, fmt.Errorf(
			"decode openai-compatible response: %w; status=%s model=%s body=%s",
			err,
			resp.Status,
			g.config.Model,
			diagnosticBody(string(respBody)),
		)
	}
	if len(completion.Choices) == 0 {
		return ChapterGenerationResult{}, fmt.Errorf(
			"openai-compatible api returned no choices; status=%s model=%s body=%s",
			resp.Status,
			g.config.Model,
			diagnosticBody(string(respBody)),
		)
	}

	content, reasoning := extractAssistantText(completion.Choices[0])
	content = strings.TrimSpace(content)
	reasoningLen := len(reasoning)
	log.Printf(
		"[e-director:llm] response status=%s request_model=%s response_model=%s choice_count=%d content_len=%d reasoning_present=%t reasoning_len=%d content_preview=%q",
		resp.Status,
		g.config.Model,
		completion.Model,
		len(completion.Choices),
		len(content),
		reasoningLen > 0,
		reasoningLen,
		diagnosticPreview(content),
	)

	if content == "" {
		return ChapterGenerationResult{}, fmt.Errorf(
			"openai-compatible api returned empty assistant content; status=%s model=%s choice_count=%d reasoning_present=%t reasoning_len=%d body=%s",
			resp.Status,
			g.config.Model,
			len(completion.Choices),
			reasoningLen > 0,
			reasoningLen,
			diagnosticBody(string(respBody)),
		)
	}
	content = stripJSONCodeFence(content)
	var result ChapterGenerationResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return ChapterGenerationResult{}, fmt.Errorf(
			"decode model JSON content: %w; model=%s reasoning_present=%t reasoning_len=%d content_preview=%q",
			err,
			g.config.Model,
			reasoningLen > 0,
			reasoningLen,
			diagnosticPreview(content),
		)
	}
	normalizeGenerationResult(&result)
	if err := validateGenerationResult(result); err != nil {
		return ChapterGenerationResult{}, err
	}
	return result, nil
}

func stripJSONCodeFence(content string) string {
	content = strings.TrimSpace(content)
	if !strings.HasPrefix(content, "```") {
		return content
	}
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSpace(content)
	if strings.HasPrefix(strings.ToLower(content), "json") {
		content = strings.TrimSpace(content[4:])
	}
	if idx := strings.LastIndex(content, "```"); idx >= 0 {
		content = strings.TrimSpace(content[:idx])
	}
	return content
}

func validateGenerationResult(result ChapterGenerationResult) error {
	if strings.TrimSpace(result.ChapterSummary) == "" {
		return errors.New("model output missing chapter_summary")
	}
	if len(result.Scenes) == 0 {
		return errors.New("model output missing scenes")
	}
	for i, scene := range result.Scenes {
		if strings.TrimSpace(scene.YAMLContent) == "" {
			return fmt.Errorf("model output scene %d missing yaml_content", i+1)
		}
	}
	return nil
}

func normalizeGenerationResult(result *ChapterGenerationResult) {
	if result == nil {
		return
	}
	if strings.TrimSpace(result.ChapterID) == "" {
		result.ChapterID = "chapter_001"
	}
	if strings.TrimSpace(result.ChapterTitle) == "" {
		result.ChapterTitle = result.ChapterID
	}
	if strings.TrimSpace(result.ChapterSummary) == "" {
		result.ChapterSummary = result.ChapterTitle
	}
	for i := range result.Scenes {
		if strings.TrimSpace(result.Scenes[i].SceneID) == "" {
			result.Scenes[i].SceneID = fmt.Sprintf("%s_scene_%03d", result.ChapterID, i+1)
		}
		if result.Scenes[i].SceneIndex <= 0 {
			result.Scenes[i].SceneIndex = i + 1
		}
		if strings.TrimSpace(result.Scenes[i].Title) == "" {
			result.Scenes[i].Title = fmt.Sprintf("Scene %d", i+1)
		}
		if strings.TrimSpace(result.Scenes[i].Summary) == "" {
			result.Scenes[i].Summary = result.Scenes[i].Title
		}
		if len(result.Scenes[i].DesignReasons) == 0 {
			result.Scenes[i].DesignReasons = defaultSceneDesignReasons(result.Scenes[i])
		}
		for j := range result.Scenes[i].DesignReasons {
			if strings.TrimSpace(result.Scenes[i].DesignReasons[j].HoverText) == "" {
				fallback := result.Scenes[i].DesignReasons[j].Description
				if strings.TrimSpace(fallback) == "" {
					fallback = result.Scenes[i].DesignReasons[j].Title
				}
				result.Scenes[i].DesignReasons[j].HoverText = fallback
			}
		}
	}
}

func defaultSceneDesignReasons(scene GeneratedScene) []domain.DesignReason {
	title := "Generated YAML"
	description := strings.TrimSpace(scene.Summary)
	if description == "" {
		description = "Scene YAML was generated from the chapter content."
	}
	return []domain.DesignReason{
		{
			ReasonID:    fmt.Sprintf("%s_reason_001", defaultString(scene.SceneID, "scene")),
			TargetPath:  "scene",
			FieldName:   "yaml_content",
			ReasonType:  "structure",
			Title:       title,
			Description: description,
			HoverText:   description,
			Confidence:  0.6,
		},
	}
}
