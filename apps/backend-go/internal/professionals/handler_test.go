package professionals

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
	entities   []Professional
	lastFilter ListFilters
	err        error
}

func (r *handlerRepository) List(ctx context.Context, clinicID string, filters ListFilters) ([]Professional, error) {
	r.lastFilter = filters
	if r.err != nil {
		return nil, r.err
	}
	return r.entities, nil
}

func (r *handlerRepository) FindByID(ctx context.Context, clinicID, professionalID string) (Professional, error) {
	if r.err != nil {
		return Professional{}, r.err
	}
	for _, entity := range r.entities {
		if entity.ID == professionalID && entity.ClinicID == clinicID {
			return entity, nil
		}
	}
	return Professional{}, ErrProfessionalNotFound
}

func (r *handlerRepository) Create(ctx context.Context, p Professional, serviceIDs []string) (Professional, error) {
	if r.err != nil {
		return Professional{}, r.err
	}
	p.ID = "professional-new"
	p.ServiceIDs = serviceIDsJSON(serviceIDs)
	p.CreatedAt = time.Now()
	p.UpdatedAt = p.CreatedAt
	return p, nil
}

func (r *handlerRepository) Update(ctx context.Context, p Professional, serviceIDs []string) (Professional, error) {
	if r.err != nil {
		return Professional{}, r.err
	}
	p.ServiceIDs = serviceIDsJSON(serviceIDs)
	p.CreatedAt = time.Now()
	p.UpdatedAt = p.CreatedAt
	return p, nil
}

func TestHandlerList(t *testing.T) {
	repo := &handlerRepository{entities: []Professional{
		{ID: "professional-1", ClinicID: "clinic-1", FullName: "Dra. Ana", WorkingHours: []byte(`{}`), IsActive: true, ServiceIDs: []byte(`[]`)},
	}}
	handler := NewHandler(NewService(repo))

	req := httptest.NewRequest(http.MethodGet, "/api/professionals?is_active=true&service_id=service-1", nil)
	req = req.WithContext(auth.ContextWithClaims(req.Context(), auth.Claims{ClinicID: "clinic-1", Role: auth.RoleAssistant}))
	rec := httptest.NewRecorder()

	handler.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Dra. Ana") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
	if repo.lastFilter.IsActive == nil || !*repo.lastFilter.IsActive || repo.lastFilter.ServiceID != "service-1" {
		t.Fatalf("filters were not parsed: %#v", repo.lastFilter)
	}
}

func TestHandlerCreate(t *testing.T) {
	repo := &handlerRepository{}
	handler := NewHandler(NewService(repo))

	body := `{"full_name":"Dra. Ana Gomez","calendar_color":"#2563EB","working_hours":{},"service_ids":["service-1"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/professionals", strings.NewReader(body))
	req = req.WithContext(auth.ContextWithClaims(req.Context(), auth.Claims{ClinicID: "clinic-1", Role: auth.RoleClinicAdmin}))
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. Body: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "professional-new") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestHandlerCreateValidationError(t *testing.T) {
	repo := &handlerRepository{}
	handler := NewHandler(NewService(repo))

	body := `{"calendar_color":"blue"}`
	req := httptest.NewRequest(http.MethodPost, "/api/professionals", strings.NewReader(body))
	req = req.WithContext(auth.ContextWithClaims(req.Context(), auth.Claims{ClinicID: "clinic-1", Role: auth.RoleClinicAdmin}))
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d. Body: %s", rec.Code, rec.Body.String())
	}
}
