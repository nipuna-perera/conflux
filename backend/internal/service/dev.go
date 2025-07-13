// Development utilities for streamlining development workflow
// Only enabled in development environment
package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"conflux/internal/models"
	"conflux/pkg/jwt"
)

// DevService provides development-specific utilities
type DevService struct {
	userService *UserService
	authService *AuthService
}

// NewDevService creates a new development service
func NewDevService(userService *UserService, authService *AuthService) *DevService {
	return &DevService{
		userService: userService,
		authService: authService,
	}
}

// GetDevToken generates a JWT token for the default development user
func (s *DevService) GetDevToken(ctx context.Context) (string, error) {
	if os.Getenv("ENVIRONMENT") != "development" {
		return "", fmt.Errorf("dev tokens only available in development environment")
	}

	// Get the development user
	user, err := s.userService.GetUserByEmail(ctx, "dev@conflux.local")
	if err != nil {
		return "", fmt.Errorf("development user not found: %w", err)
	}

	// Generate JWT token using TokenManager
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-key"
	}

	tokenManager := jwt.NewTokenManager(jwtSecret, "conflux-dev")
	token, err := tokenManager.GenerateToken(user.ID, user.Email, 24*time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to generate dev token: %w", err)
	}

	return token, nil
}

// CreateDevUser ensures the development user exists (fallback if migration doesn't run)
func (s *DevService) CreateDevUser(ctx context.Context) error {
	if os.Getenv("ENVIRONMENT") != "development" {
		return fmt.Errorf("dev user creation only available in development environment")
	}

	// Check if dev user already exists
	_, err := s.userService.GetUserByEmail(ctx, "dev@conflux.local")
	if err == nil {
		log.Println("Development user already exists")
		return nil
	}

	// Create the development user
	req := &models.RegisterRequest{
		Email:     "dev@conflux.local",
		Password:  "password123",
		FirstName: "Dev",
		LastName:  "User",
	}

	_, err = s.userService.CreateUser(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create dev user: %w", err)
	}

	log.Println("Development user created: dev@conflux.local / password123")
	return nil
}
