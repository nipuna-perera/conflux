// Development endpoints for streamlining development workflow
// Only available in development environment for security
package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"conflux/internal/service"
	"conflux/pkg/utils"
)

// DevHandler handles development-specific endpoints
type DevHandler struct {
	devService *service.DevService
}

// NewDevHandler creates a new development handler
func NewDevHandler(devService *service.DevService) *DevHandler {
	return &DevHandler{
		devService: devService,
	}
}

// GetDevToken provides a JWT token for the development user
// POST /dev/token
func (h *DevHandler) GetDevToken(w http.ResponseWriter, r *http.Request) {
	// Only allow in development environment
	if os.Getenv("ENVIRONMENT") != "development" {
		utils.ErrorResponse(w, http.StatusForbidden, "Development endpoints not available")
		return
	}

	token, err := h.devService.GetDevToken(r.Context())
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to generate dev token: "+err.Error())
		return
	}

	response := map[string]interface{}{
		"token": token,
		"user": map[string]string{
			"email":      "dev@conflux.local",
			"first_name": "Dev",
			"last_name":  "User",
		},
		"instructions": "Use this token in Authorization header: Bearer " + token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateDevUser ensures the development user exists
// POST /dev/user
func (h *DevHandler) CreateDevUser(w http.ResponseWriter, r *http.Request) {
	// Only allow in development environment
	if os.Getenv("ENVIRONMENT") != "development" {
		utils.ErrorResponse(w, http.StatusForbidden, "Development endpoints not available")
		return
	}

	err := h.devService.CreateDevUser(r.Context())
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to create dev user: "+err.Error())
		return
	}

	response := map[string]interface{}{
		"message": "Development user ready",
		"credentials": map[string]string{
			"email":    "dev@conflux.local",
			"password": "password123",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
