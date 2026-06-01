package appointments

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/shared"
)

var (
	ErrMissingClinicID          = errors.New("clinic id is required")
	ErrProfessionalIDRequired   = errors.New("professional_id is required")
	ErrServiceIDRequired        = errors.New("service_id is required")
	ErrContactNameRequired      = errors.New("contact_name is required")
	ErrStartTimeRequired        = errors.New("starts_at is required")
	ErrInvalidAppointmentTime   = errors.New("appointment end time must be after start time")
	ErrInvalidAppointmentStatus = errors.New("invalid appointment status")
	ErrInvalidConfirmation      = errors.New("invalid confirmation status")
	ErrInvalidAppointmentSource = errors.New("invalid appointment source")
)

var allowedStatuses = map[string]bool{
	"scheduled":            true,
	"confirmed":            true,
	"pending_confirmation": true,
	"rescheduled":          true,
	"no_show":              true,
	"cancelled":            true,
	"completed":            true,
	"converted_from_lead":  true,
}

var allowedConfirmationStatuses = map[string]bool{
	"pending":      true,
	"confirmed":    true,
	"not_required": true,
	"failed":       true,
}

var allowedSources = map[string]bool{
	"manual":          true,
	"whatsapp":        true,
	"instagram":       true,
	"facebook":        true,
	"web":             true,
	"llamada":         true,
	"otro":            true,
	"lead_conversion": true,
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context, clinicID string, filters ListFilters) (ListAppointmentsResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return ListAppointmentsResponse{}, ErrMissingClinicID
	}
	filters = normalizeFilters(filters)

	appointments, err := s.repository.List(ctx, clinicID, filters)
	if err != nil {
		return ListAppointmentsResponse{}, err
	}
	data := make([]AppointmentResponse, len(appointments))
	for i, appointment := range appointments {
		data[i] = mapModelToResponse(appointment)
	}
	return ListAppointmentsResponse{Data: data}, nil
}

func (s *Service) Get(ctx context.Context, clinicID, appointmentID string) (AppointmentResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return AppointmentResponse{}, ErrMissingClinicID
	}
	appointment, err := s.repository.FindByID(ctx, clinicID, strings.TrimSpace(appointmentID))
	if err != nil {
		return AppointmentResponse{}, err
	}
	return mapModelToResponse(appointment), nil
}

func (s *Service) Create(ctx context.Context, clinicID, userID string, req CreateAppointmentRequest) (AppointmentResponse, error) {
	appointment, err := buildCreateAppointment(clinicID, userID, req)
	if err != nil {
		return AppointmentResponse{}, err
	}
	created, err := s.repository.Create(ctx, appointment)
	if err != nil {
		return AppointmentResponse{}, err
	}
	return mapModelToResponse(created), nil
}

func (s *Service) Update(ctx context.Context, clinicID, appointmentID string, req UpdateAppointmentRequest) (AppointmentResponse, error) {
	appointment, err := buildUpdateAppointment(clinicID, appointmentID, req)
	if err != nil {
		return AppointmentResponse{}, err
	}
	updated, err := s.repository.Update(ctx, appointment)
	if err != nil {
		return AppointmentResponse{}, err
	}
	return mapModelToResponse(updated), nil
}

func (s *Service) UpdateStatus(ctx context.Context, clinicID, appointmentID string, req UpdateStatusRequest) (AppointmentResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return AppointmentResponse{}, ErrMissingClinicID
	}
	status := strings.TrimSpace(req.Status)
	if !allowedStatuses[status] {
		return AppointmentResponse{}, ErrInvalidAppointmentStatus
	}
	adminNote := normalizeOptionalString(req.AdminNote)
	updated, err := s.repository.UpdateStatus(ctx, clinicID, strings.TrimSpace(appointmentID), status, adminNote)
	if err != nil {
		return AppointmentResponse{}, err
	}
	return mapModelToResponse(updated), nil
}

func (s *Service) Reschedule(ctx context.Context, clinicID, appointmentID string, req RescheduleRequest) (AppointmentResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return AppointmentResponse{}, ErrMissingClinicID
	}
	startsAt, endsAt, err := resolveTimeRange(req.StartsAt, req.EndsAt, req.DurationMins)
	if err != nil {
		return AppointmentResponse{}, err
	}
	adminNote := normalizeOptionalString(req.AdminNote)
	updated, err := s.repository.Reschedule(ctx, clinicID, strings.TrimSpace(appointmentID), startsAt, endsAt, adminNote)
	if err != nil {
		return AppointmentResponse{}, err
	}
	return mapModelToResponse(updated), nil
}

func (s *Service) ConvertLead(ctx context.Context, clinicID, userID, leadID string, req ConvertLeadRequest) (ConvertLeadResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return ConvertLeadResponse{}, ErrMissingClinicID
	}
	startsAt, endsAt, err := resolveTimeRange(req.StartsAt, req.EndsAt, req.DurationMins)
	if err != nil {
		return ConvertLeadResponse{}, err
	}
	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = "pending_confirmation"
	}
	if !allowedStatuses[status] {
		return ConvertLeadResponse{}, ErrInvalidAppointmentStatus
	}
	professionalID := strings.TrimSpace(req.ProfessionalID)
	if professionalID == "" {
		return ConvertLeadResponse{}, ErrProfessionalIDRequired
	}
	serviceID := strings.TrimSpace(req.ServiceID)
	if serviceID == "" {
		return ConvertLeadResponse{}, ErrServiceIDRequired
	}
	userID = strings.TrimSpace(userID)
	var createdBy *string
	if userID != "" {
		createdBy = &userID
	}
	appointment := Appointment{
		ClinicID:           clinicID,
		ProfessionalID:     professionalID,
		ServiceID:          serviceID,
		StartsAt:           startsAt,
		EndsAt:             endsAt,
		Status:             status,
		ConfirmationStatus: "pending",
		Source:             "lead_conversion",
		AdminNotes:         normalizeOptionalString(req.AdminNotes),
		CreatedByUserID:    createdBy,
	}
	return s.repository.ConvertLead(ctx, clinicID, strings.TrimSpace(leadID), appointment, req.UpdateLeadStatus)
}

func buildCreateAppointment(clinicID, userID string, req CreateAppointmentRequest) (Appointment, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return Appointment{}, ErrMissingClinicID
	}
	startsAt, endsAt, err := resolveTimeRange(req.StartsAt, req.EndsAt, req.DurationMins)
	if err != nil {
		return Appointment{}, err
	}
	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = "pending_confirmation"
	}
	if !allowedStatuses[status] {
		return Appointment{}, ErrInvalidAppointmentStatus
	}
	source := strings.TrimSpace(req.Source)
	if source == "" {
		source = "manual"
	}
	if !allowedSources[source] {
		return Appointment{}, ErrInvalidAppointmentSource
	}
	return baseAppointment(clinicID, userID, req.ProfessionalID, req.ServiceID, req.ContactName, req.ContactPhone, req.LeadID, startsAt, endsAt, status, "pending", source, req.AdminNotes)
}

func buildUpdateAppointment(clinicID, appointmentID string, req UpdateAppointmentRequest) (Appointment, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return Appointment{}, ErrMissingClinicID
	}
	startsAt, endsAt, err := resolveTimeRange(req.StartsAt, req.EndsAt, req.DurationMins)
	if err != nil {
		return Appointment{}, err
	}
	status := strings.TrimSpace(req.Status)
	if !allowedStatuses[status] {
		return Appointment{}, ErrInvalidAppointmentStatus
	}
	confirmation := strings.TrimSpace(req.ConfirmationStatus)
	if confirmation == "" {
		confirmation = "pending"
	}
	if !allowedConfirmationStatuses[confirmation] {
		return Appointment{}, ErrInvalidConfirmation
	}
	appointment, err := baseAppointment(clinicID, "", req.ProfessionalID, req.ServiceID, req.ContactName, req.ContactPhone, nil, startsAt, endsAt, status, confirmation, "", req.AdminNotes)
	if err != nil {
		return Appointment{}, err
	}
	appointment.ID = strings.TrimSpace(appointmentID)
	return appointment, nil
}

func baseAppointment(clinicID, userID, professionalID, serviceID, contactName string, contactPhone, leadID *string, startsAt, endsAt time.Time, status, confirmation, source string, adminNotes *string) (Appointment, error) {
	professionalID = strings.TrimSpace(professionalID)
	if professionalID == "" {
		return Appointment{}, ErrProfessionalIDRequired
	}
	serviceID = strings.TrimSpace(serviceID)
	if serviceID == "" {
		return Appointment{}, ErrServiceIDRequired
	}
	contactName = strings.TrimSpace(contactName)
	if contactName == "" {
		return Appointment{}, ErrContactNameRequired
	}
	phone := normalizePhonePointer(contactPhone)
	lead := normalizeOptionalString(leadID)
	userID = strings.TrimSpace(userID)
	var createdBy *string
	if userID != "" {
		createdBy = &userID
	}
	return Appointment{
		ClinicID:           clinicID,
		ProfessionalID:     professionalID,
		LeadID:             lead,
		ServiceID:          serviceID,
		ContactName:        contactName,
		ContactPhone:       phone,
		StartsAt:           startsAt,
		EndsAt:             endsAt,
		Status:             status,
		ConfirmationStatus: confirmation,
		Source:             source,
		AdminNotes:         normalizeOptionalString(adminNotes),
		CreatedByUserID:    createdBy,
	}, nil
}

func resolveTimeRange(startsAt time.Time, explicitEnd *time.Time, durationMins *int) (time.Time, time.Time, error) {
	if startsAt.IsZero() {
		return time.Time{}, time.Time{}, ErrStartTimeRequired
	}
	var endsAt time.Time
	if explicitEnd != nil {
		endsAt = *explicitEnd
	} else if durationMins != nil && *durationMins > 0 {
		endsAt = startsAt.Add(time.Duration(*durationMins) * time.Minute)
	} else {
		endsAt = startsAt.Add(time.Hour)
	}
	if !endsAt.After(startsAt) {
		return time.Time{}, time.Time{}, ErrInvalidAppointmentTime
	}
	return startsAt, endsAt, nil
}

func normalizeFilters(filters ListFilters) ListFilters {
	filters.Date = strings.TrimSpace(filters.Date)
	filters.DateFrom = strings.TrimSpace(filters.DateFrom)
	filters.DateTo = strings.TrimSpace(filters.DateTo)
	filters.ProfessionalID = strings.TrimSpace(filters.ProfessionalID)
	filters.ServiceID = strings.TrimSpace(filters.ServiceID)
	filters.Status = strings.TrimSpace(filters.Status)
	filters.LeadID = strings.TrimSpace(filters.LeadID)
	return filters
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func normalizePhonePointer(value *string) *string {
	if value == nil {
		return nil
	}
	phone := shared.NormalizePhone(*value)
	if phone == "" {
		return nil
	}
	return &phone
}

func mapModelToResponse(a Appointment) AppointmentResponse {
	var lead *LeadSummary
	if a.LeadID != nil && a.LeadName != nil {
		lead = &LeadSummary{ID: *a.LeadID, FullName: *a.LeadName}
	}
	return AppointmentResponse{
		ID:       a.ID,
		ClinicID: a.ClinicID,
		Professional: ProfessionalSummary{
			ID:            a.ProfessionalID,
			FullName:      a.ProfessionalName,
			CalendarColor: a.ProfessionalColor,
		},
		Lead: lead,
		Service: ServiceSummary{
			ID:   a.ServiceID,
			Name: a.ServiceName,
		},
		ContactName:        a.ContactName,
		ContactPhone:       a.ContactPhone,
		StartsAt:           a.StartsAt,
		EndsAt:             a.EndsAt,
		Status:             a.Status,
		Source:             a.Source,
		ConfirmationStatus: a.ConfirmationStatus,
		AdminNotes:         a.AdminNotes,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}
}
