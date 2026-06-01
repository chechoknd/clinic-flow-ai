package dashboard

import (
	"context"
	"database/sql"
)

type Repository interface {
	Summary(ctx context.Context, clinicID string) (Summary, error)
	ScheduleSummary(ctx context.Context, clinicID string, dateStr string) (ScheduleSummary, error)
	GetActiveProfessionalsWithHours(ctx context.Context, clinicID string) ([]ProfessionalHours, error)
	GetAppointmentsForDate(ctx context.Context, clinicID string, dateStr string) ([]AppointmentSummary, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Summary(ctx context.Context, clinicID string) (Summary, error) {
	summary := Summary{LeadsByStatus: map[string]int{}}

	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM leads
		WHERE clinic_id = $1
	`, clinicID).Scan(&summary.LeadsTotal); err != nil {
		return Summary{}, err
	}

	statusRows, err := r.db.QueryContext(ctx, `
		SELECT status, COUNT(*)
		FROM leads
		WHERE clinic_id = $1
		GROUP BY status
		ORDER BY status ASC
	`, clinicID)
	if err != nil {
		return Summary{}, err
	}
	defer statusRows.Close()

	for statusRows.Next() {
		var status string
		var count int
		if err := statusRows.Scan(&status, &count); err != nil {
			return Summary{}, err
		}
		summary.LeadsByStatus[status] = count
	}
	if err := statusRows.Err(); err != nil {
		return Summary{}, err
	}

	topRows, err := r.db.QueryContext(ctx, `
		SELECT
			s.id::text,
			s.name,
			COUNT(l.id)
		FROM leads l
		JOIN clinic_services s ON s.id = l.service_id AND s.clinic_id = l.clinic_id
		WHERE l.clinic_id = $1 AND l.service_id IS NOT NULL
		GROUP BY s.id, s.name
		ORDER BY COUNT(l.id) DESC, s.name ASC
		LIMIT 5
	`, clinicID)
	if err != nil {
		return Summary{}, err
	}
	defer topRows.Close()

	for topRows.Next() {
		var service TopService
		if err := topRows.Scan(&service.ServiceID, &service.ServiceName, &service.LeadCount); err != nil {
			return Summary{}, err
		}
		summary.TopServices = append(summary.TopServices, service)
	}
	if err := topRows.Err(); err != nil {
		return Summary{}, err
	}

	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM leads
		WHERE clinic_id = $1
			AND next_action_at IS NOT NULL
			AND next_action_at >= date_trunc('day', NOW())
			AND next_action_at < date_trunc('day', NOW()) + interval '1 day'
	`, clinicID).Scan(&summary.PendingFollowUpsToday); err != nil {
		return Summary{}, err
	}

	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM leads
		WHERE clinic_id = $1
			AND next_action_at IS NOT NULL
			AND next_action_at < NOW()
	`, clinicID).Scan(&summary.OverdueFollowUps); err != nil {
		return Summary{}, err
	}

	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM leads
		WHERE clinic_id = $1
			AND next_action_at IS NOT NULL
			AND next_action_at > NOW()
	`, clinicID).Scan(&summary.UpcomingFollowUps); err != nil {
		return Summary{}, err
	}

	converted := summary.LeadsByStatus["Convertido"]
	if summary.LeadsTotal > 0 {
		summary.ConversionRate = float64(converted) / float64(summary.LeadsTotal)
	}

	return summary, nil
}

func (r *PostgresRepository) ScheduleSummary(ctx context.Context, clinicID string, dateStr string) (ScheduleSummary, error) {
	summary := ScheduleSummary{
		Date: dateStr,
	}

	// 1. Todays appointments
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM appointments
		WHERE clinic_id = $1
			AND starts_at >= $2::timestamptz
			AND starts_at < $2::timestamptz + interval '1 day'
			AND status IN ('scheduled', 'confirmed', 'pending_confirmation', 'rescheduled', 'converted_from_lead')
	`, clinicID, dateStr).Scan(&summary.TodaysAppointments)
	if err != nil {
		return ScheduleSummary{}, err
	}

	// 2. Pending confirmation
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM appointments
		WHERE clinic_id = $1
			AND starts_at >= $2::timestamptz
			AND starts_at < $2::timestamptz + interval '1 day'
			AND status = 'pending_confirmation'
	`, clinicID, dateStr).Scan(&summary.AppointmentsPendingConfirmation)
	if err != nil {
		return ScheduleSummary{}, err
	}

	// 3. Hot leads without appointment
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM leads l
		WHERE l.clinic_id = $1
			AND l.status = 'Interesado'
			AND l.id NOT IN (
				SELECT DISTINCT lead_id
				FROM appointments
				WHERE clinic_id = $1
					AND lead_id IS NOT NULL
					AND status IN ('scheduled', 'confirmed', 'pending_confirmation', 'rescheduled', 'converted_from_lead')
			)
	`, clinicID).Scan(&summary.HotLeadsWithoutAppointment)
	if err != nil {
		return ScheduleSummary{}, err
	}

	// 4. Overdue followups
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM leads
		WHERE clinic_id = $1
			AND next_action_at IS NOT NULL
			AND next_action_at < NOW()
	`, clinicID).Scan(&summary.OverdueFollowUps)
	if err != nil {
		return ScheduleSummary{}, err
	}

	// 5. Leads converted to appointments
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM appointments
		WHERE clinic_id = $1
			AND starts_at >= $2::timestamptz
			AND starts_at < $2::timestamptz + interval '1 day'
			AND source = 'lead_conversion'
	`, clinicID, dateStr).Scan(&summary.LeadsConvertedToAppointments)
	if err != nil {
		return ScheduleSummary{}, err
	}

	// 6. Appointments by professional
	profRows, err := r.db.QueryContext(ctx, `
		SELECT a.professional_id::text, p.full_name, COUNT(a.id)
		FROM appointments a
		JOIN clinic_professionals p ON p.id = a.professional_id AND p.clinic_id = a.clinic_id
		WHERE a.clinic_id = $1
			AND a.starts_at >= $2::timestamptz
			AND a.starts_at < $2::timestamptz + interval '1 day'
			AND a.status IN ('scheduled', 'confirmed', 'pending_confirmation', 'rescheduled', 'converted_from_lead', 'completed')
		GROUP BY a.professional_id, p.full_name
		ORDER BY COUNT(a.id) DESC
	`, clinicID, dateStr)
	if err != nil {
		return ScheduleSummary{}, err
	}
	defer profRows.Close()

	for profRows.Next() {
		var pc ProfessionalCountModel
		if err := profRows.Scan(&pc.ProfessionalID, &pc.ProfessionalName, &pc.AppointmentCount); err != nil {
			return ScheduleSummary{}, err
		}
		summary.AppointmentsByProfessional = append(summary.AppointmentsByProfessional, pc)
	}
	if err := profRows.Err(); err != nil {
		return ScheduleSummary{}, err
	}

	// 7. Top services by demand
	serviceRows, err := r.db.QueryContext(ctx, `
		SELECT a.service_id::text, s.name, COUNT(a.id)
		FROM appointments a
		JOIN clinic_services s ON s.id = a.service_id AND s.clinic_id = a.clinic_id
		WHERE a.clinic_id = $1
			AND a.starts_at >= $2::timestamptz
			AND a.starts_at < $2::timestamptz + interval '1 day'
			AND a.status IN ('scheduled', 'confirmed', 'pending_confirmation', 'rescheduled', 'converted_from_lead', 'completed')
		GROUP BY a.service_id, s.name
		ORDER BY COUNT(a.id) DESC
		LIMIT 5
	`, clinicID, dateStr)
	if err != nil {
		return ScheduleSummary{}, err
	}
	defer serviceRows.Close()

	for serviceRows.Next() {
		var sc ServiceCountModel
		if err := serviceRows.Scan(&sc.ServiceID, &sc.ServiceName, &sc.AppointmentCount); err != nil {
			return ScheduleSummary{}, err
		}
		summary.TopServicesByScheduleDemand = append(summary.TopServicesByScheduleDemand, sc)
	}
	if err := serviceRows.Err(); err != nil {
		return ScheduleSummary{}, err
	}

	return summary, nil
}

func (r *PostgresRepository) GetActiveProfessionalsWithHours(ctx context.Context, clinicID string) ([]ProfessionalHours, error) {
	const query = `
		SELECT id::text, working_hours
		FROM clinic_professionals
		WHERE clinic_id = $1 AND is_active = true
	`
	rows, err := r.db.QueryContext(ctx, query, clinicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var professionals []ProfessionalHours
	for rows.Next() {
		var p ProfessionalHours
		if err := rows.Scan(&p.ID, &p.WorkingHours); err != nil {
			return nil, err
		}
		professionals = append(professionals, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return professionals, nil
}

func (r *PostgresRepository) GetAppointmentsForDate(ctx context.Context, clinicID string, dateStr string) ([]AppointmentSummary, error) {
	const query = `
		SELECT starts_at, ends_at, professional_id::text
		FROM appointments
		WHERE clinic_id = $1
			AND starts_at >= $2::timestamptz
			AND starts_at < $2::timestamptz + interval '1 day'
			AND status IN ('scheduled', 'confirmed', 'pending_confirmation', 'rescheduled', 'converted_from_lead')
	`
	rows, err := r.db.QueryContext(ctx, query, clinicID, dateStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []AppointmentSummary
	for rows.Next() {
		var a AppointmentSummary
		if err := rows.Scan(&a.StartsAt, &a.EndsAt, &a.ProfessionalID); err != nil {
			return nil, err
		}
		appointments = append(appointments, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return appointments, nil
}
