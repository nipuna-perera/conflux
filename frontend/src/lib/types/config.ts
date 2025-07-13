// Configuration management types
// TypeScript interfaces for configuration templates, user configs, and related data
// Matches backend models for consistent API integration
export interface ConfigTemplate {
	id: number;
	name: string;
	display_name: string;
	description: string;
	version: string;
	category: string;
	format: ConfigFormat;
	supported_formats?: ConfigFormat[];
	default_content: string;
	schema?: string;
	variables?: ConfigVariable[];
	created_at: string;
	updated_at: string;
}

export interface ConfigVariable {
	id: number;
	template_id: number;
	name: string;
	path: string;
	type: string;
	description: string;
	default_value?: string;
	required: boolean;
	validation_rule?: string;
}

export interface UserConfig {
	id: number;
	user_id: number;
	template_id?: number;
	name: string;
	description: string;
	format: ConfigFormat;
	content: string;
	is_shared: boolean;
	created_at: string;
	updated_at: string;
	template?: ConfigTemplate;
	versions?: ConfigVersion[];
}

export interface ConfigVersion {
	id: number;
	config_id: number;
	version: number;
	content: string;
	change_note: string;
	created_by: number;
	created_at: string;
}

export interface ConfigImport {
	id: number;
	user_id: number;
	source_type: ConfigSourceType;
	source_url: string;
	status: ImportStatus;
	error_message?: string;
	config_id?: number;
	created_at: string;
	completed_at?: string;
}

export type ConfigFormat = 'yaml' | 'json' | 'toml' | 'env';

export type ConfigSourceType = 'local' | 'url' | 'github' | 'gitlab';

export type ImportStatus = 'pending' | 'processing' | 'completed' | 'failed';

export interface ConfigDiff {
	line_number: number;
	type: 'added' | 'removed' | 'modified';
	old_content: string;
	new_content: string;
}

export interface ShareRequest {
	config_id: number;
	share_with: number[];
	permissions: 'read' | 'write';
}

export interface APIKey {
	id: number;
	user_id: number;
	name: string;
	permissions: string[];
	last_used_at?: string;
	expires_at?: string;
	created_at: string;
	is_active: boolean;
}

// API request/response types
export interface CreateConfigRequest {
	template_id: number;
	name: string;
}

export interface UpdateConfigRequest {
	name?: string;
	description?: string;
	content?: string;
	change_note?: string;
	format?: ConfigFormat;
}

export interface ValidateConfigRequest {
	content: string;
	format: ConfigFormat;
	template_id?: number;
}

export interface ConvertFormatRequest {
	content: string;
	from_format: ConfigFormat;
	to_format: ConfigFormat;
}

export interface DetectFormatRequest {
	content: string;
}

export interface PaginatedResponse<T> {
	data: T[];
	pagination: {
		page: number;
		limit: number;
		total: number;
	};
}

export interface TemplatesResponse {
	templates: ConfigTemplate[];
	pagination: {
		page: number;
		limit: number;
		total: number;
	};
}

export interface ConfigsResponse {
	configs: UserConfig[];
	pagination: {
		page: number;
		limit: number;
		total: number;
	};
}

export interface VersionsResponse {
	versions: ConfigVersion[];
	pagination: {
		page: number;
		limit: number;
		total: number;
	};
}

// Utility types
export interface FormatDetectionResult {
	format: ConfigFormat;
}

export interface ConversionResult {
	content: string;
}

export interface ValidationResult {
	message: string;
}
