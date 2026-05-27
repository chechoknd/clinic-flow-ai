package dashboard

import (
	"net/http"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/auth"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, tokenManager *auth.TokenManager) {
	mux.Handle("GET /api/dashboard/summary", tokenManager.Authenticate(http.HandlerFunc(handler.Summary)))
	mux.Handle("GET /api/dashboard/actions", tokenManager.Authenticate(http.HandlerFunc(handler.PriorityActions)))
}
