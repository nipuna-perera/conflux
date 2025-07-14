// Configuration parser and format detection utilities
// Automatically detects config format and provides parsing/validation capabilities
// Supports YAML, JSON, TOML, and ENV formats with validation
package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"conflux/internal/models"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// Parser handles configuration parsing and format detection
type Parser struct{}

// NewParser creates a new configuration parser
func NewParser() *Parser {
	return &Parser{}
}

// DetectFormat attempts to automatically detect the configuration format
func (p *Parser) DetectFormat(content string) (models.ConfigFormat, error) {
	content = strings.TrimSpace(content)

	if content == "" {
		return "", fmt.Errorf("empty content")
	}

	// Try JSON first (most strict)
	if p.isValidJSON(content) {
		return models.FormatJSON, nil
	}

	// Try YAML
	if p.isValidYAML(content) {
		return models.FormatYAML, nil
	}

	// Try TOML
	if p.isValidTOML(content) {
		return models.FormatTOML, nil
	}

	// Check if it looks like ENV format
	if p.looksLikeEnv(content) {
		return models.FormatENV, nil
	}

	return "", fmt.Errorf("unable to detect configuration format")
}

// ParseConfig parses configuration content based on the specified format
func (p *Parser) ParseConfig(content string, format models.ConfigFormat) (map[string]interface{}, error) {
	switch format {
	case models.FormatJSON:
		return p.parseJSON(content)
	case models.FormatYAML:
		return p.parseYAML(content)
	case models.FormatTOML:
		return p.parseTOML(content)
	case models.FormatENV:
		return p.parseEnv(content)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// ConvertFormat converts configuration from one format to another
func (p *Parser) ConvertFormat(content string, fromFormat, toFormat models.ConfigFormat) (string, error) {
	// Parse the source format
	data, err := p.ParseConfig(content, fromFormat)
	if err != nil {
		return "", fmt.Errorf("failed to parse source format: %w", err)
	}

	// Convert to target format
	return p.SerializeConfig(data, toFormat)
}

// SerializeConfig serializes configuration data to the specified format
func (p *Parser) SerializeConfig(data map[string]interface{}, format models.ConfigFormat) (string, error) {
	switch format {
	case models.FormatJSON:
		return p.serializeJSON(data)
	case models.FormatYAML:
		return p.serializeYAML(data)
	case models.FormatTOML:
		return p.serializeTOML(data)
	case models.FormatENV:
		return p.serializeEnv(data)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

// ValidateConfig validates configuration against a JSON schema if provided
func (p *Parser) ValidateConfig(content string, format models.ConfigFormat, schema *string) error {
	// Parse the configuration
	data, err := p.ParseConfig(content, format)
	if err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// If no schema provided, just check if it's parseable
	if schema == nil {
		return nil
	}

	// TODO: Implement JSON schema validation
	// This would use a library like github.com/xeipuuv/gojsonschema
	_ = data
	return nil
}

// Private helper methods

func (p *Parser) isValidJSON(content string) bool {
	var js interface{}
	return json.Unmarshal([]byte(content), &js) == nil
}

func (p *Parser) isValidYAML(content string) bool {
	var yml interface{}
	return yaml.Unmarshal([]byte(content), &yml) == nil
}

func (p *Parser) isValidTOML(content string) bool {
	var tml interface{}
	return toml.Unmarshal([]byte(content), &tml) == nil
}

func (p *Parser) looksLikeEnv(content string) bool {
	lines := strings.Split(content, "\n")
	validLines := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "=") {
			validLines++
		} else {
			return false
		}
	}

	return validLines > 0
}

func (p *Parser) parseJSON(content string) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(content), &data)
	return data, err
}

func (p *Parser) parseYAML(content string) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := yaml.Unmarshal([]byte(content), &data)
	return data, err
}

func (p *Parser) parseTOML(content string) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := toml.Unmarshal([]byte(content), &data)
	return data, err
}

func (p *Parser) parseEnv(content string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid env line: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}

		data[key] = value
	}

	return data, nil
}

func (p *Parser) serializeJSON(data map[string]interface{}) (string, error) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	return string(bytes), err
}

func (p *Parser) serializeYAML(data map[string]interface{}) (string, error) {
	bytes, err := yaml.Marshal(data)
	return string(bytes), err
}

func (p *Parser) serializeTOML(data map[string]interface{}) (string, error) {
	var buf strings.Builder
	err := toml.NewEncoder(&buf).Encode(data)
	return buf.String(), err
}

func (p *Parser) serializeEnv(data map[string]interface{}) (string, error) {
	lines := make([]string, 0, len(data))

	for key, value := range data {
		// Convert value to string
		var valueStr string
		switch v := value.(type) {
		case string:
			valueStr = v
		case bool:
			valueStr = fmt.Sprintf("%t", v)
		case int, int64, float64:
			valueStr = fmt.Sprintf("%v", v)
		default:
			// For complex types, serialize as JSON
			bytes, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			valueStr = string(bytes)
		}

		// Quote values that contain spaces or special characters
		if strings.ContainsAny(valueStr, " \t\n\"'\\") {
			valueStr = fmt.Sprintf("%q", strings.ReplaceAll(valueStr, "\"", "\\\""))
		}

		lines = append(lines, fmt.Sprintf("%s=%s", key, valueStr))
	}

	return strings.Join(lines, "\n"), nil
}
