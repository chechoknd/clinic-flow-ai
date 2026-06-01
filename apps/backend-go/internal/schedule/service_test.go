package schedule

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRepository struct {
	prof         Professional
	profErr      error
	serviceDur   int
	activeAppts  []AppointmentRange
	apptsErr     error
}

func (r *fakeRepository) FindProfessional(ctx context.Context, clinicID, professionalID string) (Professional, error) {
	if r.profErr != nil {
		return Professional{}, r.profErr
	}
	if r.prof.ID != professionalID || r.prof.ClinicID != clinicID {
		return Professional{}, ErrProfessionalNotFound
	}
	return r.prof, nil
}

func (r *fakeRepository) GetServiceDuration(ctx context.Context, clinicID, serviceID string) (int, error) {
	return r.serviceDur, nil
}

func (r *fakeRepository) FindActiveAppointments(ctx context.Context, clinicID, professionalID string, fromTime, toTime time.Time) ([]AppointmentRange, error) {
	if r.apptsErr != nil {
		return nil, r.apptsErr
	}
	return r.activeAppts, nil
}

func TestGetAvailability(t *testing.T) {
	// Setup timezone location
	loc := time.UTC

	// Setup working hours: Monday 09:00 - 11:00
	whJSON := []byte(`{
		"monday": [{"start": "09:00", "end": "11:00"}],
		"tuesday": [],
		"wednesday": [],
		"thursday": [],
		"friday": [],
		"saturday": [],
		"sunday": []
	}`)

	prof := Professional{
		ID:           "prof-1",
		ClinicID:     "clinic-1",
		FullName:     "Dr. House",
		WorkingHours: whJSON,
		IsActive:     true,
	}

	// Active appointment on Monday from 09:30 to 10:00
	mondayDate := time.Date(2026, 6, 1, 0, 0, 0, 0, loc) // 2026-06-01 is a Monday
	appt := AppointmentRange{
		StartsAt: time.Date(2026, 6, 1, 9, 30, 0, 0, loc),
		EndsAt:   time.Date(2026, 6, 1, 10, 0, 0, 0, loc),
	}

	repo := &fakeRepository{
		prof:        prof,
		serviceDur:  30,
		activeAppts: []AppointmentRange{appt},
	}

	service := NewService(repo)

	// Check availability on Monday 2026-06-01 (should have: 09:00-09:30 and 10:00-10:30 and 10:30-11:00)
	// Since appointment is 09:30-10:00, slot candidate starts at 09:00 (ends 09:30) is free.
	// Slot candidate starting at 09:30 (ends 10:00) overlaps.
	// Slot candidate starting at 10:00 (ends 10:30) is free.
	// Slot candidate starting at 10:30 (ends 11:00) is free.
	req := AvailabilityRequest{
		ProfessionalID:  "prof-1",
		ServiceID:       "service-1",
		DateFrom:        mondayDate,
		DateTo:          mondayDate,
		DurationMinutes: 30,
	}

	res, err := service.GetAvailability(context.Background(), "clinic-1", req)
	if err != nil {
		t.Fatalf("failed to get availability: %v", err)
	}

	if len(res.Slots) != 3 {
		t.Fatalf("expected 3 available slots, got %d", len(res.Slots))
	}

	// Verify first slot: 09:00 - 09:30
	if !res.Slots[0].StartsAt.Equal(time.Date(2026, 6, 1, 9, 0, 0, 0, loc)) ||
		!res.Slots[0].EndsAt.Equal(time.Date(2026, 6, 1, 9, 30, 0, 0, loc)) {
		t.Errorf("unexpected first slot: %s to %s", res.Slots[0].StartsAt, res.Slots[0].EndsAt)
	}

	// Verify second slot: 10:00 - 10:30
	if !res.Slots[1].StartsAt.Equal(time.Date(2026, 6, 1, 10, 0, 0, 0, loc)) ||
		!res.Slots[1].EndsAt.Equal(time.Date(2026, 6, 1, 10, 30, 0, 0, loc)) {
		t.Errorf("unexpected second slot: %s to %s", res.Slots[1].StartsAt, res.Slots[1].EndsAt)
	}

	// Verify third slot: 10:30 - 11:00
	if !res.Slots[2].StartsAt.Equal(time.Date(2026, 6, 1, 10, 30, 0, 0, loc)) ||
		!res.Slots[2].EndsAt.Equal(time.Date(2026, 6, 1, 11, 0, 0, 0, loc)) {
		t.Errorf("unexpected third slot: %s to %s", res.Slots[2].StartsAt, res.Slots[2].EndsAt)
	}
}

func TestGetAvailabilityValidatesClinicID(t *testing.T) {
	service := NewService(&fakeRepository{})
	_, err := service.GetAvailability(context.Background(), "", AvailabilityRequest{})
	if !errors.Is(err, ErrMissingClinicID) {
		t.Errorf("expected ErrMissingClinicID, got %v", err)
	}
}

func TestGetAvailabilityValidatesProfessionalID(t *testing.T) {
	service := NewService(&fakeRepository{})
	_, err := service.GetAvailability(context.Background(), "clinic-1", AvailabilityRequest{})
	if !errors.Is(err, ErrProfessionalIDRequired) {
		t.Errorf("expected ErrProfessionalIDRequired, got %v", err)
	}
}
