package appointments

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/auth"
)

type handlerRepository struct {
	fakeRepository
}

func TestHandlerCreate(t *testing.T) {
	handler := NewHandler(NewService(&handlerRepository{}))
	body := `{"professional_id":"professional-1","service_id":"service-1","contact_name":"Maria Perez","starts_at":"2026-06-02T14:00:00Z","duration_minutes":60}`
	req := httptest.NewRequest(http.MethodPost, "/api/appointments", strings.NewReader(body))
	req = req.WithContext(auth.ContextWithClaims(req.Context(), auth.Claims{UserID: "user-1", ClinicID: "clinic-1", Role: auth.RoleAssistant}))
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. Body: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "appointment-new") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestHandlerStatusConflict(t *testing.T) {
	handler := NewHandler(NewService(&handlerRepository{fakeRepository: fakeRepository{err: ErrScheduleConflict}}))
	body := `{"professional_id":"professional-1","service_id":"service-1","contact_name":"Maria Perez","starts_at":"2026-06-02T14:00:00Z","duration_minutes":60}`
	req := httptest.NewRequest(http.MethodPost, "/api/appointments", strings.NewReader(body))
	req = req.WithContext(auth.ContextWithClaims(req.Context(), auth.Claims{UserID: "user-1", ClinicID: "clinic-1", Role: auth.RoleAssistant}))
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d. Body: %s", rec.Code, rec.Body.String())
	}
}

func (r *handlerRepository) List(ctx context.Context, clinicID string, filters ListFilters) ([]Appointment, error) {
	if r.err != nil {
		return nil, r.err
	}
	return []Appointment{
		{
			ID:                 "appointment-1",
			ClinicID:           clinicID,
			ProfessionalID:     "professional-1",
			ProfessionalName:   "Dra. Ana",
			ServiceID:          "service-1",
			ServiceName:        "Ortodoncia",
			ContactName:        "Maria",
			StartsAt:           time.Date(2026, 6, 2, 14, 0, 0, 0, time.UTC),
			EndsAt:             time.Date(2026, 6, 2, 15, 0, 0, 0, time.UTC),
			Status:             "scheduled",
			ConfirmationStatus: "pending",
			Source:             "manual",
		},
	}, nil
}

func TestHandlerList(t *testing.T) {
	handler := NewHandler(NewService(&handlerRepository{}))
	req := httptest.NewRequest(http.MethodGet, "/api/appointments?date=2026-06-02", nil)
	req = req.WithContext(auth.ContextWithClaims(req.Context(), auth.Claims{ClinicID: "clinic-1", Role: auth.RoleAssistant}))
	rec := httptest.NewRecorder()

	handler.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "appointment-1") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
