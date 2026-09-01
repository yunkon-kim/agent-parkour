package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OllamaProvider struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewOllamaProvider(baseURL, model string, timeout time.Duration) (*OllamaProvider, error) {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	if model == "" {
		model = "deepseek-r1:14b"
	}
	return &OllamaProvider{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		model:   model,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

func (p *OllamaProvider) Name() string {
	return "ollama"
}

func (p *OllamaProvider) Model() string {
	return p.model
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Format string `json:"format,omitempty"`
}

type ollamaResponse struct {
	Response string `json:"response"`
	Error    string `json:"error,omitempty"`
}

func (p *OllamaProvider) generateText(ctx context.Context, prompt, format string) (string, error) {
	url := fmt.Sprintf("%s/api/generate", p.baseURL)

	reqBody := ollamaRequest{
		Model:  p.model,
		Prompt: prompt,
		Stream: false,
		Format: format,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama request failed (make sure ollama is running at %s): %w", p.baseURL, err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading ollama response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var ollamaResp ollamaResponse
	if err := json.Unmarshal(respBytes, &ollamaResp); err != nil {
		return "", fmt.Errorf("parsing ollama response failed: %w", err)
	}

	if ollamaResp.Error != "" {
		return "", fmt.Errorf("ollama error: %s", ollamaResp.Error)
	}

	return ollamaResp.Response, nil
}

func (p *OllamaProvider) DecomposeRule(ctx context.Context, title, content string, maxTokens int) (*DecomposeResult, error) {
	prompt := BuildDecomposePrompt(title, content, maxTokens)
	rawOutput, err := p.generateText(ctx, prompt, "json")
	if err != nil {
		return nil, err
	}

	cleaned := cleanJSONOutput(rawOutput)
	var result DecomposeResult
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("failed to parse decomposition JSON from Ollama: %w (raw: %s)", err, rawOutput)
	}

	return &result, nil
}

func (p *OllamaProvider) GenerateDescription(ctx context.Context, content string) (string, error) {
	prompt := BuildDescriptionPrompt(content)
	rawOutput, err := p.generateText(ctx, prompt, "")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(rawOutput), nil
}
