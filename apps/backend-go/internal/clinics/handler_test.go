package clinics

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/auth"
)

type handlerRepository struct {
	clinic Clinic
	err    error
}

func (r handlerRepository) FindByID(ctx context.Context, clinicID string) (Clinic, error) {
	if r.err != nil {
		return Clinic{}, r.err
	}
	return r.clinic, nil
}

func TestHandlerCurrentRequiresClaims(t *testing.T) {
	handler := NewHandler(NewService(handlerRepository{}))
	req := httptest.NewRequest(http.MethodGet, "/api/clinics/current", nil)
	rec := httptest.NewRecorder()

	handler.Current(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestHandlerCurrentReturnsClinic(t *testing.T) {
	handler := NewHandler(NewService(handlerRepository{clinic: Clinic{
		ID:                "clinic-id",
		Name:              "Sonrisa Viva Demo",
		ClinicType:        "odontologia",
		City:              "Bogota",
		WhatsApp:          "+573001112233",
		OpeningHours:      []byte(`{}`),
		GeneralFAQ:        []byte(`[]`),
		CommunicationTone: "amable",
	}}))
	req := httptest.NewRequest(http.MethodGet, "/api/clinics/current", nil)
	req = req.WithContext(auth.ContextWithClaims(req.Context(), auth.Claims{ClinicID: "clinic-id", Role: auth.RoleClinicAdmin, UserID: "user-id"}))
	rec := httptest.NewRecorder()

	handler.Current(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Sonrisa Viva Demo") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
