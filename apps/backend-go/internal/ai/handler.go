package ai

import (
	"encoding/json"
	"net/http"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/auth"
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
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "The request could not be completed.")
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
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "The request could not be completed.")
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
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "The request could not be completed.")
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, statusCode int, code string, message string) {
	writeJSON(w, statusCode, ErrorResponse{Error: APIError{Code: code, Message: message}})
}
