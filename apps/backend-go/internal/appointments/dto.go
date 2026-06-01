package appointments

import "time"

type AppointmentResponse struct {
	ID                 string              `json:"id"`
	ClinicID           string              `json:"clinic_id"`
	Professional       ProfessionalSummary `json:"professional"`
	Lead               *LeadSummary        `json:"lead,omitempty"`
	Service            ServiceSummary      `json:"service"`
	ContactName        string              `json:"contact_name"`
	ContactPhone       *string             `json:"contact_phone,omitempty"`
	StartsAt           time.Time           `json:"starts_at"`
	EndsAt             time.Time           `json:"ends_at"`
	Status             string              `json:"status"`
	Source             string              `json:"source"`
	ConfirmationStatus string              `json:"confirmation_status"`
	AdminNotes         *string             `json:"admin_notes,omitempty"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
}

type ProfessionalSummary struct {
	ID            string  `json:"id"`
	FullName      string  `json:"full_name"`
	CalendarColor *string `json:"calendar_color,omitempty"`
}

type LeadSummary struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
}

type ServiceSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListAppointmentsResponse struct {
	Data []AppointmentResponse `json:"data"`
}

type CreateAppointmentRequest struct {
	ProfessionalID string     `json:"professional_id"`
	LeadID         *string    `json:"lead_id,omitempty"`
	ServiceID      string     `json:"service_id"`
	ContactName    string     `json:"contact_name"`
	ContactPhone   *string    `json:"contact_phone,omitempty"`
	StartsAt       time.Time  `json:"starts_at"`
	EndsAt         *time.Time `json:"ends_at,omitempty"`
	DurationMins   *int       `json:"duration_minutes,omitempty"`
	Status         string     `json:"status"`
	Source         string     `json:"source"`
	AdminNotes     *string    `json:"admin_notes,omitempty"`
}

type UpdateAppointmentRequest struct {
	ProfessionalID     string     `json:"professional_id"`
	ServiceID          string     `json:"service_id"`
	ContactName        string     `json:"contact_name"`
	ContactPhone       *string    `json:"contact_phone,omitempty"`
	StartsAt           time.Time  `json:"starts_at"`
	EndsAt             *time.Time `json:"ends_at,omitempty"`
	DurationMins       *int       `json:"duration_minutes,omitempty"`
	Status             string     `json:"status"`
	ConfirmationStatus string     `json:"confirmation_status"`
	AdminNotes         *string    `json:"admin_notes,omitempty"`
}

type UpdateStatusRequest struct {
	Status    string  `json:"status"`
	AdminNote *string `json:"admin_note,omitempty"`
}

type RescheduleRequest struct {
	StartsAt     time.Time  `json:"starts_at"`
	EndsAt       *time.Time `json:"ends_at,omitempty"`
	DurationMins *int       `json:"duration_minutes,omitempty"`
	AdminNote    *string    `json:"admin_note,omitempty"`
}

type ConvertLeadRequest struct {
	ProfessionalID   string     `json:"professional_id"`
	ServiceID        string     `json:"service_id"`
	StartsAt         time.Time  `json:"starts_at"`
	EndsAt           *time.Time `json:"ends_at,omitempty"`
	DurationMins     *int       `json:"duration_minutes,omitempty"`
	Status           string     `json:"status"`
	AdminNotes       *string    `json:"admin_notes,omitempty"`
	UpdateLeadStatus bool       `json:"update_lead_status"`
}

type ConvertLeadResponse struct {
	LeadID            string `json:"lead_id"`
	LeadStatus        string `json:"lead_status"`
	AppointmentID     string `json:"appointment_id"`
	AppointmentStatus string `json:"appointment_status"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
