package ai

import (
	"net/http"

	"github.com/chechoknd/clinic-flow-ai/apps/backend-go/internal/auth"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, tokenManager *auth.TokenManager) {
	mux.Handle("POST /api/ai/reply-suggestion", tokenManager.Authenticate(http.HandlerFunc(handler.ReplySuggestion)))
	mux.Handle("POST /api/ai/objection-handler", tokenManager.Authenticate(http.HandlerFunc(handler.ObjectionHandler)))
}
