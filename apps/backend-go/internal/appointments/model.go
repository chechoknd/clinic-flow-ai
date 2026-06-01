package appointments

import "time"

type Appointment struct {
	ID                 string
	ClinicID           string
	ProfessionalID     string
	ProfessionalName   string
	ProfessionalColor  *string
	LeadID             *string
	LeadName           *string
	ServiceID          string
	ServiceName        string
	ContactName        string
	ContactPhone       *string
	StartsAt           time.Time
	EndsAt             time.Time
	Status             string
	ConfirmationStatus string
	Source             string
	AdminNotes         *string
	CreatedByUserID    *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type ListFilters struct {
	Date           string
	DateFrom       string
	DateTo         string
	ProfessionalID string
	ServiceID      string
	Status         string
	LeadID         string
}
