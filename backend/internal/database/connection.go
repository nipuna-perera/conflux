// Database connection management and factory
// Abstracts database connection creation for MySQL and PostgreSQL
// Provides connection pooling and health checking capabilities
package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"configarr/internal/config"
)

// ConnectionFactory creates database connections based on configuration
type ConnectionFactory struct {
	config *config.Config
}

// NewConnectionFactory creates a new connection factory
func NewConnectionFactory(cfg *config.Config) *ConnectionFactory {
	return &ConnectionFactory{config: cfg}
}

// NewConnection creates a new database connection based on DB_TYPE
// Returns *sql.DB instance configured for either MySQL or PostgreSQL
func (cf *ConnectionFactory) NewConnection() (*sql.DB, error) {
	var dsn string
	var driverName string

	switch cf.config.DBType {
	case "mysql":
		driverName = "mysql"
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cf.config.DBUser,
			cf.config.DBPassword,
			cf.config.DBHost,
			cf.config.DBPort,
			cf.config.DBName,
		)
	case "postgres":
		driverName = "postgres"
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cf.config.DBHost,
			cf.config.DBPort,
			cf.config.DBUser,
			cf.config.DBPassword,
			cf.config.DBName,
		)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cf.config.DBType)
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db, nil
}

// HealthCheck verifies database connectivity
func (cf *ConnectionFactory) HealthCheck(db *sql.DB) error {
	return db.Ping()
}
