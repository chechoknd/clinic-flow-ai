package leads

import (
	"context"
	"database/sql"
	"testing"
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
}
