package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/auth"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/clinics"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/config"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/health"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/pkg/database"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		log.Printf("database readiness disabled: %v", err)
	} else {
		defer db.Close()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", health.Healthz)
	mux.HandleFunc("GET /readyz", health.Readyz(db))

	if db != nil {
		authRepository := auth.NewPostgresRepository(db)
		authService, err := auth.NewService(authRepository, cfg.JWTSecret, cfg.JWTExpiresIn)
		if err != nil {
			log.Printf("auth routes disabled: %v", err)
		} else {
			auth.RegisterRoutes(mux, auth.NewHandler(authService))

			tokenManager, err := auth.NewTokenManager(cfg.JWTSecret)
			if err != nil {
				log.Printf("protected routes disabled: %v", err)
			} else {
				clinicRepository := clinics.NewPostgresRepository(db)
				clinicService := clinics.NewService(clinicRepository)
				clinics.RegisterRoutes(mux, clinics.NewHandler(clinicService), tokenManager)
			}
		}
	}

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("clinicflow api listening on %s", cfg.HTTPAddr)
		errCh <- server.ListenAndServe()
	}()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	case sig := <-shutdownCh:
		log.Printf("received signal %s, shutting down", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("graceful shutdown failed: %v", err)
		}
	}
}
