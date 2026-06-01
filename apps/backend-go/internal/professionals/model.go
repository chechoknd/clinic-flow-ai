package professionals

import "time"

type Professional struct {
	ID              string
	ClinicID        string
	FullName        string
	RoleOrSpecialty *string
	CalendarColor   *string
	WorkingHours    []byte
	IsActive        bool
	ServiceIDs      []byte
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
