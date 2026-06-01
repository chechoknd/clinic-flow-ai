package professionals

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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

	filters, err := parseListFilters(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	res, err := h.service.List(r.Context(), claims.ClinicID, filters)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "The request could not be completed.")
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

	professionalID := r.PathValue("id")
	if professionalID == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Professional ID is required.")
		return
	}

	res, err := h.service.Get(r.Context(), claims.ClinicID, professionalID)
	if errors.Is(err, ErrProfessionalNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Professional was not found.")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "The request could not be completed.")
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

	var req CreateProfessionalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "The request body is invalid.")
		return
	}

	res, err := h.service.Create(r.Context(), claims.ClinicID, req)
	if err != nil {
		writeProfessionalError(w, err)
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

	professionalID := r.PathValue("id")
	if professionalID == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Professional ID is required.")
		return
	}

	var req UpdateProfessionalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "The request body is invalid.")
		return
	}

	res, err := h.service.Update(r.Context(), claims.ClinicID, professionalID, req)
	if err != nil {
		writeProfessionalError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func parseListFilters(r *http.Request) (ListFilters, error) {
	query := r.URL.Query()
	var filters ListFilters
	if value := query.Get("is_active"); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return ListFilters{}, errors.New("is_active must be true or false")
		}
		filters.IsActive = &parsed
	}
	filters.ServiceID = query.Get("service_id")
	return filters, nil
}

func writeProfessionalError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrProfessionalNotFound):
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Professional was not found.")
	case errors.Is(err, ErrProfessionalNameRequired), errors.Is(err, ErrInvalidCalendarColor), errors.Is(err, ErrInvalidWorkingHours):
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
	case errors.Is(err, ErrDuplicateProfessionalName):
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
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
