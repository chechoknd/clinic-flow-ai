package services

import (
	"net/http"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/auth"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, tokenManager *auth.TokenManager) {
	mux.Handle("GET /api/services", tokenManager.Authenticate(http.HandlerFunc(handler.List)))
	mux.Handle("GET /api/services/{id}", tokenManager.Authenticate(http.HandlerFunc(handler.Get)))
	mux.Handle("POST /api/services", tokenManager.Authenticate(auth.RequireRoles(auth.RoleClinicAdmin)(http.HandlerFunc(handler.Create))))
	mux.Handle("PUT /api/services/{id}", tokenManager.Authenticate(auth.RequireRoles(auth.RoleClinicAdmin)(http.HandlerFunc(handler.Update))))
	mux.Handle("DELETE /api/services/{id}", tokenManager.Authenticate(auth.RequireRoles(auth.RoleClinicAdmin)(http.HandlerFunc(handler.Delete))))
}
