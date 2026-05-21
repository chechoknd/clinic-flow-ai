package auth

import (
	"context"
	"database/sql"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

type Repository interface {
	FindByEmail(ctx context.Context, email string) (User, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) FindByEmail(ctx context.Context, email string) (User, error) {
	const query = `
		SELECT id::text, clinic_id::text, full_name, email, password_hash, role, is_active
		FROM users
		WHERE lower(email) = lower($1)
		LIMIT 1
	`

	var user User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.ClinicID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, err
	}

	return user, nil
}
