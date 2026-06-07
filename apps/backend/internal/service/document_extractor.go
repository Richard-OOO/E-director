package service

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"io"
	"path/filepath"
	"strings"

	"rsc.io/pdf"
)

const defaultImportMaxBytes int64 = 20 << 20

var (
	ErrUnsupportedDocument = errors.New("unsupported document type")
	ErrEmptyDocumentText   = errors.New("document contains no extractable text")
)

type DocumentTextExtractor interface {
	Extract(ctx context.Context, input DocumentExtractInput) (string, error)
}

type DocumentExtractInput struct {
	SourceType  string
	Filename    string
	ContentType string
	Reader      io.Reader
	MaxBytes    int64
}

type DefaultDocumentTextExtractor struct{}

func NewDocumentTextExtractor() DefaultDocumentTextExtractor {
	return DefaultDocumentTextExtractor{}
}

func (e DefaultDocumentTextExtractor) Extract(ctx context.Context, input DocumentExtractInput) (string, error) {
	if input.Reader == nil {
		return "", ErrInvalidInput
	}
	maxBytes := input.MaxBytes
	if maxBytes <= 0 {
		maxBytes = defaultImportMaxBytes
	}
	content, err := io.ReadAll(io.LimitReader(input.Reader, maxBytes+1))
	if err != nil {
		return "", err
	}
	if int64(len(content)) > maxBytes {
		return "", ErrInvalidInput
	}

	sourceType := normalizeDocumentSourceType(input.SourceType, input.Filename)
	var text string
	switch sourceType {
	case "docx":
		text, err = extractDocxText(content)
	case "pdf":
		text, err = extractPDFText(content)
	default:
		return "", ErrUnsupportedDocument
	}
	if err != nil {
		return "", err
	}
	text = normalizeExtractedText(text)
	if text == "" {
		return "", ErrEmptyDocumentText
	}
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return text, nil
	}
}

func normalizeDocumentSourceType(sourceType, filename string) string {
	sourceType = strings.TrimSpace(strings.ToLower(sourceType))
	if sourceType == "docx" || sourceType == "pdf" {
		return sourceType
	}
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
	if ext == "docx" || ext == "pdf" {
		return ext
	}
	return sourceType
}

func extractDocxText(content []byte) (string, error) {
	reader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return "", err
	}
	for _, file := range reader.File {
		if file.Name != "word/document.xml" {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()
		return extractDocxDocumentXML(rc)
	}
	return "", ErrUnsupportedDocument
}

func extractDocxDocumentXML(reader io.Reader) (string, error) {
	decoder := xml.NewDecoder(reader)
	var b strings.Builder
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		switch node := token.(type) {
		case xml.StartElement:
			switch node.Name.Local {
			case "t":
				var text string
				if err := decoder.DecodeElement(&text, &node); err != nil {
					return "", err
				}
				b.WriteString(text)
			case "tab":
				b.WriteByte('\t')
			case "br":
				b.WriteByte('\n')
			}
		case xml.EndElement:
			if node.Name.Local == "p" {
				b.WriteByte('\n')
			}
		}
	}
	return b.String(), nil
}

func extractPDFText(content []byte) (string, error) {
	reader, err := pdf.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return "", err
	}
	var b strings.Builder
	for pageIndex := 1; pageIndex <= reader.NumPage(); pageIndex++ {
		page := reader.Page(pageIndex)
		for _, text := range page.Content().Text {
			if strings.TrimSpace(text.S) == "" {
				continue
			}
			if b.Len() > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(text.S)
		}
		b.WriteByte('\n')
	}
	return b.String(), nil
}

func normalizeExtractedText(text string) string {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	cleaned := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.Join(strings.Fields(line), " ")
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}
	return strings.TrimSpace(strings.Join(cleaned, "\n"))
}
