package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestAuthenticateAddsClaimsToContext(t *testing.T) {
	manager, err := NewTokenManager("test-secret")
	if err != nil {
		t.Fatalf("new token manager: %v", err)
	}
	token, err := manager.Sign(User{ID: "user-id", ClinicID: "clinic-id", Role: RoleAssistant}, time.Now(), time.Hour)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	manager.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := ClaimsFromContext(r.Context())
		if !ok {
			t.Fatal("expected claims in context")
		}
		if claims.UserID != "user-id" || claims.ClinicID != "clinic-id" || claims.Role != RoleAssistant {
			t.Fatalf("unexpected claims: %#v", claims)
		}
		w.WriteHeader(http.StatusNoContent)
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
}

func TestAuthenticateRejectsMissingToken(t *testing.T) {
	manager, err := NewTokenManager("test-secret")
	if err != nil {
		t.Fatalf("new token manager: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()

	manager.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "UNAUTHORIZED") {
		t.Fatalf("expected unauthorized response, got %s", rec.Body.String())
	}
}

func TestAuthenticateRejectsExpiredToken(t *testing.T) {
	manager, err := NewTokenManager("test-secret")
	if err != nil {
		t.Fatalf("new token manager: %v", err)
	}
	token, err := manager.Sign(User{ID: "user-id", ClinicID: "clinic-id", Role: RoleAssistant}, time.Now().Add(-2*time.Hour), time.Hour)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	manager.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestRequireRolesAllowsConfiguredRole(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req = req.WithContext(context.WithValue(req.Context(), claimsContextKey, Claims{Role: RoleClinicAdmin}))
	rec := httptest.NewRecorder()

	RequireRoles(RoleClinicAdmin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
}

func TestRequireRolesRejectsDifferentRole(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req = req.WithContext(context.WithValue(req.Context(), claimsContextKey, Claims{Role: RoleAssistant}))
	rec := httptest.NewRecorder()

	RequireRoles(RoleClinicAdmin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}
}
