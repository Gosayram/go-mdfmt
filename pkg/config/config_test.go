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
	content := `line_width: 100
heading:
  style: "setext"
  normalize_levels: false
list:
  bullet_style: "*"
  number_style: ")"
code:
  fence_style: "~~~"
  language_detection: false
`

	tmpfile, err := os.CreateTemp("", "test-config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	cfg := Default()
	err = cfg.LoadFromFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}

	// Verify loaded values
	if cfg.LineWidth != 100 {
		t.Errorf("Expected LineWidth 100, got %d", cfg.LineWidth)
	}
	if cfg.Heading.Style != "setext" {
		t.Errorf("Expected Heading.Style 'setext', got %s", cfg.Heading.Style)
	}
	if cfg.List.BulletStyle != "*" {
		t.Errorf("Expected List.BulletStyle '*', got %s", cfg.List.BulletStyle)
	}
}

func TestLoadFromFile_NotFound(t *testing.T) {
	cfg := Default()
	err := cfg.LoadFromFile("nonexistent.yaml")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestLoadFromFile_InvalidYAML(t *testing.T) {
	content := "invalid: yaml: content:"

	tmpfile, err := os.CreateTemp("", "test-invalid-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	cfg := Default()
	err = cfg.LoadFromFile(tmpfile.Name())
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
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
	loadedCfg := Default()
	err = loadedCfg.LoadFromFile(configFile)
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

// Benchmark tests
func BenchmarkConfig_Default(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Default()
	}
}

func BenchmarkConfig_Validate(b *testing.B) {
	cfg := Default()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := cfg.Validate()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConfig_LoadFromFile(b *testing.B) {
	// Create temporary config file
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	configContent := `line_width: 80
heading:
  style: "atx"
  normalize_levels: true
list:
  bullet_style: "-"
  number_style: "."
  consistent_indentation: true
code:
  fence_style: "` + "```" + `"
  language_detection: true
whitespace:
  max_blank_lines: 2
  trim_trailing_spaces: true
  ensure_final_newline: true
files:
  extensions: [".md", ".markdown", ".mdown"]
  ignore_patterns: ["node_modules/**", ".git/**", "vendor/**"]
`

	err = os.WriteFile(tmpfile.Name(), []byte(configContent), 0644)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cfg := Default()
		err := cfg.LoadFromFile(tmpfile.Name())
		if err != nil {
			b.Fatal(err)
		}
	}
}
