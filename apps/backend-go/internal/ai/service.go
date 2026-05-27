package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/clinics"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/leads"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/services"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/shared"
)

type Service struct {
	provider       Provider
	generations    GenerationRepository
	clinicService  *clinics.Service
	serviceService *services.Service
	leadService    *leads.Service
}

func NewService(provider Provider, generations GenerationRepository, cs *clinics.Service, ss *services.Service, ls *leads.Service) *Service {
	return &Service{
		provider:       provider,
		generations:    generations,
		clinicService:  cs,
		serviceService: ss,
		leadService:    ls,
	}
}

func (s *Service) ReplySuggestion(ctx context.Context, clinicID, userID string, req ReplySuggestionRequest) (ReplySuggestionResponse, error) {
	if err := validateReplySuggestionRequest(req); err != nil {
		return ReplySuggestionResponse{}, err
	}

	aiCtx, err := s.getAIContext(ctx, clinicID, req.ServiceID, req.LeadID)
	if err != nil {
		return ReplySuggestionResponse{}, err
	}

	systemPrompt := BuildSystemPrompt(aiCtx)
	userPrompt := BuildReplySuggestionUserPrompt(req.PatientMessage, req.DesiredTone)

	rawRes, err := s.provider.Generate(ctx, systemPrompt, userPrompt)
	if err != nil {
		s.recordGeneration(ctx, generationRecordInput{
			ClinicID:        clinicID,
			UserID:          userID,
			Feature:         "reply_suggestion",
			Status:          "provider_error",
			ErrorCode:       "AI_PROVIDER_ERROR",
			InputCharCount:  len(systemPrompt) + len(userPrompt),
			OutputCharCount: len(rawRes),
		})
		return ReplySuggestionResponse{}, ProviderError{Provider: s.provider.Name(), Err: err}
	}

	var variants ReplyVariants
	if err := json.Unmarshal([]byte(rawRes), &variants); err != nil {
		s.recordGeneration(ctx, generationRecordInput{
			ClinicID:        clinicID,
			UserID:          userID,
			Feature:         "reply_suggestion",
			Status:          "response_error",
			ErrorCode:       "AI_RESPONSE_ERROR",
			InputCharCount:  len(systemPrompt) + len(userPrompt),
			OutputCharCount: len(rawRes),
		})
		return ReplySuggestionResponse{}, ResponseError{Err: err}
	}

	safetyStatus, safe := ValidateSafety(rawRes)
	if !safe {
		s.recordGeneration(ctx, generationRecordInput{
			ClinicID:        clinicID,
			UserID:          userID,
			Feature:         "reply_suggestion",
			Status:          "safety_blocked",
			SafetyStatus:    safetyStatus,
			ErrorCode:       "AI_SAFETY_BLOCKED",
			InputCharCount:  len(systemPrompt) + len(userPrompt),
			OutputCharCount: len(rawRes),
		})
		return ReplySuggestionResponse{}, SafetyBlockedError{Status: safetyStatus}
	}

	generationID := s.recordGeneration(ctx, generationRecordInput{
		ClinicID:        clinicID,
		UserID:          userID,
		Feature:         "reply_suggestion",
		Status:          "success",
		SafetyStatus:    safetyStatus,
		InputCharCount:  len(systemPrompt) + len(userPrompt),
		OutputCharCount: len(rawRes),
	})

	return ReplySuggestionResponse{
		GenerationID: generationID,
		Variants:     variants,
		SafetyStatus: safetyStatus,
	}, nil
}

func (s *Service) ObjectionHandler(ctx context.Context, clinicID, userID string, req ObjectionHandlerRequest) (ObjectionHandlerResponse, error) {
	if err := validateObjectionHandlerRequest(req); err != nil {
		return ObjectionHandlerResponse{}, err
	}

	aiCtx, err := s.getAIContext(ctx, clinicID, req.ServiceID, req.LeadID)
	if err != nil {
		return ObjectionHandlerResponse{}, err
	}

	systemPrompt := BuildSystemPrompt(aiCtx)
	userPrompt := BuildObjectionHandlerUserPrompt(req.Objection)

	rawRes, err := s.provider.Generate(ctx, systemPrompt, userPrompt)
	if err != nil {
		s.recordGeneration(ctx, generationRecordInput{
			ClinicID:        clinicID,
			UserID:          userID,
			Feature:         "objection_handler",
			Status:          "provider_error",
			ErrorCode:       "AI_PROVIDER_ERROR",
			InputCharCount:  len(systemPrompt) + len(userPrompt),
			OutputCharCount: len(rawRes),
		})
		return ObjectionHandlerResponse{}, ProviderError{Provider: s.provider.Name(), Err: err}
	}

	var res ObjectionHandlerResponse
	if err := json.Unmarshal([]byte(rawRes), &res); err != nil {
		s.recordGeneration(ctx, generationRecordInput{
			ClinicID:        clinicID,
			UserID:          userID,
			Feature:         "objection_handler",
			Status:          "response_error",
			ErrorCode:       "AI_RESPONSE_ERROR",
			InputCharCount:  len(systemPrompt) + len(userPrompt),
			OutputCharCount: len(rawRes),
		})
		return ObjectionHandlerResponse{}, ResponseError{Err: err}
	}

	safetyStatus, safe := ValidateSafety(rawRes)
	if !safe {
		s.recordGeneration(ctx, generationRecordInput{
			ClinicID:        clinicID,
			UserID:          userID,
			Feature:         "objection_handler",
			Status:          "safety_blocked",
			SafetyStatus:    safetyStatus,
			ErrorCode:       "AI_SAFETY_BLOCKED",
			InputCharCount:  len(systemPrompt) + len(userPrompt),
			OutputCharCount: len(rawRes),
		})
		return ObjectionHandlerResponse{}, SafetyBlockedError{Status: safetyStatus}
	}
	res.GenerationID = s.recordGeneration(ctx, generationRecordInput{
		ClinicID:        clinicID,
		UserID:          userID,
		Feature:         "objection_handler",
		Status:          "success",
		SafetyStatus:    safetyStatus,
		InputCharCount:  len(systemPrompt) + len(userPrompt),
		OutputCharCount: len(rawRes),
	})
	res.SafetyStatus = safetyStatus

	return res, nil
}

func (s *Service) FollowUpMessage(ctx context.Context, clinicID, userID string, req FollowUpMessageRequest) (FollowUpMessageResponse, error) {
	if err := validateFollowUpMessageRequest(req); err != nil {
		return FollowUpMessageResponse{}, err
	}

	aiCtx, err := s.getAIContext(ctx, clinicID, req.ServiceID, req.LeadID)
	if err != nil {
		return FollowUpMessageResponse{}, err
	}

	systemPrompt := BuildSystemPrompt(aiCtx)
	userPrompt := BuildFollowUpMessageUserPrompt(req.LastContactNote)

	rawRes, err := s.provider.Generate(ctx, systemPrompt, userPrompt)
	if err != nil {
		s.recordGeneration(ctx, generationRecordInput{
			ClinicID:        clinicID,
			UserID:          userID,
			Feature:         "follow_up_message",
			Status:          "provider_error",
			ErrorCode:       "AI_PROVIDER_ERROR",
			InputCharCount:  len(systemPrompt) + len(userPrompt),
			OutputCharCount: len(rawRes),
		})
		return FollowUpMessageResponse{}, ProviderError{Provider: s.provider.Name(), Err: err}
	}

	var res FollowUpMessageResponse
	if err := json.Unmarshal([]byte(rawRes), &res); err != nil {
		s.recordGeneration(ctx, generationRecordInput{
			ClinicID:        clinicID,
			UserID:          userID,
			Feature:         "follow_up_message",
			Status:          "response_error",
			ErrorCode:       "AI_RESPONSE_ERROR",
			InputCharCount:  len(systemPrompt) + len(userPrompt),
			OutputCharCount: len(rawRes),
		})
		return FollowUpMessageResponse{}, ResponseError{Err: err}
	}

	safetyStatus, safe := ValidateSafety(rawRes)
	if !safe {
		s.recordGeneration(ctx, generationRecordInput{
			ClinicID:        clinicID,
			UserID:          userID,
			Feature:         "follow_up_message",
			Status:          "safety_blocked",
			SafetyStatus:    safetyStatus,
			ErrorCode:       "AI_SAFETY_BLOCKED",
			InputCharCount:  len(systemPrompt) + len(userPrompt),
			OutputCharCount: len(rawRes),
		})
		return FollowUpMessageResponse{}, SafetyBlockedError{Status: safetyStatus}
	}
	res.GenerationID = s.recordGeneration(ctx, generationRecordInput{
		ClinicID:        clinicID,
		UserID:          userID,
		Feature:         "follow_up_message",
		Status:          "success",
		SafetyStatus:    safetyStatus,
		InputCharCount:  len(systemPrompt) + len(userPrompt),
		OutputCharCount: len(rawRes),
	})
	res.SafetyStatus = safetyStatus

	return res, nil
}

func (s *Service) getAIContext(ctx context.Context, clinicID, serviceID, leadID string) (Context, error) {
	clinicRes, err := s.clinicService.Current(ctx, clinicID)
	if err != nil {
		return Context{}, err
	}

	aiCtx := Context{
		ClinicName:        clinicRes.Name,
		ClinicType:        clinicRes.ClinicType,
		City:              clinicRes.City,
		CommunicationTone: clinicRes.CommunicationTone,
		ServiceName:       "servicio no especificado",
		ServicePriceFrom:  "no informado",
		ServiceBenefits:   "[]",
		ServiceFAQ:        "[]",
		CommonObjections:  "[]",
		PatientName:       "prospecto",
		LeadNotes:         "[]",
	}

	leadID = strings.TrimSpace(leadID)
	serviceID = strings.TrimSpace(serviceID)
	if leadID == "" && serviceID == "" {
		return aiCtx, nil
	}

	if leadID == "" {
		serviceRes, err := s.serviceService.Get(ctx, clinicID, serviceID)
		if err != nil {
			return Context{}, err
		}
		aiCtx.ServiceName = serviceRes.Name
		aiCtx.ServicePriceFrom = formatServicePriceFrom(serviceRes.PriceFrom, serviceRes.CurrencyCode)
		aiCtx.ServiceBenefits = string(serviceRes.Benefits)
		aiCtx.ServiceFAQ = string(serviceRes.FAQ)
		aiCtx.CommonObjections = string(serviceRes.CommonObjections)
		return aiCtx, nil
	}

	leadRes, err := s.leadService.Get(ctx, clinicID, leadID)
	if err != nil {
		return Context{}, err
	}
	var leadNotes []string
	for _, n := range leadRes.Notes {
		leadNotes = append(leadNotes, n.Body)
	}
	aiCtx.PatientName = leadRes.FullName
	aiCtx.LeadNotes = fmt.Sprintf("[%s]", strings.Join(leadNotes, ", "))

	if serviceID == "" && leadRes.Service != nil {
		serviceID = leadRes.Service.ID
	}
	if serviceID == "" {
		return aiCtx, nil
	}

	serviceRes, err := s.serviceService.Get(ctx, clinicID, serviceID)
	if err != nil {
		return Context{}, err
	}
	aiCtx.ServiceName = serviceRes.Name
	aiCtx.ServicePriceFrom = formatServicePriceFrom(serviceRes.PriceFrom, serviceRes.CurrencyCode)
	aiCtx.ServiceBenefits = string(serviceRes.Benefits)
	aiCtx.ServiceFAQ = string(serviceRes.FAQ)
	aiCtx.CommonObjections = string(serviceRes.CommonObjections)

	return aiCtx, nil
}

func formatServicePriceFrom(price *services.DecimalString, currencyCode string) string {
	if price == nil {
		return "no informado"
	}
	currency, ok := shared.CurrencyByCode(currencyCode)
	if !ok {
		return price.String()
	}
	return fmt.Sprintf("%s %s", currency.Code, price.String())
}

func validateReplySuggestionRequest(req ReplySuggestionRequest) error {
	message := strings.TrimSpace(req.PatientMessage)
	if message == "" {
		return ValidationError{Field: "patient_message", Message: "patient_message is required"}
	}
	if len(message) > 2000 {
		return ValidationError{Field: "patient_message", Message: "patient_message is too long"}
	}
	if len(strings.TrimSpace(req.DesiredTone)) > 120 {
		return ValidationError{Field: "desired_tone", Message: "desired_tone is too long"}
	}
	return nil
}

func validateObjectionHandlerRequest(req ObjectionHandlerRequest) error {
	objection := strings.TrimSpace(req.Objection)
	if objection == "" {
		return ValidationError{Field: "objection", Message: "objection is required"}
	}
	if len(objection) > 1000 {
		return ValidationError{Field: "objection", Message: "objection is too long"}
	}
	return nil
}

func validateFollowUpMessageRequest(req FollowUpMessageRequest) error {
	if strings.TrimSpace(req.LeadID) == "" {
		return ValidationError{Field: "lead_id", Message: "lead_id is required"}
	}
	if len(strings.TrimSpace(req.LastContactNote)) > 1000 {
		return ValidationError{Field: "last_contact_note", Message: "last_contact_note is too long"}
	}
	return nil
}

func (s *Service) AnalyzeConversation(ctx context.Context, clinicID, userID string, req AnalyzeConversationRequest) (AnalyzeConversationResponse, error) {
	if err := validateAnalyzeConversationRequest(req); err != nil {
		return AnalyzeConversationResponse{}, err
	}

	aiCtx, err := s.getAIContext(ctx, clinicID, req.ServiceID, req.LeadID)
	if err != nil {
		return AnalyzeConversationResponse{}, err
	}

	systemPrompt := BuildSystemPrompt(aiCtx)
	userPrompt := BuildAnalyzeConversationUserPrompt(strings.TrimSpace(req.ConversationText), strings.TrimSpace(req.Source))

	rawRes, err := s.provider.Generate(ctx, systemPrompt, userPrompt)
	if err != nil {
		s.recordGeneration(ctx, generationRecordInput{
			ClinicID:        clinicID,
			UserID:          userID,
			Feature:         "analyze_conversation",
			Status:          "provider_error",
			ErrorCode:       "AI_PROVIDER_ERROR",
			InputCharCount:  len(systemPrompt) + len(userPrompt),
			OutputCharCount: len(rawRes),
		})
		return AnalyzeConversationResponse{}, ProviderError{Provider: s.provider.Name(), Err: err}
	}

	var res AnalyzeConversationResponse
	if err := json.Unmarshal([]byte(rawRes), &res); err != nil {
		s.recordGeneration(ctx, generationRecordInput{
			ClinicID:        clinicID,
			UserID:          userID,
			Feature:         "analyze_conversation",
			Status:          "response_error",
			ErrorCode:       "AI_RESPONSE_ERROR",
			InputCharCount:  len(systemPrompt) + len(userPrompt),
			OutputCharCount: len(rawRes),
		})
		return AnalyzeConversationResponse{}, ResponseError{Err: err}
	}

	safetyStatus, safe := ValidateSafety(rawRes)
	if !safe {
		s.recordGeneration(ctx, generationRecordInput{
			ClinicID:        clinicID,
			UserID:          userID,
			Feature:         "analyze_conversation",
			Status:          "safety_blocked",
			SafetyStatus:    safetyStatus,
			ErrorCode:       "AI_SAFETY_BLOCKED",
			InputCharCount:  len(systemPrompt) + len(userPrompt),
			OutputCharCount: len(rawRes),
		})
		return AnalyzeConversationResponse{}, SafetyBlockedError{Status: safetyStatus}
	}
	res.AnalysisID = s.recordGeneration(ctx, generationRecordInput{
		ClinicID:        clinicID,
		UserID:          userID,
		Feature:         "analyze_conversation",
		Status:          "success",
		SafetyStatus:    safetyStatus,
		InputCharCount:  len(systemPrompt) + len(userPrompt),
		OutputCharCount: len(rawRes),
	})
	res.SafetyStatus = safetyStatus
	res.SuggestedStatus = normalizeSuggestedStatus(res.SuggestedStatus)
	res.Intent = normalizeConfidence(res.Intent)
	res.DetectedService.Confidence = normalizeConfidence(res.DetectedService.Confidence)

	return res, nil
}

type generationRecordInput struct {
	ClinicID        string
	UserID          string
	Feature         string
	Status          string
	SafetyStatus    string
	ErrorCode       string
	InputCharCount  int
	OutputCharCount int
}

func (s *Service) recordGeneration(ctx context.Context, input generationRecordInput) string {
	fallbackID := "gen-" + input.ClinicID
	if s.generations == nil {
		return fallbackID
	}

	model := strings.TrimSpace(s.provider.Model())
	if model == "" {
		model = "unknown"
	}

	id, err := s.generations.CreateGeneration(ctx, GenerationRecord{
		ClinicID:        input.ClinicID,
		UserID:          input.UserID,
		Feature:         input.Feature,
		Provider:        s.provider.Name(),
		Model:           model,
		Status:          input.Status,
		SafetyStatus:    input.SafetyStatus,
		ErrorCode:       input.ErrorCode,
		InputCharCount:  input.InputCharCount,
		OutputCharCount: input.OutputCharCount,
	})
	if err != nil {
		log.Printf("ai generation metadata record failed: clinic_id=%s feature=%s provider=%s status=%s: %v", input.ClinicID, input.Feature, s.provider.Name(), input.Status, err)
		return fallbackID
	}
	return id
}

func validateAnalyzeConversationRequest(req AnalyzeConversationRequest) error {
	conversation := strings.TrimSpace(req.ConversationText)
	if conversation == "" {
		return ValidationError{Field: "conversation_text", Message: "conversation_text is required"}
	}
	if len(conversation) < 20 {
		return ValidationError{Field: "conversation_text", Message: "conversation_text is too short"}
	}
	if len(conversation) > 8000 {
		return ValidationError{Field: "conversation_text", Message: "conversation_text is too long"}
	}
	if len(strings.TrimSpace(req.Source)) > 80 {
		return ValidationError{Field: "source", Message: "source is too long"}
	}
	return nil
}

func normalizeConfidence(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "low", "medium", "high":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return "medium"
	}
}

func normalizeSuggestedStatus(status string) string {
	switch strings.TrimSpace(status) {
	case "Nuevo", "Contactado", "Interesado", "Agendado", "No Respondio", "Perdido", "Convertido":
		return strings.TrimSpace(status)
	default:
		return "Nuevo"
	}
}
