package services

import "encoding/json"

type ServiceResponse struct {
	ID               string          `json:"id"`
	ClinicID         string          `json:"clinic_id"`
	Name             string          `json:"name"`
	Description      *string         `json:"description,omitempty"`
	DurationMinutes  *int            `json:"duration_minutes,omitempty"`
	PriceFrom        *float64        `json:"price_from,omitempty"`
	CurrencyCode     string          `json:"currency_code"`
	Benefits         json.RawMessage `json:"benefits"`
	FAQ              json.RawMessage `json:"faq"`
	CommonObjections json.RawMessage `json:"common_objections"`
	IsActive         bool            `json:"is_active"`
}

type CreateServiceRequest struct {
	Name             string          `json:"name"`
	Description      *string         `json:"description,omitempty"`
	DurationMinutes  *int            `json:"duration_minutes,omitempty"`
	PriceFrom        *float64        `json:"price_from,omitempty"`
	Benefits         json.RawMessage `json:"benefits"`
	FAQ              json.RawMessage `json:"faq"`
	CommonObjections json.RawMessage `json:"common_objections"`
}

type UpdateServiceRequest struct {
	Name             string          `json:"name"`
	Description      *string         `json:"description,omitempty"`
	DurationMinutes  *int            `json:"duration_minutes,omitempty"`
	PriceFrom        *float64        `json:"price_from,omitempty"`
	Benefits         json.RawMessage `json:"benefits"`
	FAQ              json.RawMessage `json:"faq"`
	CommonObjections json.RawMessage `json:"common_objections"`
	IsActive         bool            `json:"is_active"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
