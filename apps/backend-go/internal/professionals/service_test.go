package professionals

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRepository struct {
	entities   []Professional
	lastFilter ListFilters
	err        error
}

func (r *fakeRepository) List(ctx context.Context, clinicID string, filters ListFilters) ([]Professional, error) {
	r.lastFilter = filters
	if r.err != nil {
		return nil, r.err
	}
	return r.entities, nil
}

func (r *fakeRepository) FindByID(ctx context.Context, clinicID, professionalID string) (Professional, error) {
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

func (r *fakeRepository) Create(ctx context.Context, p Professional, serviceIDs []string) (Professional, error) {
	if r.err != nil {
		return Professional{}, r.err
	}
	p.ID = "professional-new"
	p.ServiceIDs = serviceIDsJSON(serviceIDs)
	p.CreatedAt = time.Now()
	p.UpdatedAt = p.CreatedAt
	r.entities = append(r.entities, p)
	return p, nil
}

func (r *fakeRepository) Update(ctx context.Context, p Professional, serviceIDs []string) (Professional, error) {
	if r.err != nil {
		return Professional{}, r.err
	}
	for i, entity := range r.entities {
		if entity.ID == p.ID && entity.ClinicID == p.ClinicID {
			p.ServiceIDs = serviceIDsJSON(serviceIDs)
			p.CreatedAt = entity.CreatedAt
			p.UpdatedAt = time.Now()
			r.entities[i] = p
			return p, nil
		}
	}
	return Professional{}, ErrProfessionalNotFound
}

func TestServiceCreateProfessional(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo)

	color := "#2563EB"
	role := "Ortodoncia"
	res, err := service.Create(context.Background(), "clinic-1", CreateProfessionalRequest{
		FullName:        " Dra. Ana Gomez ",
		RoleOrSpecialty: &role,
		CalendarColor:   &color,
		WorkingHours:    []byte(`{"monday":[{"start":"08:00","end":"12:00"}]}`),
		ServiceIDs:      []string{"service-1", "service-1", "service-2"},
	})
	if err != nil {
		t.Fatalf("create professional: %v", err)
	}
	if res.ID != "professional-new" || res.FullName != "Dra. Ana Gomez" {
		t.Fatalf("unexpected response: %#v", res)
	}
	if string(res.ServiceIDs) != `["service-1","service-2"]` {
		t.Fatalf("service ids not normalized: %s", string(res.ServiceIDs))
	}
}

func TestServiceCreateValidatesName(t *testing.T) {
	service := NewService(&fakeRepository{})

	_, err := service.Create(context.Background(), "clinic-1", CreateProfessionalRequest{})
	if !errors.Is(err, ErrProfessionalNameRequired) {
		t.Fatalf("expected ErrProfessionalNameRequired, got %v", err)
	}
}

func TestServiceCreateValidatesCalendarColor(t *testing.T) {
	service := NewService(&fakeRepository{})
	color := "blue"

	_, err := service.Create(context.Background(), "clinic-1", CreateProfessionalRequest{
		FullName:      "Dra. Ana Gomez",
		CalendarColor: &color,
	})
	if !errors.Is(err, ErrInvalidCalendarColor) {
		t.Fatalf("expected ErrInvalidCalendarColor, got %v", err)
	}
}

func TestServiceCreateValidatesWorkingHoursObject(t *testing.T) {
	service := NewService(&fakeRepository{})

	_, err := service.Create(context.Background(), "clinic-1", CreateProfessionalRequest{
		FullName:     "Dra. Ana Gomez",
		WorkingHours: []byte(`[]`),
	})
	if !errors.Is(err, ErrInvalidWorkingHours) {
		t.Fatalf("expected ErrInvalidWorkingHours, got %v", err)
	}
}

func TestServiceListPassesFilters(t *testing.T) {
	repo := &fakeRepository{entities: []Professional{
		{ID: "professional-1", ClinicID: "clinic-1", FullName: "Dra. Ana", WorkingHours: []byte(`{}`), ServiceIDs: []byte(`[]`)},
	}}
	service := NewService(repo)
	active := true

	res, err := service.List(context.Background(), "clinic-1", ListFilters{IsActive: &active, ServiceID: "service-1"})
	if err != nil {
		t.Fatalf("list professionals: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 professional, got %d", len(res))
	}
	if repo.lastFilter.IsActive == nil || !*repo.lastFilter.IsActive || repo.lastFilter.ServiceID != "service-1" {
		t.Fatalf("filters were not passed: %#v", repo.lastFilter)
	}
}
