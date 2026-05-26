package services

import (
	"context"
	"database/sql"
	"errors"
)

var ErrServiceNotFound = errors.New("service not found")

type Repository interface {
	ListByClinicID(ctx context.Context, clinicID string) ([]ServiceEntity, error)
	FindByID(ctx context.Context, clinicID, serviceID string) (ServiceEntity, error)
	Create(ctx context.Context, service ServiceEntity) (ServiceEntity, error)
	Update(ctx context.Context, service ServiceEntity) error
	Delete(ctx context.Context, clinicID, serviceID string) error
	CurrencyCodeByClinicID(ctx context.Context, clinicID string) (string, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) ListByClinicID(ctx context.Context, clinicID string) ([]ServiceEntity, error) {
	const query = `
		SELECT
			id::text,
			clinic_id::text,
			name,
			description,
			duration_minutes,
			price_from::text,
			(SELECT c.currency_code FROM clinics c WHERE c.id = clinic_services.clinic_id),
			benefits,
			faq,
			common_objections,
			is_active
		FROM clinic_services
		WHERE clinic_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, clinicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []ServiceEntity
	for rows.Next() {
		var s ServiceEntity
		var priceFrom sql.NullString
		err := rows.Scan(
			&s.ID,
			&s.ClinicID,
			&s.Name,
			&s.Description,
			&s.DurationMinutes,
			&priceFrom,
			&s.CurrencyCode,
			&s.Benefits,
			&s.FAQ,
			&s.CommonObjections,
			&s.IsActive,
		)
		if err != nil {
			return nil, err
		}
		if priceFrom.Valid {
			value := DecimalString(priceFrom.String)
			s.PriceFrom = &value
		}
		entities = append(entities, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, clinicID, serviceID string) (ServiceEntity, error) {
	const query = `
		SELECT
			id::text,
			clinic_id::text,
			name,
			description,
			duration_minutes,
			price_from::text,
			(SELECT c.currency_code FROM clinics c WHERE c.id = clinic_services.clinic_id),
			benefits,
			faq,
			common_objections,
			is_active
		FROM clinic_services
		WHERE id = $1 AND clinic_id = $2
		LIMIT 1
	`

	var s ServiceEntity
	var priceFrom sql.NullString
	err := r.db.QueryRowContext(ctx, query, serviceID, clinicID).Scan(
		&s.ID,
		&s.ClinicID,
		&s.Name,
		&s.Description,
		&s.DurationMinutes,
		&priceFrom,
		&s.CurrencyCode,
		&s.Benefits,
		&s.FAQ,
		&s.CommonObjections,
		&s.IsActive,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return ServiceEntity{}, ErrServiceNotFound
	}
	if err != nil {
		return ServiceEntity{}, err
	}
	if priceFrom.Valid {
		value := DecimalString(priceFrom.String)
		s.PriceFrom = &value
	}

	return s, nil
}

func (r *PostgresRepository) Create(ctx context.Context, s ServiceEntity) (ServiceEntity, error) {
	const query = `
		INSERT INTO clinic_services (
			clinic_id,
			name,
			description,
			duration_minutes,
			price_from,
			benefits,
			faq,
			common_objections,
			is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id::text
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		s.ClinicID,
		s.Name,
		s.Description,
		s.DurationMinutes,
		s.PriceFrom,
		s.Benefits,
		s.FAQ,
		s.CommonObjections,
		s.IsActive,
	).Scan(&s.ID)

	if err != nil {
		return ServiceEntity{}, err
	}

	return s, nil
}

func (r *PostgresRepository) Update(ctx context.Context, s ServiceEntity) error {
	const query = `
		UPDATE clinic_services
		SET
			name = $3,
			description = $4,
			duration_minutes = $5,
			price_from = $6,
			benefits = $7,
			faq = $8,
			common_objections = $9,
			is_active = $10
		WHERE id = $1 AND clinic_id = $2
	`

	res, err := r.db.ExecContext(
		ctx,
		query,
		s.ID,
		s.ClinicID,
		s.Name,
		s.Description,
		s.DurationMinutes,
		s.PriceFrom,
		s.Benefits,
		s.FAQ,
		s.CommonObjections,
		s.IsActive,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrServiceNotFound
	}

	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, clinicID, serviceID string) error {
	const query = `
		DELETE FROM clinic_services
		WHERE id = $1 AND clinic_id = $2
	`

	res, err := r.db.ExecContext(ctx, query, serviceID, clinicID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrServiceNotFound
	}

	return nil
}

func (r *PostgresRepository) CurrencyCodeByClinicID(ctx context.Context, clinicID string) (string, error) {
	const query = `SELECT currency_code FROM clinics WHERE id = $1 LIMIT 1`

	var currencyCode string
	err := r.db.QueryRowContext(ctx, query, clinicID).Scan(&currencyCode)
	if err != nil {
		return "", err
	}
	return currencyCode, nil
}
