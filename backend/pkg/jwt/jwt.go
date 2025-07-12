// JWT token utilities package
// Provides token generation, validation, and claims extraction
// Reusable JWT functionality across the application
package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenManager handles JWT operations
type TokenManager struct {
	secretKey []byte
	issuer    string
}

// NewTokenManager creates a new JWT token manager
func NewTokenManager(secretKey, issuer string) *TokenManager {
	return &TokenManager{
		secretKey: []byte(secretKey),
		issuer:    issuer,
	}
}

// Claims represents JWT token claims
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token for user
// Includes user ID, email, and expiration claims
func (tm *TokenManager) GenerateToken(userID int, email string, duration time.Duration) (string, error) {
	// Token generation implementation
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    tm.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tm.secretKey)
}

// ValidateToken verifies JWT signature and extracts claims
// Returns user information from valid tokens
func (tm *TokenManager) ValidateToken(tokenString string) (*Claims, error) {
	// Token validation implementation
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tm.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
