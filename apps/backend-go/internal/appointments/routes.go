package appointments

import (
	"net/http"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/auth"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, tokenManager *auth.TokenManager) {
	mux.Handle("GET /api/appointments", tokenManager.Authenticate(http.HandlerFunc(handler.List)))
	mux.Handle("GET /api/appointments/{id}", tokenManager.Authenticate(http.HandlerFunc(handler.Get)))
	mux.Handle("POST /api/appointments", tokenManager.Authenticate(http.HandlerFunc(handler.Create)))
	mux.Handle("PUT /api/appointments/{id}", tokenManager.Authenticate(http.HandlerFunc(handler.Update)))
	mux.Handle("POST /api/appointments/{id}/status", tokenManager.Authenticate(http.HandlerFunc(handler.UpdateStatus)))
	mux.Handle("POST /api/appointments/{id}/reschedule", tokenManager.Authenticate(http.HandlerFunc(handler.Reschedule)))
	mux.Handle("POST /api/leads/{id}/convert-to-appointment", tokenManager.Authenticate(http.HandlerFunc(handler.ConvertLead)))
}
