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

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/ai"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/appointments"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/auth"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/clinics"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/config"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/dashboard"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/health"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/leads"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/professionals"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/schedule"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/services"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/shared"
	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/pkg/database"
)

func main() {
	cfg := config.Load()

	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

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

				serviceRepository := services.NewPostgresRepository(db)
				serviceService := services.NewService(serviceRepository)
				services.RegisterRoutes(mux, services.NewHandler(serviceService), tokenManager)

				professionalRepository := professionals.NewPostgresRepository(db)
				professionalService := professionals.NewService(professionalRepository)
				professionals.RegisterRoutes(mux, professionals.NewHandler(professionalService), tokenManager)

				appointmentRepository := appointments.NewPostgresRepository(db)
				appointmentService := appointments.NewService(appointmentRepository)
				appointments.RegisterRoutes(mux, appointments.NewHandler(appointmentService), tokenManager)

				leadRepository := leads.NewPostgresRepository(db)
				leadService := leads.NewService(leadRepository)
				leads.RegisterRoutes(mux, leads.NewHandler(leadService), tokenManager)

				scheduleRepository := schedule.NewPostgresRepository(db)
				scheduleService := schedule.NewService(scheduleRepository)
				schedule.RegisterRoutes(mux, schedule.NewHandler(scheduleService), tokenManager)

				dashboardRepository := dashboard.NewPostgresRepository(db)
				dashboardService := dashboard.NewService(dashboardRepository)
				dashboard.RegisterRoutes(mux, dashboard.NewHandler(dashboardService), tokenManager)

				aiProvider, err := ai.NewProvider(cfg)
				if err != nil {
					log.Printf("ai module disabled: %v", err)
				} else {
					aiService := ai.NewService(aiProvider, clinicService, serviceService, leadService)
					ai.RegisterRoutes(mux, ai.NewHandler(aiService), tokenManager)
				}
			}
		}
	}

	server := &http.Server{
		Addr: cfg.HTTPAddr,
		Handler: shared.SecurityHeadersMiddleware(
			shared.CORSMiddleware(cfg.AllowedOrigins)(
				shared.NewRateLimiter(60, time.Minute).Limit(
					shared.MaxBytesMiddleware(1024 * 1024)(mux),
				),
			),
		),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       120 * time.Second,
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
