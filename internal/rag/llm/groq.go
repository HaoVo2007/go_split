package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type LLM interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type GroqLLM struct {
	BaseURL string
	Model   string
	APIKey  string
	Client  *http.Client
}

func NewGroqLLM(baseURL, model, apiKey string) *GroqLLM {
	return &GroqLLM{
		BaseURL: strings.TrimRight(baseURL, "/"),
		Model:   model,
		APIKey:  apiKey,
		Client:  &http.Client{},
	}
}

type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqChatRequest struct {
	Model    string        `json:"model"`
	Messages []groqMessage `json:"messages"`
}

type groqChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (l *GroqLLM) Generate(ctx context.Context, prompt string) (string, error) {
	payload, err := json.Marshal(groqChatRequest{
		Model: l.Model,
		Messages: []groqMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, l.BaseURL+"/openai/v1/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+l.APIKey)

	resp, err := l.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to generate response: %s: %s", resp.Status, string(body))
	}

	var result groqChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("empty llm response")
	}

	return result.Choices[0].Message.Content, nil
}
