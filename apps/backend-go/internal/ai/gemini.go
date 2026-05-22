package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type GeminiProvider struct {
	apiKey string
	model  string
	client *http.Client
}

func NewGeminiProvider(apiKey, model string) *GeminiProvider {
	if model == "" {
		model = "gemini-1.5-flash" // Sensible default for speed/cost
	}
	return &GeminiProvider{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *GeminiProvider) Name() string {
	return "gemini"
}

func (p *GeminiProvider) Generate(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	// Gemini uses a different URL structure: https://generativelanguage.googleapis.com/v1beta/models/{model}:generateContent?key={apiKey}
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", p.model, p.apiKey)

	// Combine system prompt into the instructions if the API version supports it,
	// or just prepend it to the first message for simplicity in MVP.
	fullPrompt := fmt.Sprintf("System Instructions:\n%s\n\nUser Message:\n%s", systemPrompt, userPrompt)

	requestBody := map[string]any{
		"contents": []map[string]any{
			{
				"parts": []map[string]string{
					{"text": fullPrompt},
				},
			},
		},
		"generationConfig": map[string]any{
			"temperature":      0.7,
			"topK":             40,
			"topP":             0.95,
			"maxOutputTokens":  1024,
			"responseMimeType": "application/json",
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error struct {
				Message string `json:"message"`
				Code    int    `json:"code"`
				Status  string `json:"status"`
			} `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return "", fmt.Errorf("gemini api error (%s): %s", errResp.Error.Status, errResp.Error.Message)
	}

	var res struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("no content candidates returned from gemini")
	}

	return res.Candidates[0].Content.Parts[0].Text, nil
}
