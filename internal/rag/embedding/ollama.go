package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Embedder interface {
	Embed(ctx context.Context, text string) ([]float64, error)
}

type OllamaEmbedder struct {
	BaseURL string
	Model   string
	Client  *http.Client
}

func NewOllamaEmbedder(baseURL, model string) *OllamaEmbedder {
	return &OllamaEmbedder{
		BaseURL: strings.TrimRight(baseURL, "/"),
		Model:   model,
		Client:  &http.Client{},
	}
}

type ollamaEmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type ollamaEmbedResponse struct {
	Embeddings [][]float64 `json:"embeddings"`
}

func (e *OllamaEmbedder) Embed(ctx context.Context, text string) ([]float64, error) {
	payload, err := json.Marshal(ollamaEmbedRequest{
		Model: e.Model,
		Input: text,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.BaseURL+"/api/embed", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to embed text: %s: %s", resp.Status, string(body))
	}

	var result ollamaEmbedResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if len(result.Embeddings) == 0 || len(result.Embeddings[0]) == 0 {
		return nil, fmt.Errorf("empty embedding response")
	}

	return result.Embeddings[0], nil
}
