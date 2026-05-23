package ai

import (
	"errors"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/config"
)

func NewProvider(cfg config.Config) (Provider, error) {
	switch cfg.AIProvider {
	case "mock":
		return NewMockProvider(), nil
	case "openai":
		if cfg.OpenAIKey == "" {
			return nil, errors.New("OPENAI_API_KEY is required for openai provider")
		}
		model := cfg.AIModel
		if model == "" {
			model = "gpt-4o-mini"
		}
		return NewOpenAIProvider(cfg.OpenAIKey, model), nil
	case "gemini":
		if cfg.GeminiKey == "" {
			return nil, errors.New("GEMINI_API_KEY is required for gemini provider")
		}
		return NewGeminiProvider(cfg.GeminiKey, cfg.AIModel), nil
	case "deepseek":
		if cfg.DeepSeekKey == "" {
			return nil, errors.New("DEEPSEEK_API_KEY is required for deepseek provider")
		}
		return NewDeepSeekProvider(cfg.DeepSeekKey, cfg.AIModel), nil
	default:
		return nil, errors.New("unsupported ai provider: " + cfg.AIProvider)
	}
}
