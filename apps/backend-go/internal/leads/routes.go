package leads

import (
	"net/http"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/auth"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, tokenManager *auth.TokenManager) {
	mux.Handle("GET /api/followups", tokenManager.Authenticate(http.HandlerFunc(handler.ListFollowUps)))
	mux.Handle("POST /api/followups/{id}/complete", tokenManager.Authenticate(http.HandlerFunc(handler.CompleteFollowUp)))
	mux.Handle("POST /api/followups/{id}/reschedule", tokenManager.Authenticate(http.HandlerFunc(handler.RescheduleFollowUp)))
	mux.Handle("GET /api/leads", tokenManager.Authenticate(http.HandlerFunc(handler.List)))
	mux.Handle("GET /api/leads/{id}", tokenManager.Authenticate(http.HandlerFunc(handler.Get)))
	mux.Handle("POST /api/leads", tokenManager.Authenticate(http.HandlerFunc(handler.Create)))
	mux.Handle("PUT /api/leads/{id}", tokenManager.Authenticate(http.HandlerFunc(handler.Update)))
}
