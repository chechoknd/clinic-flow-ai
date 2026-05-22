package dashboard

import (
	"context"
	"database/sql"
)

type Repository interface {
	Summary(ctx context.Context, clinicID string) (Summary, error)
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
