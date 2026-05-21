package leads

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"strings"
)

var ErrMissingClinicID = errors.New("clinic id is required")

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
		data[i] = LeadResponse{
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

	lead, notes, err := s.repository.FindByID(ctx, clinicID, leadID)
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

	return LeadDetailResponse{
		ID:           lead.ID,
		FullName:     lead.FullName,
		Phone:        lead.Phone,
		Service:      serviceDTO,
		Status:       lead.Status,
		Source:       lead.Source,
		Notes:        noteResponses,
		NextActionAt: lead.NextActionAt,
		CreatedAt:    lead.CreatedAt,
	}, nil
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
	req.Phone = strings.TrimSpace(req.Phone)
	if req.Phone == "" {
		return LeadResponse{}, errors.New("phone is required")
	}

	if req.Status == "" {
		req.Status = "Nuevo"
	}
	if req.Source == "" {
		req.Source = "whatsapp"
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

	created, err := s.repository.Create(ctx, model, req.Notes)
	if err != nil {
		return LeadResponse{}, err
	}

	return LeadResponse{
		ID:        created.ID,
		Status:    created.Status,
		CreatedAt: created.CreatedAt,
	}, nil
}

func (s *Service) Update(ctx context.Context, clinicID, leadID string, req UpdateLeadRequest) error {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return ErrMissingClinicID
	}

	if req.Status == "" {
		return errors.New("status is required")
	}

	var nextActionAt *sql.NullTime
	if req.NextActionAt != nil {
		nextActionAt = &sql.NullTime{Time: *req.NextActionAt, Valid: true}
	}

	return s.repository.Update(ctx, clinicID, leadID, req.Status, nextActionAt, req.Note)
}
