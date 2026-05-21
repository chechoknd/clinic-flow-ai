package leads

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrLeadNotFound = errors.New("lead not found")

type Repository interface {
	List(ctx context.Context, clinicID string, filter ListFilter) ([]Lead, int, error)
	FindByID(ctx context.Context, clinicID, leadID string) (Lead, []LeadNote, error)
	Create(ctx context.Context, lead Lead, initialNote string) (Lead, error)
	Update(ctx context.Context, clinicID, leadID string, status string, nextActionAt *sql.NullTime, note string) error
}

type ListFilter struct {
	Status    string
	ServiceID string
	Page      int
	PageSize  int
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

func (r *PostgresRepository) FindByID(ctx context.Context, clinicID, leadID string) (Lead, []LeadNote, error) {
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
		return Lead{}, nil, ErrLeadNotFound
	}
	if err != nil {
		return Lead{}, nil, err
	}

	const notesQuery = `
		SELECT
			id::text,
			body,
			created_at
		FROM lead_notes
		WHERE lead_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, notesQuery, leadID)
	if err != nil {
		return l, nil, err
	}
	defer rows.Close()

	var notes []LeadNote
	for rows.Next() {
		var n LeadNote
		if err := rows.Scan(&n.ID, &n.Body, &n.CreatedAt); err != nil {
			return l, nil, err
		}
		notes = append(notes, n)
	}

	return l, notes, nil
}

func (r *PostgresRepository) Create(ctx context.Context, l Lead, initialNote string) (Lead, error) {
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
		RETURNING id::text, created_at
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
	).Scan(&l.ID, &l.CreatedAt)

	if err != nil {
		return Lead{}, err
	}

	if initialNote != "" {
		const noteQuery = `INSERT INTO lead_notes (lead_id, body) VALUES ($1, $2)`
		if _, err := tx.ExecContext(ctx, noteQuery, l.ID, initialNote); err != nil {
			return Lead{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return Lead{}, err
	}

	return l, nil
}

func (r *PostgresRepository) Update(ctx context.Context, clinicID, leadID string, status string, nextActionAt *sql.NullTime, note string) error {
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
		args = append(args, nextActionAt)
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
		const noteQuery = `INSERT INTO lead_notes (lead_id, body) VALUES ($1, $2)`
		if _, err := tx.ExecContext(ctx, noteQuery, leadID, note); err != nil {
			return err
		}
	}

	return tx.Commit()
}
