// Configuration management service
// Handles CRUD operations for configuration templates and user configurations
// Provides version management, import/export, and validation capabilities
package service

import (
	"fmt"
	"time"

	"conflux/internal/models"
	"conflux/pkg/config"
)

// ConfigService provides configuration management functionality
type ConfigService struct {
	configRepo ConfigRepository
	parser     *config.Parser
}

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
}

// NewConfigService creates a new configuration service
func NewConfigService(configRepo ConfigRepository) *ConfigService {
	return &ConfigService{
		configRepo: configRepo,
		parser:     config.NewParser(),
	}
}

// Template Management

// CreateTemplate creates a new configuration template
func (s *ConfigService) CreateTemplate(template *models.ConfigTemplate) error {
	// Validate the template content
	if err := s.validateTemplateContent(template); err != nil {
		return fmt.Errorf("template validation failed: %w", err)
	}

	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()

	return s.configRepo.CreateTemplate(template)
}

// GetTemplate retrieves a configuration template by ID
func (s *ConfigService) GetTemplate(id int) (*models.ConfigTemplate, error) {
	return s.configRepo.GetTemplate(id)
}

// GetTemplates retrieves all configuration templates with optional filtering
func (s *ConfigService) GetTemplates(category, search string, page, limit int) ([]*models.ConfigTemplate, int64, error) {
	return s.configRepo.GetTemplates(category, search, page, limit)
}

// UpdateTemplate updates an existing configuration template
func (s *ConfigService) UpdateTemplate(id int, updates *models.ConfigTemplate) error {
	existing, err := s.configRepo.GetTemplate(id)
	if err != nil {
		return err
	}

	// Validate updated content
	if updates.DefaultContent != "" {
		updates.Format = existing.Format // Keep original format
		if err := s.validateTemplateContent(updates); err != nil {
			return fmt.Errorf("template validation failed: %w", err)
		}
	}

	updates.UpdatedAt = time.Now()
	return s.configRepo.UpdateTemplate(id, updates)
}

// DeleteTemplate deletes a configuration template
func (s *ConfigService) DeleteTemplate(id int) error {
	return s.configRepo.DeleteTemplate(id)
}

// User Configuration Management

// CreateUserConfig creates a new user configuration from a template
func (s *ConfigService) CreateUserConfig(userID int, templateID int, name string) (*models.UserConfig, error) {
	template, err := s.configRepo.GetTemplate(templateID)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	userConfig := &models.UserConfig{
		UserID:     userID,
		TemplateID: &templateID,
		Name:       name,
		Content:    template.DefaultContent,
		Format:     template.Format,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.configRepo.CreateUserConfig(userConfig); err != nil {
		return nil, err
	}

	// Create initial version
	if err := s.createConfigVersion(userConfig, "Initial version"); err != nil {
		return nil, fmt.Errorf("failed to create initial version: %w", err)
	}

	return userConfig, nil
}

// GetUserConfig retrieves a user configuration by ID
func (s *ConfigService) GetUserConfig(id int, userID int) (*models.UserConfig, error) {
	config, err := s.configRepo.GetUserConfig(id)
	if err != nil {
		return nil, err
	}

	if config.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to configuration")
	}

	return config, nil
}

// GetUserConfigs retrieves all configurations for a user
func (s *ConfigService) GetUserConfigs(userID int, templateID *int, page, limit int) ([]*models.UserConfig, int64, error) {
	return s.configRepo.GetUserConfigs(userID, templateID, page, limit)
}

// UpdateUserConfig updates a user configuration and creates a new version
func (s *ConfigService) UpdateUserConfig(id int, userID int, content, changeNote string, format *models.ConfigFormat) (*models.UserConfig, error) {
	config, err := s.GetUserConfig(id, userID)
	if err != nil {
		return nil, err
	}

	// Validate new content
	actualFormat := config.Format
	if format != nil {
		actualFormat = *format
	}

	if err := s.validateConfigContent(content, actualFormat); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Update configuration
	config.Content = content
	if format != nil {
		config.Format = *format
	}
	config.UpdatedAt = time.Now()

	if err := s.configRepo.UpdateUserConfig(id, config); err != nil {
		return nil, err
	}

	// Create new version
	if err := s.createConfigVersion(config, changeNote); err != nil {
		return nil, fmt.Errorf("failed to create version: %w", err)
	}

	return config, nil
}

// DeleteUserConfig deletes a user configuration
func (s *ConfigService) DeleteUserConfig(id int, userID int) error {
	config, err := s.GetUserConfig(id, userID)
	if err != nil {
		return err
	}

	return s.configRepo.DeleteUserConfig(config.ID)
}

// Version Management

// GetConfigVersions retrieves version history for a configuration
func (s *ConfigService) GetConfigVersions(configID int, userID int, page, limit int) ([]*models.ConfigVersion, int64, error) {
	// Verify user owns the configuration
	if _, err := s.GetUserConfig(configID, userID); err != nil {
		return nil, 0, err
	}

	return s.configRepo.GetConfigVersions(configID, page, limit)
}

// GetConfigVersion retrieves a specific version
func (s *ConfigService) GetConfigVersion(versionID int, userID int) (*models.ConfigVersion, error) {
	version, err := s.configRepo.GetConfigVersion(versionID)
	if err != nil {
		return nil, err
	}

	// Verify user owns the configuration
	if _, err := s.GetUserConfig(version.ConfigID, userID); err != nil {
		return nil, err
	}

	return version, nil
}

// RestoreConfigVersion restores a configuration to a previous version
func (s *ConfigService) RestoreConfigVersion(configID int, versionID int, userID int) (*models.UserConfig, error) {
	// Verify user owns the configuration
	if _, err := s.GetUserConfig(configID, userID); err != nil {
		return nil, err
	}

	version, err := s.GetConfigVersion(versionID, userID)
	if err != nil {
		return nil, err
	}

	if version.ConfigID != configID {
		return nil, fmt.Errorf("version does not belong to this configuration")
	}

	// Update configuration with version content
	return s.UpdateUserConfig(configID, userID, version.Content, fmt.Sprintf("Restored to version %d", version.Version), nil)
}

// Format Detection and Conversion

// DetectFormat automatically detects the format of configuration content
func (s *ConfigService) DetectFormat(content string) (models.ConfigFormat, error) {
	return s.parser.DetectFormat(content)
}

// ConvertFormat converts configuration from one format to another
func (s *ConfigService) ConvertFormat(content string, fromFormat, toFormat models.ConfigFormat) (string, error) {
	return s.parser.ConvertFormat(content, fromFormat, toFormat)
}

// ValidateConfig validates configuration content
func (s *ConfigService) ValidateConfig(content string, format models.ConfigFormat, templateID *int) error {
	// Basic format validation
	if err := s.validateConfigContent(content, format); err != nil {
		return err
	}

	// Template-specific validation if provided
	if templateID != nil {
		template, err := s.configRepo.GetTemplate(*templateID)
		if err != nil {
			return err
		}

		// Use template schema if available
		return s.parser.ValidateConfig(content, format, template.Schema)
	}

	return nil
}

// ImportConfig imports configuration from external source
func (s *ConfigService) ImportConfig(userID int, sourceType models.ConfigSourceType, sourceURL string) (*models.ConfigImport, error) {
	// Create import record
	importRecord := &models.ConfigImport{
		UserID:     userID,
		SourceType: sourceType,
		SourceURL:  sourceURL,
		Status:     models.ImportPending,
		CreatedAt:  time.Now(),
	}

	if err := s.configRepo.CreateImport(importRecord); err != nil {
		return nil, err
	}

	// TODO: Process import asynchronously
	// This would handle URL fetching, Git clone, file upload, etc.

	return importRecord, nil
}

// ExportConfig exports configuration in specified format
func (s *ConfigService) ExportConfig(configID int, userID int, format models.ConfigFormat) (string, error) {
	config, err := s.GetUserConfig(configID, userID)
	if err != nil {
		return "", err
	}

	if config.Format == format {
		return config.Content, nil
	}

	return s.parser.ConvertFormat(config.Content, config.Format, format)
}

// Private helper methods

func (s *ConfigService) validateTemplateContent(template *models.ConfigTemplate) error {
	return s.validateConfigContent(template.DefaultContent, template.Format)
}

func (s *ConfigService) validateConfigContent(content string, format models.ConfigFormat) error {
	_, err := s.parser.ParseConfig(content, format)
	return err
}

func (s *ConfigService) createConfigVersion(config *models.UserConfig, changeNote string) error {
	// Get the next version number
	versions, _, err := s.configRepo.GetConfigVersions(config.ID, 1, 1)
	if err != nil {
		return err
	}

	versionNumber := 1
	if len(versions) > 0 {
		versionNumber = versions[0].Version + 1
	}

	version := &models.ConfigVersion{
		ConfigID:   config.ID,
		Version:    versionNumber,
		Content:    config.Content,
		ChangeNote: changeNote,
		CreatedBy:  config.UserID,
		CreatedAt:  time.Now(),
	}

	return s.configRepo.CreateVersion(version)
}
