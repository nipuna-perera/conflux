// Password hashing utilities
// Provides secure password hashing and verification using bcrypt
// Ensures password security best practices throughout the application
package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword securely hashes a password using bcrypt
// Returns hashed password string suitable for database storage
func HashPassword(password string) (string, error) {
	// Bcrypt hashing implementation with appropriate cost
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// VerifyPassword compares password with hash
// Returns true if password matches the hash
func VerifyPassword(password, hash string) bool {
	// Password verification implementation using bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
