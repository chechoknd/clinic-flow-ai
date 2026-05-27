package dashboard

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRepository struct {
	summary Summary
	actions []ActionItem
	err     error
}

func (r fakeRepository) Summary(ctx context.Context, clinicID string) (Summary, error) {
	if r.err != nil {
		return Summary{}, r.err
	}
	return r.summary, nil
}

func (r fakeRepository) PriorityActions(ctx context.Context, clinicID string, limit int) ([]ActionItem, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.actions, nil
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

func TestServicePriorityActions(t *testing.T) {
	nextActionAt := time.Date(2026, 5, 27, 14, 30, 0, 0, time.UTC)
	createdAt := time.Date(2026, 5, 26, 10, 0, 0, 0, time.UTC)
	serviceID := "service-1"
	serviceName := "Blanqueamiento dental"
	repo := fakeRepository{actions: []ActionItem{{
		Type:         "overdue_followup",
		Tone:         "urgent",
		Priority:     100,
		LeadID:       "lead-1",
		FullName:     "Lead Demo",
		Phone:        "+573001234567",
		ServiceID:    &serviceID,
		ServiceName:  &serviceName,
		Status:       "Interesado",
		Source:       "whatsapp",
		Reason:       "Seguimiento vencido",
		NextActionAt: &nextActionAt,
		CreatedAt:    createdAt,
	}}}
	service := NewService(repo)

	res, err := service.PriorityActions(context.Background(), "clinic-1", 4)
	if err != nil {
		t.Fatalf("priority actions: %v", err)
	}
	if len(res.Data) != 1 {
		t.Fatalf("expected one action, got %#v", res.Data)
	}
	action := res.Data[0]
	if action.Type != "overdue_followup" || action.Tone != "urgent" || action.Priority != 100 {
		t.Fatalf("unexpected action classification: %#v", action)
	}
	if action.LeadID != "lead-1" || action.ServiceName == nil || *action.ServiceName != "Blanqueamiento dental" {
		t.Fatalf("unexpected action lead data: %#v", action)
	}
	if action.NextActionAt == nil || *action.NextActionAt != "2026-05-27T14:30:00Z" {
		t.Fatalf("unexpected next action timestamp: %#v", action.NextActionAt)
	}
}

func TestServicePriorityActionsRequiresClinicID(t *testing.T) {
	service := NewService(fakeRepository{})

	_, err := service.PriorityActions(context.Background(), " ", 4)
	if !errors.Is(err, ErrMissingClinicID) {
		t.Fatalf("expected ErrMissingClinicID, got %v", err)
	}
}
