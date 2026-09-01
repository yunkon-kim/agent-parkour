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

type ClaudeProvider struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

func NewClaudeProvider(apiKey, model string, timeout time.Duration) (*ClaudeProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY is required for claude provider")
	}
	if model == "" {
		model = "claude-3-7-sonnet-20250219"
	}
	return &ClaudeProvider{
		apiKey: apiKey,
		model:  model,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

func (p *ClaudeProvider) Name() string {
	return "anthropic"
}

func (p *ClaudeProvider) Model() string {
	return p.model
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []claudeMessage `json:"messages"`
}

type claudeResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (p *ClaudeProvider) generateText(ctx context.Context, prompt string) (string, error) {
	url := "https://api.anthropic.com/v1/messages"

	reqBody := claudeRequest{
		Model:     p.model,
		MaxTokens: 4096,
		Messages: []claudeMessage{
			{Role: "user", Content: prompt},
		},
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
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("claude api request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading claude response body failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("claude api returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(respBytes, &claudeResp); err != nil {
		return "", fmt.Errorf("parsing claude response failed: %w", err)
	}

	if claudeResp.Error != nil {
		return "", fmt.Errorf("claude api error (%s): %s", claudeResp.Error.Type, claudeResp.Error.Message)
	}

	if len(claudeResp.Content) == 0 {
		return "", fmt.Errorf("claude returned empty content")
	}

	return claudeResp.Content[0].Text, nil
}

func (p *ClaudeProvider) DecomposeRule(ctx context.Context, title, content string, maxTokens int) (*DecomposeResult, error) {
	prompt := BuildDecomposePrompt(title, content, maxTokens)
	rawOutput, err := p.generateText(ctx, prompt)
	if err != nil {
		return nil, err
	}

	cleaned := cleanJSONOutput(rawOutput)
	var result DecomposeResult
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("failed to parse decomposition JSON from Claude: %w (raw: %s)", err, rawOutput)
	}

	return &result, nil
}

func (p *ClaudeProvider) GenerateDescription(ctx context.Context, content string) (string, error) {
	prompt := BuildDescriptionPrompt(content)
	rawOutput, err := p.generateText(ctx, prompt)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(rawOutput), nil
}
