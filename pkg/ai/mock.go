package ai

import (
	"context"
	"fmt"
	"strings"
)

// MockProvider provides deterministic responses for offline unit tests without external API calls.
type MockProvider struct {
	model string
}

func NewMockProvider(model string) *MockProvider {
	if model == "" {
		model = "mock-model"
	}
	return &MockProvider{model: model}
}

func (p *MockProvider) Name() string {
	return "mock"
}

func (p *MockProvider) Model() string {
	return p.model
}

func (p *MockProvider) DecomposeRule(ctx context.Context, title, content string, maxTokens int) (*DecomposeResult, error) {
	// Deterministically decompose based on markdown headers (e.g. ## or ###)
	lines := strings.Split(content, "\n")
	var subRules []*SubRule

	var currentTitle string
	var currentLines []string

	flush := func() {
		if len(currentLines) == 0 {
			return
		}
		if currentTitle == "" {
			currentTitle = "Core Conventions"
		}
		id := strings.ToLower(strings.ReplaceAll(currentTitle, " ", "-"))
		subRules = append(subRules, &SubRule{
			ID:          id,
			Title:       currentTitle,
			Description: fmt.Sprintf("Modular guidelines for %s", currentTitle),
			Globs:       []string{"**/*"},
			IsSkill:     len(currentLines) > 30, // Decompose to skill if long
			Content:     strings.Join(currentLines, "\n"),
		})
		currentLines = nil
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "## ") {
			flush()
			currentTitle = strings.TrimPrefix(line, "## ")
		}
		currentLines = append(currentLines, line)
	}
	flush()

	if len(subRules) == 0 {
		subRules = append(subRules, &SubRule{
			ID:          "main",
			Title:       title,
			Description: fmt.Sprintf("Guidelines for %s", title),
			Content:     content,
		})
	}

	return &DecomposeResult{
		OriginalTitle: title,
		Summary:       fmt.Sprintf("Decomposed into %d modular sub-rules via Mock Engine", len(subRules)),
		SubRules:      subRules,
	}, nil
}

func (p *MockProvider) GenerateDescription(ctx context.Context, content string) (string, error) {
	lines := strings.Split(content, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" && !strings.HasPrefix(l, "#") && !strings.HasPrefix(l, "---") {
			if len(l) > 100 {
				return l[:97] + "...", nil
			}
			return l, nil
		}
	}
	return "Standard project guidelines and rules.", nil
}
