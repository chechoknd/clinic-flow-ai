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

func TestAIRequestValidationAllowsPartialContext(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "reply without lead or service",
			err:  validateReplySuggestionRequest(ReplySuggestionRequest{PatientMessage: "Quiero saber el precio del tratamiento"}),
		},
		{
			name: "objection with only lead",
			err:  validateObjectionHandlerRequest(ObjectionHandlerRequest{LeadID: "lead-1", Objection: "Es costoso este tratamiento ?"}),
		},
		{
			name: "follow up with only lead",
			err:  validateFollowUpMessageRequest(FollowUpMessageRequest{LeadID: "lead-1"}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err != nil {
				t.Fatalf("expected valid request, got %v", tt.err)
			}
		})
	}
}
