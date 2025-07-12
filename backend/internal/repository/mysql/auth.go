// MySQL implementation of AuthRepository interface
// Handles authentication session operations specific to MySQL database
// Implements session management and token validation for MySQL
package mysql

import (
	"context"
	"database/sql"
	"time"

	"configarr/internal/models"
)

// AuthRepository implements service.AuthRepository for MySQL
type AuthRepository struct {
	db *sql.DB
}

// NewAuthRepository creates a new MySQL auth repository
func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateSession creates a new session record in MySQL
func (r *AuthRepository) CreateSession(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO sessions (user_id, token, expires_at) 
		VALUES (?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query, userID, token, expiresAt)
	return err
}

// ValidateSession validates session token and returns user if valid
func (r *AuthRepository) ValidateSession(ctx context.Context, token string) (*models.User, error) {
	query := `
		SELECT u.id, u.email, u.password_hash, u.first_name, u.last_name, u.created_at, u.updated_at
		FROM users u
		INNER JOIN sessions s ON u.id = s.user_id
		WHERE s.token = ? AND s.expires_at > NOW()`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// InvalidateSession removes session from MySQL database
func (r *AuthRepository) InvalidateSession(ctx context.Context, token string) error {
	query := `DELETE FROM sessions WHERE token = ?`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}
