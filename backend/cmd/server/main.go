// Main entry point for the Go backend server
// Initializes database connections, sets up routes, and starts the HTTP server
// Acts as the composition root for dependency injection
package main

import (
	"log"
	"net/http"

	"conflux/internal/api"
	apiHandlers "conflux/internal/api/handlers"
	"conflux/internal/config"
	"conflux/internal/database"
	"conflux/internal/repository/mysql"
	"conflux/internal/repository/postgres"
	"conflux/internal/service"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database connection (MySQL or PostgreSQL based on config)
	dbFactory := database.NewConnectionFactory(cfg)
	db, err := dbFactory.NewConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Health check database connection
	if err := dbFactory.HealthCheck(db); err != nil {
		log.Fatal("Database health check failed:", err)
	}

	// Run database migrations
	migrator := database.NewMigrator(db, cfg.DBType)
	if err := migrator.Up(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Set up repository layer with database connection
	var userRepo service.UserRepository
	var authRepo service.AuthRepository

	switch cfg.DBType {
	case "mysql":
		userRepo = mysql.NewUserRepository(db)
		authRepo = mysql.NewAuthRepository(db)
	case "postgres":
		userRepo = postgres.NewUserRepository(db)
		authRepo = postgres.NewAuthRepository(db)
	default:
		log.Fatal("Unsupported database type:", cfg.DBType)
	}

	// Initialize service layer with repository dependencies
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, authRepo)

	// Set up API handlers with service dependencies
	healthHandler := apiHandlers.NewHealthHandler(db)
	authHandler := apiHandlers.NewAuthHandler(authService)
	userHandler := apiHandlers.NewUserHandler(userService)

	// Configure middleware chain and set up routes
	router := api.SetupRoutes(userHandler, authHandler, healthHandler)

	// Configure CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins(cfg.AllowedOrigins),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)(router)

	// Start HTTP server
	addr := cfg.Host + ":" + cfg.Port
	log.Printf("Server starting on %s", addr)
	log.Printf("Database type: %s", cfg.DBType)

	if err := http.ListenAndServe(addr, corsHandler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
