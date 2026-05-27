package leads

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"regexp"
	"strings"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/shared"
)

var ErrMissingClinicID = errors.New("clinic id is required")

var uuidPattern = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context, clinicID string, filter ListFilter) (PaginatedLeadsResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return PaginatedLeadsResponse{}, ErrMissingClinicID
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	leads, total, err := s.repository.List(ctx, clinicID, filter)
	if err != nil {
		return PaginatedLeadsResponse{}, err
	}

	data := make([]LeadResponse, len(leads))
	for i, l := range leads {
		data[i] = leadToResponse(l)
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.PageSize)))

	return PaginatedLeadsResponse{
		Data: data,
		Pagination: PaginationResponse{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *Service) Get(ctx context.Context, clinicID, leadID string) (LeadDetailResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return LeadDetailResponse{}, ErrMissingClinicID
	}

	lead, notes, insights, err := s.repository.FindByID(ctx, clinicID, leadID)
	if err != nil {
		return LeadDetailResponse{}, err
	}

	var serviceDTO *LeadServiceDTO
	if lead.ServiceID != nil && lead.ServiceName != nil {
		serviceDTO = &LeadServiceDTO{
			ID:   *lead.ServiceID,
			Name: *lead.ServiceName,
		}
	}

	noteResponses := make([]LeadNoteResponse, len(notes))
	for i, n := range notes {
		noteResponses[i] = LeadNoteResponse{
			ID:        n.ID,
			Body:      n.Body,
			CreatedAt: n.CreatedAt,
		}
	}

	insightResponses := make([]AIInsightResponse, len(insights))
	for i, insight := range insights {
		insightResponses[i] = aiInsightToResponse(insight)
	}

	return LeadDetailResponse{
		ID:           lead.ID,
		FullName:     lead.FullName,
		Phone:        lead.Phone,
		Service:      serviceDTO,
		Status:       lead.Status,
		Source:       lead.Source,
		Notes:        noteResponses,
		AIInsights:   insightResponses,
		NextActionAt: lead.NextActionAt,
		CreatedAt:    lead.CreatedAt,
	}, nil
}

func (s *Service) ListFollowUps(ctx context.Context, clinicID string, filter FollowUpFilter) (FollowUpFilterResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return FollowUpFilterResponse{}, ErrMissingClinicID
	}

	filter.Due = strings.TrimSpace(filter.Due)
	if filter.Due == "" {
		filter.Due = "all"
	}
	if filter.Due != "all" && filter.Due != "today" && filter.Due != "overdue" && filter.Due != "upcoming" {
		return FollowUpFilterResponse{}, errors.New("invalid follow-up due filter")
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	leads, total, err := s.repository.ListFollowUps(ctx, clinicID, filter)
	if err != nil {
		return FollowUpFilterResponse{}, err
	}

	data := make([]LeadResponse, len(leads))
	for i, l := range leads {
		data[i] = leadToResponse(l)
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.PageSize)))
	return FollowUpFilterResponse{
		Data: data,
		Pagination: PaginationResponse{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *Service) CompleteFollowUp(ctx context.Context, clinicID, leadID string, req CompleteFollowUpRequest) error {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return ErrMissingClinicID
	}
	leadID = strings.TrimSpace(leadID)
	if leadID == "" {
		return errors.New("lead id is required")
	}

	req.Status = strings.TrimSpace(req.Status)
	if req.Status == "" {
		req.Status = "Contactado"
	}
	if !isAllowedStatus(req.Status) {
		return errors.New("invalid status")
	}

	return s.repository.CompleteFollowUp(ctx, clinicID, leadID, req.Status, strings.TrimSpace(req.Note))
}

func (s *Service) RescheduleFollowUp(ctx context.Context, clinicID, leadID string, req RescheduleFollowUpRequest) error {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return ErrMissingClinicID
	}
	leadID = strings.TrimSpace(leadID)
	if leadID == "" {
		return errors.New("lead id is required")
	}
	if req.NextActionAt.IsZero() {
		return errors.New("next_action_at is required")
	}

	return s.repository.RescheduleFollowUp(ctx, clinicID, leadID, sql.NullTime{Time: req.NextActionAt, Valid: true}, strings.TrimSpace(req.Note))
}

func (s *Service) Create(ctx context.Context, clinicID string, req CreateLeadRequest) (LeadResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return LeadResponse{}, ErrMissingClinicID
	}

	req.FullName = strings.TrimSpace(req.FullName)
	if req.FullName == "" {
		return LeadResponse{}, errors.New("full name is required")
	}

	req.Phone = shared.NormalizePhone(req.Phone)
	if req.Phone == "" {
		return LeadResponse{}, errors.New("phone is required")
	}
	if !shared.IsValidPhone(req.Phone) {
		return LeadResponse{}, errors.New("invalid phone format (expected E.164, e.g., +573001234567)")
	}

	if req.Status == "" {
		req.Status = "Nuevo"
	}
	if !isAllowedStatus(req.Status) {
		return LeadResponse{}, errors.New("invalid status")
	}

	req.Source = strings.TrimSpace(req.Source)
	if req.Source == "" {
		req.Source = "whatsapp"
	}

	allowedSources := map[string]bool{
		"whatsapp":  true,
		"instagram": true,
		"facebook":  true,
		"web":       true,
		"llamada":   true,
		"otro":      true,
	}
	if !allowedSources[req.Source] {
		return LeadResponse{}, errors.New("invalid lead source")
	}

	model := Lead{
		ClinicID:     clinicID,
		FullName:     req.FullName,
		Phone:        req.Phone,
		ServiceID:    req.ServiceID,
		Status:       req.Status,
		Source:       req.Source,
		NextActionAt: req.NextActionAt,
	}

	insight, err := aiInsightFromRequest(req.ReviewedAIAnalysis, req.Source)
	if err != nil {
		return LeadResponse{}, err
	}

	created, err := s.repository.Create(ctx, model, strings.TrimSpace(req.Notes), insight)
	if err != nil {
		return LeadResponse{}, err
	}

	return LeadResponse{
		ID:           created.ID,
		FullName:     created.FullName,
		Phone:        created.Phone,
		ServiceID:    created.ServiceID,
		Status:       created.Status,
		Source:       created.Source,
		NextActionAt: created.NextActionAt,
		CreatedAt:    created.CreatedAt,
		UpdatedAt:    created.UpdatedAt,
	}, nil
}

func (s *Service) Update(ctx context.Context, clinicID, leadID string, req UpdateLeadRequest) error {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return ErrMissingClinicID
	}

	req.Status = strings.TrimSpace(req.Status)
	if req.Status == "" {
		return errors.New("status is required")
	}

	if !isAllowedStatus(req.Status) {
		return errors.New("invalid status")
	}

	var nextActionAt *sql.NullTime
	if req.NextActionAt != nil {
		nextActionAt = &sql.NullTime{Time: *req.NextActionAt, Valid: true}
	}

	insight, err := aiInsightFromRequest(req.ReviewedAIAnalysis, "")
	if err != nil {
		return err
	}

	return s.repository.Update(ctx, clinicID, leadID, req.Status, nextActionAt, strings.TrimSpace(req.Note), insight)
}

func leadToResponse(l Lead) LeadResponse {
	return LeadResponse{
		ID:           l.ID,
		FullName:     l.FullName,
		Phone:        l.Phone,
		ServiceID:    l.ServiceID,
		ServiceName:  l.ServiceName,
		Status:       l.Status,
		Source:       l.Source,
		NextActionAt: l.NextActionAt,
		CreatedAt:    l.CreatedAt,
		UpdatedAt:    l.UpdatedAt,
	}
}

func aiInsightFromRequest(req *ReviewedAIAnalysisRequest, defaultSource string) (*AIInsight, error) {
	if req == nil {
		return nil, nil
	}

	intent := strings.ToLower(strings.TrimSpace(req.Intent))
	if intent == "" {
		intent = "medium"
	}
	if intent != "low" && intent != "medium" && intent != "high" {
		return nil, errors.New("invalid reviewed ai analysis intent")
	}

	var generationID *string
	analysisID := strings.TrimSpace(req.AnalysisID)
	if analysisID != "" {
		if !uuidPattern.MatchString(analysisID) {
			return nil, errors.New("invalid reviewed ai analysis id")
		}
		generationID = &analysisID
	}

	source := truncate(strings.TrimSpace(req.Source), 50)
	if source == "" {
		source = truncate(strings.TrimSpace(defaultSource), 50)
	}

	objections := make([]string, 0, len(req.DetectedObjections))
	seen := map[string]bool{}
	for _, objection := range req.DetectedObjections {
		value := truncate(strings.TrimSpace(objection), 80)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		objections = append(objections, value)
		if len(objections) == 10 {
			break
		}
	}

	return &AIInsight{
		AIGenerationID:      generationID,
		Intent:              intent,
		DetectedObjections:  objections,
		CommercialSummary:   truncate(strings.TrimSpace(req.CommercialSummary), 1000),
		SuggestedNextAction: truncate(strings.TrimSpace(req.SuggestedNextAction), 500),
		Source:              source,
	}, nil
}

func aiInsightToResponse(insight AIInsight) AIInsightResponse {
	return AIInsightResponse{
		ID:                  insight.ID,
		AnalysisID:          insight.AIGenerationID,
		Intent:              insight.Intent,
		DetectedObjections:  insight.DetectedObjections,
		CommercialSummary:   insight.CommercialSummary,
		SuggestedNextAction: insight.SuggestedNextAction,
		Source:              insight.Source,
		CreatedAt:           insight.CreatedAt,
	}
}

func truncate(value string, maxLen int) string {
	runes := []rune(value)
	if len(runes) <= maxLen {
		return value
	}
	return string(runes[:maxLen])
}

func isAllowedStatus(status string) bool {
	allowedStatuses := map[string]bool{
		"Nuevo":        true,
		"Contactado":   true,
		"Interesado":   true,
		"Agendado":     true,
		"No Respondio": true,
		"Perdido":      true,
		"Convertido":   true,
	}
	return allowedStatuses[status]
}
