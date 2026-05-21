package ai

type ReplySuggestionRequest struct {
	LeadID         string `json:"lead_id"`
	ServiceID      string `json:"service_id"`
	PatientMessage string `json:"patient_message"`
	DesiredTone    string `json:"desired_tone,omitempty"`
}

type ReplySuggestionResponse struct {
	GenerationID string            `json:"generation_id"`
	Variants     ReplyVariants      `json:"variants"`
	SafetyStatus string            `json:"safety_status"`
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

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
