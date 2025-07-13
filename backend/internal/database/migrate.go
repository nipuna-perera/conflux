// Database migration management system
// Handles schema migrations for both MySQL and PostgreSQL
// Ensures database schema is up-to-date on application startup
package database

import (
	"database/sql"
	"fmt"
	"log"
)

// Migrator handles database schema migrations
type Migrator struct {
	db     *sql.DB
	dbType string
}

// NewMigrator creates a new migration manager
func NewMigrator(db *sql.DB, dbType string) *Migrator {
	return &Migrator{
		db:     db,
		dbType: dbType,
	}
}

// Up runs all pending migrations
func (m *Migrator) Up() error {
	log.Println("Running database migrations...")

	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Run migrations based on database type
	switch m.dbType {
	case "mysql":
		return m.runMySQLMigrations()
	case "postgres":
		return m.runPostgreSQLMigrations()
	default:
		return fmt.Errorf("unsupported database type: %s", m.dbType)
	}
}

// Down rolls back the last migration
func (m *Migrator) Down() error {
	log.Println("Rolling back last migration...")
	// Implementation for rollback would go here
	return nil
}

// createMigrationsTable creates the migrations tracking table
func (m *Migrator) createMigrationsTable() error {
	var query string

	switch m.dbType {
	case "mysql":
		query = `
			CREATE TABLE IF NOT EXISTS migrations (
				id INT AUTO_INCREMENT PRIMARY KEY,
				version VARCHAR(255) NOT NULL UNIQUE,
				applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`
	case "postgres":
		query = `
			CREATE TABLE IF NOT EXISTS migrations (
				id SERIAL PRIMARY KEY,
				version VARCHAR(255) NOT NULL UNIQUE,
				applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
			)`
	}

	_, err := m.db.Exec(query)
	return err
}

// runMySQLMigrations runs MySQL-specific migrations
func (m *Migrator) runMySQLMigrations() error {
	migrations := []struct {
		version string
		query   string
	}{
		{
			version: "001_create_users_table",
			query: `
				CREATE TABLE IF NOT EXISTS users (
					id INT AUTO_INCREMENT PRIMARY KEY,
					email VARCHAR(255) UNIQUE NOT NULL,
					password_hash VARCHAR(255) NOT NULL,
					first_name VARCHAR(100) NOT NULL,
					last_name VARCHAR(100) NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
					INDEX idx_email (email)
				)`,
		},
		{
			version: "002_create_sessions_table",
			query: `
				CREATE TABLE IF NOT EXISTS sessions (
					id INT AUTO_INCREMENT PRIMARY KEY,
					user_id INT NOT NULL,
					token VARCHAR(500) NOT NULL UNIQUE,
					expires_at TIMESTAMP NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
					INDEX idx_token (token),
					INDEX idx_user_id (user_id)
				)`,
		},
		{
			version: "004_seed_dev_user",
			query: `
				INSERT INTO users (email, password_hash, first_name, last_name, created_at, updated_at)
				SELECT 
					'dev@conflux.local',
					'$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
					'Dev',
					'User',
					NOW(),
					NOW()
				WHERE NOT EXISTS (
					SELECT 1 FROM users WHERE email = 'dev@conflux.local'
				)`,
		},
	}

	return m.runMigrations(migrations)
}

// runPostgreSQLMigrations runs PostgreSQL-specific migrations
func (m *Migrator) runPostgreSQLMigrations() error {
	migrations := []struct {
		version string
		query   string
	}{
		{
			version: "001_create_users_table",
			query: `
				CREATE TABLE IF NOT EXISTS users (
					id SERIAL PRIMARY KEY,
					email VARCHAR(255) UNIQUE NOT NULL,
					password_hash VARCHAR(255) NOT NULL,
					first_name VARCHAR(100) NOT NULL,
					last_name VARCHAR(100) NOT NULL,
					created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
					updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
				);
				
				CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
				
				CREATE OR REPLACE FUNCTION update_updated_at_column()
				RETURNS TRIGGER AS $$
				BEGIN
					NEW.updated_at = NOW();
					RETURN NEW;
				END;
				$$ language 'plpgsql';
				
				DROP TRIGGER IF EXISTS update_users_updated_at ON users;
				CREATE TRIGGER update_users_updated_at BEFORE UPDATE
					ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();`,
		},
		{
			version: "002_create_sessions_table",
			query: `
				CREATE TABLE IF NOT EXISTS sessions (
					id SERIAL PRIMARY KEY,
					user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
					token VARCHAR(500) NOT NULL UNIQUE,
					expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
					created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
				);
				
				CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token);
				CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);`,
		},
		{
			version: "004_seed_dev_user",
			query: `
				INSERT INTO users (email, password_hash, first_name, last_name, created_at, updated_at)
				SELECT 
					'dev@conflux.local',
					'$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
					'Dev',
					'User',
					NOW(),
					NOW()
				WHERE NOT EXISTS (
					SELECT 1 FROM users WHERE email = 'dev@conflux.local'
				)`,
		},
	}

	return m.runMigrations(migrations)
}

// runMigrations executes a list of migrations
func (m *Migrator) runMigrations(migrations []struct {
	version string
	query   string
}) error {
	for _, migration := range migrations {
		// Check if migration already applied
		var count int
		err := m.db.QueryRow("SELECT COUNT(*) FROM migrations WHERE version = ?", migration.version).Scan(&count)
		if err != nil && m.dbType == "postgres" {
			// PostgreSQL uses $1 instead of ?
			err = m.db.QueryRow("SELECT COUNT(*) FROM migrations WHERE version = $1", migration.version).Scan(&count)
		}
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if count > 0 {
			log.Printf("Migration %s already applied, skipping", migration.version)
			continue
		}

		// Run migration
		log.Printf("Applying migration: %s", migration.version)
		if _, err := m.db.Exec(migration.query); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.version, err)
		}

		// Record migration
		var insertQuery string
		if m.dbType == "mysql" {
			insertQuery = "INSERT INTO migrations (version) VALUES (?)"
			_, err = m.db.Exec(insertQuery, migration.version)
		} else {
			insertQuery = "INSERT INTO migrations (version) VALUES ($1)"
			_, err = m.db.Exec(insertQuery, migration.version)
		}

		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migration.version, err)
		}

		log.Printf("Successfully applied migration: %s", migration.version)
	}

	return nil
}
