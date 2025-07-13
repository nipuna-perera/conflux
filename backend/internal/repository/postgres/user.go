// PostgreSQL implementation of UserRepository interface
// Handles user CRUD operations specific to PostgreSQL database
// Implements SQL queries and transaction management for PostgreSQL
package postgres

import (
	"context"
	"database/sql"

	"conflux/internal/models"
)

// UserRepository implements repository.UserRepository for PostgreSQL
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user into PostgreSQL database
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, user.Email, user.Password, user.FirstName, user.LastName).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
	)

	return err
}

// GetByID retrieves user by ID from PostgreSQL
func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, created_at, updated_at 
		FROM users WHERE id = $1`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByEmail retrieves user by email from PostgreSQL
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, created_at, updated_at 
		FROM users WHERE email = $1`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// Update updates user information in PostgreSQL
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET email = $1, first_name = $2, last_name = $3, updated_at = NOW() 
		WHERE id = $4`

	_, err := r.db.ExecContext(ctx, query, user.Email, user.FirstName, user.LastName, user.ID)
	return err
}

// Delete removes user from PostgreSQL database
func (r *UserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
