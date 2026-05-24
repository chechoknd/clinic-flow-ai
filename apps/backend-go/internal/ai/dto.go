package ai

type ReplySuggestionRequest struct {
	LeadID         string `json:"lead_id"`
	ServiceID      string `json:"service_id"`
	PatientMessage string `json:"patient_message"`
	DesiredTone    string `json:"desired_tone,omitempty"`
}

type ReplySuggestionResponse struct {
	GenerationID string        `json:"generation_id"`
	Variants     ReplyVariants `json:"variants"`
	SafetyStatus string        `json:"safety_status"`
}

type ReplyVariants struct {
	Short           string `json:"short"`
	Persuasive      string `json:"persuasive"`
	Technical       string `json:"technical"`
	ClosingQuestion string `json:"closing_question"`
}

type ObjectionHandlerRequest struct {
	LeadID    string `json:"lead_id"`
	ServiceID string `json:"service_id"`
	Objection string `json:"objection"`
}

type ObjectionHandlerResponse struct {
	GenerationID        string `json:"generation_id"`
	ObjectionType       string `json:"objection_type"`
	RecommendedStrategy string `json:"recommended_strategy"`
	SuggestedMessage    string `json:"suggested_message"`
	ClosingQuestion     string `json:"closing_question"`
	SafetyStatus        string `json:"safety_status"`
}

type FollowUpMessageRequest struct {
	LeadID          string `json:"lead_id"`
	ServiceID       string `json:"service_id"`
	LastContactNote string `json:"last_contact_note,omitempty"`
}

type FollowUpMessageResponse struct {
	GenerationID      string `json:"generation_id"`
	SuggestedMessage  string `json:"suggested_message"`
	RecommendedTiming string `json:"recommended_timing"`
	NextStep          string `json:"next_step"`
	SafetyStatus      string `json:"safety_status"`
}

type AnalyzeConversationRequest struct {
	ConversationText string `json:"conversation_text"`
	Source           string `json:"source,omitempty"`
	LeadID           string `json:"lead_id,omitempty"`
	ServiceID        string `json:"service_id,omitempty"`
}

type AnalyzeConversationResponse struct {
	AnalysisID          string          `json:"analysis_id"`
	DetectedLead        DetectedLead    `json:"detected_lead"`
	DetectedService     DetectedService `json:"detected_service"`
	Intent              string          `json:"intent"`
	DetectedObjections  []string        `json:"detected_objections"`
	SuggestedStatus     string          `json:"suggested_status"`
	CommercialSummary   string          `json:"commercial_summary"`
	SuggestedReply      string          `json:"suggested_reply"`
	SuggestedNextAction string          `json:"suggested_next_action"`
	SuggestedFollowUpAt string          `json:"suggested_follow_up_at,omitempty"`
	SafetyStatus        string          `json:"safety_status"`
}

type DetectedLead struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}

type DetectedService struct {
	ServiceID   string `json:"service_id"`
	ServiceName string `json:"service_name"`
	Confidence  string `json:"confidence"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
