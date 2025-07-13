// HTTP route configuration and middleware setup
// Defines all API endpoints and applies appropriate middleware
// Acts as the routing hub for the entire API surface
package api

import (
	"net/http"

	"configarr/internal/api/handlers"
	"configarr/internal/api/middleware"

	"github.com/gorilla/mux"
)

// SetupRoutes configures all HTTP routes and middleware
// Returns configured router ready for HTTP server
func SetupRoutes(
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	healthHandler *handlers.HealthHandler,
) *mux.Router {
	router := mux.NewRouter()

	// Global middleware chain
	router.Use(middleware.Logging)
	router.Use(middleware.Recovery)

	// API routes
	api := router.PathPrefix("/api").Subrouter()

	// Health check endpoint
	api.HandleFunc("/health", healthHandler.CheckHealth).Methods("GET")

	// Public routes (no authentication required)
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", authHandler.Login).Methods("POST")
	auth.HandleFunc("/register", authHandler.Register).Methods("POST")

	// Protected routes (authentication required)
	protected := api.PathPrefix("/users").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/profile", userHandler.GetProfile).Methods("GET")
	protected.HandleFunc("/profile", userHandler.UpdateProfile).Methods("PUT")
	protected.HandleFunc("/{id}", userHandler.GetUser).Methods("GET")

	// Logout endpoint (requires auth)
	logoutHandler := middleware.AuthMiddleware(http.HandlerFunc(authHandler.Logout))
	auth.Handle("/logout", logoutHandler).Methods("POST")

	return router
}
