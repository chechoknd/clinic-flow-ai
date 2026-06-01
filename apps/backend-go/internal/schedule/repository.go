package schedule

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrProfessionalNotFound = errors.New("professional not found")
)

type Repository interface {
	FindProfessional(ctx context.Context, clinicID, professionalID string) (Professional, error)
	GetServiceDuration(ctx context.Context, clinicID, serviceID string) (int, error)
	FindActiveAppointments(ctx context.Context, clinicID, professionalID string, fromTime, toTime time.Time) ([]AppointmentRange, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) FindProfessional(ctx context.Context, clinicID, professionalID string) (Professional, error) {
	const query = `
		SELECT id::text, clinic_id::text, full_name, working_hours, is_active
		FROM clinic_professionals
		WHERE clinic_id = $1 AND id = $2 AND is_active = true
		LIMIT 1
	`
	var p Professional
	err := r.db.QueryRowContext(ctx, query, clinicID, professionalID).Scan(
		&p.ID,
		&p.ClinicID,
		&p.FullName,
		&p.WorkingHours,
		&p.IsActive,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Professional{}, ErrProfessionalNotFound
	}
	if err != nil {
		return Professional{}, err
	}
	return p, nil
}

func (r *PostgresRepository) GetServiceDuration(ctx context.Context, clinicID, serviceID string) (int, error) {
	const query = `
		SELECT duration_minutes
		FROM clinic_services
		WHERE clinic_id = $1 AND id = $2 AND is_active = true
		LIMIT 1
	`
	var duration int
	err := r.db.QueryRowContext(ctx, query, clinicID, serviceID).Scan(&duration)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil // service not found or inactive, default to 0
	}
	if err != nil {
		return 0, err
	}
	return duration, nil
}

func (r *PostgresRepository) FindActiveAppointments(ctx context.Context, clinicID, professionalID string, fromTime, toTime time.Time) ([]AppointmentRange, error) {
	const query = `
		SELECT starts_at, ends_at
		FROM appointments
		WHERE clinic_id = $1
			AND professional_id = $2
			AND status IN ('scheduled', 'confirmed', 'pending_confirmation', 'rescheduled', 'converted_from_lead')
			AND ends_at > $3
			AND starts_at < $4
		ORDER BY starts_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, clinicID, professionalID, fromTime, toTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []AppointmentRange
	for rows.Next() {
		var a AppointmentRange
		if err := rows.Scan(&a.StartsAt, &a.EndsAt); err != nil {
			return nil, err
		}
		appointments = append(appointments, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return appointments, nil
}
