// Authentication HTTP handlers
// Processes login, registration, and logout HTTP requests
// Transforms HTTP requests/responses and delegates to service layer
package handlers

import (
	"encoding/json"
	"net/http"

	"conflux/internal/models"
	"conflux/internal/service"
	"conflux/pkg/utils"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates authentication handler with service dependency
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles user login requests
// POST /auth/login - Authenticates user and returns JWT token
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// HTTP handler implementation:
	// - Parse and validate JSON request body
	// - Call authentication service
	// - Return JSON response with token or error

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, response)
}

// Register handles user registration requests
// POST /auth/register - Creates new user account
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Registration handler implementation
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Note: This would typically use UserService.CreateUser
	// For now, we'll implement basic registration logic here
	utils.ErrorResponse(w, http.StatusNotImplemented, "Registration not yet implemented")
}

// Logout handles user logout requests
// POST /auth/logout - Invalidates user session
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Extract token from Authorization header
	token := r.Header.Get("Authorization")
	if token == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Authorization header required")
		return
	}

	// Remove "Bearer " prefix
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	if err := h.authService.Logout(r.Context(), token); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}
