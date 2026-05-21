package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
)

var ErrMissingClinicID = errors.New("clinic id is required")

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context, clinicID string) ([]ServiceResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return nil, ErrMissingClinicID
	}

	entities, err := s.repository.ListByClinicID(ctx, clinicID)
	if err != nil {
		return nil, err
	}

	res := make([]ServiceResponse, len(entities))
	for i, entity := range entities {
		r, _ := mapModelToResponse(entity)
		res[i] = r
	}

	return res, nil
}

func (s *Service) Get(ctx context.Context, clinicID, serviceID string) (ServiceResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return ServiceResponse{}, ErrMissingClinicID
	}

	entity, err := s.repository.FindByID(ctx, clinicID, serviceID)
	if err != nil {
		return ServiceResponse{}, err
	}

	return mapModelToResponse(entity)
}

func (s *Service) Create(ctx context.Context, clinicID string, req CreateServiceRequest) (ServiceResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return ServiceResponse{}, ErrMissingClinicID
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return ServiceResponse{}, errors.New("service name is required")
	}

	// Default empty JSON if nil
	if req.Benefits == nil {
		req.Benefits = json.RawMessage(`[]`)
	}
	if req.FAQ == nil {
		req.FAQ = json.RawMessage(`[]`)
	}
	if req.CommonObjections == nil {
		req.CommonObjections = json.RawMessage(`[]`)
	}

	entity := ServiceEntity{
		ClinicID:         clinicID,
		Name:             req.Name,
		Description:      req.Description,
		DurationMinutes:  req.DurationMinutes,
		PriceFrom:        req.PriceFrom,
		Benefits:         []byte(req.Benefits),
		FAQ:              []byte(req.FAQ),
		CommonObjections: []byte(req.CommonObjections),
		IsActive:         true,
	}

	created, err := s.repository.Create(ctx, entity)
	if err != nil {
		return ServiceResponse{}, err
	}

	return mapModelToResponse(created)
}

func (s *Service) Update(ctx context.Context, clinicID, serviceID string, req UpdateServiceRequest) (ServiceResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return ServiceResponse{}, ErrMissingClinicID
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return ServiceResponse{}, errors.New("service name is required")
	}

	entity := ServiceEntity{
		ID:               serviceID,
		ClinicID:         clinicID,
		Name:             req.Name,
		Description:      req.Description,
		DurationMinutes:  req.DurationMinutes,
		PriceFrom:        req.PriceFrom,
		Benefits:         []byte(req.Benefits),
		FAQ:              []byte(req.FAQ),
		CommonObjections: []byte(req.CommonObjections),
		IsActive:         req.IsActive,
	}

	if err := s.repository.Update(ctx, entity); err != nil {
		return ServiceResponse{}, err
	}

	// We could return mapModelToResponse(entity) but some fields might be updated by DB (like updated_at, though not in model yet)
	// For consistency with create, let's just map it.
	return mapModelToResponse(entity)
}

func (s *Service) Delete(ctx context.Context, clinicID, serviceID string) error {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return ErrMissingClinicID
	}

	return s.repository.Delete(ctx, clinicID, serviceID)
}

func mapModelToResponse(m ServiceEntity) (ServiceResponse, error) {
	return ServiceResponse{
		ID:               m.ID,
		ClinicID:         m.ClinicID,
		Name:             m.Name,
		Description:      m.Description,
		DurationMinutes:  m.DurationMinutes,
		PriceFrom:        m.PriceFrom,
		Benefits:         json.RawMessage(m.Benefits),
		FAQ:              json.RawMessage(m.FAQ),
		CommonObjections: json.RawMessage(m.CommonObjections),
		IsActive:         m.IsActive,
	}, nil
}
