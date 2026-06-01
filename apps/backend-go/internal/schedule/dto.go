package schedule

import "time"

type AvailabilityRequest struct {
	ProfessionalID  string    `json:"professional_id"`
	ServiceID       string    `json:"service_id"`
	DateFrom        time.Time `json:"date_from"`
	DateTo          time.Time `json:"date_to"`
	DurationMinutes int       `json:"duration_minutes"`
}

type TimeSlot struct {
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
}

type AvailabilityResponse struct {
	ProfessionalID string     `json:"professional_id"`
	Slots          []TimeSlot `json:"slots"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}
