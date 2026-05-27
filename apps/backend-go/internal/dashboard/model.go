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

type ActionItem struct {
	Type         string
	Tone         string
	Priority     int
	LeadID       string
	FullName     string
	Phone        string
	ServiceID    *string
	ServiceName  *string
	Status       string
	Source       string
	Reason       string
	NextActionAt *time.Time
	CreatedAt    time.Time
}
