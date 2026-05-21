package clinics

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	clinic Clinic
	err    error
	updated bool
}

func (r *fakeRepository) FindByID(ctx context.Context, clinicID string) (Clinic, error) {
	if r.err != nil {
		return Clinic{}, r.err
	}
	return r.clinic, nil
}

func (r *fakeRepository) Update(ctx context.Context, clinic Clinic) error {
	if r.err != nil {
		return r.err
	}
	r.clinic = clinic
	r.updated = true
	return nil
}

func TestCurrentReturnsClinicResponse(t *testing.T) {
	phone := "+573001112233"
	address := "Calle 123"
	repo := &fakeRepository{clinic: Clinic{
		ID:                "clinic-id",
		Name:              "Sonrisa Viva Demo",
		ClinicType:        "odontologia",
		City:              "Bogota",
		Phone:             &phone,
		WhatsApp:          "+573001112233",
		Address:           &address,
		OpeningHours:      []byte(`{"monday_friday":"08:00-18:00"}`),
		GeneralFAQ:        []byte(`[]`),
		CommunicationTone: "amable",
	}}
	service := NewService(repo)

	res, err := service.Current(context.Background(), "clinic-id")
	if err != nil {
		t.Fatalf("current clinic: %v", err)
	}
	if res.ID != "clinic-id" || res.Name != "Sonrisa Viva Demo" || res.ClinicType != "odontologia" {
		t.Fatalf("unexpected response: %#v", res)
	}
	if string(res.OpeningHours) != `{"monday_friday":"08:00-18:00"}` {
		t.Fatalf("unexpected opening hours: %s", string(res.OpeningHours))
	}
}

func TestCurrentRequiresClinicID(t *testing.T) {
	service := NewService(&fakeRepository{})

	_, err := service.Current(context.Background(), "")
	if !errors.Is(err, ErrMissingClinicID) {
		t.Fatalf("expected ErrMissingClinicID, got %v", err)
	}
}

func TestUpdateChangesClinicData(t *testing.T) {
	repo := &fakeRepository{clinic: Clinic{
		ID:                "clinic-id",
		Name:              "Old Name",
		ClinicType:        "odontologia",
		City:              "Old City",
		WhatsApp:          "+573001112233",
		CommunicationTone: "amable",
	}}
	service := NewService(repo)

	req := UpdateClinicRequest{
		Name:              "New Name",
		City:              "New City",
		WhatsApp:          "+573009998877",
		CommunicationTone: "profesional",
		OpeningHours:      []byte(`{"all":"24h"}`),
		GeneralFAQ:        []byte(`[]`),
	}

	res, err := service.Update(context.Background(), "clinic-id", req)
	if err != nil {
		t.Fatalf("update clinic: %v", err)
	}

	if !repo.updated {
		t.Fatal("repository was not updated")
	}
	if res.Name != "New Name" || res.City != "New City" || res.CommunicationTone != "profesional" {
		t.Fatalf("unexpected response: %#v", res)
	}
}
