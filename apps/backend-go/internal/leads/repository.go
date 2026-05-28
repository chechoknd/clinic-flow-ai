package leads

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

var ErrLeadNotFound = errors.New("lead not found")

type Repository interface {
	List(ctx context.Context, clinicID string, filter ListFilter) ([]Lead, int, error)
	ListFollowUps(ctx context.Context, clinicID string, filter FollowUpFilter) ([]Lead, int, error)
	FindByID(ctx context.Context, clinicID, leadID string) (Lead, []LeadNote, []AIInsight, error)
	Create(ctx context.Context, lead Lead, initialNote string, insight *AIInsight) (Lead, error)
	Update(ctx context.Context, clinicID, leadID string, status string, nextActionAt *sql.NullTime, note string, contactOutcome string, insight *AIInsight) error
	CompleteFollowUp(ctx context.Context, clinicID, leadID, status, note, contactOutcome string) error
	RescheduleFollowUp(ctx context.Context, clinicID, leadID string, nextActionAt sql.NullTime, note string) error
}

type ListFilter struct {
	Status    string
	ServiceID string
	Page      int
	PageSize  int
}

type FollowUpFilter struct {
	Due      string
	Page     int
	PageSize int
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) List(ctx context.Context, clinicID string, filter ListFilter) ([]Lead, int, error) {
	where := "WHERE l.clinic_id = $1"
	args := []any{clinicID}
	argIdx := 2

	if filter.Status != "" {
		where += fmt.Sprintf(" AND l.status = $%d", argIdx)
		args = append(args, filter.Status)
		argIdx++
	}

	if filter.ServiceID != "" {
		where += fmt.Sprintf(" AND l.service_id = $%d", argIdx)
		args = append(args, filter.ServiceID)
		argIdx++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM leads l %s", where)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	query := fmt.Sprintf(`
		SELECT
			l.id::text,
			l.full_name,
			l.phone,
			l.service_id::text,
			s.name as service_name,
			l.status,
			l.source,
			l.next_action_at,
			l.created_at,
			l.updated_at
		FROM leads l
		LEFT JOIN clinic_services s ON l.service_id = s.id
		%s
		ORDER BY l.created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argIdx, argIdx+1)

	args = append(args, filter.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var leads []Lead
	for rows.Next() {
		var l Lead
		err := rows.Scan(
			&l.ID,
			&l.FullName,
			&l.Phone,
			&l.ServiceID,
			&l.ServiceName,
			&l.Status,
			&l.Source,
			&l.NextActionAt,
			&l.CreatedAt,
			&l.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		leads = append(leads, l)
	}

	return leads, total, nil
}

func (r *PostgresRepository) ListFollowUps(ctx context.Context, clinicID string, filter FollowUpFilter) ([]Lead, int, error) {
	where := "WHERE l.clinic_id = $1 AND l.next_action_at IS NOT NULL"
	args := []any{clinicID}
	argIdx := 2

	switch filter.Due {
	case "today":
		where += fmt.Sprintf(" AND l.next_action_at <= date_trunc('day', NOW()) + interval '1 day' AND l.next_action_at >= date_trunc('day', NOW())")
	case "overdue":
		where += fmt.Sprintf(" AND l.next_action_at < NOW()")
	case "upcoming":
		where += fmt.Sprintf(" AND l.next_action_at > NOW()")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM leads l %s", where)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	query := fmt.Sprintf(`
		SELECT
			l.id::text,
			l.full_name,
			l.phone,
			l.service_id::text,
			s.name as service_name,
			l.status,
			l.source,
			l.next_action_at,
			l.created_at,
			l.updated_at
		FROM leads l
		LEFT JOIN clinic_services s ON l.service_id = s.id
		%s
		ORDER BY l.next_action_at ASC, l.created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argIdx, argIdx+1)

	args = append(args, filter.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var leads []Lead
	for rows.Next() {
		var l Lead
		err := rows.Scan(
			&l.ID,
			&l.FullName,
			&l.Phone,
			&l.ServiceID,
			&l.ServiceName,
			&l.Status,
			&l.Source,
			&l.NextActionAt,
			&l.CreatedAt,
			&l.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		leads = append(leads, l)
	}

	return leads, total, nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, clinicID, leadID string) (Lead, []LeadNote, []AIInsight, error) {
	const leadQuery = `
		SELECT
			l.id::text,
			l.clinic_id::text,
			l.full_name,
			l.phone,
			l.service_id::text,
			s.name as service_name,
			l.status,
			l.source,
			l.next_action_at,
			l.created_at,
			l.updated_at
		FROM leads l
		LEFT JOIN clinic_services s ON l.service_id = s.id
		WHERE l.id = $1 AND l.clinic_id = $2
	`

	var l Lead
	err := r.db.QueryRowContext(ctx, leadQuery, leadID, clinicID).Scan(
		&l.ID,
		&l.ClinicID,
		&l.FullName,
		&l.Phone,
		&l.ServiceID,
		&l.ServiceName,
		&l.Status,
		&l.Source,
		&l.NextActionAt,
		&l.CreatedAt,
		&l.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Lead{}, nil, nil, ErrLeadNotFound
	}
	if err != nil {
		return Lead{}, nil, nil, err
	}

	const notesQuery = `
		SELECT
			id::text,
			body,
			contact_outcome,
			created_at
		FROM lead_notes
		WHERE lead_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, notesQuery, leadID)
	if err != nil {
		return l, nil, nil, err
	}
	defer rows.Close()

	var notes []LeadNote
	for rows.Next() {
		var n LeadNote
		var contactOutcome sql.NullString
		if err := rows.Scan(&n.ID, &n.Body, &contactOutcome, &n.CreatedAt); err != nil {
			return l, nil, nil, err
		}
		if contactOutcome.Valid {
			n.ContactOutcome = &contactOutcome.String
		}
		notes = append(notes, n)
	}

	insights, err := r.listAIInsights(ctx, clinicID, leadID)
	if err != nil {
		return l, notes, nil, err
	}

	return l, notes, insights, nil
}

func (r *PostgresRepository) Create(ctx context.Context, l Lead, initialNote string, insight *AIInsight) (Lead, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Lead{}, err
	}
	defer tx.Rollback()

	const leadQuery = `
		INSERT INTO leads (
			clinic_id,
			full_name,
			phone,
			service_id,
			status,
			source,
			next_action_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id::text, created_at, updated_at
	`

	err = tx.QueryRowContext(
		ctx,
		leadQuery,
		l.ClinicID,
		l.FullName,
		l.Phone,
		l.ServiceID,
		l.Status,
		l.Source,
		l.NextActionAt,
	).Scan(&l.ID, &l.CreatedAt, &l.UpdatedAt)

	if err != nil {
		return Lead{}, err
	}

	if initialNote != "" {
		const noteQuery = `INSERT INTO lead_notes (lead_id, body) VALUES ($1, $2)`
		if _, err := tx.ExecContext(ctx, noteQuery, l.ID, initialNote); err != nil {
			return Lead{}, err
		}
	}

	if insight != nil {
		insight.ClinicID = l.ClinicID
		insight.LeadID = l.ID
		if err := insertAIInsightTx(ctx, tx, *insight); err != nil {
			return Lead{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return Lead{}, err
	}

	return l, nil
}

func (r *PostgresRepository) Update(ctx context.Context, clinicID, leadID string, status string, nextActionAt *sql.NullTime, note string, contactOutcome string, insight *AIInsight) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := "UPDATE leads SET status = $3"
	args := []any{leadID, clinicID, status}
	argIdx := 4

	if nextActionAt != nil {
		query += fmt.Sprintf(", next_action_at = $%d", argIdx)
		args = append(args, *nextActionAt)
		argIdx++
	}

	query += fmt.Sprintf(" WHERE id = $1 AND clinic_id = $2")

	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrLeadNotFound
	}

	if note != "" {
		const noteQuery = `INSERT INTO lead_notes (lead_id, body, contact_outcome) VALUES ($1, $2, NULLIF($3, ''))`
		if _, err := tx.ExecContext(ctx, noteQuery, leadID, note, contactOutcome); err != nil {
			return err
		}
	}

	if insight != nil {
		insight.ClinicID = clinicID
		insight.LeadID = leadID
		if err := insertAIInsightTx(ctx, tx, *insight); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) listAIInsights(ctx context.Context, clinicID, leadID string) ([]AIInsight, error) {
	const query = `
		SELECT
			id::text,
			clinic_id::text,
			lead_id::text,
			ai_generation_id::text,
			intent,
			detected_objections,
			COALESCE(commercial_summary, ''),
			COALESCE(suggested_next_action, ''),
			COALESCE(source, ''),
			created_at
		FROM lead_ai_insights
		WHERE clinic_id = $1 AND lead_id = $2
		ORDER BY created_at DESC
		LIMIT 10
	`

	rows, err := r.db.QueryContext(ctx, query, clinicID, leadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var insights []AIInsight
	for rows.Next() {
		var insight AIInsight
		var generationID sql.NullString
		var objectionsRaw []byte
		if err := rows.Scan(
			&insight.ID,
			&insight.ClinicID,
			&insight.LeadID,
			&generationID,
			&insight.Intent,
			&objectionsRaw,
			&insight.CommercialSummary,
			&insight.SuggestedNextAction,
			&insight.Source,
			&insight.CreatedAt,
		); err != nil {
			return nil, err
		}
		if generationID.Valid {
			insight.AIGenerationID = &generationID.String
		}
		if len(objectionsRaw) > 0 {
			if err := json.Unmarshal(objectionsRaw, &insight.DetectedObjections); err != nil {
				return nil, err
			}
		}
		insights = append(insights, insight)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return insights, nil
}

func insertAIInsightTx(ctx context.Context, tx *sql.Tx, insight AIInsight) error {
	objections, err := json.Marshal(insight.DetectedObjections)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO lead_ai_insights (
			clinic_id,
			lead_id,
			ai_generation_id,
			intent,
			detected_objections,
			commercial_summary,
			suggested_next_action,
			source
		) VALUES ($1, $2, NULLIF($3, '')::uuid, $4, $5::jsonb, NULLIF($6, ''), NULLIF($7, ''), NULLIF($8, ''))
	`,
		insight.ClinicID,
		insight.LeadID,
		stringValue(insight.AIGenerationID),
		insight.Intent,
		string(objections),
		insight.CommercialSummary,
		insight.SuggestedNextAction,
		insight.Source,
	)
	return err
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func (r *PostgresRepository) CompleteFollowUp(ctx context.Context, clinicID, leadID, status, note, contactOutcome string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `
		UPDATE leads
		SET status = $3, next_action_at = NULL
		WHERE id = $1 AND clinic_id = $2
	`, leadID, clinicID, status)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrLeadNotFound
	}

	if note != "" {
		const noteQuery = `INSERT INTO lead_notes (lead_id, body, contact_outcome) VALUES ($1, $2, NULLIF($3, ''))`
		if _, err := tx.ExecContext(ctx, noteQuery, leadID, note, contactOutcome); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) RescheduleFollowUp(ctx context.Context, clinicID, leadID string, nextActionAt sql.NullTime, note string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `
		UPDATE leads
		SET next_action_at = $3
		WHERE id = $1 AND clinic_id = $2
	`, leadID, clinicID, nextActionAt)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrLeadNotFound
	}

	if note != "" {
		const noteQuery = `INSERT INTO lead_notes (lead_id, body) VALUES ($1, $2)`
		if _, err := tx.ExecContext(ctx, noteQuery, leadID, note); err != nil {
			return err
		}
	}

	return tx.Commit()
}
