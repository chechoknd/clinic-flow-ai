package professionals

import (
	"net/http"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/auth"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, tokenManager *auth.TokenManager) {
	mux.Handle("GET /api/professionals", tokenManager.Authenticate(http.HandlerFunc(handler.List)))
	mux.Handle("GET /api/professionals/{id}", tokenManager.Authenticate(http.HandlerFunc(handler.Get)))
	mux.Handle("POST /api/professionals", tokenManager.Authenticate(auth.RequireRoles(auth.RoleClinicAdmin)(http.HandlerFunc(handler.Create))))
	mux.Handle("PUT /api/professionals/{id}", tokenManager.Authenticate(auth.RequireRoles(auth.RoleClinicAdmin)(http.HandlerFunc(handler.Update))))
}
