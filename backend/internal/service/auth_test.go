package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"conflux/internal/models"
	"conflux/pkg/utils"
)

// MockAuthRepository is a mock implementation of the AuthRepository interface.
// It is used for testing purposes to simulate the behavior of an authentication repository.
// 
// Fields:
// - sessions: A map that stores session tokens and their corresponding session data.
// - createSessionErr: An error to simulate failures in the CreateSession method.
// - validateSessionErr: An error to simulate failures in the ValidateSession method.
// - invalidateSessionErr: An error to simulate failures in the InvalidateSession method.
// - userForSession: A predefined user to return during session validation, if set.
//
// Usage:
// - Use NewMockAuthRepository to create an instance of this mock.
// - Set the error fields (e.g., createSessionErr) to simulate specific error scenarios.
// - Use the sessions map to inspect or manipulate session data during tests.
type MockAuthRepository struct {
	sessions             map[string]*models.Session
	createSessionErr     error
	validateSessionErr   error
	invalidateSessionErr error
	userForSession       *models.User
}

// NewMockAuthRepository creates a new mock auth repository
func NewMockAuthRepository() *MockAuthRepository {
	return &MockAuthRepository{
		sessions: make(map[string]*models.Session),
	}
}

// CreateSession implements AuthRepository.CreateSession
func (m *MockAuthRepository) CreateSession(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	if m.createSessionErr != nil {
		return m.createSessionErr
	}

	session := &models.Session{
		ID:        len(m.sessions) + 1,
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	m.sessions[token] = session
	return nil
}

// ValidateSession implements AuthRepository.ValidateSession
func (m *MockAuthRepository) ValidateSession(ctx context.Context, token string) (*models.User, error) {
	if m.validateSessionErr != nil {
		return nil, m.validateSessionErr
	}

	session, exists := m.sessions[token]
	if !exists {
		return nil, errors.New("session not found")
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("session expired")
	}

	if m.userForSession != nil {
		return m.userForSession, nil
	}

	// Return a default user if none set
	return &models.User{
		ID:    session.UserID,
		Email: "test@example.com",
	}, nil
}

// InvalidateSession implements AuthRepository.InvalidateSession
func (m *MockAuthRepository) InvalidateSession(ctx context.Context, token string) error {
	if m.invalidateSessionErr != nil {
		return m.invalidateSessionErr
	}

	delete(m.sessions, token)
	return nil
}

// Helper methods for testing
func (m *MockAuthRepository) SetCreateSessionError(err error) {
	m.createSessionErr = err
}

func (m *MockAuthRepository) SetValidateSessionError(err error) {
	m.validateSessionErr = err
}

func (m *MockAuthRepository) SetInvalidateSessionError(err error) {
	m.invalidateSessionErr = err
}

func (m *MockAuthRepository) SetUserForSession(user *models.User) {
	m.userForSession = user
}

func (m *MockAuthRepository) HasSession(token string) bool {
	_, exists := m.sessions[token]
	return exists
}

func (m *MockAuthRepository) SessionCount() int {
	return len(m.sessions)
}

// Tests for AuthService
func TestNewAuthService(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockAuthRepo := NewMockAuthRepository()

	authService := NewAuthService(mockUserRepo, mockAuthRepo)

	if authService == nil {
		t.Fatal("NewAuthService returned nil")
	}
	if authService.userRepo != mockUserRepo {
		t.Error("AuthService not initialized with correct user repository")
	}
	if authService.authRepo != mockAuthRepo {
		t.Error("AuthService not initialized with correct auth repository")
	}
	if authService.tokenManager == nil {
		t.Error("AuthService token manager not initialized")
	}
}

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name             string
		loginReq         *models.LoginRequest
		setupUser        *models.User
		userRepoErr      error
		authRepoErr      error
		wantErr          bool
		errorContains    string
		validateResponse bool
	}{
		{
			name: "successful login",
			loginReq: &models.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupUser: &models.User{
				ID:        1,
				Email:     "test@example.com",
				Password:  mustHashPassword("password123"),
				FirstName: "John",
				LastName:  "Doe",
			},
			wantErr:          false,
			validateResponse: true,
		},
		{
			name: "user not found",
			loginReq: &models.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			userRepoErr:   errors.New("user not found"),
			wantErr:       true,
			errorContains: "invalid credentials",
		},
		{
			name: "incorrect password",
			loginReq: &models.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupUser: &models.User{
				ID:       1,
				Email:    "test@example.com",
				Password: mustHashPassword("password123"),
			},
			wantErr:       true,
			errorContains: "invalid credentials",
		},
		{
			name: "session creation error",
			loginReq: &models.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupUser: &models.User{
				ID:       1,
				Email:    "test@example.com",
				Password: mustHashPassword("password123"),
			},
			authRepoErr:   errors.New("database error"),
			wantErr:       true,
			errorContains: "failed to create session",
		},
		{
			name: "empty email",
			loginReq: &models.LoginRequest{
				Email:    "",
				Password: "password123",
			},
			wantErr:       true, // GetByEmail will fail with empty email
			errorContains: "invalid credentials",
		},
		{
			name: "empty password",
			loginReq: &models.LoginRequest{
				Email:    "test@example.com",
				Password: "",
			},
			setupUser: &models.User{
				ID:       1,
				Email:    "test@example.com",
				Password: mustHashPassword("password123"),
			},
			wantErr:       true,
			errorContains: "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := NewMockUserRepository()
			mockAuthRepo := NewMockAuthRepository()
			authService := NewAuthService(mockUserRepo, mockAuthRepo)

			// Set up user if needed
			if tt.setupUser != nil {
				err := mockUserRepo.Create(context.Background(), tt.setupUser)
				if err != nil {
					t.Fatalf("Failed to create test user: %v", err)
				}
			}

			// Set up repository errors
			if tt.userRepoErr != nil {
				mockUserRepo.SetGetByEmailError(tt.userRepoErr)
			}
			if tt.authRepoErr != nil {
				mockAuthRepo.SetCreateSessionError(tt.authRepoErr)
			}

			// Execute test
			response, err := authService.Login(context.Background(), tt.loginReq)

			// Validate results
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("error %q should contain %q", err.Error(), tt.errorContains)
				}
				if response != nil {
					t.Error("expected nil response on error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tt.validateResponse && response != nil {
					// Validate response structure
					if response.Token == "" {
						t.Error("response should contain non-empty token")
					}
					if response.ExpiresIn <= 0 {
						t.Error("response should contain positive expires_in")
					}
					if response.User == nil {
						t.Error("response should contain user")
					} else {
						if response.User.Password != "" {
							t.Error("response user password should be sanitized")
						}
						if response.User.Email != tt.loginReq.Email {
							t.Errorf("response user email mismatch: got %q, want %q",
								response.User.Email, tt.loginReq.Email)
						}
					}

					// Verify session was created
					if !mockAuthRepo.HasSession(response.Token) {
						t.Error("session should be created in repository")
					}
				}
			}
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	tests := []struct {
		name           string
		setupUser      *models.User
		generateToken  bool
		token          string
		authRepoErr    error
		userIDMismatch bool
		wantErr        bool
		errorContains  string
	}{
		{
			name: "successful token validation",
			setupUser: &models.User{
				ID:    1,
				Email: "test@example.com",
			},
			generateToken: true,
			wantErr:       false,
		},
		{
			name:          "invalid token format",
			token:         "invalid.token.format",
			wantErr:       true,
			errorContains: "invalid token",
		},
		{
			name:          "empty token",
			token:         "",
			wantErr:       true,
			errorContains: "invalid token",
		},
		{
			name: "session not found",
			setupUser: &models.User{
				ID:    1,
				Email: "test@example.com",
			},
			generateToken: true,
			authRepoErr:   errors.New("session not found"),
			wantErr:       true,
			errorContains: "session not found or expired",
		},
		{
			name: "user ID mismatch",
			setupUser: &models.User{
				ID:    1,
				Email: "test@example.com",
			},
			generateToken:  true,
			userIDMismatch: true,
			wantErr:        true,
			errorContains:  "token user mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := NewMockUserRepository()
			mockAuthRepo := NewMockAuthRepository()
			authService := NewAuthService(mockUserRepo, mockAuthRepo)

			var token string
			if tt.generateToken && tt.setupUser != nil {
				// Generate a valid token for the user
				var err error
				token, err = authService.tokenManager.GenerateToken(
					tt.setupUser.ID,
					tt.setupUser.Email,
					time.Hour,
				)
				if err != nil {
					t.Fatalf("Failed to generate token: %v", err)
				}

				// Create session
				err = mockAuthRepo.CreateSession(
					context.Background(),
					tt.setupUser.ID,
					token,
					time.Now().Add(time.Hour),
				)
				if err != nil {
					t.Fatalf("Failed to create session: %v", err)
				}

				// Set up user for session validation
				sessionUser := *tt.setupUser
				if tt.userIDMismatch {
					sessionUser.ID = tt.setupUser.ID + 999 // Different ID
				}
				mockAuthRepo.SetUserForSession(&sessionUser)
			} else if tt.token != "" {
				token = tt.token
			}

			// Set up repository errors
			if tt.authRepoErr != nil {
				mockAuthRepo.SetValidateSessionError(tt.authRepoErr)
			}

			// Execute test
			user, err := authService.ValidateToken(context.Background(), token)

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
					if user.Password != "" {
						t.Error("returned user password should be sanitized")
					}
					if user.ID != tt.setupUser.ID {
						t.Errorf("user ID mismatch: got %d, want %d", user.ID, tt.setupUser.ID)
					}
				}
			}
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	tests := []struct {
		name          string
		token         string
		repoErr       error
		wantErr       bool
		errorContains string
	}{
		{
			name:    "successful logout",
			token:   "valid_token_123",
			wantErr: false,
		},
		{
			name:          "repository error",
			token:         "valid_token_123",
			repoErr:       errors.New("database error"),
			wantErr:       true,
			errorContains: "database error",
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: false, // Repository handles empty token
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := NewMockUserRepository()
			mockAuthRepo := NewMockAuthRepository()
			authService := NewAuthService(mockUserRepo, mockAuthRepo)

			// Set up repository error
			if tt.repoErr != nil {
				mockAuthRepo.SetInvalidateSessionError(tt.repoErr)
			}

			// Execute test
			err := authService.Logout(context.Background(), tt.token)

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

func TestAuthService_Integration(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockAuthRepo := NewMockAuthRepository()
	authService := NewAuthService(mockUserRepo, mockAuthRepo)

	ctx := context.Background()

	// 1. Create a test user
	testUser := &models.User{
		Email:     "integration@example.com",
		Password:  mustHashPassword("testpassword123"),
		FirstName: "Integration",
		LastName:  "Test",
	}
	err := mockUserRepo.Create(ctx, testUser)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// 2. Test login
	loginReq := &models.LoginRequest{
		Email:    "integration@example.com",
		Password: "testpassword123",
	}

	authResponse, err := authService.Login(ctx, loginReq)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if authResponse.Token == "" {
		t.Error("Login should return a token")
	}
	if authResponse.User.Password != "" {
		t.Error("Login response user password should be sanitized")
	}

	// 3. Test token validation
	// Set up correct user for session validation
	mockAuthRepo.SetUserForSession(&models.User{
		ID:    testUser.ID,
		Email: testUser.Email,
	})

	user, err := authService.ValidateToken(ctx, authResponse.Token)
	if err != nil {
		t.Fatalf("Token validation failed: %v", err)
	}
	if user.Email != testUser.Email {
		t.Errorf("Validated user email mismatch: got %q, want %q",
			user.Email, testUser.Email)
	}
	if user.Password != "" {
		t.Error("Validated user password should be sanitized")
	}

	// 4. Test logout
	err = authService.Logout(ctx, authResponse.Token)
	if err != nil {
		t.Fatalf("Logout failed: %v", err)
	}

	// 5. Verify token is invalidated
	_, err = authService.ValidateToken(ctx, authResponse.Token)
	if err == nil {
		t.Error("Token validation should fail after logout")
	}

	// 6. Test login with wrong password
	wrongLoginReq := &models.LoginRequest{
		Email:    "integration@example.com",
		Password: "wrongpassword",
	}
	_, err = authService.Login(ctx, wrongLoginReq)
	if err == nil {
		t.Error("Login with wrong password should fail")
	}
	if !strings.Contains(err.Error(), "invalid credentials") {
		t.Errorf("Error should mention invalid credentials: %v", err)
	}
}

// Helper function to hash password for testing
func mustHashPassword(password string) string {
	hash, err := utils.HashPassword(password)
	if err != nil {
		panic("Failed to hash password in test: " + err.Error())
	}
	return hash
}

// Benchmark tests
func BenchmarkAuthService_Login(b *testing.B) {
	mockUserRepo := NewMockUserRepository()
	mockAuthRepo := NewMockAuthRepository()
	authService := NewAuthService(mockUserRepo, mockAuthRepo)

	// Set up test user
	testUser := &models.User{
		Email:    "bench@example.com",
		Password: mustHashPassword("password123"),
	}
	err := mockUserRepo.Create(context.Background(), testUser)
	if err != nil {
		b.Fatal(err)
	}

	loginReq := &models.LoginRequest{
		Email:    "bench@example.com",
		Password: "password123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := authService.Login(context.Background(), loginReq)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAuthService_ValidateToken(b *testing.B) {
	mockUserRepo := NewMockUserRepository()
	mockAuthRepo := NewMockAuthRepository()
	authService := NewAuthService(mockUserRepo, mockAuthRepo)

	// Set up test user and token
	testUser := &models.User{
		ID:       1,
		Email:    "bench@example.com",
		Password: mustHashPassword("password123"),
	}
	err := mockUserRepo.Create(context.Background(), testUser)
	if err != nil {
		b.Fatal(err)
	}

	loginReq := &models.LoginRequest{
		Email:    "bench@example.com",
		Password: "password123",
	}

	authResponse, err := authService.Login(context.Background(), loginReq)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := authService.ValidateToken(context.Background(), authResponse.Token)
		if err != nil {
			b.Fatal(err)
		}
	}
}
