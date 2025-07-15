package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"conflux/internal/models"
)

// MockUserRepository is a mock implementation of the UserRepository interface,
// designed for use in unit tests. It allows developers to simulate various
// behaviors of a user repository by configuring its error fields (e.g., createErr,
// getByIDErr) to return specific errors during method calls.
//
// Usage:
// - Use NewMockUserRepository() to create a new instance.
// - Configure the mock's behavior by setting the error fields (e.g., createErr).
// - Call the methods (e.g., Create, GetByID) as you would with a real UserRepository.
//
// Behavioral Characteristics:
// - The mock stores users in memory and assigns unique IDs starting from 1.
// - Methods return copies of users to prevent unintended mutations.
// - If an error field is set (e.g., createErr), the corresponding method will
//   return that error instead of performing its normal operation.
type MockUserRepository struct {
	users         map[int]*models.User
	emailToUserID map[string]int
	nextID        int
	createErr     error
	getByIDErr    error
	getByEmailErr error
	updateErr     error
	deleteErr     error
}

// NewMockUserRepository creates a new mock user repository
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:         make(map[int]*models.User),
		emailToUserID: make(map[string]int),
		nextID:        1,
	}
}

// Create implements UserRepository.Create
func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	if m.createErr != nil {
		return m.createErr
	}

	// Assign ID and store user
	user.ID = m.nextID
	m.nextID++

	userCopy := *user
	m.users[user.ID] = &userCopy
	m.emailToUserID[user.Email] = user.ID

	return nil
}

// GetByID implements UserRepository.GetByID
func (m *MockUserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}

	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	// Return a copy to avoid mutations
	userCopy := *user
	return &userCopy, nil
}

// GetByEmail implements UserRepository.GetByEmail
func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.getByEmailErr != nil {
		return nil, m.getByEmailErr
	}

	userID, exists := m.emailToUserID[email]
	if !exists {
		return nil, errors.New("user not found")
	}

	user := m.users[userID]
	// Return a copy to avoid mutations
	userCopy := *user
	return &userCopy, nil
}

// Update implements UserRepository.Update
func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	if m.updateErr != nil {
		return m.updateErr
	}

	if _, exists := m.users[user.ID]; !exists {
		return errors.New("user not found")
	}

	// Update the stored user
	userCopy := *user
	m.users[user.ID] = &userCopy
	m.emailToUserID[user.Email] = user.ID

	return nil
}

// Delete implements UserRepository.Delete
func (m *MockUserRepository) Delete(ctx context.Context, id int) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}

	user, exists := m.users[id]
	if !exists {
		return errors.New("user not found")
	}

	delete(m.users, id)
	delete(m.emailToUserID, user.Email)

	return nil
}

// Helper methods for testing
func (m *MockUserRepository) SetCreateError(err error) {
	m.createErr = err
}

func (m *MockUserRepository) SetGetByIDError(err error) {
	m.getByIDErr = err
}

func (m *MockUserRepository) SetGetByEmailError(err error) {
	m.getByEmailErr = err
}

func (m *MockUserRepository) SetUpdateError(err error) {
	m.updateErr = err
}

func (m *MockUserRepository) SetDeleteError(err error) {
	m.deleteErr = err
}

func (m *MockUserRepository) UserCount() int {
	return len(m.users)
}

func (m *MockUserRepository) HasUser(id int) bool {
	_, exists := m.users[id]
	return exists
}

func (m *MockUserRepository) HasUserByEmail(email string) bool {
	_, exists := m.emailToUserID[email]
	return exists
}

// Tests for UserService
func TestNewUserService(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewUserService(mockRepo)

	if service == nil {
		t.Fatal("NewUserService returned nil")
	}
	if service.userRepo != mockRepo {
		t.Error("UserService not initialized with correct repository")
	}
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name          string
		request       *models.RegisterRequest
		existingEmail string
		repoCreateErr error
		repoEmailErr  error
		wantErr       bool
		errorContains string
	}{
		{
			name: "successful user creation",
			request: &models.RegisterRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			wantErr: false,
		},
		{
			name: "email already exists",
			request: &models.RegisterRequest{
				Email:     "existing@example.com",
				Password:  "password123",
				FirstName: "Jane",
				LastName:  "Doe",
			},
			existingEmail: "existing@example.com",
			wantErr:       true,
			errorContains: "email already exists",
		},
		{
			name: "repository create error",
			request: &models.RegisterRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			repoCreateErr: errors.New("database error"),
			wantErr:       true,
			errorContains: "failed to create user",
		},
		{
			name: "empty password",
			request: &models.RegisterRequest{
				Email:     "test@example.com",
				Password:  "",
				FirstName: "John",
				LastName:  "Doe",
			},
			wantErr: false, // HashPassword handles empty passwords
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockUserRepository()
			service := NewUserService(mockRepo)

			// Set up existing user if needed
			if tt.existingEmail != "" {
				existingUser := &models.User{
					Email: tt.existingEmail,
				}
				err := mockRepo.Create(context.Background(), existingUser)
				if err != nil {
					t.Fatalf("Failed to create existing user: %v", err)
				}
			}

			// Set up repository errors
			if tt.repoCreateErr != nil {
				mockRepo.SetCreateError(tt.repoCreateErr)
			}

			// Execute test
			user, err := service.CreateUser(context.Background(), tt.request)

			// Validate results
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("error %q should contain %q", err.Error(), tt.errorContains)
				}
				if user != nil {
					t.Error("expected nil user on error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if user == nil {
					t.Error("expected user, got nil")
				} else {
					// Validate user fields
					if user.Email != tt.request.Email {
						t.Errorf("email mismatch: got %q, want %q", user.Email, tt.request.Email)
					}
					if user.FirstName != tt.request.FirstName {
						t.Errorf("first name mismatch: got %q, want %q", user.FirstName, tt.request.FirstName)
					}
					if user.LastName != tt.request.LastName {
						t.Errorf("last name mismatch: got %q, want %q", user.LastName, tt.request.LastName)
					}
					if user.Password != "" {
						t.Error("password should be sanitized (empty) in returned user")
					}
					if user.ID == 0 {
						t.Error("user ID should be assigned")
					}
				}
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		setupUser     *models.User
		repoErr       error
		wantErr       bool
		errorContains string
	}{
		{
			name:   "successful user retrieval",
			userID: 1,
			setupUser: &models.User{
				Email:     "test@example.com",
				Password:  "hashed_password",
				FirstName: "John",
				LastName:  "Doe",
			},
			wantErr: false,
		},
		{
			name:          "user not found",
			userID:        999,
			wantErr:       true,
			errorContains: "failed to get user",
		},
		{
			name:          "repository error",
			userID:        1,
			repoErr:       errors.New("database error"),
			wantErr:       true,
			errorContains: "failed to get user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockUserRepository()
			service := NewUserService(mockRepo)

			// Set up user if needed
			if tt.setupUser != nil {
				err := mockRepo.Create(context.Background(), tt.setupUser)
				if err != nil {
					t.Fatalf("Failed to create test user: %v", err)
				}
			}

			// Set up repository error
			if tt.repoErr != nil {
				mockRepo.SetGetByIDError(tt.repoErr)
			}

			// Execute test
			user, err := service.GetUserByID(context.Background(), tt.userID)

			// Validate results
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("error %q should contain %q", err.Error(), tt.errorContains)
				}
				if user != nil {
					t.Error("expected nil user on error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if user == nil {
					t.Error("expected user, got nil")
				} else {
					// Validate password is sanitized
					if user.Password != "" {
						t.Error("password should be sanitized (empty) in returned user")
					}
					// Validate other fields
					if user.Email != tt.setupUser.Email {
						t.Errorf("email mismatch: got %q, want %q", user.Email, tt.setupUser.Email)
					}
				}
			}
		})
	}
}

func TestUserService_GetUserByEmail(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		setupUser     *models.User
		repoErr       error
		wantErr       bool
		errorContains string
	}{
		{
			name:  "successful user retrieval",
			email: "test@example.com",
			setupUser: &models.User{
				Email:     "test@example.com",
				Password:  "hashed_password",
				FirstName: "John",
				LastName:  "Doe",
			},
			wantErr: false,
		},
		{
			name:          "user not found",
			email:         "nonexistent@example.com",
			wantErr:       true,
			errorContains: "failed to get user",
		},
		{
			name:          "repository error",
			email:         "test@example.com",
			repoErr:       errors.New("database error"),
			wantErr:       true,
			errorContains: "failed to get user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockUserRepository()
			service := NewUserService(mockRepo)

			// Set up user if needed
			if tt.setupUser != nil {
				err := mockRepo.Create(context.Background(), tt.setupUser)
				if err != nil {
					t.Fatalf("Failed to create test user: %v", err)
				}
			}

			// Set up repository error
			if tt.repoErr != nil {
				mockRepo.SetGetByEmailError(tt.repoErr)
			}

			// Execute test
			user, err := service.GetUserByEmail(context.Background(), tt.email)

			// Validate results
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("error %q should contain %q", err.Error(), tt.errorContains)
				}
				if user != nil {
					t.Error("expected nil user on error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if user == nil {
					t.Error("expected user, got nil")
				} else {
					// Note: GetUserByEmail does NOT sanitize password (for auth purposes)
					if user.Password == "" {
						t.Error("password should NOT be sanitized in GetUserByEmail")
					}
					// Validate other fields
					if user.Email != tt.setupUser.Email {
						t.Errorf("email mismatch: got %q, want %q", user.Email, tt.setupUser.Email)
					}
				}
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	tests := []struct {
		name          string
		user          *models.User
		setupUser     *models.User
		repoErr       error
		wantErr       bool
		errorContains string
	}{
		{
			name: "successful user update",
			user: &models.User{
				ID:        1,
				Email:     "updated@example.com",
				FirstName: "Jane",
				LastName:  "Updated",
			},
			setupUser: &models.User{
				Email:     "original@example.com",
				FirstName: "John",
				LastName:  "Original",
			},
			wantErr: false,
		},
		{
			name: "repository error",
			user: &models.User{
				ID:        1,
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
			},
			setupUser: &models.User{
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
			},
			repoErr:       errors.New("database error"),
			wantErr:       true,
			errorContains: "failed to update user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockUserRepository()
			service := NewUserService(mockRepo)

			// Set up user if needed
			if tt.setupUser != nil {
				err := mockRepo.Create(context.Background(), tt.setupUser)
				if err != nil {
					t.Fatalf("Failed to create test user: %v", err)
				}
				// Set the ID for update
				tt.user.ID = 1
			}

			// Set up repository error
			if tt.repoErr != nil {
				mockRepo.SetUpdateError(tt.repoErr)
			}

			// Execute test
			err := service.UpdateUser(context.Background(), tt.user)

			// Validate results
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("error %q should contain %q", err.Error(), tt.errorContains)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// Integration test
func TestUserService_Integration(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewUserService(mockRepo)

	// Test complete user lifecycle
	ctx := context.Background()

	// 1. Create user
	registerReq := &models.RegisterRequest{
		Email:     "integration@example.com",
		Password:  "password123",
		FirstName: "Integration",
		LastName:  "Test",
	}

	createdUser, err := service.CreateUser(ctx, registerReq)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	if createdUser.ID == 0 {
		t.Error("Created user should have an ID")
	}

	// 2. Get user by ID
	retrievedUser, err := service.GetUserByID(ctx, createdUser.ID)
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}
	if retrievedUser.Email != registerReq.Email {
		t.Errorf("Retrieved user email mismatch: got %q, want %q",
			retrievedUser.Email, registerReq.Email)
	}
	if retrievedUser.Password != "" {
		t.Error("Retrieved user password should be sanitized")
	}

	// 3. Get user by email (for auth)
	authUser, err := service.GetUserByEmail(ctx, registerReq.Email)
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}
	if authUser.Password == "" {
		t.Error("Auth user password should NOT be sanitized")
	}

	// 4. Update user
	authUser.FirstName = "Updated"
	err = service.UpdateUser(ctx, authUser)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// 5. Verify update
	updatedUser, err := service.GetUserByID(ctx, authUser.ID)
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}
	if updatedUser.FirstName != "Updated" {
		t.Errorf("User not updated: got %q, want %q",
			updatedUser.FirstName, "Updated")
	}

	// 6. Try to create duplicate email
	duplicateReq := &models.RegisterRequest{
		Email:     registerReq.Email,
		Password:  "different_password",
		FirstName: "Duplicate",
		LastName:  "User",
	}
	_, err = service.CreateUser(ctx, duplicateReq)
	if err == nil {
		t.Error("Expected error for duplicate email")
	}
	if !strings.Contains(err.Error(), "email already exists") {
		t.Errorf("Error should mention duplicate email: %v", err)
	}
}

// Benchmark tests
func BenchmarkUserService_CreateUser(b *testing.B) {
	mockRepo := NewMockUserRepository()
	service := NewUserService(mockRepo)

	registerReq := &models.RegisterRequest{
		Email:     "bench@example.com",
		Password:  "password123",
		FirstName: "Benchmark",
		LastName:  "User",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use unique email for each iteration to avoid conflicts
		registerReq.Email = fmt.Sprintf("bench%d@example.com", i)
		_, err := service.CreateUser(context.Background(), registerReq)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUserService_GetUserByID(b *testing.B) {
	mockRepo := NewMockUserRepository()
	service := NewUserService(mockRepo)

	// Create a test user
	registerReq := &models.RegisterRequest{
		Email:     "bench@example.com",
		Password:  "password123",
		FirstName: "Benchmark",
		LastName:  "User",
	}
	user, err := service.CreateUser(context.Background(), registerReq)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GetUserByID(context.Background(), user.ID)
		if err != nil {
			b.Fatal(err)
		}
	}
}
