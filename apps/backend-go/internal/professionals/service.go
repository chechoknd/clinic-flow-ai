package professionals

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

var (
	ErrMissingClinicID           = errors.New("clinic id is required")
	ErrProfessionalNameRequired  = errors.New("professional name is required")
	ErrInvalidCalendarColor      = errors.New("calendar color must be a hex color like #2563EB")
	ErrInvalidWorkingHours       = errors.New("working hours must be a valid JSON object")
	ErrDuplicateProfessionalName = errors.New("professional name already exists in this clinic")
)

var colorPattern = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context, clinicID string, filters ListFilters) ([]ProfessionalResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return nil, ErrMissingClinicID
	}
	filters.ServiceID = strings.TrimSpace(filters.ServiceID)

	entities, err := s.repository.List(ctx, clinicID, filters)
	if err != nil {
		return nil, err
	}

	res := make([]ProfessionalResponse, len(entities))
	for i, entity := range entities {
		res[i] = mapModelToResponse(entity)
	}

	return res, nil
}

func (s *Service) Get(ctx context.Context, clinicID, professionalID string) (ProfessionalResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return ProfessionalResponse{}, ErrMissingClinicID
	}

	entity, err := s.repository.FindByID(ctx, clinicID, strings.TrimSpace(professionalID))
	if err != nil {
		return ProfessionalResponse{}, err
	}

	return mapModelToResponse(entity), nil
}

func (s *Service) Create(ctx context.Context, clinicID string, req CreateProfessionalRequest) (ProfessionalResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return ProfessionalResponse{}, ErrMissingClinicID
	}

	entity, serviceIDs, err := buildCreateEntity(clinicID, req)
	if err != nil {
		return ProfessionalResponse{}, err
	}

	created, err := s.repository.Create(ctx, entity, serviceIDs)
	if err != nil {
		return ProfessionalResponse{}, mapRepositoryError(err)
	}

	return mapModelToResponse(created), nil
}

func (s *Service) Update(ctx context.Context, clinicID, professionalID string, req UpdateProfessionalRequest) (ProfessionalResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return ProfessionalResponse{}, ErrMissingClinicID
	}

	entity, serviceIDs, err := buildUpdateEntity(clinicID, strings.TrimSpace(professionalID), req)
	if err != nil {
		return ProfessionalResponse{}, err
	}

	updated, err := s.repository.Update(ctx, entity, serviceIDs)
	if err != nil {
		return ProfessionalResponse{}, mapRepositoryError(err)
	}

	return mapModelToResponse(updated), nil
}

func buildCreateEntity(clinicID string, req CreateProfessionalRequest) (Professional, []string, error) {
	fullName := strings.TrimSpace(req.FullName)
	if fullName == "" {
		return Professional{}, nil, ErrProfessionalNameRequired
	}

	workingHours, err := normalizeWorkingHours(req.WorkingHours)
	if err != nil {
		return Professional{}, nil, err
	}

	role := normalizeOptionalString(req.RoleOrSpecialty)
	color, err := normalizeCalendarColor(req.CalendarColor)
	if err != nil {
		return Professional{}, nil, err
	}

	return Professional{
		ClinicID:        clinicID,
		FullName:        fullName,
		RoleOrSpecialty: role,
		CalendarColor:   color,
		WorkingHours:    workingHours,
		IsActive:        true,
	}, normalizeServiceIDs(req.ServiceIDs), nil
}

func buildUpdateEntity(clinicID, professionalID string, req UpdateProfessionalRequest) (Professional, []string, error) {
	fullName := strings.TrimSpace(req.FullName)
	if fullName == "" {
		return Professional{}, nil, ErrProfessionalNameRequired
	}

	workingHours, err := normalizeWorkingHours(req.WorkingHours)
	if err != nil {
		return Professional{}, nil, err
	}

	role := normalizeOptionalString(req.RoleOrSpecialty)
	color, err := normalizeCalendarColor(req.CalendarColor)
	if err != nil {
		return Professional{}, nil, err
	}

	return Professional{
		ID:              professionalID,
		ClinicID:        clinicID,
		FullName:        fullName,
		RoleOrSpecialty: role,
		CalendarColor:   color,
		WorkingHours:    workingHours,
		IsActive:        req.IsActive,
	}, normalizeServiceIDs(req.ServiceIDs), nil
}

func normalizeWorkingHours(raw json.RawMessage) ([]byte, error) {
	if len(raw) == 0 {
		return []byte(`{}`), nil
	}

	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, ErrInvalidWorkingHours
	}

	return []byte(raw), nil
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func normalizeCalendarColor(value *string) (*string, error) {
	normalized := normalizeOptionalString(value)
	if normalized == nil {
		return nil, nil
	}
	if !colorPattern.MatchString(*normalized) {
		return nil, ErrInvalidCalendarColor
	}
	return normalized, nil
}

func normalizeServiceIDs(serviceIDs []string) []string {
	seen := make(map[string]struct{}, len(serviceIDs))
	normalized := make([]string, 0, len(serviceIDs))
	for _, serviceID := range serviceIDs {
		serviceID = strings.TrimSpace(serviceID)
		if serviceID == "" {
			continue
		}
		if _, exists := seen[serviceID]; exists {
			continue
		}
		seen[serviceID] = struct{}{}
		normalized = append(normalized, serviceID)
	}
	return normalized
}

func mapModelToResponse(m Professional) ProfessionalResponse {
	workingHours := json.RawMessage(m.WorkingHours)
	if len(workingHours) == 0 {
		workingHours = json.RawMessage(`{}`)
	}
	serviceIDs := json.RawMessage(m.ServiceIDs)
	if len(serviceIDs) == 0 {
		serviceIDs = json.RawMessage(`[]`)
	}

	return ProfessionalResponse{
		ID:              m.ID,
		ClinicID:        m.ClinicID,
		FullName:        m.FullName,
		RoleOrSpecialty: m.RoleOrSpecialty,
		CalendarColor:   m.CalendarColor,
		WorkingHours:    workingHours,
		IsActive:        m.IsActive,
		ServiceIDs:      serviceIDs,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func mapRepositoryError(err error) error {
	message := err.Error()
	if strings.Contains(message, "clinic_professionals_name_per_clinic_unique") {
		return ErrDuplicateProfessionalName
	}
	return err
}
