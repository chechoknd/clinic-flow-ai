package leads

import (
	"context"
	"database/sql"
	"testing"
	"time"
)

type fakeRepository struct {
	leads []Lead
	notes []LeadNote
	err   error
}

func (r *fakeRepository) List(ctx context.Context, clinicID string, filter ListFilter) ([]Lead, int, error) {
	if r.err != nil {
		return nil, 0, r.err
	}
	return r.leads, len(r.leads), nil
}

func (r *fakeRepository) ListFollowUps(ctx context.Context, clinicID string, filter FollowUpFilter) ([]Lead, int, error) {
	if r.err != nil {
		return nil, 0, r.err
	}
	var filtered []Lead
	for _, l := range r.leads {
		if l.ClinicID == clinicID && l.NextActionAt != nil {
			filtered = append(filtered, l)
		}
	}
	return filtered, len(filtered), nil
}

func (r *fakeRepository) FindByID(ctx context.Context, clinicID, leadID string) (Lead, []LeadNote, error) {
	if r.err != nil {
		return Lead{}, nil, r.err
	}
	for _, l := range r.leads {
		if l.ID == leadID && l.ClinicID == clinicID {
			var leadNotes []LeadNote
			for _, n := range r.notes {
				if n.LeadID == leadID {
					leadNotes = append(leadNotes, n)
				}
			}
			return l, leadNotes, nil
		}
	}
	return Lead{}, nil, ErrLeadNotFound
}

func (r *fakeRepository) Create(ctx context.Context, l Lead, initialNote string) (Lead, error) {
	if r.err != nil {
		return Lead{}, r.err
	}
	l.ID = "new-id"
	l.CreatedAt = time.Date(2026, 5, 21, 10, 0, 0, 0, time.UTC)
	l.UpdatedAt = l.CreatedAt
	r.leads = append(r.leads, l)
	if initialNote != "" {
		r.notes = append(r.notes, LeadNote{ID: "note-id", LeadID: l.ID, Body: initialNote})
	}
	return l, nil
}

func (r *fakeRepository) Update(ctx context.Context, clinicID, leadID string, status string, nextActionAt *sql.NullTime, note string) error {
	if r.err != nil {
		return r.err
	}
	for i, l := range r.leads {
		if l.ID == leadID && l.ClinicID == clinicID {
			r.leads[i].Status = status
			if nextActionAt != nil && nextActionAt.Valid {
				r.leads[i].NextActionAt = &nextActionAt.Time
			}
			if note != "" {
				r.notes = append(r.notes, LeadNote{ID: "note-id", LeadID: leadID, Body: note})
			}
			return nil
		}
	}
	return ErrLeadNotFound
}

func (r *fakeRepository) CompleteFollowUp(ctx context.Context, clinicID, leadID, status, note string) error {
	if r.err != nil {
		return r.err
	}
	for i, l := range r.leads {
		if l.ID == leadID && l.ClinicID == clinicID {
			r.leads[i].Status = status
			r.leads[i].NextActionAt = nil
			if note != "" {
				r.notes = append(r.notes, LeadNote{ID: "note-id", LeadID: leadID, Body: note})
			}
			return nil
		}
	}
	return ErrLeadNotFound
}

func (r *fakeRepository) RescheduleFollowUp(ctx context.Context, clinicID, leadID string, nextActionAt sql.NullTime, note string) error {
	if r.err != nil {
		return r.err
	}
	for i, l := range r.leads {
		if l.ID == leadID && l.ClinicID == clinicID {
			if nextActionAt.Valid {
				r.leads[i].NextActionAt = &nextActionAt.Time
			}
			if note != "" {
				r.notes = append(r.notes, LeadNote{ID: "note-id", LeadID: leadID, Body: note})
			}
			return nil
		}
	}
	return ErrLeadNotFound
}

func TestLeadServiceCreate(t *testing.T) {
	repo := &fakeRepository{}
	s := NewService(repo)

	req := CreateLeadRequest{
		FullName: "Maria Perez",
		Phone:    "+573001112233",
		Notes:    "Interesada en blanqueamiento",
	}

	res, err := s.Create(context.Background(), "clinic-1", req)
	if err != nil {
		t.Fatalf("create lead: %v", err)
	}
	if res.ID != "new-id" || res.Status != "Nuevo" {
		t.Fatalf("unexpected response: %#v", res)
	}
	if res.FullName != "Maria Perez" || res.Phone != "+573001112233" || res.Source != "whatsapp" {
		t.Fatalf("create response omitted lead fields: %#v", res)
	}
	if res.CreatedAt.IsZero() || res.UpdatedAt.IsZero() {
		t.Fatalf("create response omitted timestamps: %#v", res)
	}
}

func TestLeadServiceListFollowUps(t *testing.T) {
	nextActionAt := time.Date(2026, 5, 22, 15, 0, 0, 0, time.UTC)
	repo := &fakeRepository{leads: []Lead{
		{ID: "lead-1", ClinicID: "clinic-1", FullName: "Maria Perez", Phone: "+573001112233", Status: "Interesado", Source: "whatsapp", NextActionAt: &nextActionAt},
		{ID: "lead-2", ClinicID: "clinic-1", FullName: "Carlos Perez", Phone: "+573001112244", Status: "Nuevo", Source: "whatsapp"},
	}}
	s := NewService(repo)

	res, err := s.ListFollowUps(context.Background(), "clinic-1", FollowUpFilter{})
	if err != nil {
		t.Fatalf("list follow-ups: %v", err)
	}
	if len(res.Data) != 1 || res.Data[0].ID != "lead-1" {
		t.Fatalf("unexpected follow-ups: %#v", res.Data)
	}
}

func TestLeadServiceCompleteFollowUp(t *testing.T) {
	nextActionAt := time.Date(2026, 5, 22, 15, 0, 0, 0, time.UTC)
	repo := &fakeRepository{leads: []Lead{{ID: "lead-1", ClinicID: "clinic-1", FullName: "Maria Perez", Phone: "+573001112233", Status: "Interesado", Source: "whatsapp", NextActionAt: &nextActionAt}}}
	s := NewService(repo)

	err := s.CompleteFollowUp(context.Background(), "clinic-1", "lead-1", CompleteFollowUpRequest{Note: "Se envio seguimiento."})
	if err != nil {
		t.Fatalf("complete follow-up: %v", err)
	}
	if repo.leads[0].NextActionAt != nil || repo.leads[0].Status != "Contactado" {
		t.Fatalf("follow-up not completed: %#v", repo.leads[0])
	}
	if len(repo.notes) != 1 {
		t.Fatalf("expected completion note")
	}
}

func TestLeadServiceRescheduleFollowUp(t *testing.T) {
	repo := &fakeRepository{leads: []Lead{{ID: "lead-1", ClinicID: "clinic-1", FullName: "Maria Perez", Phone: "+573001112233", Status: "Interesado", Source: "whatsapp"}}}
	s := NewService(repo)
	nextActionAt := time.Date(2026, 5, 23, 15, 0, 0, 0, time.UTC)

	err := s.RescheduleFollowUp(context.Background(), "clinic-1", "lead-1", RescheduleFollowUpRequest{NextActionAt: nextActionAt})
	if err != nil {
		t.Fatalf("reschedule follow-up: %v", err)
	}
	if repo.leads[0].NextActionAt == nil || !repo.leads[0].NextActionAt.Equal(nextActionAt) {
		t.Fatalf("follow-up not rescheduled: %#v", repo.leads[0])
	}
}
