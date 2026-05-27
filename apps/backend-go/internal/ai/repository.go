package ai

import (
	"context"
	"database/sql"
)

type GenerationRepository interface {
	CreateGeneration(ctx context.Context, record GenerationRecord) (string, error)
}

type PostgresGenerationRepository struct {
	db *sql.DB
}

func NewPostgresGenerationRepository(db *sql.DB) *PostgresGenerationRepository {
	return &PostgresGenerationRepository{db: db}
}

func (r *PostgresGenerationRepository) CreateGeneration(ctx context.Context, record GenerationRecord) (string, error) {
	const query = `
		INSERT INTO ai_generations (
			clinic_id,
			user_id,
			feature,
			provider,
			model,
			status,
			safety_status,
			error_code,
			input_char_count,
			output_char_count
		) VALUES ($1, NULLIF($2, '')::uuid, $3, $4, $5, $6, NULLIF($7, ''), NULLIF($8, ''), $9, $10)
		RETURNING id::text
	`

	var id string
	err := r.db.QueryRowContext(
		ctx,
		query,
		record.ClinicID,
		record.UserID,
		record.Feature,
		record.Provider,
		record.Model,
		record.Status,
		record.SafetyStatus,
		record.ErrorCode,
		record.InputCharCount,
		record.OutputCharCount,
	).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}
