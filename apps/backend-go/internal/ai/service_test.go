package ai

import (
	"strings"
	"testing"
)

func TestBuildSystemPrompt(t *testing.T) {
	ctx := Context{
		ClinicName:       "Sonrisa Viva",
		ClinicType:       "Odontología",
		City:             "Bogotá",
		ServicePriceFrom: "250000",
	}
	prompt := BuildSystemPrompt(ctx)
	if !strings.Contains(prompt, "Sonrisa Viva") || !strings.Contains(prompt, "Bogotá") || !strings.Contains(prompt, "250000") {
		t.Errorf("prompt missing context info: %s", prompt)
	}
}

func TestBuildReplySuggestionUserPrompt(t *testing.T) {
	prompt := BuildReplySuggestionUserPrompt("¿Precio?", "profesional")
	if !strings.Contains(prompt, "¿Precio?") || !strings.Contains(prompt, "profesional") {
		t.Errorf("prompt missing user message or tone: %s", prompt)
	}
}
