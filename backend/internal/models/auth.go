// Authentication-related data models
// Defines structures for login requests, JWT tokens, and sessions
// Provides validation for authentication payloads
package models

import (
	"strings"
	"time"
)

// LoginRequest represents user login credentials
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validates the login request
func (lr *LoginRequest) Validate() error {
	lr.Email = strings.TrimSpace(lr.Email)
	lr.Password = strings.TrimSpace(lr.Password)
	// Add validation logic here
	return nil
}

// RegisterRequest represents user registration data
type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Validate validates the registration request
func (rr *RegisterRequest) Validate() error {
	rr.Email = strings.TrimSpace(rr.Email)
	rr.FirstName = strings.TrimSpace(rr.FirstName)
	rr.LastName = strings.TrimSpace(rr.LastName)
	// Add validation logic here
	return nil
}

// AuthResponse represents successful authentication response
type AuthResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
	User      *User  `json:"user"`
}

// Session represents a user session
type Session struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
