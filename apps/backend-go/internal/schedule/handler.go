package schedule

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/auth"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAvailability(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required.")
		return
	}

	query := r.URL.Query()
	professionalID := query.Get("professional_id")
	if professionalID == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "professional_id is required.")
		return
	}

	dateFromStr := query.Get("date_from")
	dateToStr := query.Get("date_to")

	dateFrom, err := parseDate(dateFromStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", fmt.Sprintf("invalid date_from: %v", err))
		return
	}

	dateTo, err := parseDate(dateToStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", fmt.Sprintf("invalid date_to: %v", err))
		return
	}

	serviceID := query.Get("service_id")

	durationMins := 0
	if val := query.Get("duration_minutes"); val != "" {
		parsed, err := strconv.Atoi(val)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "duration_minutes must be an integer.")
			return
		}
		durationMins = parsed
	}

	req := AvailabilityRequest{
		ProfessionalID:  professionalID,
		ServiceID:       serviceID,
		DateFrom:        dateFrom,
		DateTo:          dateTo,
		DurationMinutes: durationMins,
	}

	res, err := h.service.GetAvailability(r.Context(), claims.ClinicID, req)
	if errors.Is(err, ErrProfessionalNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Professional not found or is inactive.")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", fmt.Sprintf("Could not calculate availability: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func parseDate(val string) (time.Time, error) {
	val = strings.TrimSpace(val)
	if val == "" {
		return time.Time{}, errors.New("date value is empty")
	}
	// Try ISO 8601 / RFC 3339 format (Javascript toISOString)
	if t, err := time.Parse(time.RFC3339, val); err == nil {
		return t, nil
	}
	// Try standard date YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", val); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("use YYYY-MM-DD or RFC3339 format")
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, statusCode int, code string, message string) {
	writeJSON(w, statusCode, ErrorResponse{Error: APIError{Code: code, Message: message}})
}
