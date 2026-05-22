package dashboard

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	summary Summary
	err     error
}

func (r fakeRepository) Summary(ctx context.Context, clinicID string) (Summary, error) {
	if r.err != nil {
		return Summary{}, r.err
	}
	return r.summary, nil
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
