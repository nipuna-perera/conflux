// User business logic service layer
// Orchestrates user operations between API handlers and repository
// Implements business rules, validation, and data transformation
package service

import (
	"context"
	"fmt"

	"conflux/internal/models"
	"conflux/pkg/utils"
)

// UserRepository defines data access methods for user entities
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
}

// UserService handles user business logic
type UserService struct {
	userRepo UserRepository
}

// NewUserService creates a new user service with repository dependency
func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser handles user registration business logic
// Validates input, hashes password, and creates user record
func (s *UserService) CreateUser(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	// Validate registration request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user record
	user := &models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Return sanitized user data (without password)
	user.Password = ""
	return user, nil
}

// GetUserByID retrieves user by ID with business rule validation
func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Sanitize password
	user.Password = ""
	return user, nil
}

// GetUserByEmail retrieves user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Don't sanitize password here as it might be needed for authentication
	return user, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	if err := user.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
