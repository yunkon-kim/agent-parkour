package ai

import (
	"fmt"
	"strings"
)

// NewProvider creates a unified AIProvider instance matching the given configuration.
func NewProvider(cfg *Config) (AIProvider, error) {
	if cfg == nil {
		cfg = LoadConfig()
	}

	providerName := strings.ToLower(strings.TrimSpace(cfg.Provider))

	switch providerName {
	case "gemini", "google":
		return NewGeminiProvider(cfg.APIKey, cfg.Model, cfg.Timeout)

	case "anthropic", "claude":
		return NewClaudeProvider(cfg.APIKey, cfg.Model, cfg.Timeout)

	case "openai", "gpt":
		return NewOpenAIProvider(cfg.APIKey, cfg.Model, cfg.BaseURL, cfg.Timeout)

	case "ollama", "local":
		return NewOllamaProvider(cfg.BaseURL, cfg.Model, cfg.Timeout)

	case "mock", "offline", "test":
		return NewMockProvider(cfg.Model), nil

	default:
		return nil, fmt.Errorf("unsupported AI provider: %q (supported: gemini, claude, openai, ollama, mock)", providerName)
	}
}
