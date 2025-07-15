package utils

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "testpassword123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false, // bcrypt can hash empty strings
		},
		{
			name:     "long password",
			password: "this_is_a_very_long_password_that_should_still_work_fine_123456789",
			wantErr:  false,
		},
		{
			name:     "password with special characters",
			password: "p@ssw0rd!@#$%^&*()",
			wantErr:  false,
		},
		{
			name:     "unicode password",
			password: "пароль123",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify hash is not empty
				if hash == "" {
					t.Error("HashPassword() returned empty hash")
				}

				// Verify hash is different from password
				if hash == tt.password {
					t.Error("HashPassword() returned password unchanged")
				}

				// Verify hash can be used with bcrypt.CompareHashAndPassword
				if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(tt.password)); err != nil {
					t.Errorf("Generated hash cannot be verified: %v", err)
				}
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	// Generate test hash
	testPassword := "testpassword123"
	testHash, err := HashPassword(testPassword)
	if err != nil {
		t.Fatalf("Failed to generate test hash: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "correct password",
			password: testPassword,
			hash:     testHash,
			want:     true,
		},
		{
			name:     "incorrect password",
			password: "wrongpassword",
			hash:     testHash,
			want:     false,
		},
		{
			name:     "empty password with valid hash",
			password: "",
			hash:     testHash,
			want:     false,
		},
		{
			name:     "valid password with empty hash",
			password: testPassword,
			hash:     "",
			want:     false,
		},
		{
			name:     "both empty",
			password: "",
			hash:     "",
			want:     false,
		},
		{
			name:     "invalid hash format",
			password: testPassword,
			hash:     "invalid_hash",
			want:     false,
		},
		{
			name:     "case sensitive password",
			password: "TestPassword123",
			hash:     testHash,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := VerifyPassword(tt.password, tt.hash)
			if got != tt.want {
				t.Errorf("VerifyPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Benchmark tests
func BenchmarkHashPassword(b *testing.B) {
	password := "benchmarkpassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := HashPassword(password)
		if err != nil {
			b.Fatalf("HashPassword() error = %v", err)
		}
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	password := "benchmarkpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		b.Fatalf("Failed to generate test hash: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		VerifyPassword(password, hash)
	}
}
