package jwt

import (
	"testing"
	"time"
)

func TestNewTokenManager(t *testing.T) {
	tests := []struct {
		name      string
		secretKey string
		issuer    string
	}{
		{
			name:      "valid token manager",
			secretKey: "test-secret-key",
			issuer:    "test-issuer",
		},
		{
			name:      "empty secret key",
			secretKey: "",
			issuer:    "test-issuer",
		},
		{
			name:      "empty issuer",
			secretKey: "test-secret-key",
			issuer:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := NewTokenManager(tt.secretKey, tt.issuer)
			if tm == nil {
				t.Fatal("NewTokenManager returned nil")
			}
			if string(tm.secretKey) != tt.secretKey {
				t.Errorf("expected secretKey %s, got %s", tt.secretKey, string(tm.secretKey))
			}
			if tm.issuer != tt.issuer {
				t.Errorf("expected issuer %s, got %s", tt.issuer, tm.issuer)
			}
		})
	}
}

func TestTokenManager_GenerateToken(t *testing.T) {
	tm := NewTokenManager("test-secret-key", "test-issuer")

	tests := []struct {
		name     string
		userID   int
		email    string
		duration time.Duration
		wantErr  bool
	}{
		{
			name:     "valid token generation",
			userID:   123,
			email:    "test@example.com",
			duration: time.Hour,
			wantErr:  false,
		},
		{
			name:     "zero user ID",
			userID:   0,
			email:    "test@example.com",
			duration: time.Hour,
			wantErr:  false,
		},
		{
			name:     "empty email",
			userID:   123,
			email:    "",
			duration: time.Hour,
			wantErr:  false,
		},
		{
			name:     "zero duration",
			userID:   123,
			email:    "test@example.com",
			duration: 0,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := tm.GenerateToken(tt.userID, tt.email, tt.duration)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if token != "" {
					t.Error("expected empty token on error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if token == "" {
					t.Error("expected non-empty token")
				}

				// Verify token structure (should have 3 parts separated by dots)
				parts := splitToken(token)
				if len(parts) != 3 {
					t.Errorf("expected JWT to have 3 parts, got %d", len(parts))
				}
			}
		})
	}
}

func TestTokenManager_ValidateToken(t *testing.T) {
	tm := NewTokenManager("test-secret-key", "test-issuer")

	// Generate a valid token for testing
	validToken, err := tm.GenerateToken(123, "test@example.com", time.Hour)
	if err != nil {
		t.Fatalf("failed to generate valid token: %v", err)
	}

	// Generate an expired token
	expiredToken, err := tm.GenerateToken(456, "expired@example.com", -time.Hour)
	if err != nil {
		t.Fatalf("failed to generate expired token: %v", err)
	}

	// Generate token with different secret for signature mismatch test
	wrongSecretTM := NewTokenManager("wrong-secret", "test-issuer")
	wrongSignatureToken, err := wrongSecretTM.GenerateToken(789, "wrong@example.com", time.Hour)
	if err != nil {
		t.Fatalf("failed to generate wrong signature token: %v", err)
	}

	tests := []struct {
		name           string
		tokenString    string
		wantErr        bool
		expectedUserID int
		expectedEmail  string
	}{
		{
			name:           "valid token",
			tokenString:    validToken,
			wantErr:        false,
			expectedUserID: 123,
			expectedEmail:  "test@example.com",
		},
		{
			name:        "expired token",
			tokenString: expiredToken,
			wantErr:     true,
		},
		{
			name:        "invalid signature",
			tokenString: wrongSignatureToken,
			wantErr:     true,
		},
		{
			name:        "malformed token",
			tokenString: "invalid.token",
			wantErr:     true,
		},
		{
			name:        "empty token",
			tokenString: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := tm.ValidateToken(tt.tokenString)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if claims != nil {
					t.Error("expected nil claims on error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if claims == nil {
					t.Error("expected non-nil claims")
					return
				}
				if claims.UserID != tt.expectedUserID {
					t.Errorf("expected UserID %d, got %d", tt.expectedUserID, claims.UserID)
				}
				if claims.Email != tt.expectedEmail {
					t.Errorf("expected Email %s, got %s", tt.expectedEmail, claims.Email)
				}
				if claims.Issuer != tm.issuer {
					t.Errorf("expected Issuer %s, got %s", tm.issuer, claims.Issuer)
				}
			}
		})
	}
}

func TestTokenManager_TokenRoundTrip(t *testing.T) {
	tm := NewTokenManager("test-secret-key-for-roundtrip", "roundtrip-issuer")

	tests := []struct {
		name     string
		userID   int
		email    string
		duration time.Duration
	}{
		{
			name:     "basic roundtrip",
			userID:   100,
			email:    "roundtrip@example.com",
			duration: time.Hour,
		},
		{
			name:     "special characters in email",
			userID:   200,
			email:    "user+test@example-domain.co.uk",
			duration: time.Minute * 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate token
			token, err := tm.GenerateToken(tt.userID, tt.email, tt.duration)
			if err != nil {
				t.Fatalf("failed to generate token: %v", err)
			}
			if token == "" {
				t.Fatal("got empty token")
			}

			// Validate token
			claims, err := tm.ValidateToken(token)
			if err != nil {
				t.Fatalf("failed to validate token: %v", err)
			}
			if claims == nil {
				t.Fatal("got nil claims")
			}

			// Verify claims match input
			if claims.UserID != tt.userID {
				t.Errorf("expected UserID %d, got %d", tt.userID, claims.UserID)
			}
			if claims.Email != tt.email {
				t.Errorf("expected Email %s, got %s", tt.email, claims.Email)
			}
			if claims.Issuer != tm.issuer {
				t.Errorf("expected Issuer %s, got %s", tm.issuer, claims.Issuer)
			}

			// Verify timing
			now := time.Now()
			if claims.IssuedAt == nil || claims.IssuedAt.Time.After(now) {
				t.Error("IssuedAt should be before or equal to now")
			}
			if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(now) {
				t.Error("ExpiresAt should be after now")
			}
		})
	}
}

func TestTokenManager_DifferentSecretKeys(t *testing.T) {
	tm1 := NewTokenManager("secret-key-1", "issuer-1")
	tm2 := NewTokenManager("secret-key-2", "issuer-2")

	// Generate token with tm1
	token, err := tm1.GenerateToken(123, "test@example.com", time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// Try to validate with tm1 (should succeed)
	claims, err := tm1.ValidateToken(token)
	if err != nil {
		t.Errorf("validation with same secret should succeed: %v", err)
	}
	if claims == nil {
		t.Error("expected claims from same secret validation")
	} else if claims.UserID != 123 {
		t.Errorf("expected UserID 123, got %d", claims.UserID)
	}

	// Try to validate with tm2 (should fail due to different secret)
	claims, err = tm2.ValidateToken(token)
	if err == nil {
		t.Error("validation with different secret should fail")
	}
	if claims != nil {
		t.Error("expected nil claims with different secret")
	}
}

// Helper function to split JWT token into parts
func splitToken(token string) []string {
	return strings.Split(token, ".")
}

// Benchmark tests
func BenchmarkTokenManager_GenerateToken(b *testing.B) {
	tm := NewTokenManager("benchmark-secret-key", "benchmark-issuer")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tm.GenerateToken(123, "benchmark@example.com", time.Hour)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTokenManager_ValidateToken(b *testing.B) {
	tm := NewTokenManager("benchmark-secret-key", "benchmark-issuer")

	// Pre-generate token for validation benchmark
	token, err := tm.GenerateToken(123, "benchmark@example.com", time.Hour)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tm.ValidateToken(token)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTokenManager_GenerateAndValidate(b *testing.B) {
	tm := NewTokenManager("benchmark-secret-key", "benchmark-issuer")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		token, err := tm.GenerateToken(123, "benchmark@example.com", time.Hour)
		if err != nil {
			b.Fatal(err)
		}

		_, err = tm.ValidateToken(token)
		if err != nil {
			b.Fatal(err)
		}
	}
}
