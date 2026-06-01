package professionals

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
)

var ErrProfessionalNotFound = errors.New("professional not found")

type Repository interface {
	List(ctx context.Context, clinicID string, filters ListFilters) ([]Professional, error)
	FindByID(ctx context.Context, clinicID, professionalID string) (Professional, error)
	Create(ctx context.Context, professional Professional, serviceIDs []string) (Professional, error)
	Update(ctx context.Context, professional Professional, serviceIDs []string) (Professional, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) List(ctx context.Context, clinicID string, filters ListFilters) ([]Professional, error) {
	query := `
		SELECT
			p.id::text,
			p.clinic_id::text,
			p.full_name,
			p.role_or_specialty,
			p.calendar_color,
			p.working_hours,
			p.is_active,
			COALESCE(jsonb_agg(ps.service_id::text) FILTER (WHERE ps.service_id IS NOT NULL), '[]'::jsonb) AS service_ids,
			p.created_at,
			p.updated_at
		FROM clinic_professionals p
		LEFT JOIN professional_services ps
			ON ps.clinic_id = p.clinic_id AND ps.professional_id = p.id
		WHERE p.clinic_id = $1
	`

	args := []any{clinicID}
	if filters.IsActive != nil {
		args = append(args, *filters.IsActive)
		query += ` AND p.is_active = $2`
	}
	if filters.ServiceID != "" {
		args = append(args, filters.ServiceID)
		query += ` AND EXISTS (
			SELECT 1
			FROM professional_services ps_filter
			WHERE ps_filter.clinic_id = p.clinic_id
				AND ps_filter.professional_id = p.id
				AND ps_filter.service_id = $` + placeholder(len(args)) + `
		)`
	}

	query += `
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var professionals []Professional
	for rows.Next() {
		p, err := scanProfessional(rows)
		if err != nil {
			return nil, err
		}
		professionals = append(professionals, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return professionals, nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, clinicID, professionalID string) (Professional, error) {
	const query = `
		SELECT
			p.id::text,
			p.clinic_id::text,
			p.full_name,
			p.role_or_specialty,
			p.calendar_color,
			p.working_hours,
			p.is_active,
			COALESCE(jsonb_agg(ps.service_id::text) FILTER (WHERE ps.service_id IS NOT NULL), '[]'::jsonb) AS service_ids,
			p.created_at,
			p.updated_at
		FROM clinic_professionals p
		LEFT JOIN professional_services ps
			ON ps.clinic_id = p.clinic_id AND ps.professional_id = p.id
		WHERE p.id = $1 AND p.clinic_id = $2
		GROUP BY p.id
		LIMIT 1
	`

	p, err := scanProfessional(r.db.QueryRowContext(ctx, query, professionalID, clinicID))
	if errors.Is(err, sql.ErrNoRows) {
		return Professional{}, ErrProfessionalNotFound
	}
	if err != nil {
		return Professional{}, err
	}

	return p, nil
}

func (r *PostgresRepository) Create(ctx context.Context, p Professional, serviceIDs []string) (Professional, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Professional{}, err
	}
	defer tx.Rollback()

	const query = `
		INSERT INTO clinic_professionals (
			clinic_id,
			full_name,
			role_or_specialty,
			calendar_color,
			working_hours,
			is_active
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id::text, created_at, updated_at
	`
	err = tx.QueryRowContext(
		ctx,
		query,
		p.ClinicID,
		p.FullName,
		p.RoleOrSpecialty,
		p.CalendarColor,
		p.WorkingHours,
		p.IsActive,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return Professional{}, err
	}

	if err := replaceServices(ctx, tx, p.ClinicID, p.ID, serviceIDs); err != nil {
		return Professional{}, err
	}
	if err := tx.Commit(); err != nil {
		return Professional{}, err
	}

	p.ServiceIDs = serviceIDsJSON(serviceIDs)
	return p, nil
}

func (r *PostgresRepository) Update(ctx context.Context, p Professional, serviceIDs []string) (Professional, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Professional{}, err
	}
	defer tx.Rollback()

	const query = `
		UPDATE clinic_professionals
		SET
			full_name = $3,
			role_or_specialty = $4,
			calendar_color = $5,
			working_hours = $6,
			is_active = $7
		WHERE id = $1 AND clinic_id = $2
		RETURNING created_at, updated_at
	`
	err = tx.QueryRowContext(
		ctx,
		query,
		p.ID,
		p.ClinicID,
		p.FullName,
		p.RoleOrSpecialty,
		p.CalendarColor,
		p.WorkingHours,
		p.IsActive,
	).Scan(&p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Professional{}, ErrProfessionalNotFound
	}
	if err != nil {
		return Professional{}, err
	}

	if err := replaceServices(ctx, tx, p.ClinicID, p.ID, serviceIDs); err != nil {
		return Professional{}, err
	}
	if err := tx.Commit(); err != nil {
		return Professional{}, err
	}

	p.ServiceIDs = serviceIDsJSON(serviceIDs)
	return p, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanProfessional(row rowScanner) (Professional, error) {
	var p Professional
	err := row.Scan(
		&p.ID,
		&p.ClinicID,
		&p.FullName,
		&p.RoleOrSpecialty,
		&p.CalendarColor,
		&p.WorkingHours,
		&p.IsActive,
		&p.ServiceIDs,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return Professional{}, err
	}
	return p, nil
}

func replaceServices(ctx context.Context, tx *sql.Tx, clinicID, professionalID string, serviceIDs []string) error {
	if _, err := tx.ExecContext(ctx, `
		DELETE FROM professional_services
		WHERE clinic_id = $1 AND professional_id = $2
	`, clinicID, professionalID); err != nil {
		return err
	}

	for _, serviceID := range serviceIDs {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO professional_services (clinic_id, professional_id, service_id)
			VALUES ($1, $2, $3)
		`, clinicID, professionalID, serviceID); err != nil {
			return err
		}
	}

	return nil
}

func serviceIDsJSON(serviceIDs []string) []byte {
	if serviceIDs == nil {
		serviceIDs = []string{}
	}
	payload, _ := json.Marshal(serviceIDs)
	return payload
}

func placeholder(position int) string {
	return strconv.Itoa(position)
}
