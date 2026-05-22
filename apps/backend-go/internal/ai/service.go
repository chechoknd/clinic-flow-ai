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
	aiCtx, err := s.getAIContext(ctx, clinicID, req.ServiceID, req.LeadID)
	if err != nil {
		return ReplySuggestionResponse{}, err
	}

	systemPrompt := BuildSystemPrompt(aiCtx)
	userPrompt := BuildReplySuggestionUserPrompt(req.PatientMessage, req.DesiredTone)

	rawRes, err := s.provider.Generate(ctx, systemPrompt, userPrompt)
	if err != nil {
		return ReplySuggestionResponse{}, err
	}

	var variants ReplyVariants
	if err := json.Unmarshal([]byte(rawRes), &variants); err != nil {
		// Fallback if AI didn't return valid JSON
		return ReplySuggestionResponse{}, fmt.Errorf("failed to parse AI response: %w", err)
	}

	// Basic safety check on all variants
	safetyStatus, _ := ValidateSafety(rawRes)

	return ReplySuggestionResponse{
		GenerationID: "gen-" + clinicID, // Future: UUID
		Variants:     variants,
		SafetyStatus: safetyStatus,
	}, nil
}

func (s *Service) ObjectionHandler(ctx context.Context, clinicID string, req ObjectionHandlerRequest) (ObjectionHandlerResponse, error) {
	aiCtx, err := s.getAIContext(ctx, clinicID, req.ServiceID, req.LeadID)
	if err != nil {
		return ObjectionHandlerResponse{}, err
	}

	systemPrompt := BuildSystemPrompt(aiCtx)
	userPrompt := BuildObjectionHandlerUserPrompt(req.Objection)

	rawRes, err := s.provider.Generate(ctx, systemPrompt, userPrompt)
	if err != nil {
		return ObjectionHandlerResponse{}, err
	}

	var res ObjectionHandlerResponse
	if err := json.Unmarshal([]byte(rawRes), &res); err != nil {
		return ObjectionHandlerResponse{}, fmt.Errorf("failed to parse AI response: %w", err)
	}

	safetyStatus, _ := ValidateSafety(rawRes)
	res.GenerationID = "gen-" + clinicID
	res.SafetyStatus = safetyStatus

	return res, nil
}

func (s *Service) FollowUpMessage(ctx context.Context, clinicID string, req FollowUpMessageRequest) (FollowUpMessageResponse, error) {
	aiCtx, err := s.getAIContext(ctx, clinicID, req.ServiceID, req.LeadID)
	if err != nil {
		return FollowUpMessageResponse{}, err
	}

	systemPrompt := BuildSystemPrompt(aiCtx)
	userPrompt := BuildFollowUpMessageUserPrompt(req.LastContactNote)

	rawRes, err := s.provider.Generate(ctx, systemPrompt, userPrompt)
	if err != nil {
		return FollowUpMessageResponse{}, err
	}

	var res FollowUpMessageResponse
	if err := json.Unmarshal([]byte(rawRes), &res); err != nil {
		return FollowUpMessageResponse{}, fmt.Errorf("failed to parse AI response: %w", err)
	}

	safetyStatus, _ := ValidateSafety(rawRes)
	res.GenerationID = "gen-" + clinicID
	res.SafetyStatus = safetyStatus

	return res, nil
}

func (s *Service) getAIContext(ctx context.Context, clinicID, serviceID, leadID string) (Context, error) {
	clinicRes, err := s.clinicService.Current(ctx, clinicID)
	if err != nil {
		return Context{}, err
	}

	serviceRes, err := s.serviceService.Get(ctx, clinicID, serviceID)
	if err != nil {
		return Context{}, err
	}

	leadRes, err := s.leadService.Get(ctx, clinicID, leadID)
	if err != nil {
		return Context{}, err
	}

	var leadNotes []string
	for _, n := range leadRes.Notes {
		leadNotes = append(leadNotes, n.Body)
	}

	return Context{
		ClinicName:        clinicRes.Name,
		ClinicType:        clinicRes.ClinicType,
		City:              clinicRes.City,
		CommunicationTone: clinicRes.CommunicationTone,
		ServiceName:       serviceRes.Name,
		ServicePriceFrom:  formatServicePriceFrom(serviceRes.PriceFrom),
		ServiceBenefits:   string(serviceRes.Benefits),
		ServiceFAQ:        string(serviceRes.FAQ),
		CommonObjections:  string(serviceRes.CommonObjections),
		PatientName:       leadRes.FullName,
		LeadNotes:         fmt.Sprintf("[%s]", strings.Join(leadNotes, ", ")),
	}, nil
}

func formatServicePriceFrom(price *float64) string {
	if price == nil {
		return "no informado"
	}
	return strconv.FormatFloat(*price, 'f', -1, 64)
}
