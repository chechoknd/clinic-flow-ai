package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Pinger interface {
	PingContext(ctx context.Context) error
}

type response struct {
	Status string `json:"status"`
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, response{Status: "ok"})
}

func Readyz(db Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			writeJSON(w, http.StatusServiceUnavailable, response{Status: "not_ready"})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, response{Status: "not_ready"})
			return
		}

		writeJSON(w, http.StatusOK, response{Status: "ready"})
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, payload response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
