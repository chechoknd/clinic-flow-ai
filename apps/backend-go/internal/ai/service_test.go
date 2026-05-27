package ai

import (
	"context"
	"strings"
	"testing"
)

type fakeGenerationRepository struct {
	record GenerationRecord
}

func (r *fakeGenerationRepository) CreateGeneration(ctx context.Context, record GenerationRecord) (string, error) {
	r.record = record
	return "generation-1", nil
}

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

func TestBuildAnalyzeConversationUserPrompt(t *testing.T) {
	prompt := BuildAnalyzeConversationUserPrompt("Paciente: Quiero saber precio", "whatsapp")
	if !strings.Contains(prompt, "Paciente: Quiero saber precio") || !strings.Contains(prompt, "detected_lead") {
		t.Fatalf("prompt missing conversation or expected JSON shape: %s", prompt)
	}
	if !strings.Contains(prompt, "No diagnostiques") {
		t.Fatalf("prompt missing safety rules: %s", prompt)
	}
}

func TestAnalyzeConversationValidation(t *testing.T) {
	if err := validateAnalyzeConversationRequest(AnalyzeConversationRequest{ConversationText: "Paciente pregunta por blanqueamiento y precio"}); err != nil {
		t.Fatalf("expected valid request, got %v", err)
	}

	err := validateAnalyzeConversationRequest(AnalyzeConversationRequest{ConversationText: ""})
	if err == nil || !strings.Contains(err.Error(), "conversation_text") {
		t.Fatalf("expected conversation_text validation error, got %v", err)
	}
}

func TestRecordGenerationStoresMetadataOnly(t *testing.T) {
	repo := &fakeGenerationRepository{}
	service := &Service{
		provider:    NewMockProvider(),
		generations: repo,
	}

	id := service.recordGeneration(context.Background(), generationRecordInput{
		ClinicID:        "clinic-1",
		UserID:          "user-1",
		Feature:         "reply_suggestion",
		Status:          "success",
		SafetyStatus:    "passed",
		InputCharCount:  120,
		OutputCharCount: 80,
	})

	if id != "generation-1" {
		t.Fatalf("unexpected generation id: %s", id)
	}
	if repo.record.ClinicID != "clinic-1" || repo.record.UserID != "user-1" {
		t.Fatalf("unexpected tenant/user metadata: %#v", repo.record)
	}
	if repo.record.Provider != "mock" || repo.record.Model != "mock" {
		t.Fatalf("unexpected provider metadata: %#v", repo.record)
	}
	if repo.record.Feature != "reply_suggestion" || repo.record.Status != "success" || repo.record.SafetyStatus != "passed" {
		t.Fatalf("unexpected generation status metadata: %#v", repo.record)
	}
	if repo.record.InputCharCount != 120 || repo.record.OutputCharCount != 80 {
		t.Fatalf("unexpected char count metadata: %#v", repo.record)
	}
}
