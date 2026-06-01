package dashboard

import "time"

type StatusCount struct {
	Status string
	Count  int
}

type TopService struct {
	ServiceID   string
	ServiceName string
	LeadCount   int
}

type Summary struct {
	LeadsTotal            int
	LeadsByStatus         map[string]int
	TopServices           []TopService
	PendingFollowUpsToday int
	OverdueFollowUps      int
	UpcomingFollowUps     int
	ConversionRate        float64
}

type ScheduleSummary struct {
	Date                            string
	TodaysAppointments              int
	AppointmentsPendingConfirmation int
	AvailableSlots                  int
	HotLeadsWithoutAppointment      int
	OverdueFollowUps                int
	AppointmentsByProfessional      []ProfessionalCountModel
	LeadsConvertedToAppointments    int
	TopServicesByScheduleDemand     []ServiceCountModel
}

type ProfessionalCountModel struct {
	ProfessionalID   string
	ProfessionalName string
	AppointmentCount int
}

type ServiceCountModel struct {
	ServiceID        string
	ServiceName      string
	AppointmentCount int
}

type ProfessionalHours struct {
	ID           string
	WorkingHours []byte
}

type AppointmentSummary struct {
	StartsAt       time.Time
	EndsAt         time.Time
	ProfessionalID string
}
