package professionals

import (
	"encoding/json"
	"time"
)

type ProfessionalResponse struct {
	ID              string          `json:"id"`
	ClinicID        string          `json:"clinic_id"`
	FullName        string          `json:"full_name"`
	RoleOrSpecialty *string         `json:"role_or_specialty,omitempty"`
	CalendarColor   *string         `json:"calendar_color,omitempty"`
	WorkingHours    json.RawMessage `json:"working_hours"`
	IsActive        bool            `json:"is_active"`
	ServiceIDs      json.RawMessage `json:"service_ids"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type CreateProfessionalRequest struct {
	FullName        string          `json:"full_name"`
	RoleOrSpecialty *string         `json:"role_or_specialty,omitempty"`
	CalendarColor   *string         `json:"calendar_color,omitempty"`
	WorkingHours    json.RawMessage `json:"working_hours"`
	ServiceIDs      []string        `json:"service_ids"`
}

type UpdateProfessionalRequest struct {
	FullName        string          `json:"full_name"`
	RoleOrSpecialty *string         `json:"role_or_specialty,omitempty"`
	CalendarColor   *string         `json:"calendar_color,omitempty"`
	WorkingHours    json.RawMessage `json:"working_hours"`
	IsActive        bool            `json:"is_active"`
	ServiceIDs      []string        `json:"service_ids"`
}

type ListFilters struct {
	IsActive  *bool
	ServiceID string
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
