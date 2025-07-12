// Health check handler for monitoring and load balancer probes
// Provides endpoint to verify service availability and database connectivity
// Essential for container orchestration and monitoring systems
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// HealthHandler provides health check endpoints
type HealthHandler struct {
	db *sql.DB
}

// NewHealthHandler creates a new health check handler
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// CheckHealth returns service health status
// GET /health - Returns 200 OK if service is healthy
func (h *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	// Health check implementation:
	// - Verify database connectivity
	// - Check critical dependencies
	// - Return appropriate HTTP status

	response := map[string]interface{}{
		"status": "healthy",
		"checks": map[string]string{},
	}

	// Check database connectivity
	if h.db != nil {
		if err := h.db.Ping(); err != nil {
			response["status"] = "unhealthy"
			response["checks"].(map[string]string)["database"] = "failed: " + err.Error()
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			response["checks"].(map[string]string)["database"] = "healthy"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
