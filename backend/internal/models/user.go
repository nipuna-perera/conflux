// User data model and business logic
// Defines the User entity structure and validation rules
// Used across service and repository layers for type safety
package models

import (
	"strings"
	"time"
)

// User represents a user entity in the system
type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password_hash"` // Hidden from JSON
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Validate performs business rule validation on user data
func (u *User) Validate() error {
	// Trim whitespace
	u.Email = strings.TrimSpace(u.Email)
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)

	// Basic validation rules would go here
	// For now, we'll keep it simple
	return nil
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}
