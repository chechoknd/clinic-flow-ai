package clinics

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

func (s *Service) Current(ctx context.Context, clinicID string) (ClinicResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return ClinicResponse{}, ErrMissingClinicID
	}

	clinic, err := s.repository.FindByID(ctx, clinicID)
	if err != nil {
		return ClinicResponse{}, err
	}

	return ClinicResponse{
		ID:                clinic.ID,
		Name:              clinic.Name,
		ClinicType:        clinic.ClinicType,
		City:              clinic.City,
		Phone:             clinic.Phone,
		WhatsApp:          clinic.WhatsApp,
		Address:           clinic.Address,
		OpeningHours:      json.RawMessage(clinic.OpeningHours),
		GeneralFAQ:        json.RawMessage(clinic.GeneralFAQ),
		CommunicationTone: clinic.CommunicationTone,
	}, nil
}
