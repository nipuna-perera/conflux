// MySQL implementation of UserRepository interface
// Handles user CRUD operations specific to MySQL database
// Implements SQL queries and transaction management for MySQL
package mysql

import (
	"context"
	"database/sql"
	"time"

	"configarr/internal/models"
)

// UserRepository implements repository.UserRepository for MySQL
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new MySQL user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user into MySQL database
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name) 
		VALUES (?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, user.Email, user.Password, user.FirstName, user.LastName)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return nil
}

// GetByID retrieves user by ID from MySQL
func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, created_at, updated_at 
		FROM users WHERE id = ?`

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

// GetByEmail retrieves user by email from MySQL
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, created_at, updated_at 
		FROM users WHERE email = ?`

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

// Update updates user information in MySQL
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET email = ?, first_name = ?, last_name = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, user.Email, user.FirstName, user.LastName, user.ID)
	return err
}

// Delete removes user from MySQL database
func (r *UserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
