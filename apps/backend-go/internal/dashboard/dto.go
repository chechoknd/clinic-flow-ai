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

type ActionsResponse struct {
	Data []ActionItemDTO `json:"data"`
}

type ActionItemDTO struct {
	Type         string  `json:"type"`
	Tone         string  `json:"tone"`
	Priority     int     `json:"priority"`
	LeadID       string  `json:"lead_id"`
	FullName     string  `json:"full_name"`
	Phone        string  `json:"phone"`
	ServiceID    *string `json:"service_id,omitempty"`
	ServiceName  *string `json:"service_name,omitempty"`
	Status       string  `json:"status"`
	Source       string  `json:"source"`
	Reason       string  `json:"reason"`
	NextActionAt *string `json:"next_action_at,omitempty"`
	CreatedAt    string  `json:"created_at"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
