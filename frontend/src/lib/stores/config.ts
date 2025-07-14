// Configuration state management store
// Manages templates, user configs, and configuration editor state
// Provides reactive data for configuration management components
import { writable, derived } from 'svelte/store';
import type {
	ConfigTemplate,
	UserConfig,
	ConfigVersion,
	ConfigFormat
} from '../types/config';
import { configAPI } from '../utils/configApi';

// Type for the editor store state
type EditorState = {
	config: UserConfig | null;
	content: string;
	format: ConfigFormat;
	isDirty: boolean;
	isValid: boolean;
	validationError: string | null;
	versions: ConfigVersion[];
	loading: boolean;
	saving: boolean;
	error: string | null;
};

// Templates store
function createTemplatesStore() {
	const { subscribe, set, update } = writable<{
		templates: ConfigTemplate[];
		loading: boolean;
		error: string | null;
		total: number;
		page: number;
		limit: number;
	}>({
		templates: [],
		loading: false,
		error: null,
		total: 0,
		page: 1,
		limit: 20
	});

	return {
		subscribe,
		async loadTemplates(params?: {
			category?: string;
			search?: string;
			page?: number;
			limit?: number;
		}) {
			update(state => ({ ...state, loading: true, error: null }));
			
			try {
				const response = await configAPI.getTemplates(params);
				update(state => ({
					...state,
					templates: response.templates,
					total: response.pagination.total,
					page: response.pagination.page,
					limit: response.pagination.limit,
					loading: false
				}));
			} catch (error) {
				update(state => ({
					...state,
					error: error instanceof Error ? error.message : 'Failed to load templates',
					loading: false
				}));
			}
		},
		async createTemplate(template: Omit<ConfigTemplate, 'id' | 'created_at' | 'updated_at'>) {
			try {
				const newTemplate = await configAPI.createTemplate(template);
				update(state => ({
					...state,
					templates: [newTemplate, ...state.templates]
				}));
				return newTemplate;
			} catch (error) {
				update(state => ({
					...state,
					error: error instanceof Error ? error.message : 'Failed to create template'
				}));
				throw error;
			}
		},
		clear: () => set({
			templates: [],
			loading: false,
			error: null,
			total: 0,
			page: 1,
			limit: 20
		})
	};
}

// User configurations store
function createConfigsStore() {
	const { subscribe, set, update } = writable<{
		configs: UserConfig[];
		loading: boolean;
		error: string | null;
		total: number;
		page: number;
		limit: number;
	}>({
		configs: [],
		loading: false,
		error: null,
		total: 0,
		page: 1,
		limit: 20
	});

	return {
		subscribe,
		async loadConfigs(params?: {
			template_id?: number;
			page?: number;
			limit?: number;
		}) {
			update(state => ({ ...state, loading: true, error: null }));
			
			try {
				const response = await configAPI.getUserConfigs(params);
				update(state => ({
					...state,
					configs: response.configs,
					total: response.pagination.total,
					page: response.pagination.page,
					limit: response.pagination.limit,
					loading: false
				}));
			} catch (error) {
				update(state => ({
					...state,
					error: error instanceof Error ? error.message : 'Failed to load configurations',
					loading: false
				}));
			}
		},
		async createConfig(templateId: number, name: string) {
			try {
				const newConfig = await configAPI.createUserConfig({ template_id: templateId, name });
				update(state => ({
					...state,
					configs: [newConfig, ...state.configs]
				}));
				return newConfig;
			} catch (error) {
				update(state => ({
					...state,
					error: error instanceof Error ? error.message : 'Failed to create configuration'
				}));
				throw error;
			}
		},
		async updateConfig(id: number, content: string, changeNote: string, format?: ConfigFormat) {
			try {
				const updatedConfig = await configAPI.updateUserConfig(id, {
					content,
					change_note: changeNote,
					format
				});
				
				update(state => ({
					...state,
					configs: state.configs.map(config => 
						config.id === id ? updatedConfig : config
					)
				}));
				return updatedConfig;
			} catch (error) {
				update(state => ({
					...state,
					error: error instanceof Error ? error.message : 'Failed to update configuration'
				}));
				throw error;
			}
		},
		async deleteConfig(id: number) {
			try {
				await configAPI.deleteUserConfig(id);
				update(state => ({
					...state,
					configs: state.configs.filter(config => config.id !== id)
				}));
			} catch (error) {
				update(state => ({
					...state,
					error: error instanceof Error ? error.message : 'Failed to delete configuration'
				}));
				throw error;
			}
		},
		clear: () => set({
			configs: [],
			loading: false,
			error: null,
			total: 0,
			page: 1,
			limit: 20
		})
	};
}

// Configuration editor store
function createEditorStore() {
	const { subscribe, set, update } = writable<EditorState>({
		config: null,
		content: '',
		format: 'yaml',
		isDirty: false,
		isValid: true,
		validationError: null,
		versions: [],
		loading: false,
		saving: false,
		error: null
	});

	return {
		subscribe,
		async loadConfig(id: number) {
			update(state => ({ ...state, loading: true, error: null }));
			
			try {
				const config = await configAPI.getUserConfig(id);
				update(state => ({
					...state,
					config,
					content: config.content,
					format: config.format,
					isDirty: false,
					loading: false
				}));
				
				// Load versions
				const versionsResponse = await configAPI.getConfigVersions(id);
				update(state => ({
					...state,
					versions: versionsResponse.versions
				}));
			} catch (error) {
				update(state => ({
					...state,
					error: error instanceof Error ? error.message : 'Failed to load configuration',
					loading: false
				}));
			}
		},
		updateContent(content: string) {
			update(state => ({
				...state,
				content,
				isDirty: content !== state.config?.content,
				validationError: null
			}));
		},
		async changeFormat(newFormat: ConfigFormat) {
			const state = await new Promise<EditorState>(resolve => {
				const unsubscribe = subscribe(resolve);
				unsubscribe();
			});
			
			if (state.format !== newFormat && state.content) {
				try {
					const result = await configAPI.convertFormat({
						content: state.content,
						from_format: state.format,
						to_format: newFormat
					});
					
					update(s => ({
						...s,
						content: result.content,
						format: newFormat,
						isDirty: true
					}));
				} catch (error) {
					update(s => ({
						...s,
						error: error instanceof Error ? error.message : 'Failed to convert format'
					}));
				}
			} else {
				update(s => ({ ...s, format: newFormat }));
			}
		},
		async detectFormat() {
			const state = await new Promise<EditorState>(resolve => {
				const unsubscribe = subscribe(resolve);
				unsubscribe();
			});
			
			if (state.content) {
				try {
					const result = await configAPI.detectFormat({ content: state.content });
					update(s => ({ ...s, format: result.format }));
				} catch (error) {
					update(s => ({
						...s,
						error: error instanceof Error ? error.message : 'Failed to detect format'
					}));
				}
			}
		},
		async validateConfig() {
			const state = await new Promise<EditorState>(resolve => {
				const unsubscribe = subscribe(resolve);
				unsubscribe();
			});
			
			if (state.content && state.config) {
				try {
					await configAPI.validateConfig({
						content: state.content,
						format: state.format,
						template_id: state.config.template_id || undefined
					});
					
					update(s => ({
						...s,
						isValid: true,
						validationError: null
					}));
				} catch (error) {
					update(s => ({
						...s,
						isValid: false,
						validationError: error instanceof Error ? error.message : 'Validation failed'
					}));
				}
			}
		},
		async saveConfig(changeNote: string) {
			const state = await new Promise<EditorState>(resolve => {
				const unsubscribe = subscribe(resolve);
				unsubscribe();
			});
			
			if (!state.config || !state.isDirty) return;
			
			update(s => ({ ...s, saving: true, error: null }));
			
			try {
				const updatedConfig = await configAPI.updateUserConfig(state.config.id, {
					content: state.content,
					change_note: changeNote,
					format: state.format
				});
				
				// Reload versions
				const versionsResponse = await configAPI.getConfigVersions(state.config.id);
				
				update(s => ({
					...s,
					config: updatedConfig,
					isDirty: false,
					saving: false,
					versions: versionsResponse.versions
				}));
			} catch (error) {
				update(s => ({
					...s,
					error: error instanceof Error ? error.message : 'Failed to save configuration',
					saving: false
				}));
				throw error;
			}
		},
		async restoreVersion(versionId: number) {
			const state = await new Promise<EditorState>(resolve => {
				const unsubscribe = subscribe(resolve);
				unsubscribe();
			});
			
			if (!state.config) return;
			
			try {
				const restoredConfig = await configAPI.restoreConfigVersion(state.config.id, versionId);
				const versionsResponse = await configAPI.getConfigVersions(state.config.id);
				
				update(s => ({
					...s,
					config: restoredConfig,
					content: restoredConfig.content,
					format: restoredConfig.format,
					isDirty: false,
					versions: versionsResponse.versions
				}));
			} catch (error) {
				update(s => ({
					...s,
					error: error instanceof Error ? error.message : 'Failed to restore version'
				}));
				throw error;
			}
		},
		clear: () => set({
			config: null,
			content: '',
			format: 'yaml',
			isDirty: false,
			isValid: true,
			validationError: null,
			versions: [],
			loading: false,
			saving: false,
			error: null
		})
	};
}

// Store instances
export const templatesStore = createTemplatesStore();
export const configsStore = createConfigsStore();
export const editorStore = createEditorStore();

// Derived stores
export const hasUnsavedChanges = derived(
	editorStore,
	$editor => $editor.isDirty
);
