package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.LineWidth != 80 {
		t.Errorf("Expected LineWidth to be 80, got %d", cfg.LineWidth)
	}

	if cfg.Heading.Style != "atx" {
		t.Errorf("Expected Heading.Style to be 'atx', got %s", cfg.Heading.Style)
	}

	if !cfg.Heading.NormalizeLevels {
		t.Error("Expected Heading.NormalizeLevels to be true")
	}

	if cfg.List.BulletStyle != "-" {
		t.Errorf("Expected List.BulletStyle to be '-', got %s", cfg.List.BulletStyle)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "valid default config",
			config:  Default(),
			wantErr: false,
		},
		{
			name: "invalid line width",
			config: &Config{
				LineWidth:  0,
				Heading:    HeadingConfig{Style: "atx"},
				List:       ListConfig{BulletStyle: "-", NumberStyle: "."},
				Code:       CodeConfig{FenceStyle: "```"},
				Whitespace: WhitespaceConfig{MaxBlankLines: 2},
			},
			wantErr: true,
		},
		{
			name: "invalid heading style",
			config: &Config{
				LineWidth:  80,
				Heading:    HeadingConfig{Style: "invalid"},
				List:       ListConfig{BulletStyle: "-", NumberStyle: "."},
				Code:       CodeConfig{FenceStyle: "```"},
				Whitespace: WhitespaceConfig{MaxBlankLines: 2},
			},
			wantErr: true,
		},
		{
			name: "invalid bullet style",
			config: &Config{
				LineWidth:  80,
				Heading:    HeadingConfig{Style: "atx"},
				List:       ListConfig{BulletStyle: "invalid", NumberStyle: "."},
				Code:       CodeConfig{FenceStyle: "```"},
				Whitespace: WhitespaceConfig{MaxBlankLines: 2},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadFromFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test.yaml")

	configContent := `line_width: 100
heading:
  style: setext
  normalize_levels: false
list:
  bullet_style: "*"
  number_style: ")"
code:
  fence_style: "~~~"
whitespace:
  max_blank_lines: 1
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	cfg, err := LoadFromFile(configFile)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	if cfg.LineWidth != 100 {
		t.Errorf("Expected LineWidth to be 100, got %d", cfg.LineWidth)
	}

	if cfg.Heading.Style != "setext" {
		t.Errorf("Expected Heading.Style to be 'setext', got %s", cfg.Heading.Style)
	}

	if cfg.Heading.NormalizeLevels {
		t.Error("Expected Heading.NormalizeLevels to be false")
	}

	if cfg.List.BulletStyle != "*" {
		t.Errorf("Expected List.BulletStyle to be '*', got %s", cfg.List.BulletStyle)
	}

	if cfg.List.NumberStyle != ")" {
		t.Errorf("Expected List.NumberStyle to be ')', got %s", cfg.List.NumberStyle)
	}

	if cfg.Code.FenceStyle != "~~~" {
		t.Errorf("Expected Code.FenceStyle to be '~~~', got %s", cfg.Code.FenceStyle)
	}

	if cfg.Whitespace.MaxBlankLines != 1 {
		t.Errorf("Expected Whitespace.MaxBlankLines to be 1, got %d", cfg.Whitespace.MaxBlankLines)
	}
}

func TestLoadFromFileNotFound(t *testing.T) {
	_, err := LoadFromFile("nonexistent.yaml")
	if err == nil {
		t.Error("Expected error when loading nonexistent file")
	}
}

func TestSaveToFile(t *testing.T) {
	cfg := Default()
	cfg.LineWidth = 120

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "save_test.yaml")

	err := cfg.SaveToFile(configFile)
	if err != nil {
		t.Fatalf("SaveToFile() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Load it back and verify
	loadedCfg, err := LoadFromFile(configFile)
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if loadedCfg.LineWidth != 120 {
		t.Errorf("Expected LineWidth to be 120, got %d", loadedCfg.LineWidth)
	}
}

func TestIsMarkdownFile(t *testing.T) {
	cfg := Default()

	tests := []struct {
		filename string
		expected bool
	}{
		{"README.md", true},
		{"doc.markdown", true},
		{"file.mdown", true},
		{"script.js", false},
		{"style.css", false},
		{"README.MD", true}, // case insensitive
		{"file.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := cfg.IsMarkdownFile(tt.filename)
			if result != tt.expected {
				t.Errorf("IsMarkdownFile(%s) = %v, expected %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestShouldIgnore(t *testing.T) {
	cfg := Default()

	tests := []struct {
		filename string
		expected bool
	}{
		{"README.md", false},
		{"node_modules/package.json", true},
		{".git/config", true},
		{"docs/guide.md", false},
		{"node_modules/lib/index.js", true},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := cfg.ShouldIgnore(tt.filename)
			if result != tt.expected {
				t.Errorf("ShouldIgnore(%s) = %v, expected %v", tt.filename, result, tt.expected)
			}
		})
	}
}
