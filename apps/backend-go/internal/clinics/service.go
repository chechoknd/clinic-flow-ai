package clinics

import (
	"context"
	"encoding/json"
	"errors"
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

func (s *Service) Current(ctx context.Context, clinicID string) (ClinicResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return ClinicResponse{}, ErrMissingClinicID
	}

	clinic, err := s.repository.FindByID(ctx, clinicID)
	if err != nil {
		return ClinicResponse{}, err
	}

	return mapClinicToResponse(clinic), nil
}

func (s *Service) Update(ctx context.Context, clinicID string, req UpdateClinicRequest) (ClinicResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return ClinicResponse{}, ErrMissingClinicID
	}

	// Basic validation
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return ClinicResponse{}, errors.New("clinic name is required")
	}
	req.City = strings.TrimSpace(req.City)
	if req.City == "" {
		return ClinicResponse{}, errors.New("city is required")
	}
	req.CountryCode = shared.NormalizeCountryCode(req.CountryCode)
	req.CurrencyCode = shared.NormalizeCurrencyCode(req.CurrencyCode)
	if req.CountryCode == "" {
		req.CountryCode = "CO"
	}
	if req.CurrencyCode == "" {
		currency, ok := shared.CurrencyForCountry(req.CountryCode)
		if !ok {
			return ClinicResponse{}, errors.New("invalid country code")
		}
		req.CurrencyCode = currency.Code
	}
	if !shared.IsValidCountryCurrency(req.CountryCode, req.CurrencyCode) {
		return ClinicResponse{}, errors.New("invalid country currency")
	}
	req.WhatsApp = shared.NormalizePhone(req.WhatsApp)
	if req.WhatsApp == "" {
		return ClinicResponse{}, errors.New("whatsapp is required")
	}
	if !shared.IsValidPhone(req.WhatsApp) {
		return ClinicResponse{}, errors.New("invalid whatsapp format (expected E.164, e.g., +573001234567)")
	}

	if req.Phone != nil {
		normalized := shared.NormalizePhone(*req.Phone)
		if normalized != "" && !shared.IsValidPhone(normalized) {
			return ClinicResponse{}, errors.New("invalid phone format (expected E.164)")
		}
		req.Phone = &normalized
	}

	allowedTones := map[string]bool{
		"amable":      true,
		"profesional": true,
		"cercano":     true,
		"juvenil":     true,
		"elegante":    true,
	}
	if !allowedTones[req.CommunicationTone] {
		return ClinicResponse{}, errors.New("invalid communication tone")
	}

	clinic := Clinic{
		ID:                clinicID,
		Name:              req.Name,
		City:              req.City,
		CountryCode:       req.CountryCode,
		CurrencyCode:      req.CurrencyCode,
		Phone:             req.Phone,
		WhatsApp:          req.WhatsApp,
		Address:           req.Address,
		OpeningHours:      []byte(req.OpeningHours),
		GeneralFAQ:        []byte(req.GeneralFAQ),
		CommunicationTone: req.CommunicationTone,
	}

	if err := s.repository.Update(ctx, clinic); err != nil {
		return ClinicResponse{}, err
	}

	// Return the updated state (we could also re-fetch if we want to be sure of DB state)
	// For simplicity, we return what we just sent, but with the full fields.
	// Note: ClinicType is not updatable via this endpoint for now.
	updated, err := s.repository.FindByID(ctx, clinicID)
	if err != nil {
		return ClinicResponse{}, err
	}

	return mapClinicToResponse(updated), nil
}

func mapClinicToResponse(clinic Clinic) ClinicResponse {
	currency, _ := shared.CurrencyByCode(clinic.CurrencyCode)

	return ClinicResponse{
		ID:                clinic.ID,
		Name:              clinic.Name,
		ClinicType:        clinic.ClinicType,
		City:              clinic.City,
		CountryCode:       clinic.CountryCode,
		CurrencyCode:      clinic.CurrencyCode,
		Currency:          currency,
		Phone:             clinic.Phone,
		WhatsApp:          clinic.WhatsApp,
		Address:           clinic.Address,
		OpeningHours:      json.RawMessage(clinic.OpeningHours),
		GeneralFAQ:        json.RawMessage(clinic.GeneralFAQ),
		CommunicationTone: clinic.CommunicationTone,
	}
}
