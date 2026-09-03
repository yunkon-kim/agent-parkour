package ai

import (
	"bufio"
	"os"
	"strings"
	"time"
)

// Config holds the AI provider settings loaded from .env or environment variables.
type Config struct {
	Enabled   bool          `json:"enabled"`
	Provider  string        `json:"provider"` // gemini, anthropic, openai, ollama, mock
	Model     string        `json:"model"`
	APIKey    string        `json:"api_key,omitempty"`
	BaseURL   string        `json:"base_url,omitempty"`
	Timeout   time.Duration `json:"timeout"`
	MaxTokens int           `json:"max_tokens"`
}

// LoadConfig loads the configuration from .env file (if exists) and environment variables.
func LoadConfig() *Config {
	// 1. Try to load from .env file in current directory or parent directory
	loadDotEnv(".env")

	enabledStr := os.Getenv("PARKOUR_AI_ENABLED")
	if enabledStr == "" {
		enabledStr = os.Getenv("TOKEN_HOP_AI_ENABLED")
	}
	enabled := strings.ToLower(enabledStr) == "true" || enabledStr == "1"

	provider := strings.ToLower(strings.TrimSpace(os.Getenv("PARKOUR_AI_PROVIDER")))
	if provider == "" {
		provider = strings.ToLower(strings.TrimSpace(os.Getenv("TOKEN_HOP_AI_PROVIDER")))
	}
	if provider == "" {
		provider = "gemini"
	}

	model := strings.TrimSpace(os.Getenv("PARKOUR_AI_MODEL"))
	if model == "" {
		model = strings.TrimSpace(os.Getenv("TOKEN_HOP_AI_MODEL"))
	}

	cfg := &Config{
		Enabled:   enabled,
		Provider:  provider,
		Model:     model,
		Timeout:   30 * time.Second,
		MaxTokens: 4096,
	}

	// Resolve provider-specific keys and defaults
	switch provider {
	case "gemini", "google":
		cfg.Provider = "gemini"
		cfg.APIKey = os.Getenv("GEMINI_API_KEY")
		if cfg.Model == "" {
			cfg.Model = "gemini-2.5-pro"
		}
	case "anthropic", "claude":
		cfg.Provider = "anthropic"
		cfg.APIKey = os.Getenv("ANTHROPIC_API_KEY")
		if cfg.Model == "" {
			cfg.Model = "claude-3-7-sonnet-20250219"
		}
	case "openai", "gpt":
		cfg.Provider = "openai"
		cfg.APIKey = os.Getenv("OPENAI_API_KEY")
		if cfg.Model == "" {
			cfg.Model = "gpt-4o"
		}
	case "ollama", "local":
		cfg.Provider = "ollama"
		cfg.BaseURL = os.Getenv("OLLAMA_BASE_URL")
		if cfg.BaseURL == "" {
			cfg.BaseURL = "http://localhost:11434"
		}
		if cfg.Model == "" {
			cfg.Model = os.Getenv("OLLAMA_MODEL")
			if cfg.Model == "" {
				cfg.Model = "deepseek-r1:14b"
			}
		}
	case "mock":
		cfg.Provider = "mock"
		cfg.Model = "mock-model"
	}

	return cfg
}

// loadDotEnv parses a simple KEY=VALUE .env file without external dependencies.
func loadDotEnv(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		// Strip surrounding quotes
		if (strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"")) ||
			(strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'")) {
			val = val[1 : len(val)-1]
		}

		// Set env only if not already set by system environment
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}
