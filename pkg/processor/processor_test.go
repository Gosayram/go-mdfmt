package processor

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/Gosayram/go-mdfmt/pkg/config"
)

func TestNewFileProcessor(t *testing.T) {
	cfg := config.Default()
	processor := NewFileProcessor(cfg, true)

	if processor.config != cfg {
		t.Error("Expected config to be set correctly")
	}
	if !processor.verbose {
		t.Error("Expected verbose to be true")
	}
}

func TestIsMarkdownFile(t *testing.T) {
	cfg := config.Default()
	processor := NewFileProcessor(cfg, false)

	tests := []struct {
		path     string
		expected bool
	}{
		{"README.md", true},
		{"doc.markdown", true},
		{"file.mdown", true},
		{"script.js", false},
		{"style.css", false},
		{"README.MD", true},
		{"file.txt", false},
		{"test.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := processor.isMarkdownFile(tt.path)
			if result != tt.expected {
				t.Errorf("isMarkdownFile(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestShouldIgnoreFile(t *testing.T) {
	cfg := config.Default()
	processor := NewFileProcessor(cfg, false)

	tests := []struct {
		path     string
		expected bool
	}{
		{"README.md", false},
		{"node_modules/package.json", true},
		{".git/config", true},
		{"docs/guide.md", false},
		{"node_modules/lib/index.js", true},
		{"vendor/github.com/pkg/errors/errors.go", true},
		{"regular/file.md", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := processor.shouldIgnoreFile(tt.path)
			if result != tt.expected {
				t.Errorf("shouldIgnoreFile(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestFindFiles(t *testing.T) {
	// Create a temporary directory structure for testing
	tmpDir, err := os.MkdirTemp("", "mdfmt-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	testFiles := []string{
		"README.md",
		"docs/guide.md",
		"docs/api.markdown",
		"src/main.go",
		"node_modules/package.json",
		".git/config",
	}

	for _, file := range testFiles {
		fullPath := filepath.Join(tmpDir, file)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}

	cfg := config.Default()
	processor := NewFileProcessor(cfg, false)

	// Test finding files in the temp directory
	files, err := processor.FindFiles([]string{tmpDir})
	if err != nil {
		t.Fatalf("FindFiles failed: %v", err)
	}

	// We should find 3 markdown files (README.md, docs/guide.md, docs/api.markdown)
	// but not the ones in ignored directories
	expectedCount := 3
	if len(files) != expectedCount {
		t.Errorf("Expected %d files, got %d", expectedCount, len(files))
		for _, file := range files {
			t.Logf("Found file: %s", file.RelativePath)
		}
	}

	// Check that all found files are markdown files
	for _, file := range files {
		if !processor.isMarkdownFile(file.Path) {
			t.Errorf("Non-markdown file found: %s", file.Path)
		}
	}
}

func TestReadWriteFile(t *testing.T) {
	fp := NewFileProcessor(config.Default(), false)

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test-*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// Test data
	testContent := []byte("# Test Content\n\nThis is a test.")

	// Write content
	err = fp.writeFile(tmpfile.Name(), testContent)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Read content back
	readContent, err := fp.readFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	// Compare content
	if !bytes.Equal(testContent, readContent) {
		t.Errorf("Content mismatch. Expected %q, got %q", testContent, readContent)
	}
}

func TestBackupFile(t *testing.T) {
	fp := NewFileProcessor(config.Default(), false)

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test-*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer os.Remove(tmpfile.Name() + ".backup")

	// Write some content
	testContent := []byte("# Original Content")
	err = fp.writeFile(tmpfile.Name(), testContent)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create backup
	err = fp.BackupFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("BackupFile failed: %v", err)
	}

	// Verify backup exists and has same content
	backupContent, err := fp.readFile(tmpfile.Name() + ".backup")
	if err != nil {
		t.Fatalf("Failed to read backup file: %v", err)
	}

	if !bytes.Equal(testContent, backupContent) {
		t.Errorf("Backup content mismatch. Expected %q, got %q", testContent, backupContent)
	}
}

func TestProcessFiles(t *testing.T) {
	cfg := config.Default()
	processor := NewFileProcessor(cfg, false)

	// Create test files
	files := []FileInfo{
		{Path: "test1.md", RelativePath: "test1.md", Size: 100},
		{Path: "test2.md", RelativePath: "test2.md", Size: 200},
		{Path: "test3.md", RelativePath: "test3.md", Size: 300},
	}

	// Mock processor function with atomic counter
	var processCount int64
	mockProcessor := func(file FileInfo) ProcessingResult {
		atomic.AddInt64(&processCount, 1)
		return ProcessingResult{
			File:      file,
			Success:   true,
			Error:     nil,
			Changed:   false,
			BytesRead: file.Size,
		}
	}

	// Process files
	results := processor.ProcessFiles(files, mockProcessor)

	// Check results
	if len(results) != len(files) {
		t.Errorf("Expected %d results, got %d", len(files), len(results))
	}

	if atomic.LoadInt64(&processCount) != int64(len(files)) {
		t.Errorf("Expected processor to be called %d times, got %d", len(files), atomic.LoadInt64(&processCount))
	}

	// Check all results are successful
	for i, result := range results {
		if !result.Success {
			t.Errorf("Result %d should be successful", i)
		}
		if result.Error != nil {
			t.Errorf("Result %d should not have error: %v", i, result.Error)
		}
	}
}

func TestMinFunction(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{1, 2, 1},
		{5, 3, 3},
		{10, 10, 10},
		{0, 1, 0},
		{-1, 5, -1},
	}

	for _, tt := range tests {
		result := min(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("min(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
		}
	}
}

// Benchmark tests
func BenchmarkFileProcessor_FindFiles(b *testing.B) {
	cfg := config.Default()
	processor := NewFileProcessor(cfg, false)

	// Create temporary directory structure
	tempDir := b.TempDir()

	// Create more test files to simulate real project
	testFiles := []string{
		"README.md",
		"docs/guide.md",
		"docs/api.markdown",
		"docs/tutorial.md",
		"docs/examples/basic.md",
		"docs/examples/advanced.md",
		"src/README.md",
		"test.txt",
		"script.js",
		"nested/deep/file.md",
		"nested/deep/another.markdown",
		"project/docs/spec.md",
	}

	for _, file := range testFiles {
		fullPath := filepath.Join(tempDir, file)
		dir := filepath.Dir(fullPath)
		os.MkdirAll(dir, 0755)
		os.WriteFile(fullPath, []byte("# Test\n\nContent with more text"), 0644)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := processor.FindFiles([]string{tempDir})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFileProcessor_FindFilesLargeProject(b *testing.B) {
	cfg := config.Default()
	processor := NewFileProcessor(cfg, false)

	// Create temporary directory structure
	tempDir := b.TempDir()

	// Simulate a large project with many files
	for i := 0; i < 50; i++ {
		for j := 0; j < 10; j++ {
			file := fmt.Sprintf("module%d/docs/file%d.md", i, j)
			fullPath := filepath.Join(tempDir, file)
			dir := filepath.Dir(fullPath)
			os.MkdirAll(dir, 0755)
			content := fmt.Sprintf("# Module %d File %d\n\nContent here", i, j)
			os.WriteFile(fullPath, []byte(content), 0644)
		}
	}

	// Add some non-markdown files
	for i := 0; i < 100; i++ {
		file := fmt.Sprintf("src/file%d.go", i)
		fullPath := filepath.Join(tempDir, file)
		dir := filepath.Dir(fullPath)
		os.MkdirAll(dir, 0755)
		os.WriteFile(fullPath, []byte("package main"), 0644)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := processor.FindFiles([]string{tempDir})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFileProcessor_FileOperations(b *testing.B) {
	cfg := config.Default()
	processor := NewFileProcessor(cfg, false)

	// Larger, more realistic markdown content
	content := []byte(`# Test Document

This is a test paragraph with some **bold** and *italic* text.

## Section 1

- Item 1 with detailed description
- Item 2 with more content
- Item 3 with even more text

### Subsection

1. Ordered item 1
2. Ordered item 2
3. Ordered item 3

## Section 2

` + "```go\n" + `func example() {
    fmt.Println("Hello, world!")
    return true
}
` + "```" + `

> This is a blockquote with important information
> that spans multiple lines.

## Final Section

Final paragraph with [link](https://example.com) and more text.
`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create temp file
		tmpfile, err := os.CreateTemp("", "bench-*.md")
		if err != nil {
			b.Fatal(err)
		}

		// Test file operations
		err = processor.writeFile(tmpfile.Name(), content)
		if err != nil {
			b.Fatal(err)
		}

		_, err = processor.readFile(tmpfile.Name())
		if err != nil {
			b.Fatal(err)
		}

		// Cleanup
		os.Remove(tmpfile.Name())
	}
}
