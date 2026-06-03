package clinics

import (
	"context"
	"database/sql"
	"errors"
)

var ErrClinicNotFound = errors.New("clinic not found")

type Repository interface {
	FindByID(ctx context.Context, clinicID string) (Clinic, error)
	Update(ctx context.Context, clinic Clinic) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) FindByID(ctx context.Context, clinicID string) (Clinic, error) {
	const query = `
		SELECT
			id::text,
			name,
			clinic_type,
			city,
			country_code,
			currency_code,
			phone,
			whatsapp,
			address,
			opening_hours,
			general_faq,
			communication_tone
		FROM clinics
		WHERE id = $1
		LIMIT 1
	`

	var clinic Clinic
	err := r.db.QueryRowContext(ctx, query, clinicID).Scan(
		&clinic.ID,
		&clinic.Name,
		&clinic.ClinicType,
		&clinic.City,
		&clinic.CountryCode,
		&clinic.CurrencyCode,
		&clinic.Phone,
		&clinic.WhatsApp,
		&clinic.Address,
		&clinic.OpeningHours,
		&clinic.GeneralFAQ,
		&clinic.CommunicationTone,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Clinic{}, ErrClinicNotFound
	}
	if err != nil {
		return Clinic{}, err
	}

	return clinic, nil
}

func (r *PostgresRepository) Update(ctx context.Context, clinic Clinic) error {
	const query = `
		UPDATE clinics
		SET
			name = $2,
			city = $3,
			country_code = $4,
			currency_code = $5,
			phone = $6,
			whatsapp = $7,
			address = $8,
			opening_hours = $9,
			general_faq = $10,
			communication_tone = $11
		WHERE id = $1
	`

	res, err := r.db.ExecContext(
		ctx,
		query,
		clinic.ID,
		clinic.Name,
		clinic.City,
		clinic.CountryCode,
		clinic.CurrencyCode,
		clinic.Phone,
		clinic.WhatsApp,
		clinic.Address,
		clinic.OpeningHours,
		clinic.GeneralFAQ,
		clinic.CommunicationTone,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrClinicNotFound
	}

	return nil
}
