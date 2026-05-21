package health

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakePinger struct {
	err error
}

func (p fakePinger) PingContext(ctx context.Context) error {
	return p.err
}

func TestHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	Healthz(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if rec.Body.String() != "{\"status\":\"ok\"}\n" {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestReadyzWhenDatabaseIsReady(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	Readyz(fakePinger{})(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if rec.Body.String() != "{\"status\":\"ready\"}\n" {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestReadyzWhenDatabaseIsMissing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	Readyz(nil)(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, rec.Code)
	}
	if rec.Body.String() != "{\"status\":\"not_ready\"}\n" {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
