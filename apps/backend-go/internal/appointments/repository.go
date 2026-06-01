package appointments

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrAppointmentNotFound = errors.New("appointment not found")
	ErrScheduleConflict    = errors.New("professional has another active appointment in this time range")
	ErrLeadNotFound        = errors.New("lead not found")
	ErrInvalidService      = errors.New("professional cannot provide selected service")
)

type Repository interface {
	List(ctx context.Context, clinicID string, filters ListFilters) ([]Appointment, error)
	FindByID(ctx context.Context, clinicID, appointmentID string) (Appointment, error)
	Create(ctx context.Context, appointment Appointment) (Appointment, error)
	Update(ctx context.Context, appointment Appointment) (Appointment, error)
	UpdateStatus(ctx context.Context, clinicID, appointmentID, status string, adminNote *string) (Appointment, error)
	Reschedule(ctx context.Context, clinicID, appointmentID string, startsAt, endsAt time.Time, adminNote *string) (Appointment, error)
	ConvertLead(ctx context.Context, clinicID, leadID string, appointment Appointment, updateLeadStatus bool) (ConvertLeadResponse, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) List(ctx context.Context, clinicID string, filters ListFilters) ([]Appointment, error) {
	where := "WHERE a.clinic_id = $1"
	args := []any{clinicID}
	argIdx := 2

	if filters.Date != "" {
		where += fmt.Sprintf(" AND a.starts_at >= $%d::date AND a.starts_at < $%d::date + interval '1 day'", argIdx, argIdx)
		args = append(args, filters.Date)
		argIdx++
	}
	if filters.DateFrom != "" {
		where += fmt.Sprintf(" AND a.starts_at >= $%d::date", argIdx)
		args = append(args, filters.DateFrom)
		argIdx++
	}
	if filters.DateTo != "" {
		where += fmt.Sprintf(" AND a.starts_at < $%d::date + interval '1 day'", argIdx)
		args = append(args, filters.DateTo)
		argIdx++
	}
	if filters.ProfessionalID != "" {
		where += fmt.Sprintf(" AND a.professional_id = $%d", argIdx)
		args = append(args, filters.ProfessionalID)
		argIdx++
	}
	if filters.ServiceID != "" {
		where += fmt.Sprintf(" AND a.service_id = $%d", argIdx)
		args = append(args, filters.ServiceID)
		argIdx++
	}
	if filters.Status != "" {
		where += fmt.Sprintf(" AND a.status = $%d", argIdx)
		args = append(args, filters.Status)
		argIdx++
	}
	if filters.LeadID != "" {
		where += fmt.Sprintf(" AND a.lead_id = $%d", argIdx)
		args = append(args, filters.LeadID)
	}

	query := appointmentSelectQuery(where + " ORDER BY a.starts_at ASC")
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []Appointment
	for rows.Next() {
		a, err := scanAppointment(rows)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return appointments, nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, clinicID, appointmentID string) (Appointment, error) {
	query := appointmentSelectQuery("WHERE a.id = $1 AND a.clinic_id = $2 LIMIT 1")
	a, err := scanAppointment(r.db.QueryRowContext(ctx, query, appointmentID, clinicID))
	if errors.Is(err, sql.ErrNoRows) {
		return Appointment{}, ErrAppointmentNotFound
	}
	if err != nil {
		return Appointment{}, err
	}
	return a, nil
}

func (r *PostgresRepository) Create(ctx context.Context, a Appointment) (Appointment, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Appointment{}, err
	}
	defer tx.Rollback()

	if err := validateLinkedRows(ctx, tx, a); err != nil {
		return Appointment{}, err
	}
	if err := ensureNoOverlap(ctx, tx, a.ClinicID, a.ProfessionalID, "", a.StartsAt, a.EndsAt); err != nil {
		return Appointment{}, err
	}

	const query = `
		INSERT INTO appointments (
			clinic_id,
			professional_id,
			lead_id,
			service_id,
			contact_name,
			contact_phone,
			starts_at,
			ends_at,
			status,
			confirmation_status,
			source,
			admin_notes,
			created_by_user_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id::text
	`
	err = tx.QueryRowContext(
		ctx,
		query,
		a.ClinicID,
		a.ProfessionalID,
		a.LeadID,
		a.ServiceID,
		a.ContactName,
		a.ContactPhone,
		a.StartsAt,
		a.EndsAt,
		a.Status,
		a.ConfirmationStatus,
		a.Source,
		a.AdminNotes,
		a.CreatedByUserID,
	).Scan(&a.ID)
	if err != nil {
		return Appointment{}, err
	}
	if err := tx.Commit(); err != nil {
		return Appointment{}, err
	}
	return r.FindByID(ctx, a.ClinicID, a.ID)
}

func (r *PostgresRepository) Update(ctx context.Context, a Appointment) (Appointment, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Appointment{}, err
	}
	defer tx.Rollback()

	if err := appointmentExists(ctx, tx, a.ClinicID, a.ID); err != nil {
		return Appointment{}, err
	}
	if err := validateLinkedRows(ctx, tx, a); err != nil {
		return Appointment{}, err
	}
	if err := ensureNoOverlap(ctx, tx, a.ClinicID, a.ProfessionalID, a.ID, a.StartsAt, a.EndsAt); err != nil {
		return Appointment{}, err
	}

	res, err := tx.ExecContext(ctx, `
		UPDATE appointments
		SET
			professional_id = $3,
			service_id = $4,
			contact_name = $5,
			contact_phone = $6,
			starts_at = $7,
			ends_at = $8,
			status = $9,
			confirmation_status = $10,
			admin_notes = $11
		WHERE id = $1 AND clinic_id = $2
	`, a.ID, a.ClinicID, a.ProfessionalID, a.ServiceID, a.ContactName, a.ContactPhone, a.StartsAt, a.EndsAt, a.Status, a.ConfirmationStatus, a.AdminNotes)
	if err != nil {
		return Appointment{}, err
	}
	if rows, err := res.RowsAffected(); err != nil {
		return Appointment{}, err
	} else if rows == 0 {
		return Appointment{}, ErrAppointmentNotFound
	}
	if err := tx.Commit(); err != nil {
		return Appointment{}, err
	}
	return r.FindByID(ctx, a.ClinicID, a.ID)
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, clinicID, appointmentID, status string, adminNote *string) (Appointment, error) {
	res, err := r.db.ExecContext(ctx, `
		UPDATE appointments
		SET
			status = $3::varchar,
			confirmation_status = CASE WHEN $3::varchar = 'confirmed' THEN 'confirmed' ELSE confirmation_status END,
			admin_notes = COALESCE($4::text, admin_notes)
		WHERE id = $1 AND clinic_id = $2
	`, appointmentID, clinicID, status, nullableString(adminNote))
	if err != nil {
		return Appointment{}, err
	}
	if rows, err := res.RowsAffected(); err != nil {
		return Appointment{}, err
	} else if rows == 0 {
		return Appointment{}, ErrAppointmentNotFound
	}
	return r.FindByID(ctx, clinicID, appointmentID)
}

func (r *PostgresRepository) Reschedule(ctx context.Context, clinicID, appointmentID string, startsAt, endsAt time.Time, adminNote *string) (Appointment, error) {
	current, err := r.FindByID(ctx, clinicID, appointmentID)
	if err != nil {
		return Appointment{}, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Appointment{}, err
	}
	defer tx.Rollback()

	if err := ensureNoOverlap(ctx, tx, clinicID, current.ProfessionalID, appointmentID, startsAt, endsAt); err != nil {
		return Appointment{}, err
	}

	res, err := tx.ExecContext(ctx, `
		UPDATE appointments
		SET
			starts_at = $3,
			ends_at = $4,
			status = 'rescheduled',
			confirmation_status = 'pending',
			admin_notes = COALESCE($5::text, admin_notes)
		WHERE id = $1 AND clinic_id = $2
	`, appointmentID, clinicID, startsAt, endsAt, nullableString(adminNote))
	if err != nil {
		return Appointment{}, err
	}
	if rows, err := res.RowsAffected(); err != nil {
		return Appointment{}, err
	} else if rows == 0 {
		return Appointment{}, ErrAppointmentNotFound
	}
	if err := tx.Commit(); err != nil {
		return Appointment{}, err
	}
	return r.FindByID(ctx, clinicID, appointmentID)
}

func (r *PostgresRepository) ConvertLead(ctx context.Context, clinicID, leadID string, a Appointment, updateLeadStatus bool) (ConvertLeadResponse, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return ConvertLeadResponse{}, err
	}
	defer tx.Rollback()

	lead, err := findLeadForConversion(ctx, tx, clinicID, leadID)
	if err != nil {
		return ConvertLeadResponse{}, err
	}
	a.ClinicID = clinicID
	a.LeadID = &lead.ID
	if strings.TrimSpace(a.ContactName) == "" {
		a.ContactName = lead.FullName
	}
	if a.ContactPhone == nil || strings.TrimSpace(*a.ContactPhone) == "" {
		a.ContactPhone = &lead.Phone
	}

	if err := validateLinkedRows(ctx, tx, a); err != nil {
		return ConvertLeadResponse{}, err
	}
	if err := ensureNoOverlap(ctx, tx, a.ClinicID, a.ProfessionalID, "", a.StartsAt, a.EndsAt); err != nil {
		return ConvertLeadResponse{}, err
	}

	var appointmentID string
	err = tx.QueryRowContext(ctx, `
		INSERT INTO appointments (
			clinic_id,
			professional_id,
			lead_id,
			service_id,
			contact_name,
			contact_phone,
			starts_at,
			ends_at,
			status,
			confirmation_status,
			source,
			admin_notes,
			created_by_user_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'pending', 'lead_conversion', $10, $11)
		RETURNING id::text
	`, a.ClinicID, a.ProfessionalID, a.LeadID, a.ServiceID, a.ContactName, a.ContactPhone, a.StartsAt, a.EndsAt, a.Status, a.AdminNotes, a.CreatedByUserID).Scan(&appointmentID)
	if err != nil {
		return ConvertLeadResponse{}, err
	}

	leadStatus := lead.Status
	if updateLeadStatus {
		leadStatus = "Agendado"
		if _, err := tx.ExecContext(ctx, `
			UPDATE leads
			SET status = $3
			WHERE id = $1 AND clinic_id = $2
		`, leadID, clinicID, leadStatus); err != nil {
			return ConvertLeadResponse{}, err
		}
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO lead_notes (lead_id, body)
		VALUES ($1, $2)
	`, leadID, fmt.Sprintf("Cita creada para %s.", a.StartsAt.Format("2006-01-02 15:04"))); err != nil {
		return ConvertLeadResponse{}, err
	}

	if err := tx.Commit(); err != nil {
		return ConvertLeadResponse{}, err
	}

	return ConvertLeadResponse{
		LeadID:            leadID,
		LeadStatus:        leadStatus,
		AppointmentID:     appointmentID,
		AppointmentStatus: a.Status,
	}, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func appointmentSelectQuery(where string) string {
	return `
		SELECT
			a.id::text,
			a.clinic_id::text,
			a.professional_id::text,
			p.full_name,
			p.calendar_color,
			a.lead_id::text,
			l.full_name,
			a.service_id::text,
			s.name,
			a.contact_name,
			a.contact_phone,
			a.starts_at,
			a.ends_at,
			a.status,
			a.confirmation_status,
			a.source,
			a.admin_notes,
			a.created_by_user_id::text,
			a.created_at,
			a.updated_at
		FROM appointments a
		INNER JOIN clinic_professionals p ON p.clinic_id = a.clinic_id AND p.id = a.professional_id
		INNER JOIN clinic_services s ON s.clinic_id = a.clinic_id AND s.id = a.service_id
		LEFT JOIN leads l ON l.clinic_id = a.clinic_id AND l.id = a.lead_id
		` + where
}

func scanAppointment(row rowScanner) (Appointment, error) {
	var a Appointment
	err := row.Scan(
		&a.ID,
		&a.ClinicID,
		&a.ProfessionalID,
		&a.ProfessionalName,
		&a.ProfessionalColor,
		&a.LeadID,
		&a.LeadName,
		&a.ServiceID,
		&a.ServiceName,
		&a.ContactName,
		&a.ContactPhone,
		&a.StartsAt,
		&a.EndsAt,
		&a.Status,
		&a.ConfirmationStatus,
		&a.Source,
		&a.AdminNotes,
		&a.CreatedByUserID,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		return Appointment{}, err
	}
	return a, nil
}

func validateLinkedRows(ctx context.Context, tx *sql.Tx, a Appointment) error {
	if err := validateProfessionalService(ctx, tx, a.ClinicID, a.ProfessionalID, a.ServiceID); err != nil {
		return err
	}
	if a.LeadID != nil && strings.TrimSpace(*a.LeadID) != "" {
		var exists bool
		if err := tx.QueryRowContext(ctx, `
			SELECT EXISTS(
				SELECT 1 FROM leads
				WHERE id = $1 AND clinic_id = $2
			)
		`, *a.LeadID, a.ClinicID).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return ErrLeadNotFound
		}
	}
	if a.CreatedByUserID != nil && strings.TrimSpace(*a.CreatedByUserID) != "" {
		var exists bool
		if err := tx.QueryRowContext(ctx, `
			SELECT EXISTS(
				SELECT 1 FROM users
				WHERE id = $1 AND clinic_id = $2
			)
		`, *a.CreatedByUserID, a.ClinicID).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return errors.New("user does not belong to clinic")
		}
	}
	return nil
}

func validateProfessionalService(ctx context.Context, tx *sql.Tx, clinicID, professionalID, serviceID string) error {
	var valid bool
	if err := tx.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM clinic_professionals p
			INNER JOIN clinic_services s ON s.clinic_id = p.clinic_id
			WHERE p.clinic_id = $1
				AND p.id = $2
				AND p.is_active = TRUE
				AND s.id = $3
				AND s.is_active = TRUE
				AND (
					NOT EXISTS (
						SELECT 1 FROM professional_services ps_any
						WHERE ps_any.clinic_id = p.clinic_id
							AND ps_any.professional_id = p.id
					)
					OR EXISTS (
						SELECT 1 FROM professional_services ps
						WHERE ps.clinic_id = p.clinic_id
							AND ps.professional_id = p.id
							AND ps.service_id = s.id
					)
				)
		)
	`, clinicID, professionalID, serviceID).Scan(&valid); err != nil {
		return err
	}
	if !valid {
		return ErrInvalidService
	}
	return nil
}

func ensureNoOverlap(ctx context.Context, tx *sql.Tx, clinicID, professionalID, excludeAppointmentID string, startsAt, endsAt time.Time) error {
	var exists bool
	args := []any{clinicID, professionalID, startsAt, endsAt}
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM appointments
			WHERE clinic_id = $1
				AND professional_id = $2
				AND status IN ('scheduled', 'confirmed', 'pending_confirmation', 'rescheduled', 'converted_from_lead')
				AND starts_at < $4
				AND ends_at > $3
	`
	if excludeAppointmentID != "" {
		args = append(args, excludeAppointmentID)
		query += ` AND id <> $5`
	}
	query += `)`
	if err := tx.QueryRowContext(ctx, query, args...).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return ErrScheduleConflict
	}
	return nil
}

func appointmentExists(ctx context.Context, tx *sql.Tx, clinicID, appointmentID string) error {
	var exists bool
	if err := tx.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM appointments
			WHERE id = $1 AND clinic_id = $2
		)
	`, appointmentID, clinicID).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return ErrAppointmentNotFound
	}
	return nil
}

type leadForConversion struct {
	ID       string
	FullName string
	Phone    string
	Status   string
}

func findLeadForConversion(ctx context.Context, tx *sql.Tx, clinicID, leadID string) (leadForConversion, error) {
	var lead leadForConversion
	err := tx.QueryRowContext(ctx, `
		SELECT id::text, full_name, phone, status
		FROM leads
		WHERE id = $1 AND clinic_id = $2
	`, leadID, clinicID).Scan(&lead.ID, &lead.FullName, &lead.Phone, &lead.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return leadForConversion{}, ErrLeadNotFound
	}
	if err != nil {
		return leadForConversion{}, err
	}
	return lead, nil
}

func nullableString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *value, Valid: true}
}
