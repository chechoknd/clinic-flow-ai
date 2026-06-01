package appointments

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

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required.")
		return
	}
	res, err := h.service.List(r.Context(), claims.ClinicID, parseFilters(r))
	if err != nil {
		writeAppointmentError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required.")
		return
	}
	res, err := h.service.Get(r.Context(), claims.ClinicID, r.PathValue("id"))
	if err != nil {
		writeAppointmentError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required.")
		return
	}
	var req CreateAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "The request body is invalid.")
		return
	}
	res, err := h.service.Create(r.Context(), claims.ClinicID, claims.UserID, req)
	if err != nil {
		writeAppointmentError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required.")
		return
	}
	var req UpdateAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "The request body is invalid.")
		return
	}
	res, err := h.service.Update(r.Context(), claims.ClinicID, r.PathValue("id"), req)
	if err != nil {
		writeAppointmentError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required.")
		return
	}
	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "The request body is invalid.")
		return
	}
	res, err := h.service.UpdateStatus(r.Context(), claims.ClinicID, r.PathValue("id"), req)
	if err != nil {
		writeAppointmentError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) Reschedule(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required.")
		return
	}
	var req RescheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "The request body is invalid.")
		return
	}
	res, err := h.service.Reschedule(r.Context(), claims.ClinicID, r.PathValue("id"), req)
	if err != nil {
		writeAppointmentError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) ConvertLead(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required.")
		return
	}
	var req ConvertLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "The request body is invalid.")
		return
	}
	res, err := h.service.ConvertLead(r.Context(), claims.ClinicID, claims.UserID, r.PathValue("id"), req)
	if err != nil {
		writeAppointmentError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

func parseFilters(r *http.Request) ListFilters {
	query := r.URL.Query()
	return ListFilters{
		Date:           query.Get("date"),
		DateFrom:       query.Get("date_from"),
		DateTo:         query.Get("date_to"),
		ProfessionalID: query.Get("professional_id"),
		ServiceID:      query.Get("service_id"),
		Status:         query.Get("status"),
		LeadID:         query.Get("lead_id"),
	}
}

func writeAppointmentError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrAppointmentNotFound), errors.Is(err, ErrLeadNotFound):
		writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
	case errors.Is(err, ErrScheduleConflict):
		writeError(w, http.StatusConflict, "VALIDATION_ERROR", err.Error())
	case errors.Is(err, ErrInvalidService), errors.Is(err, ErrMissingClinicID), errors.Is(err, ErrProfessionalIDRequired), errors.Is(err, ErrServiceIDRequired), errors.Is(err, ErrContactNameRequired), errors.Is(err, ErrStartTimeRequired), errors.Is(err, ErrInvalidAppointmentTime), errors.Is(err, ErrInvalidAppointmentStatus), errors.Is(err, ErrInvalidConfirmation), errors.Is(err, ErrInvalidAppointmentSource):
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
	default:
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
