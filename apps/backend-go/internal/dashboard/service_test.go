package dashboard

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	summary     Summary
	schedSum    ScheduleSummary
	activeProfs []ProfessionalHours
	appts       []AppointmentSummary
	err         error
}

func (r fakeRepository) Summary(ctx context.Context, clinicID string) (Summary, error) {
	if r.err != nil {
		return Summary{}, r.err
	}
	return r.summary, nil
}

func (r fakeRepository) ScheduleSummary(ctx context.Context, clinicID string, dateStr string) (ScheduleSummary, error) {
	if r.err != nil {
		return ScheduleSummary{}, r.err
	}
	return r.schedSum, nil
}

func (r fakeRepository) GetActiveProfessionalsWithHours(ctx context.Context, clinicID string) ([]ProfessionalHours, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.activeProfs, nil
}

func (r fakeRepository) GetAppointmentsForDate(ctx context.Context, clinicID string, dateStr string) ([]AppointmentSummary, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.appts, nil
}

func TestServiceSummary(t *testing.T) {
	repo := fakeRepository{summary: Summary{
		LeadsTotal: 7,
		LeadsByStatus: map[string]int{
			"Nuevo":      2,
			"Convertido": 1,
		},
		TopServices: []TopService{{
			ServiceID:   "service-1",
			ServiceName: "Blanqueamiento dental",
			LeadCount:   4,
		}},
		PendingFollowUpsToday: 3,
		OverdueFollowUps:      1,
		UpcomingFollowUps:     5,
		ConversionRate:        1.0 / 7.0,
	}}
	service := NewService(repo)

	res, err := service.Summary(context.Background(), "clinic-1")
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if res.LeadsTotal != 7 || res.LeadsByStatus["Convertido"] != 1 {
		t.Fatalf("unexpected lead metrics: %#v", res)
	}
	if len(res.TopServices) != 1 || res.TopServices[0].ServiceName != "Blanqueamiento dental" {
		t.Fatalf("unexpected top services: %#v", res.TopServices)
	}
	if res.PendingFollowUpsToday != 3 || res.OverdueFollowUps != 1 || res.UpcomingFollowUps != 5 {
		t.Fatalf("unexpected follow-up metrics: %#v", res)
	}
}

func TestServiceSummaryRequiresClinicID(t *testing.T) {
	service := NewService(fakeRepository{})

	_, err := service.Summary(context.Background(), " ")
	if !errors.Is(err, ErrMissingClinicID) {
		t.Fatalf("expected ErrMissingClinicID, got %v", err)
	}
}

func TestServiceScheduleSummary(t *testing.T) {
	sched := ScheduleSummary{
		Date:                            "2026-06-01",
		TodaysAppointments:              5,
		AppointmentsPendingConfirmation: 2,
		HotLeadsWithoutAppointment:      3,
		OverdueFollowUps:                1,
		LeadsConvertedToAppointments:    2,
		AppointmentsByProfessional: []ProfessionalCountModel{
			{ProfessionalID: "prof-1", ProfessionalName: "Dr. House", AppointmentCount: 3},
		},
		TopServicesByScheduleDemand: []ServiceCountModel{
			{ServiceID: "serv-1", ServiceName: "Ortodoncia", AppointmentCount: 4},
		},
	}

	repo := fakeRepository{
		schedSum:    sched,
		activeProfs: []ProfessionalHours{},
		appts:       []AppointmentSummary{},
	}
	service := NewService(repo)

	res, err := service.ScheduleSummary(context.Background(), "clinic-1", "2026-06-01")
	if err != nil {
		t.Fatalf("schedule summary failed: %v", err)
	}

	if res.Date != "2026-06-01" || res.TodaysAppointments != 5 || res.AvailableSlots != 0 {
		t.Fatalf("unexpected schedule summary: %#v", res)
	}
}

