// Package main provides the command-line interface for the mdfmt markdown formatter.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Gosayram/go-mdfmt/internal/version"
	"github.com/Gosayram/go-mdfmt/pkg/config"
	"github.com/Gosayram/go-mdfmt/pkg/formatter"
	"github.com/Gosayram/go-mdfmt/pkg/parser"
	"github.com/Gosayram/go-mdfmt/pkg/processor"
	"github.com/Gosayram/go-mdfmt/pkg/renderer"
)

const (
	// ExitCodeError indicates an error occurred
	ExitCodeError = 2
	// OutputFilePermissions defines the file permissions for output files
	OutputFilePermissions = 0o600
)

var (
	flagWrite = flag.Bool("w", false, "write result to (source) file instead of stdout")
	flagCheck = flag.Bool("c", false,
		"don't write the files back, just return the status. "+
			"Return code 0 if nothing would change, 1 if some files would be reformatted")
	flagCheck2 = flag.Bool("check", false,
		"don't write the files back, just return the status. "+
			"Return code 0 if nothing would change, 1 if some files would be reformatted")
	flagList    = flag.Bool("l", false, "list files whose formatting differs from mdfmt's")
	flagDiff    = flag.Bool("d", false, "display diffs instead of rewriting files")
	flagVerbose = flag.Bool("v", false, "verbose output")
	flagVersion = flag.Bool("version", false, "print version information")
	flagHelp    = flag.Bool("h", false, "show help")
	flagConfig  = flag.String("config", "", "path to configuration file")
)

// ProcessingArgs contains arguments for file processing
type ProcessingArgs struct {
	write   bool
	check   bool
	list    bool
	diff    bool
	verbose bool
}

func main() {
	flag.Parse()

	if *flagHelp {
		printUsage()
		return
	}

	if *flagVersion {
		fmt.Println(version.GetFullVersionInfo())
		return
	}

	// Get configuration
	cfg, err := loadConfig(*flagConfig, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(ExitCodeError)
	}

	// Get file paths
	paths := flag.Args()
	if len(paths) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No input files specified\n")
		os.Exit(ExitCodeError)
	}

	// Process files
	if err := processFiles(paths, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(ExitCodeError)
	}
}

// printUsage prints the usage information
func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: mdfmt [options] [files...]

mdfmt formats Markdown files according to consistent style rules.

Options:
  -w, --write           Write formatted content back to files
  -d, --diff            Show diff of changes without writing
  -c, --check           Check if files are formatted (exit 1 if not)
      --line-width N    Maximum line width for text reflow (default: from config)
      --config FILE     Path to configuration file
  -v, --verbose         Verbose output
      --version         Show version information
  -h, --help            Show this help message

Examples:
  mdfmt README.md                    Format README.md to stdout
  mdfmt --write docs/               Format all .md files in docs/
  mdfmt --check --diff *.md         Check formatting and show diffs
  mdfmt --line-width 100 --write .  Format with 100-char line width

For more information, visit: https://github.com/Gosayram/go-mdfmt
`)
}

// loadConfig loads the configuration from file or defaults
func loadConfig(configPath string, lineWidth int) (*config.Config, error) {
	cfg := config.Default()

	if configPath != "" {
		// Load from specified config file
		if err := cfg.LoadFromFile(configPath); err != nil {
			return nil, fmt.Errorf("failed to load config from %s: %w", configPath, err)
		}
	} else {
		// Try to find config file automatically
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %w", err)
		}

		configFile, err := config.FindConfigFile(wd)
		if err == nil {
			if err := cfg.LoadFromFile(configFile); err != nil {
				return nil, fmt.Errorf("failed to load config from %s: %w", configFile, err)
			}
		}
		// If no config file found, use defaults (already set above)
	}

	// Override line width if specified
	if lineWidth > 0 {
		cfg.LineWidth = lineWidth
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// processFiles processes the specified files
func processFiles(paths []string, cfg *config.Config) error {
	fp := processor.NewFileProcessor(cfg, *flagVerbose)

	files, err := fp.FindFiles(paths)
	if err != nil {
		return fmt.Errorf("failed to find files: %w", err)
	}

	if len(files) == 0 {
		if *flagVerbose {
			fmt.Println("No markdown files found")
		}
		return nil
	}

	// Create processing arguments
	args := &ProcessingArgs{
		write:   *flagWrite,
		check:   *flagCheck || *flagCheck2,
		list:    *flagList,
		diff:    *flagDiff,
		verbose: *flagVerbose,
	}

	var hasChanges bool
	for _, file := range files {
		changed, err := processFile(file, cfg, args)
		if err != nil {
			return fmt.Errorf("error processing %s: %w", file.Path, err)
		}
		if changed {
			hasChanges = true
		}
	}

	// Handle check mode exit code
	if args.check && hasChanges {
		os.Exit(1)
	}

	return nil
}

// processFile processes a single file
func processFile(file processor.FileInfo, cfg *config.Config, args *ProcessingArgs) (bool, error) {
	content, err := os.ReadFile(file.Path)
	if err != nil {
		return false, fmt.Errorf("failed to read file: %w", err)
	}

	formatted, err := formatMarkdownContent(content, cfg)
	if err != nil {
		return false, err
	}

	changed := hasContentChanged(content, formatted)

	if args.verbose && changed {
		fmt.Printf("File %s will be reformatted\n", file.Path)
	}

	if err := handleFileOutput(file.Path, formatted, changed, args); err != nil {
		return false, err
	}

	return changed, nil
}

// formatMarkdownContent processes markdown content through parse -> format -> render pipeline
func formatMarkdownContent(content []byte, cfg *config.Config) (string, error) {
	p := parser.DefaultParser()
	doc, err := p.Parse(content)
	if err != nil {
		return "", fmt.Errorf("failed to parse markdown: %w", err)
	}

	engine := formatter.New()
	engine.RegisterDefaults()

	if formatErr := engine.Format(doc, cfg); formatErr != nil {
		return "", fmt.Errorf("failed to format document: %w", formatErr)
	}

	mdRenderer := renderer.New()
	formatted, err := mdRenderer.Render(doc, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to render document: %w", err)
	}

	return formatted, nil
}

// hasContentChanged checks if the content has been modified after formatting
func hasContentChanged(original []byte, formatted string) bool {
	originalContent := strings.TrimSpace(string(original))
	formattedContent := strings.TrimSpace(formatted)
	return originalContent != formattedContent
}

// handleFileOutput handles different output modes based on processing arguments
func handleFileOutput(filePath, formatted string, changed bool, args *ProcessingArgs) error {
	switch {
	case args.write:
		return handleWriteMode(filePath, formatted, changed, args)
	case args.check:
		return handleCheckMode(filePath, changed, args)
	case args.list:
		return handleListMode(filePath, changed)
	case args.diff:
		return handleDiffMode(filePath, changed)
	default:
		return handleStdoutMode(formatted)
	}
}

// handleWriteMode writes formatted content back to file
func handleWriteMode(filePath, formatted string, changed bool, args *ProcessingArgs) error {
	if changed {
		if err := os.WriteFile(filePath, []byte(formatted), OutputFilePermissions); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		if args.verbose {
			fmt.Printf("Formatted: %s\n", filePath)
		}
	} else if args.verbose {
		fmt.Printf("Already formatted: %s\n", filePath)
	}
	return nil
}

// handleCheckMode handles check mode output
func handleCheckMode(filePath string, changed bool, args *ProcessingArgs) error {
	if changed && args.verbose {
		fmt.Printf("would reformat %s\n", filePath)
	}
	return nil
}

// handleListMode handles list mode output
func handleListMode(filePath string, changed bool) error {
	if changed {
		fmt.Println(filePath)
	}
	return nil
}

// handleDiffMode handles diff mode output
func handleDiffMode(filePath string, changed bool) error {
	if changed {
		fmt.Printf("--- %s\n+++ %s\n", filePath, filePath)
		fmt.Println("File would be reformatted")
	}
	return nil
}

// handleStdoutMode writes formatted content to stdout
func handleStdoutMode(formatted string) error {
	fmt.Print(formatted)
	return nil
}
