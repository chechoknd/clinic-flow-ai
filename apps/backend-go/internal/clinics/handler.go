package clinics

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/auth"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Current(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required.")
		return
	}

	res, err := h.service.Current(r.Context(), claims.ClinicID)
	if errors.Is(err, ErrMissingClinicID) {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required.")
		return
	}
	if errors.Is(err, ErrClinicNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Clinic was not found.")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "The request could not be completed.")
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required.")
		return
	}

	var req UpdateClinicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "The request body is invalid.")
		return
	}

	res, err := h.service.Update(r.Context(), claims.ClinicID, req)
	if errors.Is(err, ErrMissingClinicID) {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required.")
		return
	}
	if errors.Is(err, ErrClinicNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Clinic was not found.")
		return
	}
	if err != nil {
		// Basic validation error mapping (could be improved with specific error types)
		if err.Error() == "clinic name is required" || err.Error() == "city is required" || err.Error() == "whatsapp is required" {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}

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
