// Package processor provides file processing functionality for markdown formatting.
package processor

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Gosayram/go-mdfmt/pkg/config"
)

// Constants
const (
	// FilePermissions defines the file permissions for written files
	FilePermissions = 0o600
)

// FileProcessor handles file operations and batch processing
type FileProcessor struct {
	config  *config.Config
	verbose bool
}

// NewFileProcessor creates a new file processor instance
func NewFileProcessor(cfg *config.Config, verbose bool) *FileProcessor {
	return &FileProcessor{
		config:  cfg,
		verbose: verbose,
	}
}

// FileInfo contains information about a file to be processed
type FileInfo struct {
	Path         string
	RelativePath string
	IsDirectory  bool
	Size         int64
}

// ProcessingResult contains the result of processing a file
type ProcessingResult struct {
	File      FileInfo
	Success   bool
	Error     error
	Changed   bool
	BytesRead int64
}

// FindFiles recursively finds all Markdown files in the given paths
func (fp *FileProcessor) FindFiles(paths []string) ([]FileInfo, error) {
	var files []FileInfo
	seen := make(map[string]bool)

	for _, path := range paths {
		err := fp.findFilesInPath(path, &files, seen)
		if err != nil {
			return nil, fmt.Errorf("error processing path %s: %w", path, err)
		}
	}

	return files, nil
}

// findFilesInPath recursively finds files in a single path
func (fp *FileProcessor) findFilesInPath(path string, files *[]FileInfo, seen map[string]bool) error {
	// Clean and resolve the path
	cleanPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve path %s: %w", path, err)
	}

	// Check if we've already processed this path
	if seen[cleanPath] {
		return nil
	}
	seen[cleanPath] = true

	// Get file info
	info, err := os.Stat(cleanPath)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", cleanPath, err)
	}

	if info.IsDir() {
		return fp.findFilesInDirectory(cleanPath, files, seen)
	}

	// Check if it's a Markdown file
	if fp.isMarkdownFile(cleanPath) && !fp.shouldIgnoreFile(cleanPath) {
		relPath, _ := filepath.Rel(".", cleanPath)
		*files = append(*files, FileInfo{
			Path:         cleanPath,
			RelativePath: relPath,
			IsDirectory:  false,
			Size:         info.Size(),
		})
	}

	return nil
}

// findFilesInDirectory finds files in a directory
func (fp *FileProcessor) findFilesInDirectory(dir string, files *[]FileInfo, seen map[string]bool) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if fp.verbose {
				fmt.Fprintf(os.Stderr, "Warning: skipping %s: %v\n", path, err)
			}
			return nil // Skip files we can't access
		}

		// Skip if already seen
		cleanPath, err := filepath.Abs(path)
		if err != nil {
			return nil
		}
		if seen[cleanPath] {
			return nil
		}

		// Check if we should ignore this path
		if fp.shouldIgnoreFile(path) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// If it's a Markdown file, add it
		if !d.IsDir() && fp.isMarkdownFile(path) {
			info, err := d.Info()
			if err != nil {
				return nil
			}

			relPath, _ := filepath.Rel(".", path)
			*files = append(*files, FileInfo{
				Path:         path,
				RelativePath: relPath,
				IsDirectory:  false,
				Size:         info.Size(),
			})
		}

		return nil
	})
}

// isMarkdownFile checks if a file is a Markdown file based on extension
func (fp *FileProcessor) isMarkdownFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, validExt := range fp.config.Files.Extensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

// shouldIgnoreFile checks if a file should be ignored based on patterns
func (fp *FileProcessor) shouldIgnoreFile(path string) bool {
	return fp.config.ShouldIgnore(path)
}

// ProcessFiles processes multiple files concurrently
func (fp *FileProcessor) ProcessFiles(files []FileInfo, processor func(FileInfo) ProcessingResult) []ProcessingResult {
	const maxWorkers = 8
	workers := minInt(maxWorkers, len(files))
	if workers == 0 {
		return nil
	}

	jobs := make(chan FileInfo, len(files))
	results := make(chan ProcessingResult, len(files))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range jobs {
				results <- processor(file)
			}
		}()
	}

	// Send jobs
	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	// Wait for completion
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var allResults []ProcessingResult
	for result := range results {
		allResults = append(allResults, result)
	}

	return allResults
}

// readFile reads content from a file.
func (fp *FileProcessor) readFile(path string) ([]byte, error) {
	if fp.verbose {
		fmt.Printf("Reading file: %s\n", path)
	}
	content, err := os.ReadFile(path) // #nosec G304 - path is validated through file discovery
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}
	return content, nil
}

// writeFile writes content to a file.
func (fp *FileProcessor) writeFile(path string, content []byte) error {
	if fp.verbose {
		fmt.Printf("Writing file: %s\n", path)
	}
	err := os.WriteFile(path, content, FilePermissions)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}
	return nil
}

// BackupFile creates a backup of a file before modification
func (fp *FileProcessor) BackupFile(path string) error {
	content, err := fp.readFile(path)
	if err != nil {
		return err
	}

	backupPath := path + ".backup"
	return fp.writeFile(backupPath, content)
}

// minInt returns the minimum of two integers.
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
