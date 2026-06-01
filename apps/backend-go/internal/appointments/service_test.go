package appointments

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRepository struct {
	appointments []Appointment
	convert      ConvertLeadResponse
	err          error
}

func (r *fakeRepository) List(ctx context.Context, clinicID string, filters ListFilters) ([]Appointment, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.appointments, nil
}

func (r *fakeRepository) FindByID(ctx context.Context, clinicID, appointmentID string) (Appointment, error) {
	if r.err != nil {
		return Appointment{}, r.err
	}
	for _, appointment := range r.appointments {
		if appointment.ID == appointmentID && appointment.ClinicID == clinicID {
			return appointment, nil
		}
	}
	return Appointment{}, ErrAppointmentNotFound
}

func (r *fakeRepository) Create(ctx context.Context, appointment Appointment) (Appointment, error) {
	if r.err != nil {
		return Appointment{}, r.err
	}
	appointment.ID = "appointment-new"
	appointment.ProfessionalName = "Dra. Ana"
	appointment.ServiceName = "Ortodoncia"
	appointment.CreatedAt = time.Now()
	appointment.UpdatedAt = appointment.CreatedAt
	r.appointments = append(r.appointments, appointment)
	return appointment, nil
}

func (r *fakeRepository) Update(ctx context.Context, appointment Appointment) (Appointment, error) {
	if r.err != nil {
		return Appointment{}, r.err
	}
	appointment.ProfessionalName = "Dra. Ana"
	appointment.ServiceName = "Ortodoncia"
	return appointment, nil
}

func (r *fakeRepository) UpdateStatus(ctx context.Context, clinicID, appointmentID, status string, adminNote *string) (Appointment, error) {
	if r.err != nil {
		return Appointment{}, r.err
	}
	return Appointment{ID: appointmentID, ClinicID: clinicID, ProfessionalID: "professional-1", ProfessionalName: "Dra. Ana", ServiceID: "service-1", ServiceName: "Ortodoncia", ContactName: "Maria", Status: status, ConfirmationStatus: "confirmed", Source: "manual"}, nil
}

func (r *fakeRepository) Reschedule(ctx context.Context, clinicID, appointmentID string, startsAt, endsAt time.Time, adminNote *string) (Appointment, error) {
	if r.err != nil {
		return Appointment{}, r.err
	}
	return Appointment{ID: appointmentID, ClinicID: clinicID, ProfessionalID: "professional-1", ProfessionalName: "Dra. Ana", ServiceID: "service-1", ServiceName: "Ortodoncia", ContactName: "Maria", StartsAt: startsAt, EndsAt: endsAt, Status: "rescheduled", ConfirmationStatus: "pending", Source: "manual"}, nil
}

func (r *fakeRepository) ConvertLead(ctx context.Context, clinicID, leadID string, appointment Appointment, updateLeadStatus bool) (ConvertLeadResponse, error) {
	if r.err != nil {
		return ConvertLeadResponse{}, r.err
	}
	if r.convert.AppointmentID != "" {
		return r.convert, nil
	}
	return ConvertLeadResponse{LeadID: leadID, LeadStatus: "Agendado", AppointmentID: "appointment-new", AppointmentStatus: appointment.Status}, nil
}

func TestServiceCreateAppointment(t *testing.T) {
	service := NewService(&fakeRepository{})
	startsAt := time.Date(2026, 6, 2, 14, 0, 0, 0, time.UTC)
	duration := 60

	res, err := service.Create(context.Background(), "clinic-1", "user-1", CreateAppointmentRequest{
		ProfessionalID: "professional-1",
		ServiceID:      "service-1",
		ContactName:    " Maria Perez ",
		StartsAt:       startsAt,
		DurationMins:   &duration,
	})
	if err != nil {
		t.Fatalf("create appointment: %v", err)
	}
	if res.ID != "appointment-new" || res.ContactName != "Maria Perez" || res.Status != "pending_confirmation" {
		t.Fatalf("unexpected response: %#v", res)
	}
	if !res.EndsAt.Equal(startsAt.Add(time.Hour)) {
		t.Fatalf("unexpected end time: %s", res.EndsAt)
	}
}

func TestServiceCreateValidatesContactName(t *testing.T) {
	service := NewService(&fakeRepository{})
	_, err := service.Create(context.Background(), "clinic-1", "user-1", CreateAppointmentRequest{
		ProfessionalID: "professional-1",
		ServiceID:      "service-1",
		StartsAt:       time.Now(),
	})
	if !errors.Is(err, ErrContactNameRequired) {
		t.Fatalf("expected ErrContactNameRequired, got %v", err)
	}
}

func TestServiceCreateValidatesTimeRange(t *testing.T) {
	service := NewService(&fakeRepository{})
	startsAt := time.Date(2026, 6, 2, 14, 0, 0, 0, time.UTC)
	endsAt := startsAt.Add(-time.Hour)
	_, err := service.Create(context.Background(), "clinic-1", "user-1", CreateAppointmentRequest{
		ProfessionalID: "professional-1",
		ServiceID:      "service-1",
		ContactName:    "Maria",
		StartsAt:       startsAt,
		EndsAt:         &endsAt,
	})
	if !errors.Is(err, ErrInvalidAppointmentTime) {
		t.Fatalf("expected ErrInvalidAppointmentTime, got %v", err)
	}
}

func TestServiceConvertLead(t *testing.T) {
	service := NewService(&fakeRepository{})
	duration := 60
	res, err := service.ConvertLead(context.Background(), "clinic-1", "user-1", "lead-1", ConvertLeadRequest{
		ProfessionalID:   "professional-1",
		ServiceID:        "service-1",
		StartsAt:         time.Date(2026, 6, 2, 14, 0, 0, 0, time.UTC),
		DurationMins:     &duration,
		UpdateLeadStatus: true,
	})
	if err != nil {
		t.Fatalf("convert lead: %v", err)
	}
	if res.AppointmentID != "appointment-new" || res.LeadStatus != "Agendado" {
		t.Fatalf("unexpected response: %#v", res)
	}
}
