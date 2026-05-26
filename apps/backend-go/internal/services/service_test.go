package services

import (
	"context"
	"testing"
)

type fakeRepository struct {
	entities     []ServiceEntity
	err          error
	currencyCode string
}

func (r *fakeRepository) ListByClinicID(ctx context.Context, clinicID string) ([]ServiceEntity, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.entities, nil
}

func (r *fakeRepository) FindByID(ctx context.Context, clinicID, serviceID string) (ServiceEntity, error) {
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

func (r *fakeRepository) Create(ctx context.Context, s ServiceEntity) (ServiceEntity, error) {
	if r.err != nil {
		return ServiceEntity{}, r.err
	}
	s.ID = "new-id"
	r.entities = append(r.entities, s)
	return s, nil
}

func (r *fakeRepository) Update(ctx context.Context, s ServiceEntity) error {
	if r.err != nil {
		return r.err
	}
	for i, e := range r.entities {
		if e.ID == s.ID && e.ClinicID == s.ClinicID {
			r.entities[i] = s
			return nil
		}
	}
	return ErrServiceNotFound
}

func (r *fakeRepository) CurrencyCodeByClinicID(ctx context.Context, clinicID string) (string, error) {
	if r.err != nil {
		return "", r.err
	}
	if r.currencyCode != "" {
		return r.currencyCode, nil
	}
	return "COP", nil
}

func (r *fakeRepository) Delete(ctx context.Context, clinicID, serviceID string) error {
	if r.err != nil {
		return r.err
	}
	for i, e := range r.entities {
		if e.ID == serviceID && e.ClinicID == clinicID {
			r.entities = append(r.entities[:i], r.entities[i+1:]...)
			return nil
		}
	}
	return ErrServiceNotFound
}

func TestServiceList(t *testing.T) {
	repo := &fakeRepository{entities: []ServiceEntity{
		{ID: "1", ClinicID: "clinic-1", Name: "Service 1", CurrencyCode: "COP"},
		{ID: "2", ClinicID: "clinic-1", Name: "Service 2", CurrencyCode: "COP"},
	}}
	s := NewService(repo)

	res, err := s.List(context.Background(), "clinic-1")
	if err != nil {
		t.Fatalf("list services: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 services, got %d", len(res))
	}
}

func TestServiceCreate(t *testing.T) {
	repo := &fakeRepository{}
	s := NewService(repo)

	req := CreateServiceRequest{Name: "New Service"}
	res, err := s.Create(context.Background(), "clinic-1", req)
	if err != nil {
		t.Fatalf("create service: %v", err)
	}
	if res.ID != "new-id" || res.Name != "New Service" {
		t.Fatalf("unexpected response: %#v", res)
	}
}

func TestServiceCreateRejectsDecimalsForZeroDecimalCurrency(t *testing.T) {
	price := 120000.50
	repo := &fakeRepository{currencyCode: "COP"}
	s := NewService(repo)

	_, err := s.Create(context.Background(), "clinic-1", CreateServiceRequest{Name: "Limpieza", PriceFrom: &price})
	if err == nil {
		t.Fatal("expected decimal validation error")
	}
}

func TestServiceCreateAllowsTwoDecimalsForTwoDecimalCurrency(t *testing.T) {
	price := 120.50
	repo := &fakeRepository{currencyCode: "PEN"}
	s := NewService(repo)

	_, err := s.Create(context.Background(), "clinic-1", CreateServiceRequest{Name: "Consulta", PriceFrom: &price})
	if err != nil {
		t.Fatalf("expected price to be accepted: %v", err)
	}
}
