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

type OpenAIProvider struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

func NewOpenAIProvider(apiKey, model, baseURL string, timeout time.Duration) (*OpenAIProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is required for openai provider")
	}
	if model == "" {
		model = "gpt-4o"
	}
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &OpenAIProvider{
		apiKey:  apiKey,
		model:   model,
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

func (p *OpenAIProvider) Name() string {
	return "openai"
}

func (p *OpenAIProvider) Model() string {
	return p.model
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIRequest struct {
	Model    string          `json:"model"`
	Messages []openAIMessage `json:"messages"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

func (p *OpenAIProvider) generateText(ctx context.Context, prompt string) (string, error) {
	url := fmt.Sprintf("%s/chat/completions", p.baseURL)

	reqBody := openAIRequest{
		Model: p.model,
		Messages: []openAIMessage{
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("openai api request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading openai response body failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openai api returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var openAIResp openAIResponse
	if err := json.Unmarshal(respBytes, &openAIResp); err != nil {
		return "", fmt.Errorf("parsing openai response failed: %w", err)
	}

	if openAIResp.Error != nil {
		return "", fmt.Errorf("openai api error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("openai returned empty choices")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

func (p *OpenAIProvider) DecomposeRule(ctx context.Context, title, content string, maxTokens int) (*DecomposeResult, error) {
	prompt := BuildDecomposePrompt(title, content, maxTokens)
	rawOutput, err := p.generateText(ctx, prompt)
	if err != nil {
		return nil, err
	}

	cleaned := cleanJSONOutput(rawOutput)
	var result DecomposeResult
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("failed to parse decomposition JSON from OpenAI: %w (raw: %s)", err, rawOutput)
	}

	return &result, nil
}

func (p *OpenAIProvider) GenerateDescription(ctx context.Context, content string) (string, error) {
	prompt := BuildDescriptionPrompt(content)
	rawOutput, err := p.generateText(ctx, prompt)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(rawOutput), nil
}
