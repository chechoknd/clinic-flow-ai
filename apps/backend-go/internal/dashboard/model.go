package dashboard

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
