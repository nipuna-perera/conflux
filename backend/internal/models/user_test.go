package models

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name string
		user User
		want error
	}{
		{
			name: "valid user",
			user: User{
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
			},
			want: nil,
		},
		{
			name: "user with whitespace",
			user: User{
				Email:     "  test@example.com  ",
				FirstName: "  John  ",
				LastName:  "  Doe  ",
			},
			want: nil,
		},
		{
			name: "empty user",
			user: User{},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalEmail := tt.user.Email
			originalFirstName := tt.user.FirstName
			originalLastName := tt.user.LastName

			err := tt.user.Validate()

			if err != tt.want {
				t.Errorf("User.Validate() error = %v, want %v", err, tt.want)
			}

			// Check that whitespace was trimmed
			if tt.user.Email != trimmedString(originalEmail) {
				t.Errorf("Email not trimmed properly: got %q, want %q",
					tt.user.Email, trimmedString(originalEmail))
			}
			if tt.user.FirstName != trimmedString(originalFirstName) {
				t.Errorf("FirstName not trimmed properly: got %q, want %q",
					tt.user.FirstName, trimmedString(originalFirstName))
			}
			if tt.user.LastName != trimmedString(originalLastName) {
				t.Errorf("LastName not trimmed properly: got %q, want %q",
					tt.user.LastName, trimmedString(originalLastName))
			}
		})
	}
}

func TestUser_FullName(t *testing.T) {
	tests := []struct {
		name     string
		user     User
		expected string
	}{
		{
			name: "normal full name",
			user: User{
				FirstName: "John",
				LastName:  "Doe",
			},
			expected: "John Doe",
		},
		{
			name: "empty first name",
			user: User{
				FirstName: "",
				LastName:  "Doe",
			},
			expected: " Doe",
		},
		{
			name: "empty last name",
			user: User{
				FirstName: "John",
				LastName:  "",
			},
			expected: "John ",
		},
		{
			name: "both names empty",
			user: User{
				FirstName: "",
				LastName:  "",
			},
			expected: " ",
		},
		{
			name: "names with spaces",
			user: User{
				FirstName: "John Michael",
				LastName:  "Van Doe",
			},
			expected: "John Michael Van Doe",
		},
		{
			name: "single character names",
			user: User{
				FirstName: "J",
				LastName:  "D",
			},
			expected: "J D",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.FullName()
			if got != tt.expected {
				t.Errorf("User.FullName() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestUser_JSONSerialization(t *testing.T) {
	user := User{
		ID:        123,
		Email:     "test@example.com",
		Password:  "hashed_password", // This should be hidden
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user to JSON: %v", err)
	}

	// Convert to string for easier testing
	jsonStr := string(jsonData)

	// Check that password is not included in JSON (has json:"-" tag)
	if containsSubstring(jsonStr, "password") || containsSubstring(jsonStr, "hashed_password") {
		t.Error("Password should not be included in JSON output")
	}

	// Check that other fields are included
	expectedFields := []string{
		"\"id\":123",
		"\"email\":\"test@example.com\"",
		"\"first_name\":\"John\"",
		"\"last_name\":\"Doe\"",
		"\"created_at\":",
		"\"updated_at\":",
	}

	for _, field := range expectedFields {
		if !containsSubstring(jsonStr, field) {
			t.Errorf("Expected field %q not found in JSON: %s", field, jsonStr)
		}
	}

	// Test JSON unmarshaling
	var unmarshaled User
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal user from JSON: %v", err)
	}

	// Verify unmarshaled data (password should remain empty since it's not in JSON)
	if unmarshaled.ID != user.ID {
		t.Errorf("ID mismatch: got %d, want %d", unmarshaled.ID, user.ID)
	}
	if unmarshaled.Email != user.Email {
		t.Errorf("Email mismatch: got %q, want %q", unmarshaled.Email, user.Email)
	}
	if unmarshaled.Password != "" {
		t.Errorf("Password should be empty after unmarshaling, got %q", unmarshaled.Password)
	}
	if unmarshaled.FirstName != user.FirstName {
		t.Errorf("FirstName mismatch: got %q, want %q", unmarshaled.FirstName, user.FirstName)
	}
	if unmarshaled.LastName != user.LastName {
		t.Errorf("LastName mismatch: got %q, want %q", unmarshaled.LastName, user.LastName)
	}
}

func TestUser_StructFields(t *testing.T) {
	now := time.Now()
	user := User{
		ID:        42,
		Email:     "user@example.com",
		Password:  "secret_hash",
		FirstName: "Jane",
		LastName:  "Smith",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Test field assignments and retrieval
	if user.ID != 42 {
		t.Errorf("ID = %d, want 42", user.ID)
	}
	if user.Email != "user@example.com" {
		t.Errorf("Email = %q, want %q", user.Email, "user@example.com")
	}
	if user.Password != "secret_hash" {
		t.Errorf("Password = %q, want %q", user.Password, "secret_hash")
	}
	if user.FirstName != "Jane" {
		t.Errorf("FirstName = %q, want %q", user.FirstName, "Jane")
	}
	if user.LastName != "Smith" {
		t.Errorf("LastName = %q, want %q", user.LastName, "Smith")
	}
	if !user.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", user.CreatedAt, now)
	}
	if !user.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt = %v, want %v", user.UpdatedAt, now)
	}
}

func TestUser_ZeroValue(t *testing.T) {
	var user User

	// Test zero values
	if user.ID != 0 {
		t.Errorf("Zero value ID = %d, want 0", user.ID)
	}
	if user.Email != "" {
		t.Errorf("Zero value Email = %q, want empty string", user.Email)
	}
	if user.Password != "" {
		t.Errorf("Zero value Password = %q, want empty string", user.Password)
	}
	if user.FirstName != "" {
		t.Errorf("Zero value FirstName = %q, want empty string", user.FirstName)
	}
	if user.LastName != "" {
		t.Errorf("Zero value LastName = %q, want empty string", user.LastName)
	}
	if !user.CreatedAt.IsZero() {
		t.Errorf("Zero value CreatedAt = %v, want zero time", user.CreatedAt)
	}
	if !user.UpdatedAt.IsZero() {
		t.Errorf("Zero value UpdatedAt = %v, want zero time", user.UpdatedAt)
	}

	// FullName should work with zero values
	fullName := user.FullName()
	if fullName != " " {
		t.Errorf("FullName() with zero values = %q, want %q", fullName, " ")
	}

	// Validate should work with zero values
	err := user.Validate()
	if err != nil {
		t.Errorf("Validate() with zero values = %v, want nil", err)
	}
}

// Helper functions
func trimmedString(s string) string {
	return strings.TrimSpace(s)
}

func containsSubstring(s, substr string) bool {
	return strings.Contains(s, substr)
}

// Benchmark tests
func BenchmarkUser_FullName(b *testing.B) {
	user := User{
		FirstName: "John",
		LastName:  "Doe",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = user.FullName()
	}
}

func BenchmarkUser_Validate(b *testing.B) {
	user := User{
		Email:     "  test@example.com  ",
		FirstName: "  John  ",
		LastName:  "  Doe  ",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset the values since Validate modifies them
		user.Email = "  test@example.com  "
		user.FirstName = "  John  "
		user.LastName = "  Doe  "
		_ = user.Validate()
	}
}

func BenchmarkUser_JSONMarshal(b *testing.B) {
	user := User{
		ID:        123,
		Email:     "test@example.com",
		Password:  "hashed_password",
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(user)
		if err != nil {
			b.Fatal(err)
		}
	}
}
