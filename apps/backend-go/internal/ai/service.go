package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/clinics"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/leads"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/services"
)

type Service struct {
	provider       Provider
	clinicService  *clinics.Service
	serviceService *services.Service
	leadService    *leads.Service
}

func NewService(provider Provider, cs *clinics.Service, ss *services.Service, ls *leads.Service) *Service {
	return &Service{
		provider:       provider,
		clinicService:  cs,
		serviceService: ss,
		leadService:    ls,
	}
}

func (s *Service) ReplySuggestion(ctx context.Context, clinicID string, req ReplySuggestionRequest) (ReplySuggestionResponse, error) {
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
		return ReplySuggestionResponse{}, ProviderError{Provider: s.provider.Name(), Err: err}
	}

	var variants ReplyVariants
	if err := json.Unmarshal([]byte(rawRes), &variants); err != nil {
		return ReplySuggestionResponse{}, ResponseError{Err: err}
	}

	safetyStatus, safe := ValidateSafety(rawRes)
	if !safe {
		return ReplySuggestionResponse{}, SafetyBlockedError{Status: safetyStatus}
	}

	return ReplySuggestionResponse{
		GenerationID: "gen-" + clinicID, // Future: UUID
		Variants:     variants,
		SafetyStatus: safetyStatus,
	}, nil
}

func (s *Service) ObjectionHandler(ctx context.Context, clinicID string, req ObjectionHandlerRequest) (ObjectionHandlerResponse, error) {
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
		return ObjectionHandlerResponse{}, ProviderError{Provider: s.provider.Name(), Err: err}
	}

	var res ObjectionHandlerResponse
	if err := json.Unmarshal([]byte(rawRes), &res); err != nil {
		return ObjectionHandlerResponse{}, ResponseError{Err: err}
	}

	safetyStatus, safe := ValidateSafety(rawRes)
	if !safe {
		return ObjectionHandlerResponse{}, SafetyBlockedError{Status: safetyStatus}
	}
	res.GenerationID = "gen-" + clinicID
	res.SafetyStatus = safetyStatus

	return res, nil
}

func (s *Service) FollowUpMessage(ctx context.Context, clinicID string, req FollowUpMessageRequest) (FollowUpMessageResponse, error) {
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
		return FollowUpMessageResponse{}, ProviderError{Provider: s.provider.Name(), Err: err}
	}

	var res FollowUpMessageResponse
	if err := json.Unmarshal([]byte(rawRes), &res); err != nil {
		return FollowUpMessageResponse{}, ResponseError{Err: err}
	}

	safetyStatus, safe := ValidateSafety(rawRes)
	if !safe {
		return FollowUpMessageResponse{}, SafetyBlockedError{Status: safetyStatus}
	}
	res.GenerationID = "gen-" + clinicID
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
		aiCtx.ServicePriceFrom = formatServicePriceFrom(serviceRes.PriceFrom)
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
	aiCtx.ServicePriceFrom = formatServicePriceFrom(serviceRes.PriceFrom)
	aiCtx.ServiceBenefits = string(serviceRes.Benefits)
	aiCtx.ServiceFAQ = string(serviceRes.FAQ)
	aiCtx.CommonObjections = string(serviceRes.CommonObjections)

	return aiCtx, nil
}

func formatServicePriceFrom(price *float64) string {
	if price == nil {
		return "no informado"
	}
	return strconv.FormatFloat(*price, 'f', -1, 64)
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
