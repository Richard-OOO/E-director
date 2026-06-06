package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Richard-OOO/E-director/apps/backend/internal/domain"
)

type ChapterSplitInput struct {
	ProjectID string
	JobID     string
	Title     string
	Language  string
	Content   string
}

type ChapterSplitter struct {
	FallbackChunkSize int
}

func NewChapterSplitter() ChapterSplitter {
	return ChapterSplitter{FallbackChunkSize: 4000}
}

func (s ChapterSplitter) Split(input ChapterSplitInput) ([]domain.GenerationChapter, error) {
	content := strings.TrimSpace(input.Content)
	if content == "" {
		return nil, ErrInvalidInput
	}
	chunks := splitByHeadings(content)
	if len(chunks) == 0 {
		chunks = splitByLength(content, s.chunkSize())
	}
	chapters := make([]domain.GenerationChapter, 0, len(chunks))
	now := time.Now().UTC()
	for i, chunk := range chunks {
		chapterID := fmt.Sprintf("chapter_%03d", i+1)
		chapters = append(chapters, domain.GenerationChapter{
			ID:           newGenerationID("chapter", input.JobID, chapterID),
			JobID:        input.JobID,
			ProjectID:    input.ProjectID,
			ChapterID:    chapterID,
			ChapterIndex: i + 1,
			Title:        chunk.title,
			Content:      chunk.content,
			ContentHash:  hashText(chunk.content),
			Status:       domain.GenerationChapterPending,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
	}
	return chapters, nil
}

func (s ChapterSplitter) chunkSize() int {
	if s.FallbackChunkSize > 0 {
		return s.FallbackChunkSize
	}
	return 4000
}

type chapterChunk struct {
	title   string
	content string
}

var chapterHeadingPattern = regexp.MustCompile(`(?m)^\s*((第[\d一二三四五六七八九十百千万零〇两]+[章节回部卷][^\n]*)|(Chapter\s+[0-9IVXLCDM]+[^\n]*)|(CHAPTER\s+[0-9IVXLCDM]+[^\n]*))\s*$`)

func splitByHeadings(content string) []chapterChunk {
	matches := chapterHeadingPattern.FindAllStringSubmatchIndex(content, -1)
	if len(matches) == 0 {
		return nil
	}
	chunks := make([]chapterChunk, 0, len(matches))
	for i, match := range matches {
		title := strings.TrimSpace(content[match[2]:match[3]])
		start := match[1]
		end := len(content)
		if i+1 < len(matches) {
			end = matches[i+1][0]
		}
		body := strings.TrimSpace(content[start:end])
		if body == "" {
			body = title
		}
		chunks = append(chunks, chapterChunk{title: title, content: strings.TrimSpace(title + "\n" + body)})
	}
	return chunks
}

func splitByLength(content string, chunkSize int) []chapterChunk {
	runes := []rune(content)
	chunks := make([]chapterChunk, 0, len(runes)/chunkSize+1)
	for start, index := 0, 1; start < len(runes); start, index = start+chunkSize, index+1 {
		end := start + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, chapterChunk{title: fmt.Sprintf("Chapter %d", index), content: strings.TrimSpace(string(runes[start:end]))})
	}
	return chunks
}

func hashText(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func newGenerationID(parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, ":")))
	return parts[0] + "_" + hex.EncodeToString(sum[:8])
}
