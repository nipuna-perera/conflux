-- Drop configuration management tables

-- Drop triggers first
DROP TRIGGER IF EXISTS update_user_configs_updated_at ON user_configs;
DROP TRIGGER IF EXISTS update_config_templates_updated_at ON config_templates;

-- Drop the update function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS config_imports;
DROP TABLE IF EXISTS config_versions;
DROP TABLE IF EXISTS user_configs;
DROP TABLE IF EXISTS config_variables;
DROP TABLE IF EXISTS config_templates;
