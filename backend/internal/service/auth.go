// Authentication service layer
// Handles login, token generation, and session management
// Implements JWT token creation and validation logic
package service

import (
	"context"
	"fmt"
	"time"

	"configarr/internal/models"
	"configarr/pkg/jwt"
	"configarr/pkg/utils"
)

// AuthRepository defines data access methods for authentication
type AuthRepository interface {
	CreateSession(ctx context.Context, userID int, token string, expiresAt time.Time) error
	ValidateSession(ctx context.Context, token string) (*models.User, error)
	InvalidateSession(ctx context.Context, token string) error
}

// AuthService handles authentication business logic
type AuthService struct {
	userRepo     UserRepository
	authRepo     AuthRepository
	tokenManager *jwt.TokenManager
}

// NewAuthService creates authentication service with dependencies
func NewAuthService(userRepo UserRepository, authRepo AuthRepository) *AuthService {
	// Initialize token manager with a default secret (should come from config)
	tokenManager := jwt.NewTokenManager("default-secret", "configarr")

	return &AuthService{
		userRepo:     userRepo,
		authRepo:     authRepo,
		tokenManager: tokenManager,
	}
}

// Login authenticates user credentials and returns JWT token
// Validates credentials, generates JWT, creates session record
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	// Validate login request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if !utils.VerifyPassword(req.Password, user.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	duration := time.Hour * 24 // 24 hours
	token, err := s.tokenManager.GenerateToken(user.ID, user.Email, duration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create session record
	expiresAt := time.Now().Add(duration)
	if err := s.authRepo.CreateSession(ctx, user.ID, token, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Sanitize user data
	user.Password = ""

	return &models.AuthResponse{
		Token:     token,
		ExpiresIn: int(duration.Seconds()),
		User:      user,
	}, nil
}

// ValidateToken verifies JWT token and returns user information
func (s *AuthService) ValidateToken(ctx context.Context, token string) (*models.User, error) {
	// Validate JWT token
	claims, err := s.tokenManager.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Validate session in database
	user, err := s.authRepo.ValidateSession(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("session not found or expired: %w", err)
	}

	// Ensure token claims match user
	if user.ID != claims.UserID {
		return nil, fmt.Errorf("token user mismatch")
	}

	// Sanitize user data
	user.Password = ""
	return user, nil
}

// Logout invalidates user session
func (s *AuthService) Logout(ctx context.Context, token string) error {
	return s.authRepo.InvalidateSession(ctx, token)
}
