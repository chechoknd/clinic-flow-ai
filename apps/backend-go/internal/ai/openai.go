package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type OpenAIProvider struct {
	apiKey       string
	model        string
	baseURL      string
	providerName string
	client       *http.Client
}

func NewOpenAIProvider(apiKey, model string) *OpenAIProvider {
	return newOpenAICompatibleProvider("openai", "https://api.openai.com/v1", apiKey, model)
}

func NewDeepSeekProvider(apiKey, model string) *OpenAIProvider {
	if model == "" {
		model = "deepseek-chat"
	}
	return newOpenAICompatibleProvider("deepseek", "https://api.deepseek.com", apiKey, model)
}

func newOpenAICompatibleProvider(providerName, baseURL, apiKey, model string) *OpenAIProvider {
	return &OpenAIProvider{
		apiKey:       apiKey,
		model:        model,
		baseURL:      baseURL,
		providerName: providerName,
		client:       &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *OpenAIProvider) Name() string {
	return p.providerName
}

func (p *OpenAIProvider) Generate(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	url := p.baseURL + "/chat/completions"

	requestBody := map[string]any{
		"model": p.model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"temperature":     0.7,
		"response_format": map[string]string{"type": "json_object"},
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
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return "", fmt.Errorf("%s api error (%d): %s", p.providerName, resp.StatusCode, errResp.Error.Message)
	}

	var res struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned from %s", p.providerName)
	}

	return res.Choices[0].Message.Content, nil
}
