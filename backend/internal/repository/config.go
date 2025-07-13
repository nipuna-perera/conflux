// Configuration repository interface
// Defines the contract for configuration data access layer
// Supports templates, user configs, versions, and imports
package repository

import (
	"conflux/internal/models"
)

// ConfigRepository defines the interface for configuration data access
type ConfigRepository interface {
	// Template management
	CreateTemplate(template *models.ConfigTemplate) error
	GetTemplate(id int) (*models.ConfigTemplate, error)
	GetTemplates(category, search string, page, limit int) ([]*models.ConfigTemplate, int64, error)
	UpdateTemplate(id int, updates *models.ConfigTemplate) error
	DeleteTemplate(id int) error

	// User configuration management
	CreateUserConfig(config *models.UserConfig) error
	GetUserConfig(id int) (*models.UserConfig, error)
	GetUserConfigs(userID int, templateID *int, page, limit int) ([]*models.UserConfig, int64, error)
	UpdateUserConfig(id int, config *models.UserConfig) error
	DeleteUserConfig(id int) error

	// Version management
	CreateVersion(version *models.ConfigVersion) error
	GetConfigVersion(id int) (*models.ConfigVersion, error)
	GetConfigVersions(configID int, page, limit int) ([]*models.ConfigVersion, int64, error)

	// Import management
	CreateImport(importRecord *models.ConfigImport) error
	GetImport(id int) (*models.ConfigImport, error)
	UpdateImport(id int, updates *models.ConfigImport) error

	// Template variables
	CreateVariable(variable *models.ConfigVariable) error
	GetTemplateVariables(templateID int) ([]*models.ConfigVariable, error)
	UpdateVariable(id int, updates *models.ConfigVariable) error
	DeleteVariable(id int) error
}
