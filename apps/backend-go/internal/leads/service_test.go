package leads

import (
	"context"
	"database/sql"
	"testing"
	"time"
)

type fakeRepository struct {
	leads    []Lead
	notes    []LeadNote
	insights []AIInsight
	err      error
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

func (r *fakeRepository) FindByID(ctx context.Context, clinicID, leadID string) (Lead, []LeadNote, []AIInsight, error) {
	if r.err != nil {
		return Lead{}, nil, nil, r.err
	}
	for _, l := range r.leads {
		if l.ID == leadID && l.ClinicID == clinicID {
			var leadNotes []LeadNote
			for _, n := range r.notes {
				if n.LeadID == leadID {
					leadNotes = append(leadNotes, n)
				}
			}
			var leadInsights []AIInsight
			for _, insight := range r.insights {
				if insight.LeadID == leadID && insight.ClinicID == clinicID {
					leadInsights = append(leadInsights, insight)
				}
			}
			return l, leadNotes, leadInsights, nil
		}
	}
	return Lead{}, nil, nil, ErrLeadNotFound
}

func (r *fakeRepository) Create(ctx context.Context, l Lead, initialNote string, insight *AIInsight) (Lead, error) {
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
	if insight != nil {
		insight.ID = "insight-id"
		insight.ClinicID = l.ClinicID
		insight.LeadID = l.ID
		insight.CreatedAt = l.CreatedAt
		r.insights = append(r.insights, *insight)
	}
	return l, nil
}

func (r *fakeRepository) Update(ctx context.Context, clinicID, leadID string, status string, nextActionAt *sql.NullTime, note string, contactOutcome string, insight *AIInsight) error {
	if r.err != nil {
		return r.err
	}
	for i, l := range r.leads {
		if l.ID == leadID && l.ClinicID == clinicID {
			r.leads[i].Status = status
			if nextActionAt != nil {
				if nextActionAt.Valid {
					r.leads[i].NextActionAt = &nextActionAt.Time
				} else {
					r.leads[i].NextActionAt = nil
				}
			}
			if note != "" {
				newNote := LeadNote{ID: "note-id", LeadID: leadID, Body: note}
				if contactOutcome != "" {
					newNote.ContactOutcome = &contactOutcome
				}
				r.notes = append(r.notes, newNote)
			}
			if insight != nil {
				insight.ID = "insight-id"
				insight.ClinicID = clinicID
				insight.LeadID = leadID
				insight.CreatedAt = time.Date(2026, 5, 21, 10, 0, 0, 0, time.UTC)
				r.insights = append(r.insights, *insight)
			}
			return nil
		}
	}
	return ErrLeadNotFound
}

func (r *fakeRepository) CompleteFollowUp(ctx context.Context, clinicID, leadID, status, note, contactOutcome string) error {
	if r.err != nil {
		return r.err
	}
	for i, l := range r.leads {
		if l.ID == leadID && l.ClinicID == clinicID {
			r.leads[i].Status = status
			r.leads[i].NextActionAt = nil
			if note != "" {
				newNote := LeadNote{ID: "note-id", LeadID: leadID, Body: note}
				if contactOutcome != "" {
					newNote.ContactOutcome = &contactOutcome
				}
				r.notes = append(r.notes, newNote)
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

func TestLeadServiceCreatePersistsReviewedAIAnalysis(t *testing.T) {
	repo := &fakeRepository{}
	s := NewService(repo)
	analysisID := "11111111-1111-4111-8111-111111111111"

	_, err := s.Create(context.Background(), "clinic-1", CreateLeadRequest{
		FullName: "Maria Perez",
		Phone:    "+573001112233",
		Source:   "whatsapp",
		ReviewedAIAnalysis: &ReviewedAIAnalysisRequest{
			AnalysisID:          analysisID,
			Intent:              "high",
			DetectedObjections:  []string{"precio", "precio", "agenda"},
			CommercialSummary:   "Quiere aclarar precio.",
			SuggestedNextAction: "Enviar opciones y agendar seguimiento.",
		},
	})
	if err != nil {
		t.Fatalf("create lead with reviewed analysis: %v", err)
	}
	if len(repo.insights) != 1 {
		t.Fatalf("expected one insight, got %d", len(repo.insights))
	}
	insight := repo.insights[0]
	if insight.Intent != "high" || insight.AIGenerationID == nil || *insight.AIGenerationID != analysisID {
		t.Fatalf("unexpected insight metadata: %#v", insight)
	}
	if len(insight.DetectedObjections) != 2 {
		t.Fatalf("expected deduplicated objections: %#v", insight.DetectedObjections)
	}
}

func TestLeadServiceGetReturnsAIInsights(t *testing.T) {
	createdAt := time.Date(2026, 5, 21, 10, 0, 0, 0, time.UTC)
	repo := &fakeRepository{
		leads: []Lead{{ID: "lead-1", ClinicID: "clinic-1", FullName: "Maria Perez", Phone: "+573001112233", Status: "Interesado", Source: "whatsapp"}},
		insights: []AIInsight{{
			ID:                  "insight-1",
			ClinicID:            "clinic-1",
			LeadID:              "lead-1",
			Intent:              "high",
			DetectedObjections:  []string{"precio"},
			CommercialSummary:   "Interes comercial alto.",
			SuggestedNextAction: "Enviar propuesta.",
			Source:              "whatsapp",
			CreatedAt:           createdAt,
		}},
	}
	s := NewService(repo)

	res, err := s.Get(context.Background(), "clinic-1", "lead-1")
	if err != nil {
		t.Fatalf("get lead: %v", err)
	}
	if len(res.AIInsights) != 1 || res.AIInsights[0].Intent != "high" {
		t.Fatalf("unexpected insight response: %#v", res.AIInsights)
	}
}

func TestLeadServiceGetReturnsContactOutcomes(t *testing.T) {
	createdAt := time.Date(2026, 5, 21, 10, 0, 0, 0, time.UTC)
	outcome := "asked_price"
	repo := &fakeRepository{
		leads: []Lead{{ID: "lead-1", ClinicID: "clinic-1", FullName: "Maria Perez", Phone: "+573001112233", Status: "Interesado", Source: "whatsapp"}},
		notes: []LeadNote{{
			ID:             "note-1",
			LeadID:         "lead-1",
			Body:           "Pidio precio antes de decidir.",
			ContactOutcome: &outcome,
			CreatedAt:      createdAt,
		}},
	}
	s := NewService(repo)

	res, err := s.Get(context.Background(), "clinic-1", "lead-1")
	if err != nil {
		t.Fatalf("get lead: %v", err)
	}
	if len(res.Notes) != 1 || res.Notes[0].ContactOutcome == nil || *res.Notes[0].ContactOutcome != "asked_price" {
		t.Fatalf("unexpected contact outcome response: %#v", res.Notes)
	}
}

func TestLeadServiceUpdateCanClearNextAction(t *testing.T) {
	nextActionAt := time.Date(2026, 5, 24, 15, 0, 0, 0, time.UTC)
	repo := &fakeRepository{leads: []Lead{{
		ID:           "lead-1",
		ClinicID:     "clinic-1",
		FullName:     "Maria Perez",
		Phone:        "+573001112233",
		Status:       "Interesado",
		Source:       "whatsapp",
		NextActionAt: &nextActionAt,
	}}}
	s := NewService(repo)

	err := s.Update(context.Background(), "clinic-1", "lead-1", UpdateLeadRequest{
		Status:            "Perdido",
		Note:              "No continua por ahora.",
		ContactOutcome:    "lost_price",
		ClearNextActionAt: true,
	})
	if err != nil {
		t.Fatalf("update lead: %v", err)
	}
	if repo.leads[0].NextActionAt != nil {
		t.Fatalf("expected next action to be cleared: %#v", repo.leads[0])
	}
	if len(repo.notes) != 1 || repo.notes[0].ContactOutcome == nil || *repo.notes[0].ContactOutcome != "lost_price" {
		t.Fatalf("expected contact outcome note: %#v", repo.notes)
	}
}

func TestLeadServiceUpdateRejectsInvalidContactOutcome(t *testing.T) {
	repo := &fakeRepository{leads: []Lead{{ID: "lead-1", ClinicID: "clinic-1", FullName: "Maria Perez", Phone: "+573001112233", Status: "Interesado", Source: "whatsapp"}}}
	s := NewService(repo)

	err := s.Update(context.Background(), "clinic-1", "lead-1", UpdateLeadRequest{
		Status:         "Contactado",
		Note:           "Seguimiento comercial.",
		ContactOutcome: "diagnosis_requested",
	})
	if err == nil || err.Error() != "invalid contact outcome" {
		t.Fatalf("expected invalid contact outcome error, got %v", err)
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

	err := s.CompleteFollowUp(context.Background(), "clinic-1", "lead-1", CompleteFollowUpRequest{Note: "Se envio seguimiento.", ContactOutcome: "follow_up_requested"})
	if err != nil {
		t.Fatalf("complete follow-up: %v", err)
	}
	if repo.leads[0].NextActionAt != nil || repo.leads[0].Status != "Contactado" {
		t.Fatalf("follow-up not completed: %#v", repo.leads[0])
	}
	if len(repo.notes) != 1 {
		t.Fatalf("expected completion note")
	}
	if repo.notes[0].ContactOutcome == nil || *repo.notes[0].ContactOutcome != "follow_up_requested" {
		t.Fatalf("expected completion contact outcome: %#v", repo.notes)
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
