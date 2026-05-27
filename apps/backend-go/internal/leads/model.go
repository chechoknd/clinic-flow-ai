package leads

import (
	"time"
)

type Lead struct {
	ID           string
	ClinicID     string
	FullName     string
	Phone        string
	ServiceID    *string
	ServiceName  *string
	Status       string
	Source       string
	NextActionAt *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type LeadNote struct {
	ID        string
	LeadID    string
	Body      string
	CreatedAt time.Time
}

type AIInsight struct {
	ID                  string
	ClinicID            string
	LeadID              string
	AIGenerationID      *string
	Intent              string
	DetectedObjections  []string
	CommercialSummary   string
	SuggestedNextAction string
	Source              string
	CreatedAt           time.Time
}
