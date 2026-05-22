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
