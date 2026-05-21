package migrations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var ErrNoMigrationDir = errors.New("migration directory does not exist")

type Migration struct {
	Version string
	Path    string
	SQL     string
}

type AppliedMigration struct {
	Version string
	Skipped bool
}

func Run(ctx context.Context, db *sql.DB, dir string) ([]AppliedMigration, error) {
	migrations, err := Load(dir)
	if err != nil {
		return nil, err
	}
	if err := ensureSchemaMigrations(ctx, db); err != nil {
		return nil, err
	}

	results := make([]AppliedMigration, 0, len(migrations))
	for _, migration := range migrations {
		applied, err := alreadyApplied(ctx, db, migration.Version)
		if err != nil {
			return nil, err
		}
		if applied {
			results = append(results, AppliedMigration{Version: migration.Version, Skipped: true})
			continue
		}
		if err := apply(ctx, db, migration); err != nil {
			return nil, err
		}
		results = append(results, AppliedMigration{Version: migration.Version})
	}

	return results, nil
}

func Load(dir string) ([]Migration, error) {
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNoMigrationDir
	}
	if err != nil {
		return nil, err
	}

	migrations := make([]Migration, 0)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		version := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		migrations = append(migrations, Migration{Version: version, Path: path, SQL: string(content)})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})
	return migrations, nil
}

func ensureSchemaMigrations(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func alreadyApplied(ctx context.Context, db *sql.DB, version string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)`, version).Scan(&exists)
	return exists, err
}

func apply(ctx context.Context, db *sql.DB, migration Migration) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, migration.SQL); err != nil {
		return fmt.Errorf("apply migration %s: %w", migration.Version, err)
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations(version) VALUES ($1)`, migration.Version); err != nil {
		return fmt.Errorf("record migration %s: %w", migration.Version, err)
	}
	return tx.Commit()
}
