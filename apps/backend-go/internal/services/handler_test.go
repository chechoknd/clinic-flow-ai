package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/auth"
)

type handlerRepository struct {
	entities []ServiceEntity
	err      error
}

func (r *handlerRepository) ListByClinicID(ctx context.Context, clinicID string) ([]ServiceEntity, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.entities, nil
}

func (r *handlerRepository) FindByID(ctx context.Context, clinicID, serviceID string) (ServiceEntity, error) {
	if r.err != nil {
		return ServiceEntity{}, r.err
	}
	for _, e := range r.entities {
		if e.ID == serviceID && e.ClinicID == clinicID {
			return e, nil
		}
	}
	return ServiceEntity{}, ErrServiceNotFound
}

func (r *handlerRepository) Create(ctx context.Context, s ServiceEntity) (ServiceEntity, error) {
	if r.err != nil {
		return ServiceEntity{}, r.err
	}
	s.ID = "new-id"
	s.CurrencyCode = "COP"
	r.entities = append(r.entities, s)
	return s, nil
}

func (r *handlerRepository) Update(ctx context.Context, s ServiceEntity) error {
	return r.err
}

func (r *handlerRepository) Delete(ctx context.Context, clinicID, serviceID string) error {
	return r.err
}

func (r *handlerRepository) CurrencyCodeByClinicID(ctx context.Context, clinicID string) (string, error) {
	if r.err != nil {
		return "", r.err
	}
	return "COP", nil
}

func TestHandlerList(t *testing.T) {
	repo := &handlerRepository{entities: []ServiceEntity{
		{ID: "1", ClinicID: "clinic-1", Name: "Service 1", CurrencyCode: "COP", IsActive: true, Benefits: []byte("[]"), FAQ: []byte("[]"), CommonObjections: []byte("[]")},
	}}
	h := NewHandler(NewService(repo))

	req := httptest.NewRequest(http.MethodGet, "/api/services", nil)
	req = req.WithContext(auth.ContextWithClaims(req.Context(), auth.Claims{ClinicID: "clinic-1", Role: auth.RoleAssistant}))
	rec := httptest.NewRecorder()

	h.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Service 1") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestHandlerCreate(t *testing.T) {
	repo := &handlerRepository{}
	h := NewHandler(NewService(repo))

	body := `{"name":"New Service"}`
	req := httptest.NewRequest(http.MethodPost, "/api/services", strings.NewReader(body))
	req = req.WithContext(auth.ContextWithClaims(req.Context(), auth.Claims{ClinicID: "clinic-1", Role: auth.RoleClinicAdmin}))
	rec := httptest.NewRecorder()

	h.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. Body: %s", rec.Code, rec.Body.String())
	}
}
