package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jstaud/edgelink-go/pkg/models"
)

// Store interface - anything that can retrieve latest readings
type Store interface {
	Latest(id string) (models.Reading, bool)
}

// NewREST creates a new HTTP router with our API endpoints
func NewREST(store Store) http.Handler {
	// Chi is a lightweight HTTP router for Go
	r := chi.NewRouter()
	
	// Health check endpoint - returns 200 OK if service is running
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"status":"ok","service":"edgelink-go"}`))
	})
	
	// Get latest reading for a specific device
	r.Get("/api/readings/{id}", func(w http.ResponseWriter, req *http.Request) {
		// Extract device ID from URL path
		id := chi.URLParam(req, "id")
		
		// Try to get latest reading for this device
		if reading, exists := store.Latest(id); exists {
			// Found it! Send as JSON
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(reading)
			return
		}
		
		// Device not found
		http.NotFound(w, req)
	})
	
	// Placeholder for metrics endpoint (prometheus temporarily disabled)
	r.Get("/metrics", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("# Metrics temporarily disabled\n"))
	})
	
	return r
}
