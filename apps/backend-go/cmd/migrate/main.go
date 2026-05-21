package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"time"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/config"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/migrations"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/pkg/database"
)

func main() {
	dir := flag.String("dir", "../../database/migrations", "directory containing SQL migration files")
	flag.Parse()

	cfg := config.Load()
	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	results, err := migrations.Run(ctx, db, *dir)
	if errors.Is(err, migrations.ErrNoMigrationDir) {
		log.Fatalf("migration directory not found: %s", *dir)
	}
	if err != nil {
		log.Fatalf("run migrations: %v", err)
	}
	if len(results) == 0 {
		log.Println("no migrations found")
		return
	}
	for _, result := range results {
		if result.Skipped {
			log.Printf("skipped migration %s", result.Version)
			continue
		}
		log.Printf("applied migration %s", result.Version)
	}
}
