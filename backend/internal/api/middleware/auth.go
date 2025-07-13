// JWT authentication middleware
// Validates JWT tokens and extracts user information from requests
// Protects API endpoints requiring authentication
package middleware

import (
	"context"
	"net/http"
	"strings"

	"conflux/pkg/jwt"
	"conflux/pkg/utils"
)

// UserContextKey is the key for storing user in context
type UserContextKey string

const UserKey UserContextKey = "user"

// AuthMiddleware validates JWT tokens from Authorization header
// Extracts user information and adds to request context
// Returns 401 Unauthorized for invalid or missing tokens
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// JWT validation implementation:
		// - Extract token from Authorization header
		// - Validate JWT signature and expiration
		// - Extract user claims and add to context
		// - Call next handler or return 401

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		// Check Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Bearer token required")
			return
		}

		token := authHeader[7:] // Remove "Bearer " prefix

		// Validate token
		tokenManager := jwt.NewTokenManager("default-secret", "conflux") // Should come from config
		claims, err := tokenManager.ValidateToken(token)
		if err != nil {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Add user information to context
		ctx := context.WithValue(r.Context(), UserKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuthMiddleware validates tokens when present
// Used for endpoints that work with or without authentication
func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Optional authentication implementation
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token := authHeader[7:]
			tokenManager := jwt.NewTokenManager("default-secret", "conflux")
			if claims, err := tokenManager.ValidateToken(token); err == nil {
				ctx := context.WithValue(r.Context(), UserKey, claims)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}
