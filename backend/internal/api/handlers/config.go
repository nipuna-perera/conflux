// Configuration API handlers
// Provides REST endpoints for managing configuration templates and user configs
// Supports CRUD operations, version management, and format conversion
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"conflux/internal/models"
	"conflux/internal/service"
	"conflux/pkg/utils"

	"github.com/gorilla/mux"
)

// ConfigHandler handles configuration-related HTTP requests
type ConfigHandler struct {
	configService *service.ConfigService
}

// NewConfigHandler creates a new configuration handler
func NewConfigHandler(configService *service.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		configService: configService,
	}
}

// Template Endpoints

// GetTemplates handles GET /api/templates
func (h *ConfigHandler) GetTemplates(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	category := r.URL.Query().Get("category")
	search := r.URL.Query().Get("search")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	templates, total, err := h.configService.GetTemplates(category, search, page, limit)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve templates")
		return
	}

	response := map[string]interface{}{
		"templates": templates,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	utils.JSONResponse(w, http.StatusOK, response)
}

// GetTemplate handles GET /api/templates/{id}
func (h *ConfigHandler) GetTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid template ID")
		return
	}

	template, err := h.configService.GetTemplate(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Template not found")
		return
	}

	utils.JSONResponse(w, http.StatusOK, template)
}

// CreateTemplate handles POST /api/templates
func (h *ConfigHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var template models.ConfigTemplate
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.configService.CreateTemplate(&template); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Failed to create template: "+err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusCreated, template)
}

// UpdateTemplate handles PUT /api/templates/{id}
func (h *ConfigHandler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid template ID")
		return
	}

	var updates models.ConfigTemplate
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.configService.UpdateTemplate(id, &updates); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Failed to update template: "+err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "Template updated successfully"})
}

// DeleteTemplate handles DELETE /api/templates/{id}
func (h *ConfigHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid template ID")
		return
	}

	if err := h.configService.DeleteTemplate(id); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete template")
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "Template deleted successfully"})
}

// User Configuration Endpoints

// GetUserConfigs handles GET /api/configs
func (h *ConfigHandler) GetUserConfigs(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := getUserIDFromContext(r)
	if userID == 0 {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse query parameters
	var templateID *int
	if templateIDStr := r.URL.Query().Get("template_id"); templateIDStr != "" {
		if tid, err := strconv.Atoi(templateIDStr); err == nil {
			templateID = &tid
		}
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	configs, total, err := h.configService.GetUserConfigs(userID, templateID, page, limit)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve configurations")
		return
	}

	response := map[string]interface{}{
		"configs": configs,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	utils.JSONResponse(w, http.StatusOK, response)
}

// GetUserConfig handles GET /api/configs/{id}
func (h *ConfigHandler) GetUserConfig(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid configuration ID")
		return
	}

	config, err := h.configService.GetUserConfig(id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			utils.ErrorResponse(w, http.StatusForbidden, "Unauthorized access")
		} else {
			utils.ErrorResponse(w, http.StatusNotFound, "Configuration not found")
		}
		return
	}

	utils.JSONResponse(w, http.StatusOK, config)
}

// CreateUserConfig handles POST /api/configs
func (h *ConfigHandler) CreateUserConfig(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		TemplateID int    `json:"template_id"`
		Name       string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Configuration name is required")
		return
	}

	config, err := h.configService.CreateUserConfig(userID, req.TemplateID, req.Name)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Failed to create configuration: "+err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusCreated, config)
}

// UpdateUserConfig handles PUT /api/configs/{id}
func (h *ConfigHandler) UpdateUserConfig(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid configuration ID")
		return
	}

	var req struct {
		Content    string               `json:"content"`
		ChangeNote string               `json:"change_note"`
		Format     *models.ConfigFormat `json:"format,omitempty"`
	}

	if decodeErr := json.NewDecoder(r.Body).Decode(&req); decodeErr != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Content == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Content is required")
		return
	}

	config, err := h.configService.UpdateUserConfig(id, userID, req.Content, req.ChangeNote, req.Format)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			utils.ErrorResponse(w, http.StatusForbidden, "Unauthorized access")
		} else {
			utils.ErrorResponse(w, http.StatusBadRequest, "Failed to update configuration: "+err.Error())
		}
		return
	}

	utils.JSONResponse(w, http.StatusOK, config)
}

// DeleteUserConfig handles DELETE /api/configs/{id}
func (h *ConfigHandler) DeleteUserConfig(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid configuration ID")
		return
	}

	if err := h.configService.DeleteUserConfig(id, userID); err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			utils.ErrorResponse(w, http.StatusForbidden, "Unauthorized access")
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete configuration")
		}
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "Configuration deleted successfully"})
}

// Version Management Endpoints

// GetConfigVersions handles GET /api/configs/{id}/versions
func (h *ConfigHandler) GetConfigVersions(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	configID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid configuration ID")
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	versions, total, err := h.configService.GetConfigVersions(configID, userID, page, limit)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			utils.ErrorResponse(w, http.StatusForbidden, "Unauthorized access")
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve versions")
		}
		return
	}

	response := map[string]interface{}{
		"versions": versions,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	utils.JSONResponse(w, http.StatusOK, response)
}

// RestoreConfigVersion handles POST /api/configs/{id}/versions/{version_id}/restore
func (h *ConfigHandler) RestoreConfigVersion(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	configID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid configuration ID")
		return
	}

	versionID, err := strconv.Atoi(vars["version_id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid version ID")
		return
	}

	config, err := h.configService.RestoreConfigVersion(configID, versionID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			utils.ErrorResponse(w, http.StatusForbidden, "Unauthorized access")
		} else {
			utils.ErrorResponse(w, http.StatusBadRequest, "Failed to restore version: "+err.Error())
		}
		return
	}

	utils.JSONResponse(w, http.StatusOK, config)
}

// Utility Endpoints

// DetectFormat handles POST /api/configs/detect-format
func (h *ConfigHandler) DetectFormat(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	format, err := h.configService.DetectFormat(req.Content)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Unable to detect format: "+err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"format": string(format)})
}

// ConvertFormat handles POST /api/configs/convert
func (h *ConfigHandler) ConvertFormat(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Content    string              `json:"content"`
		FromFormat models.ConfigFormat `json:"from_format"`
		ToFormat   models.ConfigFormat `json:"to_format"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	converted, err := h.configService.ConvertFormat(req.Content, req.FromFormat, req.ToFormat)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Conversion failed: "+err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"content": converted})
}

// ValidateConfig handles POST /api/configs/validate
func (h *ConfigHandler) ValidateConfig(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Content    string              `json:"content"`
		Format     models.ConfigFormat `json:"format"`
		TemplateID *int                `json:"template_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.configService.ValidateConfig(req.Content, req.Format, req.TemplateID); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "Configuration is valid"})
}

// ExportConfig handles GET /api/configs/{id}/export?format=yaml
func (h *ConfigHandler) ExportConfig(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	configID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid configuration ID")
		return
	}

	format := models.ConfigFormat(r.URL.Query().Get("format"))
	if format == "" {
		format = models.FormatYAML // Default format
	}

	content, err := h.configService.ExportConfig(configID, userID, format)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			utils.ErrorResponse(w, http.StatusForbidden, "Unauthorized access")
		} else {
			utils.ErrorResponse(w, http.StatusBadRequest, "Export failed: "+err.Error())
		}
		return
	}

	// Set appropriate content type
	contentType := "text/plain"
	switch format {
	case models.FormatJSON:
		contentType = "application/json"
	case models.FormatYAML:
		contentType = "application/x-yaml"
	case models.FormatTOML:
		contentType = "application/toml"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename=config."+string(format))
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(content)); err != nil {
		// Log error - response headers are already written so we can't send error response
		// In production, you might want to log this error properly
		return
	}
}

// Helper function to extract user ID from request context
func getUserIDFromContext(r *http.Request) int {
	if userID, ok := r.Context().Value("user_id").(int); ok {
		return userID
	}
	return 0
}
