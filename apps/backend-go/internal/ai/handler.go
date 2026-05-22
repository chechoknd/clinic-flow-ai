package ai

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/auth"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/clinics"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/leads"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/services"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ReplySuggestion(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required.")
		return
	}

	var req ReplySuggestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "The request body is invalid.")
		return
	}

	res, err := h.service.ReplySuggestion(r.Context(), claims.ClinicID, req)
	if err != nil {
		handleAIError(w, "reply_suggestion", err)
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) FollowUpMessage(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required.")
		return
	}

	var req FollowUpMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "The request body is invalid.")
		return
	}

	res, err := h.service.FollowUpMessage(r.Context(), claims.ClinicID, req)
	if err != nil {
		handleAIError(w, "follow_up_message", err)
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) ObjectionHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required.")
		return
	}

	var req ObjectionHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "The request body is invalid.")
		return
	}

	res, err := h.service.ObjectionHandler(r.Context(), claims.ClinicID, req)
	if err != nil {
		handleAIError(w, "objection_handler", err)
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func handleAIError(w http.ResponseWriter, operation string, err error) {
	switch {
	case errors.Is(err, ErrValidation):
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	case errors.Is(err, clinics.ErrClinicNotFound), errors.Is(err, services.ErrServiceNotFound), errors.Is(err, leads.ErrLeadNotFound):
		writeError(w, http.StatusNotFound, "NOT_FOUND", "The requested AI context was not found.")
	case errors.Is(err, ErrAISafetyBlocked):
		log.Printf("ai safety blocked during %s: %v", operation, err)
		writeError(w, http.StatusUnprocessableEntity, "AI_SAFETY_BLOCKED", "The AI response did not pass safety validation.")
	case errors.Is(err, ErrAIProvider):
		log.Printf("ai provider error during %s: %v", operation, err)
		writeError(w, http.StatusBadGateway, "AI_PROVIDER_ERROR", "The AI provider could not complete the request.")
	case errors.Is(err, ErrAIResponse):
		log.Printf("ai response parse error during %s: %v", operation, err)
		writeError(w, http.StatusBadGateway, "AI_PROVIDER_ERROR", "The AI provider returned an invalid response.")
	default:
		log.Printf("ai internal error during %s: %v", operation, err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "The request could not be completed.")
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, statusCode int, code string, message string) {
	writeJSON(w, statusCode, ErrorResponse{Error: APIError{Code: code, Message: message}})
}
