// Configuration API client
// Handles all configuration-related API calls with type safety
// Provides methods for templates, user configs, versions, and utilities
import { apiClient } from './api';
import type {
	ConfigTemplate,
	UserConfig,
	ConfigVersion,
	ConfigImport,
	CreateConfigRequest,
	UpdateConfigRequest,
	ValidateConfigRequest,
	ConvertFormatRequest,
	DetectFormatRequest,
	TemplatesResponse,
	ConfigsResponse,
	VersionsResponse,
	FormatDetectionResult,
	ConversionResult,
	ValidationResult,
	ConfigFormat,
	ConfigSourceType
} from '../types/config';

class ConfigAPI {
	// Template Management
	async getTemplates(params?: {
		category?: string;
		search?: string;
		page?: number;
		limit?: number;
	}): Promise<TemplatesResponse> {
		const searchParams = new URLSearchParams();
		if (params?.category) searchParams.set('category', params.category);
		if (params?.search) searchParams.set('search', params.search);
		if (params?.page) searchParams.set('page', params.page.toString());
		if (params?.limit) searchParams.set('limit', params.limit.toString());

		const url = `/templates${searchParams.toString() ? '?' + searchParams.toString() : ''}`;
		return apiClient.request(url);
	}

	async getTemplate(id: number): Promise<ConfigTemplate> {
		return apiClient.request(`/templates/${id}`);
	}

	async createTemplate(template: Omit<ConfigTemplate, 'id' | 'created_at' | 'updated_at'>): Promise<ConfigTemplate> {
		return apiClient.request('/templates', {
			method: 'POST',
			body: JSON.stringify(template)
		});
	}

	async updateTemplate(id: number, updates: Partial<ConfigTemplate>): Promise<{ message: string }> {
		return apiClient.request(`/templates/${id}`, {
			method: 'PUT',
			body: JSON.stringify(updates)
		});
	}

	async deleteTemplate(id: number): Promise<{ message: string }> {
		return apiClient.request(`/templates/${id}`, {
			method: 'DELETE'
		});
	}

	// User Configuration Management
	async getUserConfigs(params?: {
		template_id?: number;
		page?: number;
		limit?: number;
	}): Promise<ConfigsResponse> {
		const searchParams = new URLSearchParams();
		if (params?.template_id) searchParams.set('template_id', params.template_id.toString());
		if (params?.page) searchParams.set('page', params.page.toString());
		if (params?.limit) searchParams.set('limit', params.limit.toString());

		const url = `/configs${searchParams.toString() ? '?' + searchParams.toString() : ''}`;
		return apiClient.request(url);
	}

	async getUserConfig(id: number): Promise<UserConfig> {
		return apiClient.request(`/configs/${id}`);
	}

	async createUserConfig(request: CreateConfigRequest): Promise<UserConfig> {
		return apiClient.request('/configs', {
			method: 'POST',
			body: JSON.stringify(request)
		});
	}

	async updateUserConfig(id: number, request: UpdateConfigRequest): Promise<UserConfig> {
		return apiClient.request(`/configs/${id}`, {
			method: 'PUT',
			body: JSON.stringify(request)
		});
	}

	async deleteUserConfig(id: number): Promise<{ message: string }> {
		return apiClient.request(`/configs/${id}`, {
			method: 'DELETE'
		});
	}

	// Version Management
	async getConfigVersions(configId: number, params?: {
		page?: number;
		limit?: number;
	}): Promise<VersionsResponse> {
		const searchParams = new URLSearchParams();
		if (params?.page) searchParams.set('page', params.page.toString());
		if (params?.limit) searchParams.set('limit', params.limit.toString());

		const url = `/configs/${configId}/versions${searchParams.toString() ? '?' + searchParams.toString() : ''}`;
		return apiClient.request(url);
	}

	async restoreConfigVersion(configId: number, versionId: number): Promise<UserConfig> {
		return apiClient.request(`/configs/${configId}/versions/${versionId}/restore`, {
			method: 'POST'
		});
	}

	// Utility Functions
	async detectFormat(request: DetectFormatRequest): Promise<FormatDetectionResult> {
		return apiClient.request('/configs/detect-format', {
			method: 'POST',
			body: JSON.stringify(request)
		});
	}

	async convertFormat(request: ConvertFormatRequest): Promise<ConversionResult> {
		return apiClient.request('/configs/convert', {
			method: 'POST',
			body: JSON.stringify(request)
		});
	}

	async validateConfig(request: ValidateConfigRequest): Promise<ValidationResult> {
		return apiClient.request('/configs/validate', {
			method: 'POST',
			body: JSON.stringify(request)
		});
	}

	async exportConfig(configId: number, format: ConfigFormat): Promise<Blob> {
		const token = localStorage.getItem('auth_token');
		const response = await fetch(`/api/configs/${configId}/export?format=${format}`, {
			method: 'GET',
			headers: {
				'Authorization': `Bearer ${token}`
			}
		});

		if (!response.ok) {
			throw new Error(`Export failed: ${response.statusText}`);
		}

		return response.blob();
	}

	// Import Management
	async importConfig(sourceType: ConfigSourceType, sourceUrl: string): Promise<ConfigImport> {
		return apiClient.request('/configs/import', {
			method: 'POST',
			body: JSON.stringify({
				source_type: sourceType,
				source_url: sourceUrl
			})
		});
	}

	async getImportStatus(importId: number): Promise<ConfigImport> {
		return apiClient.request(`/configs/imports/${importId}`);
	}
}

export const configAPI = new ConfigAPI();
