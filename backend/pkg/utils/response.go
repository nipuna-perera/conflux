// HTTP response utilities
// Provides standardized JSON response helpers for API consistency
// Ensures uniform error handling and response formatting
package utils

import (
	"encoding/json"
	"net/http"
)

// JSONResponse sends JSON response with appropriate headers
// Standardizes response format across all API endpoints
func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	// JSON response implementation
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Log error but don't try to send another response as headers are already written
		// In a production environment, you might want to log this error
		return
	}
}

// ErrorResponse sends standardized error response
// Provides consistent error format for client consumption
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	// Error response implementation
	response := map[string]interface{}{
		"error":   true,
		"message": message,
		"status":  statusCode,
	}
	JSONResponse(w, statusCode, response)
}

// SuccessResponse sends standardized success response
func SuccessResponse(w http.ResponseWriter, data interface{}) {
	response := map[string]interface{}{
		"error": false,
		"data":  data,
	}
	JSONResponse(w, http.StatusOK, response)
}
