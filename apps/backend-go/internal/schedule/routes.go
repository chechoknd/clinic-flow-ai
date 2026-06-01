package schedule

import (
	"net/http"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/auth"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, tokenManager *auth.TokenManager) {
	mux.Handle("GET /api/schedule/availability", tokenManager.Authenticate(http.HandlerFunc(handler.GetAvailability)))
}
