package dashboard

import (
	"context"
	"database/sql"
)

type Repository interface {
	Summary(ctx context.Context, clinicID string) (Summary, error)
	PriorityActions(ctx context.Context, clinicID string, limit int) ([]ActionItem, error)
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

func (r *PostgresRepository) PriorityActions(ctx context.Context, clinicID string, limit int) ([]ActionItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		WITH latest_insights AS (
			SELECT DISTINCT ON (lead_id)
				lead_id,
				intent,
				detected_objections,
				created_at
			FROM lead_ai_insights
			WHERE clinic_id = $1
			ORDER BY lead_id, created_at DESC
		),
		priority_leads AS (
			SELECT
				'overdue_followup' AS type,
				'urgent' AS tone,
				100 AS priority,
				l.id::text AS lead_id,
				l.full_name,
				l.phone,
				l.service_id::text,
				s.name AS service_name,
				l.status,
				l.source,
				'Seguimiento vencido' AS reason,
				l.next_action_at,
				l.created_at
			FROM leads l
			LEFT JOIN clinic_services s ON s.id = l.service_id AND s.clinic_id = l.clinic_id
			WHERE l.clinic_id = $1
				AND l.next_action_at IS NOT NULL
				AND l.next_action_at < NOW()

			UNION ALL

			SELECT
				'today_followup' AS type,
				'today' AS tone,
				80 AS priority,
				l.id::text AS lead_id,
				l.full_name,
				l.phone,
				l.service_id::text,
				s.name AS service_name,
				l.status,
				l.source,
				'Seguimiento para hoy' AS reason,
				l.next_action_at,
				l.created_at
			FROM leads l
			LEFT JOIN clinic_services s ON s.id = l.service_id AND s.clinic_id = l.clinic_id
			WHERE l.clinic_id = $1
				AND l.next_action_at IS NOT NULL
				AND l.next_action_at >= NOW()
				AND l.next_action_at < date_trunc('day', NOW()) + interval '1 day'

			UNION ALL

			SELECT
				'high_intent' AS type,
				'intent' AS tone,
				70 AS priority,
				l.id::text AS lead_id,
				l.full_name,
				l.phone,
				l.service_id::text,
				s.name AS service_name,
				l.status,
				l.source,
				'Alta intencion detectada por Inbox AI' AS reason,
				l.next_action_at,
				l.created_at
			FROM leads l
			JOIN latest_insights i ON i.lead_id = l.id AND i.intent = 'high'
			LEFT JOIN clinic_services s ON s.id = l.service_id AND s.clinic_id = l.clinic_id
			WHERE l.clinic_id = $1
				AND l.status NOT IN ('Convertido', 'Perdido')

			UNION ALL

			SELECT
				'detected_objection' AS type,
				'objection' AS tone,
				65 AS priority,
				l.id::text AS lead_id,
				l.full_name,
				l.phone,
				l.service_id::text,
				s.name AS service_name,
				l.status,
				l.source,
				'Objecion comercial detectada por Inbox AI' AS reason,
				l.next_action_at,
				l.created_at
			FROM leads l
			JOIN latest_insights i ON i.lead_id = l.id AND jsonb_array_length(i.detected_objections) > 0
			LEFT JOIN clinic_services s ON s.id = l.service_id AND s.clinic_id = l.clinic_id
			WHERE l.clinic_id = $1
				AND l.status NOT IN ('Convertido', 'Perdido')

			UNION ALL

			SELECT
				'new_lead' AS type,
				'new' AS tone,
				60 AS priority,
				l.id::text AS lead_id,
				l.full_name,
				l.phone,
				l.service_id::text,
				s.name AS service_name,
				l.status,
				l.source,
				'Nuevo lead sin primer contacto registrado' AS reason,
				l.next_action_at,
				l.created_at
			FROM leads l
			LEFT JOIN clinic_services s ON s.id = l.service_id AND s.clinic_id = l.clinic_id
			WHERE l.clinic_id = $1
				AND l.status = 'Nuevo'
		)
		SELECT
			type,
			tone,
			priority,
			lead_id,
			full_name,
			phone,
			service_id,
			service_name,
			status,
			source,
			reason,
			next_action_at,
			created_at
		FROM (
			SELECT
				priority_leads.*,
				ROW_NUMBER() OVER (
					PARTITION BY lead_id
					ORDER BY priority DESC, next_action_at ASC NULLS LAST, created_at DESC
				) AS row_number
			FROM priority_leads
		) deduplicated
		WHERE row_number = 1
		ORDER BY priority DESC, next_action_at ASC NULLS LAST, created_at DESC
		LIMIT $2
	`, clinicID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []ActionItem
	for rows.Next() {
		var action ActionItem
		if err := rows.Scan(
			&action.Type,
			&action.Tone,
			&action.Priority,
			&action.LeadID,
			&action.FullName,
			&action.Phone,
			&action.ServiceID,
			&action.ServiceName,
			&action.Status,
			&action.Source,
			&action.Reason,
			&action.NextActionAt,
			&action.CreatedAt,
		); err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return actions, nil
}
