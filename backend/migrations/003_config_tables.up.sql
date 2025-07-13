-- Configuration management schema migration
-- Creates tables for templates, user configs, versions, and imports
-- Supports multiple formats and version history

-- Configuration templates table
CREATE TABLE IF NOT EXISTS config_templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    version VARCHAR(50) NOT NULL DEFAULT '1.0.0',
    category VARCHAR(100) NOT NULL,
    format VARCHAR(20) NOT NULL CHECK (format IN ('yaml', 'json', 'toml', 'env')),
    default_content TEXT NOT NULL,
    schema TEXT, -- JSON schema for validation
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster template lookups
CREATE INDEX IF NOT EXISTS idx_config_templates_category ON config_templates(category);
CREATE INDEX IF NOT EXISTS idx_config_templates_name ON config_templates(name);

-- Configuration template variables table
CREATE TABLE IF NOT EXISTS config_variables (
    id SERIAL PRIMARY KEY,
    template_id INTEGER NOT NULL REFERENCES config_templates(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    path VARCHAR(500) NOT NULL, -- JSON path like "delay" or "torrentDir"
    type VARCHAR(50) NOT NULL DEFAULT 'string',
    description TEXT,
    default_value TEXT,
    required BOOLEAN DEFAULT false,
    validation_rule TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for template variables
CREATE INDEX IF NOT EXISTS idx_config_variables_template ON config_variables(template_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_config_variables_template_name ON config_variables(template_id, name);

-- User configurations table
CREATE TABLE IF NOT EXISTS user_configs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    template_id INTEGER REFERENCES config_templates(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    format VARCHAR(20) NOT NULL CHECK (format IN ('yaml', 'json', 'toml', 'env')),
    content TEXT NOT NULL,
    is_shared BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure unique config names per user
    UNIQUE(user_id, name)
);

-- Indexes for user configs
CREATE INDEX IF NOT EXISTS idx_user_configs_user ON user_configs(user_id);
CREATE INDEX IF NOT EXISTS idx_user_configs_template ON user_configs(template_id);
CREATE INDEX IF NOT EXISTS idx_user_configs_shared ON user_configs(is_shared) WHERE is_shared = true;

-- Configuration versions table for version history
CREATE TABLE IF NOT EXISTS config_versions (
    id SERIAL PRIMARY KEY,
    config_id INTEGER NOT NULL REFERENCES user_configs(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    content TEXT NOT NULL,
    change_note TEXT,
    created_by INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure unique version numbers per config
    UNIQUE(config_id, version)
);

-- Index for config versions
CREATE INDEX IF NOT EXISTS idx_config_versions_config ON config_versions(config_id, version DESC);
CREATE INDEX IF NOT EXISTS idx_config_versions_created_by ON config_versions(created_by);

-- Configuration imports table
CREATE TABLE IF NOT EXISTS config_imports (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source_type VARCHAR(50) NOT NULL CHECK (source_type IN ('local', 'url', 'github', 'gitlab')),
    source_url TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    error_message TEXT,
    config_id INTEGER REFERENCES user_configs(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE
);

-- Index for imports
CREATE INDEX IF NOT EXISTS idx_config_imports_user ON config_imports(user_id);
CREATE INDEX IF NOT EXISTS idx_config_imports_status ON config_imports(status);

-- API keys table for programmatic access
CREATE TABLE IF NOT EXISTS api_keys (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    key_hash VARCHAR(255) NOT NULL UNIQUE,
    permissions JSON NOT NULL DEFAULT '[]',
    last_used_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for API keys
CREATE INDEX IF NOT EXISTS idx_api_keys_user ON api_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_hash ON api_keys(key_hash) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_api_keys_expires ON api_keys(expires_at) WHERE expires_at IS NOT NULL;

-- Insert default cross-seed template
INSERT INTO config_templates (name, display_name, description, category, format, default_content) VALUES
('cross-seed', 'Cross-Seed', 'Automatic cross-seeding configuration for torrent clients', 'torrenting', 'yaml', 
'# Cross-seed configuration
delay: 30
outputDir: "/downloads/torrents"
torrentDir: "/watch/folders"
duplicateCategories: true

# Torrent client settings
torznab:
  - name: "prowlarr"
    url: "http://prowlarr:9696/1/api"
    apikey: "your-api-key"

# Action settings
action: "inject"
includeEpisodes: false
includeSingleEpisodes: true
includeNonVideos: false

# Matching settings
matchMode: "safe"
skipRecheck: false
maxDataDepth: 1

# Logging
verbose: false')
ON CONFLICT (name) DO NOTHING;

-- Create a function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for updated_at
CREATE TRIGGER update_config_templates_updated_at
    BEFORE UPDATE ON config_templates
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_configs_updated_at
    BEFORE UPDATE ON user_configs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
