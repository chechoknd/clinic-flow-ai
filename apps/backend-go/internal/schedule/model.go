package schedule

import "time"

type Professional struct {
	ID           string
	ClinicID     string
	FullName     string
	WorkingHours []byte
	IsActive     bool
}

type AppointmentRange struct {
	StartsAt time.Time
	EndsAt   time.Time
}

type DayWorkingInterval struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type WorkingHours struct {
	Monday    []DayWorkingInterval `json:"monday"`
	Tuesday   []DayWorkingInterval `json:"tuesday"`
	Wednesday []DayWorkingInterval `json:"wednesday"`
	Thursday  []DayWorkingInterval `json:"thursday"`
	Friday    []DayWorkingInterval `json:"friday"`
	Saturday  []DayWorkingInterval `json:"saturday"`
	Sunday    []DayWorkingInterval `json:"sunday"`
}
