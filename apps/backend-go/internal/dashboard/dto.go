package dashboard

type SummaryResponse struct {
	LeadsTotal            int             `json:"leads_total"`
	LeadsByStatus         map[string]int  `json:"leads_by_status"`
	TopServices           []TopServiceDTO `json:"top_services"`
	PendingFollowUpsToday int             `json:"pending_followups_today"`
	OverdueFollowUps      int             `json:"overdue_followups"`
	UpcomingFollowUps     int             `json:"upcoming_followups"`
	ConversionRate        float64         `json:"conversion_rate"`
}

type TopServiceDTO struct {
	ServiceID   string `json:"service_id"`
	ServiceName string `json:"service_name"`
	LeadCount   int    `json:"lead_count"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ScheduleSummaryResponse struct {
	Date                             string                    `json:"date"`
	TodaysAppointments               int                       `json:"todays_appointments"`
	AppointmentsPendingConfirmation int                       `json:"appointments_pending_confirmation"`
	AvailableSlots                   int                       `json:"available_slots"`
	HotLeadsWithoutAppointment       int                       `json:"hot_leads_without_appointment"`
	OverdueFollowUps                 int                       `json:"overdue_followups"`
	AppointmentsByProfessional       []ProfessionalCount       `json:"appointments_by_professional"`
	LeadsConvertedToAppointments     int                       `json:"leads_converted_to_appointments"`
	TopServicesByScheduleDemand      []ServiceCount            `json:"top_services_by_schedule_demand"`
}

type ProfessionalCount struct {
	ProfessionalID   string `json:"professional_id"`
	ProfessionalName string `json:"professional_name"`
	AppointmentCount int    `json:"appointment_count"`
}

type ServiceCount struct {
	ServiceID        string `json:"service_id"`
	ServiceName      string `json:"service_name"`
	AppointmentCount int    `json:"appointment_count"`
}
