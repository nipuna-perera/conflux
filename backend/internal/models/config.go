// Configuration management data models
// Defines structures for config templates, user configurations, and version history
// Supports multiple config formats and validation schemas
package models

import (
	"time"
)

// ConfigFormat represents supported configuration formats
type ConfigFormat string

const (
	FormatYAML ConfigFormat = "yaml"
	FormatJSON ConfigFormat = "json"
	FormatTOML ConfigFormat = "toml"
	FormatENV  ConfigFormat = "env"
)

// ConfigTemplate represents a default configuration template for an application
type ConfigTemplate struct {
	ID               int              `json:"id" db:"id"`
	Name             string           `json:"name" db:"name"`                 // e.g., "cross-seed"
	DisplayName      string           `json:"display_name" db:"display_name"` // e.g., "Cross-Seed"
	Description      string           `json:"description" db:"description"`
	Version          string           `json:"version" db:"version"`     // Template version
	Category         string           `json:"category" db:"category"`   // e.g., "torrenting", "media"
	Format           ConfigFormat     `json:"format" db:"format"`       // Primary format
	SupportedFormats []ConfigFormat   `json:"supported_formats" db:"-"` // All supported formats
	DefaultContent   string           `json:"default_content" db:"default_content"`
	Schema           *string          `json:"schema,omitempty" db:"schema"` // JSON schema for validation
	Variables        []ConfigVariable `json:"variables" db:"-"`             // Template variables
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at" db:"updated_at"`
}

// ConfigVariable represents a variable in a configuration template
type ConfigVariable struct {
	ID             int     `json:"id" db:"id"`
	TemplateID     int     `json:"template_id" db:"template_id"`
	Name           string  `json:"name" db:"name"` // e.g., "DELAY"
	Path           string  `json:"path" db:"path"` // JSON path, e.g., "delay"
	Type           string  `json:"type" db:"type"` // "string", "number", "boolean", "array"
	Description    string  `json:"description" db:"description"`
	DefaultValue   *string `json:"default_value,omitempty" db:"default_value"`
	Required       bool    `json:"required" db:"required"`
	ValidationRule *string `json:"validation_rule,omitempty" db:"validation_rule"` // Regex or constraint
}

// UserConfig represents a user's configuration instance
type UserConfig struct {
	ID          int          `json:"id" db:"id"`
	UserID      int          `json:"user_id" db:"user_id"`
	TemplateID  *int         `json:"template_id,omitempty" db:"template_id"` // Null for custom configs
	Name        string       `json:"name" db:"name"`                         // User-defined name
	Description string       `json:"description" db:"description"`
	Format      ConfigFormat `json:"format" db:"format"`
	Content     string       `json:"content" db:"content"` // Current content
	IsShared    bool         `json:"is_shared" db:"is_shared"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`

	// Relationships
	Template *ConfigTemplate `json:"template,omitempty" db:"-"`
	Versions []ConfigVersion `json:"versions,omitempty" db:"-"`
}

// ConfigVersion represents a version in the configuration history
type ConfigVersion struct {
	ID         int       `json:"id" db:"id"`
	ConfigID   int       `json:"config_id" db:"config_id"`
	Version    int       `json:"version" db:"version"` // Incremental version number
	Content    string    `json:"content" db:"content"`
	ChangeNote string    `json:"change_note" db:"change_note"` // User-provided change description
	CreatedBy  int       `json:"created_by" db:"created_by"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// ConfigImport represents an import operation from external sources
type ConfigImport struct {
	ID           int              `json:"id" db:"id"`
	UserID       int              `json:"user_id" db:"user_id"`
	SourceType   ConfigSourceType `json:"source_type" db:"source_type"`
	SourceURL    string           `json:"source_url" db:"source_url"`
	Status       ImportStatus     `json:"status" db:"status"`
	ErrorMessage *string          `json:"error_message,omitempty" db:"error_message"`
	ConfigID     *int             `json:"config_id,omitempty" db:"config_id"` // Result config ID
	CreatedAt    time.Time        `json:"created_at" db:"created_at"`
	CompletedAt  *time.Time       `json:"completed_at,omitempty" db:"completed_at"`
}

// ConfigSourceType represents the source of an imported configuration
type ConfigSourceType string

const (
	SourceLocal  ConfigSourceType = "local"  // File upload
	SourceURL    ConfigSourceType = "url"    // Direct URL
	SourceGitHub ConfigSourceType = "github" // GitHub repository
	SourceGitLab ConfigSourceType = "gitlab" // GitLab repository
)

// ImportStatus represents the status of a configuration import
type ImportStatus string

const (
	ImportPending    ImportStatus = "pending"
	ImportProcessing ImportStatus = "processing"
	ImportCompleted  ImportStatus = "completed"
	ImportFailed     ImportStatus = "failed"
)

// ConfigDiff represents differences between two configuration versions
type ConfigDiff struct {
	LineNumber int    `json:"line_number"`
	Type       string `json:"type"` // "added", "removed", "modified"
	OldContent string `json:"old_content"`
	NewContent string `json:"new_content"`
}

// ShareRequest represents a request to share a configuration
type ShareRequest struct {
	ConfigID    int    `json:"config_id"`
	ShareWith   []int  `json:"share_with"`  // User IDs to share with
	Permissions string `json:"permissions"` // "read", "write"
}

// APIKey represents an API key for programmatic access
type APIKey struct {
	ID          int        `json:"id" db:"id"`
	UserID      int        `json:"user_id" db:"user_id"`
	Name        string     `json:"name" db:"name"`     // User-defined name
	KeyHash     string     `json:"-" db:"key_hash"`    // Hashed API key
	Permissions []string   `json:"permissions" db:"-"` // Stored as JSON in DB
	LastUsedAt  *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	IsActive    bool       `json:"is_active" db:"is_active"`
}
