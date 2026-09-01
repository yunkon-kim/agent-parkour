package ai

import (
	"context"
)

// AIProvider defines the standard unified interface across all LLM backends (Gemini, Claude, OpenAI, Ollama).
type AIProvider interface {
	// Name returns the provider identifier (e.g. "gemini", "claude", "openai", "ollama", "mock")
	Name() string

	// Model returns the active model identifier (e.g. "gemini-2.5-pro", "claude-3-7-sonnet-20250219")
	Model() string

	// DecomposeRule semantically analyzes and splits an oversized prompt/rule into modular JIT sub-rules / skills.
	DecomposeRule(ctx context.Context, title, content string, maxTokens int) (*DecomposeResult, error)

	// GenerateDescription creates a concise 1-sentence activation description for Cursor / Antigravity rule triggers.
	GenerateDescription(ctx context.Context, content string) (string, error)
}

// SubRule represents a single decomposed modular sub-rule or JIT skill.
type SubRule struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Globs       []string `json:"globs,omitempty"`
	IsSkill     bool     `json:"is_skill"` // true if it should be an on-demand JIT skill package
	Content     string   `json:"content"`
}

// DecomposeResult holds the outcome of a semantic rule decomposition.
type DecomposeResult struct {
	OriginalTitle string     `json:"original_title"`
	Summary       string     `json:"summary"`
	SubRules      []*SubRule `json:"sub_rules"`
}
