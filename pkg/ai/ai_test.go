package ai

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestMockProviderDecomposition(t *testing.T) {
	provider := NewMockProvider("mock-v1")
	if provider.Name() != "mock" {
		t.Fatalf("expected provider name 'mock', got %s", provider.Name())
	}

	content := `# Global Rules
## Go Guidelines
Use zerolog and table-driven tests.

## UI Design System
Use CSS variables and modern typography.
`

	result, err := provider.DecomposeRule(context.Background(), "Project Rules", content, 400)
	if err != nil {
		t.Fatalf("unexpected decomposition error: %v", err)
	}

	if len(result.SubRules) < 2 {
		t.Fatalf("expected at least 2 sub-rules, got %d", len(result.SubRules))
	}

	desc, err := provider.GenerateDescription(context.Background(), content)
	if err != nil {
		t.Fatalf("unexpected description error: %v", err)
	}
	if desc == "" {
		t.Fatalf("expected non-empty description")
	}
}

func TestConfigLoader(t *testing.T) {
	os.Setenv("PARKOUR_AI_ENABLED", "true")
	os.Setenv("PARKOUR_AI_PROVIDER", "gemini")
	os.Setenv("GEMINI_API_KEY", "test-key-123")
	os.Setenv("PARKOUR_AI_MODEL", "gemini-2.5-pro")

	cfg := LoadConfig()
	if !cfg.Enabled {
		t.Fatalf("expected enabled to be true")
	}
	if cfg.Provider != "gemini" {
		t.Fatalf("expected provider 'gemini', got %s", cfg.Provider)
	}
	if cfg.APIKey != "test-key-123" {
		t.Fatalf("expected api key 'test-key-123', got %s", cfg.APIKey)
	}
}

func TestFactoryCreation(t *testing.T) {
	// 1. Mock
	mockP, err := NewProvider(&Config{Provider: "mock"})
	if err != nil || mockP.Name() != "mock" {
		t.Fatalf("failed to create mock provider: %v", err)
	}

	// 2. Gemini with missing key should error
	_, err = NewProvider(&Config{Provider: "gemini", APIKey: ""})
	if err == nil {
		t.Fatalf("expected error for missing gemini api key")
	}

	// 3. Claude with key
	claudeP, err := NewProvider(&Config{Provider: "claude", APIKey: "test-claude-key"})
	if err != nil || claudeP.Name() != "anthropic" {
		t.Fatalf("failed to create claude provider: %v", err)
	}

	// 4. OpenAI with key
	openAIP, err := NewProvider(&Config{Provider: "openai", APIKey: "test-openai-key"})
	if err != nil || openAIP.Name() != "openai" {
		t.Fatalf("failed to create openai provider: %v", err)
	}

	// 5. Ollama
	ollamaP, err := NewProvider(&Config{Provider: "ollama", BaseURL: "http://localhost:11434"})
	if err != nil || ollamaP.Name() != "ollama" {
		t.Fatalf("failed to create ollama provider: %v", err)
	}
}

func TestOpenAIProviderMockHTTP(t *testing.T) {
	// Mock OpenAI API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"choices": [
				{
					"message": {
						"content": "{\"original_title\":\"API Rules\",\"summary\":\"split\",\"sub_rules\":[{\"id\":\"api\",\"title\":\"API Guidelines\",\"description\":\"desc\",\"is_skill\":false,\"content\":\"details\"}]}"
					}
				}
			]
		}`))
	}))
	defer server.Close()

	provider, err := NewOpenAIProvider("test-key", "gpt-4o", server.URL, 5*time.Second)
	if err != nil {
		t.Fatalf("failed to create openai provider: %v", err)
	}

	result, err := provider.DecomposeRule(context.Background(), "API Rules", "some content", 400)
	if err != nil {
		t.Fatalf("unexpected decompose error: %v", err)
	}

	if len(result.SubRules) != 1 || result.SubRules[0].ID != "api" {
		t.Fatalf("unexpected decompose result: %+v", result)
	}
}
