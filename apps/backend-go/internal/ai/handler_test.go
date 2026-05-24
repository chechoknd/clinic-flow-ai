package ai

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/auth"
)

func TestReplySuggestionValidationError(t *testing.T) {
	handler := NewHandler(&Service{})
	req := httptest.NewRequest(http.MethodPost, "/api/ai/reply-suggestion", strings.NewReader(`{}`))
	req = req.WithContext(auth.ContextWithClaims(req.Context(), auth.Claims{ClinicID: "clinic-1", UserID: "user-1", Role: auth.RoleAssistant}))
	w := httptest.NewRecorder()

	handler.ReplySuggestion(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "VALIDATION_ERROR") || !strings.Contains(w.Body.String(), "patient_message") {
		t.Fatalf("unexpected response body: %s", w.Body.String())
	}
}

func TestObjectionHandlerValidationError(t *testing.T) {
	handler := NewHandler(&Service{})
	req := httptest.NewRequest(http.MethodPost, "/api/ai/objection-handler", strings.NewReader(`{"lead_id":"lead-1","service_id":"service-1"}`))
	req = req.WithContext(auth.ContextWithClaims(req.Context(), auth.Claims{ClinicID: "clinic-1", UserID: "user-1", Role: auth.RoleAssistant}))
	w := httptest.NewRecorder()

	handler.ObjectionHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "VALIDATION_ERROR") || !strings.Contains(w.Body.String(), "objection") {
		t.Fatalf("unexpected response body: %s", w.Body.String())
	}
}

func TestAITypedErrors(t *testing.T) {
	if !errors.Is(ValidationError{Field: "lead_id", Message: "required"}, ErrValidation) {
		t.Fatal("ValidationError should match ErrValidation")
	}
	if !errors.Is(ProviderError{Provider: "deepseek", Err: errors.New("timeout")}, ErrAIProvider) {
		t.Fatal("ProviderError should match ErrAIProvider")
	}
	if !errors.Is(ResponseError{Err: errors.New("invalid json")}, ErrAIResponse) {
		t.Fatal("ResponseError should match ErrAIResponse")
	}
	if !errors.Is(SafetyBlockedError{Status: "passed_with_warnings"}, ErrAISafetyBlocked) {
		t.Fatal("SafetyBlockedError should match ErrAISafetyBlocked")
	}
}
