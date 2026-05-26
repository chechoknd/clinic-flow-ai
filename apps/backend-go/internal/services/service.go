package services

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strings"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/shared"
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

	if err := s.validatePriceForClinic(ctx, clinicID, req.PriceFrom); err != nil {
		return ServiceResponse{}, err
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

	created, err = s.repository.FindByID(ctx, clinicID, created.ID)
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

	if err := s.validatePriceForClinic(ctx, clinicID, req.PriceFrom); err != nil {
		return ServiceResponse{}, err
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

	updated, err := s.repository.FindByID(ctx, clinicID, serviceID)
	if err != nil {
		return ServiceResponse{}, err
	}

	return mapModelToResponse(updated)
}

func (s *Service) Delete(ctx context.Context, clinicID, serviceID string) error {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return ErrMissingClinicID
	}

	return s.repository.Delete(ctx, clinicID, serviceID)
}

func (s *Service) validatePriceForClinic(ctx context.Context, clinicID string, price *float64) error {
	if price == nil {
		return nil
	}
	if *price < 0 {
		return errors.New("price_from must be greater than or equal to zero")
	}

	currencyCode, err := s.repository.CurrencyCodeByClinicID(ctx, clinicID)
	if err != nil {
		return err
	}
	currency, ok := shared.CurrencyByCode(currencyCode)
	if !ok {
		return errors.New("unsupported clinic currency")
	}

	factor := math.Pow10(currency.DecimalDigits)
	if math.Round(*price*factor) != *price*factor {
		return errors.New("price_from has too many decimal places for clinic currency")
	}
	return nil
}

func mapModelToResponse(m ServiceEntity) (ServiceResponse, error) {
	return ServiceResponse{
		ID:               m.ID,
		ClinicID:         m.ClinicID,
		Name:             m.Name,
		Description:      m.Description,
		DurationMinutes:  m.DurationMinutes,
		PriceFrom:        m.PriceFrom,
		CurrencyCode:     m.CurrencyCode,
		Benefits:         json.RawMessage(m.Benefits),
		FAQ:              json.RawMessage(m.FAQ),
		CommonObjections: json.RawMessage(m.CommonObjections),
		IsActive:         m.IsActive,
	}, nil
}
