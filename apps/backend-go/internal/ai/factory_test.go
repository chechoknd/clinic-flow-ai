package ai

import (
	"strings"
	"testing"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/config"
)

func TestNewProviderDeepSeek(t *testing.T) {
	provider, err := NewProvider(config.Config{
		AIProvider:  "deepseek",
		DeepSeekKey: "test-key",
	})
	if err != nil {
		t.Fatalf("NewProvider returned error: %v", err)
	}

	if provider.Name() != "deepseek" {
		t.Fatalf("expected provider name deepseek, got %s", provider.Name())
	}

	openAICompatible, ok := provider.(*OpenAIProvider)
	if !ok {
		t.Fatalf("expected OpenAIProvider, got %T", provider)
	}
	if openAICompatible.baseURL != "https://api.deepseek.com" {
		t.Fatalf("expected DeepSeek base URL, got %s", openAICompatible.baseURL)
	}
	if openAICompatible.model != "deepseek-chat" {
		t.Fatalf("expected default DeepSeek model, got %s", openAICompatible.model)
	}
}

func TestNewProviderDeepSeekUsesConfiguredModel(t *testing.T) {
	provider, err := NewProvider(config.Config{
		AIProvider:  "deepseek",
		AIModel:     "deepseek-reasoner",
		DeepSeekKey: "test-key",
	})
	if err != nil {
		t.Fatalf("NewProvider returned error: %v", err)
	}

	openAICompatible := provider.(*OpenAIProvider)
	if openAICompatible.model != "deepseek-reasoner" {
		t.Fatalf("expected configured model, got %s", openAICompatible.model)
	}
}

func TestNewProviderDeepSeekRequiresKey(t *testing.T) {
	_, err := NewProvider(config.Config{AIProvider: "deepseek"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "DEEPSEEK_API_KEY") {
		t.Fatalf("expected DEEPSEEK_API_KEY error, got %v", err)
	}
}
