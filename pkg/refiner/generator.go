package refiner

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yunkon-kim/token-hop/pkg/audit"
	"github.com/yunkon-kim/token-hop/pkg/ir"
)

// Generator builds target-specific AI refinement prompts
type Generator struct {
	MaxTokens int
}

// NewGenerator creates a new Generator
func NewGenerator(maxTokens int) *Generator {
	if maxTokens <= 0 {
		maxTokens = 400
	}
	return &Generator{MaxTokens: maxTokens}
}

// GenerateFromDocs builds a refinement prompt from a list of UA-IR documents
func (g *Generator) GenerateFromDocs(targetPlatform string, docs []*ir.UADocument, customGuidance string) (string, error) {
	tmpl := GetTemplate(targetPlatform)

	var files []FileItem
	totalTokens := 0

	for _, doc := range docs {
		tokens := audit.EstimateTokens(doc.Payload.MarkdownBody)
		chars := len(doc.Payload.MarkdownBody)
		totalTokens += tokens

		filePath := doc.Payload.RawSource
		if filePath == "" {
			filePath = fmt.Sprintf("%s.md", doc.Metadata.ID)
		}

		files = append(files, FileItem{
			FilePath:    filePath,
			EntityType:  doc.Metadata.Type,
			FileContent: doc.Payload.MarkdownBody,
			Characters:  chars,
			Tokens:      tokens,
			IsOversized: tokens > g.MaxTokens,
		})
	}

	ctx := PromptContext{
		TargetPlatform: targetPlatform,
		PlatformName:   ResolvePlatformName(targetPlatform),
		Files:          files,
		TotalTokens:    totalTokens,
		TotalFiles:     len(files),
		CustomGuidance: customGuidance,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx); err != nil {
		return "", fmt.Errorf("failed to render refinement prompt template: %w", err)
	}

	return buf.String(), nil
}

// GenerateFromFile builds a refinement prompt from a single file path
func (g *Generator) GenerateFromFile(targetPlatform string, filePath string, customGuidance string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	content := string(data)
	tokens := audit.EstimateTokens(content)
	chars := len(content)
	entityType := InferEntityType(filePath)

	fileItem := FileItem{
		FilePath:    filePath,
		EntityType:  entityType,
		FileContent: content,
		Characters:  chars,
		Tokens:      tokens,
		IsOversized: tokens > g.MaxTokens,
	}

	ctx := PromptContext{
		TargetPlatform: targetPlatform,
		PlatformName:   ResolvePlatformName(targetPlatform),
		Files:          []FileItem{fileItem},
		TotalTokens:    tokens,
		TotalFiles:     1,
		CustomGuidance: customGuidance,
	}

	tmpl := GetTemplate(targetPlatform)
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return buf.String(), nil
}

// GenerateFromDirectory scans a directory and builds a combined refinement prompt
func (g *Generator) GenerateFromDirectory(targetPlatform string, dirPath string, customGuidance string) (string, error) {
	var files []FileItem
	totalTokens := 0

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".md" && ext != ".mdc" && ext != ".txt" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		content := string(data)
		tokens := audit.EstimateTokens(content)
		chars := len(content)
		totalTokens += tokens

		files = append(files, FileItem{
			FilePath:    path,
			EntityType:  InferEntityType(path),
			FileContent: content,
			Characters:  chars,
			Tokens:      tokens,
			IsOversized: tokens > g.MaxTokens,
		})
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to scan directory %s: %w", dirPath, err)
	}

	ctx := PromptContext{
		TargetPlatform: targetPlatform,
		PlatformName:   ResolvePlatformName(targetPlatform),
		RootDir:        dirPath,
		Files:          files,
		TotalTokens:    totalTokens,
		TotalFiles:     len(files),
		CustomGuidance: customGuidance,
	}

	tmpl := GetTemplate(targetPlatform)
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return buf.String(), nil
}

// ResolvePlatformName returns the display name of a target platform
func ResolvePlatformName(platform string) string {
	switch strings.ToLower(strings.TrimSpace(platform)) {
	case "antigravity", "gemini", "google":
		return "Google Antigravity"
	case "cursor", "cursor-ai":
		return "Cursor AI"
	case "claude", "claude-code", "anthropic":
		return "Claude Code"
	case "copilot", "github":
		return "GitHub Copilot"
	default:
		return strings.ToUpper(platform)
	}
}

// InferEntityType infers entity type from file path heuristics
func InferEntityType(path string) ir.EntityType {
	lower := strings.ToLower(path)
	if strings.Contains(lower, "skill") {
		return ir.TypeSkill
	}
	if strings.Contains(lower, "workflow") || strings.Contains(lower, "prompt") {
		return ir.TypeWorkflow
	}
	if strings.Contains(lower, "agent") {
		return ir.TypeAgent
	}
	if strings.Contains(lower, "instruction") || strings.Contains(lower, "ssot") || strings.Contains(lower, "agents.md") || strings.Contains(lower, "claude.md") {
		return ir.TypeInstruction
	}
	return ir.TypeRule
}
