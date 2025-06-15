package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the configuration for mdfmt
type Config struct {
	// LineWidth is the maximum line width for text reflow
	LineWidth int `yaml:"line_width" json:"line_width"`

	// Heading configuration
	Heading HeadingConfig `yaml:"heading" json:"heading"`

	// List configuration
	List ListConfig `yaml:"list" json:"list"`

	// Code block configuration
	Code CodeConfig `yaml:"code" json:"code"`

	// Whitespace configuration
	Whitespace WhitespaceConfig `yaml:"whitespace" json:"whitespace"`

	// File processing configuration
	Files FilesConfig `yaml:"files" json:"files"`
}

// HeadingConfig contains heading formatting options
type HeadingConfig struct {
	// Style defines the heading style: "atx" (#) or "setext" (===)
	Style string `yaml:"style" json:"style"`
	// NormalizeLevels fixes heading level jumps
	NormalizeLevels bool `yaml:"normalize_levels" json:"normalize_levels"`
}

// ListConfig contains list formatting options
type ListConfig struct {
	// BulletStyle defines the bullet character: "-", "*", or "+"
	BulletStyle string `yaml:"bullet_style" json:"bullet_style"`
	// NumberStyle defines the numbering style: "." or ")"
	NumberStyle string `yaml:"number_style" json:"number_style"`
	// ConsistentIndentation ensures consistent indentation
	ConsistentIndentation bool `yaml:"consistent_indentation" json:"consistent_indentation"`
}

// CodeConfig contains code block formatting options
type CodeConfig struct {
	// FenceStyle defines the fence style: "```" or "~~~"
	FenceStyle string `yaml:"fence_style" json:"fence_style"`
	// LanguageDetection enables automatic language detection
	LanguageDetection bool `yaml:"language_detection" json:"language_detection"`
}

// WhitespaceConfig contains whitespace handling options
type WhitespaceConfig struct {
	// MaxBlankLines defines maximum consecutive blank lines
	MaxBlankLines int `yaml:"max_blank_lines" json:"max_blank_lines"`
	// TrimTrailingSpaces removes trailing spaces
	TrimTrailingSpaces bool `yaml:"trim_trailing_spaces" json:"trim_trailing_spaces"`
	// EnsureFinalNewline ensures files end with a newline
	EnsureFinalNewline bool `yaml:"ensure_final_newline" json:"ensure_final_newline"`
}

// FilesConfig contains file processing options
type FilesConfig struct {
	// Extensions defines which file extensions to process
	Extensions []string `yaml:"extensions" json:"extensions"`
	// IgnorePatterns defines glob patterns to ignore
	IgnorePatterns []string `yaml:"ignore_patterns" json:"ignore_patterns"`
}

// Default returns the default configuration
func Default() *Config {
	return &Config{
		LineWidth: 80,
		Heading: HeadingConfig{
			Style:           "atx",
			NormalizeLevels: true,
		},
		List: ListConfig{
			BulletStyle:           "-",
			NumberStyle:           ".",
			ConsistentIndentation: true,
		},
		Code: CodeConfig{
			FenceStyle:        "```",
			LanguageDetection: true,
		},
		Whitespace: WhitespaceConfig{
			MaxBlankLines:      2,
			TrimTrailingSpaces: true,
			EnsureFinalNewline: true,
		},
		Files: FilesConfig{
			Extensions:     []string{".md", ".markdown", ".mdown"},
			IgnorePatterns: []string{"node_modules/**", ".git/**"},
		},
	}
}

// LoadFromFile loads configuration from a YAML file
func LoadFromFile(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", filename, err)
	}

	config := Default()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", filename, err)
	}

	return config, nil
}

// SaveToFile saves configuration to a YAML file
func (c *Config) SaveToFile(filename string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", filename, err)
	}

	return nil
}

// FindConfigFile searches for configuration files in common locations
func FindConfigFile(startDir string) (string, error) {
	configNames := []string{
		".mdfmt.yaml",
		".mdfmt.yml",
		".mdfmt.json",
		"mdfmt.yaml",
		"mdfmt.yml",
		"mdfmt.json",
	}

	dir := startDir
	for {
		for _, name := range configNames {
			path := filepath.Join(dir, name)
			if _, err := os.Stat(path); err == nil {
				return path, nil
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached root directory
		}
		dir = parent
	}

	return "", fmt.Errorf("no configuration file found")
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.LineWidth < 1 {
		return fmt.Errorf("line_width must be greater than 0")
	}

	if c.Heading.Style != "atx" && c.Heading.Style != "setext" {
		return fmt.Errorf("heading.style must be 'atx' or 'setext'")
	}

	if !contains([]string{"-", "*", "+"}, c.List.BulletStyle) {
		return fmt.Errorf("list.bullet_style must be '-', '*', or '+'")
	}

	if !contains([]string{".", ")"}, c.List.NumberStyle) {
		return fmt.Errorf("list.number_style must be '.' or ')'")
	}

	if !contains([]string{"```", "~~~"}, c.Code.FenceStyle) {
		return fmt.Errorf("code.fence_style must be '```' or '~~~'")
	}

	if c.Whitespace.MaxBlankLines < 0 {
		return fmt.Errorf("whitespace.max_blank_lines must be >= 0")
	}

	return nil
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// IsMarkdownFile checks if a file is a markdown file based on extension
func (c *Config) IsMarkdownFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return contains(c.Files.Extensions, ext)
}

// ShouldIgnore checks if a file should be ignored based on patterns
func (c *Config) ShouldIgnore(filename string) bool {
	for _, pattern := range c.Files.IgnorePatterns {
		// Handle directory patterns with **
		if strings.Contains(pattern, "**") {
			// Simple glob matching for directory patterns
			cleanPattern := strings.TrimSuffix(pattern, "/**")
			if strings.HasPrefix(filename, cleanPattern) {
				return true
			}
		} else {
			if matched, _ := filepath.Match(pattern, filename); matched {
				return true
			}
		}
	}
	return false
}
