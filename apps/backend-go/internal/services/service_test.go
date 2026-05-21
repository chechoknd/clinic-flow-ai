package services

import (
	"context"
	"testing"
)

type fakeRepository struct {
	entities []ServiceEntity
	err      error
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
		{ID: "1", ClinicID: "clinic-1", Name: "Service 1"},
		{ID: "2", ClinicID: "clinic-1", Name: "Service 2"},
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
