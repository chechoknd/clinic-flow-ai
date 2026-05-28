package leads

import (
	"time"
)

type LeadResponse struct {
	ID           string     `json:"id"`
	FullName     string     `json:"full_name"`
	Phone        string     `json:"phone"`
	ServiceID    *string    `json:"service_id,omitempty"`
	ServiceName  *string    `json:"service_name,omitempty"`
	Status       string     `json:"status"`
	Source       string     `json:"source"`
	NextActionAt *time.Time `json:"next_action_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type LeadDetailResponse struct {
	ID           string              `json:"id"`
	FullName     string              `json:"full_name"`
	Phone        string              `json:"phone"`
	Service      *LeadServiceDTO     `json:"service,omitempty"`
	Status       string              `json:"status"`
	Source       string              `json:"source"`
	Notes        []LeadNoteResponse  `json:"notes"`
	AIInsights   []AIInsightResponse `json:"ai_insights,omitempty"`
	NextActionAt *time.Time          `json:"next_action_at,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
}

type LeadServiceDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type LeadNoteResponse struct {
	ID             string    `json:"id"`
	Body           string    `json:"body"`
	ContactOutcome *string   `json:"contact_outcome,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type AIInsightResponse struct {
	ID                  string    `json:"id"`
	AnalysisID          *string   `json:"analysis_id,omitempty"`
	Intent              string    `json:"intent"`
	DetectedObjections  []string  `json:"detected_objections"`
	CommercialSummary   string    `json:"commercial_summary,omitempty"`
	SuggestedNextAction string    `json:"suggested_next_action,omitempty"`
	Source              string    `json:"source,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
}

type ReviewedAIAnalysisRequest struct {
	AnalysisID          string   `json:"analysis_id,omitempty"`
	Intent              string   `json:"intent"`
	DetectedObjections  []string `json:"detected_objections,omitempty"`
	CommercialSummary   string   `json:"commercial_summary,omitempty"`
	SuggestedNextAction string   `json:"suggested_next_action,omitempty"`
	Source              string   `json:"source,omitempty"`
}

type CreateLeadRequest struct {
	FullName           string                     `json:"full_name"`
	Phone              string                     `json:"phone"`
	ServiceID          *string                    `json:"service_id,omitempty"`
	Status             string                     `json:"status"`
	Source             string                     `json:"source"`
	Notes              string                     `json:"notes,omitempty"`
	NextActionAt       *time.Time                 `json:"next_action_at,omitempty"`
	ReviewedAIAnalysis *ReviewedAIAnalysisRequest `json:"reviewed_ai_analysis,omitempty"`
}

type UpdateLeadRequest struct {
	Status             string                     `json:"status"`
	Note               string                     `json:"note,omitempty"`
	ContactOutcome     string                     `json:"contact_outcome,omitempty"`
	NextActionAt       *time.Time                 `json:"next_action_at,omitempty"`
	ClearNextActionAt  bool                       `json:"clear_next_action_at,omitempty"`
	ReviewedAIAnalysis *ReviewedAIAnalysisRequest `json:"reviewed_ai_analysis,omitempty"`
}

type FollowUpFilterResponse struct {
	Data       []LeadResponse     `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

type CompleteFollowUpRequest struct {
	Status         string `json:"status,omitempty"`
	Note           string `json:"note,omitempty"`
	ContactOutcome string `json:"contact_outcome,omitempty"`
}

type RescheduleFollowUpRequest struct {
	NextActionAt time.Time `json:"next_action_at"`
	Note         string    `json:"note,omitempty"`
}

type PaginatedLeadsResponse struct {
	Data       []LeadResponse     `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

type PaginationResponse struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
