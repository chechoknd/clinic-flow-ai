package database

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var ErrMissingURL = errors.New("database url is required")

func Open(databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		return nil, ErrMissingURL
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}
