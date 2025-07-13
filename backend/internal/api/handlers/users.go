// User management HTTP handlers
// Processes user CRUD operations via HTTP requests
// Handles user profile management and user listing endpoints
package handlers

import (
	"net/http"
	"strconv"

	"configarr/internal/service"
	"configarr/pkg/utils"

	"github.com/gorilla/mux"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates user handler with service dependency
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile handles user profile retrieval
// GET /users/profile - Returns current user's profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Profile retrieval implementation:
	// - Extract user ID from JWT token
	// - Call user service to get profile
	// - Return JSON response

	// For now, we'll return a placeholder response
	// In a real implementation, we'd extract user ID from the JWT token context
	utils.ErrorResponse(w, http.StatusNotImplemented, "Profile retrieval not yet implemented")
}

// UpdateProfile handles user profile updates
// PUT /users/profile - Updates current user's profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Profile update implementation
	utils.ErrorResponse(w, http.StatusNotImplemented, "Profile update not yet implemented")
}

// GetUser handles user retrieval by ID
// GET /users/{id} - Returns user information (admin only)
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// User retrieval implementation
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		utils.ErrorResponse(w, http.StatusBadRequest, "User ID required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	utils.JSONResponse(w, http.StatusOK, user)
}
