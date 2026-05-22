package leads

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

func (h *Handler) ListFollowUps(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required.")
		return
	}

	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("page_size"))

	res, err := h.service.ListFollowUps(r.Context(), claims.ClinicID, FollowUpFilter{
		Due:      query.Get("due"),
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		if err.Error() == "invalid follow-up due filter" {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "The request could not be completed.")
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) CompleteFollowUp(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required.")
		return
	}

	leadID := r.PathValue("id")
	if leadID == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Lead ID is required.")
		return
	}

	var req CompleteFollowUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "The request body is invalid.")
		return
	}

	err := h.service.CompleteFollowUp(r.Context(), claims.ClinicID, leadID, req)
	if errors.Is(err, ErrLeadNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Lead was not found.")
		return
	}
	if err != nil {
		if err.Error() == "invalid status" {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "The request could not be completed.")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": leadID, "status": "completed"})
}

func (h *Handler) RescheduleFollowUp(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required.")
		return
	}

	leadID := r.PathValue("id")
	if leadID == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Lead ID is required.")
		return
	}

	var req RescheduleFollowUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "The request body is invalid.")
		return
	}

	err := h.service.RescheduleFollowUp(r.Context(), claims.ClinicID, leadID, req)
	if errors.Is(err, ErrLeadNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Lead was not found.")
		return
	}
	if err != nil {
		if err.Error() == "next_action_at is required" {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "The request could not be completed.")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": leadID, "status": "rescheduled"})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required.")
		return
	}

	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("page_size"))

	filter := ListFilter{
		Status:    query.Get("status"),
		ServiceID: query.Get("service_id"),
		Page:      page,
		PageSize:  pageSize,
	}

	res, err := h.service.List(r.Context(), claims.ClinicID, filter)
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

	leadID := r.PathValue("id")
	if leadID == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Lead ID is required.")
		return
	}

	res, err := h.service.Get(r.Context(), claims.ClinicID, leadID)
	if errors.Is(err, ErrLeadNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Lead was not found.")
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

	var req CreateLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "The request body is invalid.")
		return
	}

	res, err := h.service.Create(r.Context(), claims.ClinicID, req)
	if err != nil {
		if err.Error() == "full name is required" || err.Error() == "phone is required" {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "The request could not be completed.")
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

	leadID := r.PathValue("id")
	if leadID == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Lead ID is required.")
		return
	}

	var req UpdateLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "The request body is invalid.")
		return
	}

	err := h.service.Update(r.Context(), claims.ClinicID, leadID, req)
	if errors.Is(err, ErrLeadNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Lead was not found.")
		return
	}
	if err != nil {
		if err.Error() == "status is required" {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "The request could not be completed.")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": leadID, "status": req.Status})
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, statusCode int, code string, message string) {
	writeJSON(w, statusCode, ErrorResponse{Error: APIError{Code: code, Message: message}})
}
