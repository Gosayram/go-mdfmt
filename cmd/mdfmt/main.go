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
	// ExitCodeChangesNeeded indicates files need formatting (for check mode)
	ExitCodeChangesNeeded = 1
	// OutputFilePermissions defines the file permissions for output files
	OutputFilePermissions = 0o600
)

var (
	// Main operation flags
	flagWrite = flag.Bool("w", false, "write formatted content back to files")
	flagCheck = flag.Bool("c", false, "check if files are formatted correctly (exit 1 if not)")
	flagList  = flag.Bool("l", false, "list files that need formatting")
	flagDiff  = flag.Bool("d", false, "show diff of changes without writing files")

	// Long versions of operation flags
	flagWriteLong = flag.Bool("write", false, "write formatted content back to files")
	flagCheckLong = flag.Bool("check", false, "check if files are formatted correctly (exit 1 if not)")
	flagListLong  = flag.Bool("list", false, "list files that need formatting")
	flagDiffLong  = flag.Bool("diff", false, "show diff of changes without writing files")

	// Configuration flags
	flagConfig = flag.String("config", "", "path to configuration file")

	// Output flags
	flagVerbose = flag.Bool("v", false, "verbose output")
	flagQuiet   = flag.Bool("q", false, "quiet mode (suppress non-error output)")

	// Long versions of output flags
	flagVerboseLong = flag.Bool("verbose", false, "verbose output")
	flagQuietLong   = flag.Bool("quiet", false, "quiet mode (suppress non-error output)")

	// Information flags
	flagVersion  = flag.Bool("version", false, "print version information")
	flagHelp     = flag.Bool("h", false, "show help message")
	flagHelpLong = flag.Bool("help", false, "show help message")
)

// ProcessingArgs contains arguments for file processing
type ProcessingArgs struct {
	write   bool
	check   bool
	list    bool
	diff    bool
	verbose bool
	quiet   bool
}

func main() {
	// Custom usage function
	flag.Usage = printUsage
	flag.Parse()

	if *flagHelp || *flagHelpLong {
		printUsage()
		return
	}

	if *flagVersion {
		fmt.Println(version.GetFullVersionInfo())
		return
	}

	// Validate flag combinations
	if err := validateFlags(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Run 'mdfmt -h' for usage information.\n")
		os.Exit(ExitCodeError)
	}

	// Get configuration
	cfg, err := loadConfig(*flagConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(ExitCodeError)
	}

	// Get file paths
	paths := flag.Args()
	if len(paths) == 0 {
		if !*flagQuiet {
			fmt.Fprintf(os.Stderr, "Error: No input files or directories specified\n")
			fmt.Fprintf(os.Stderr, "Run 'mdfmt -h' for usage information.\n")
		}
		os.Exit(ExitCodeError)
	}

	// Process files
	if err := processFiles(paths, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(ExitCodeError)
	}
}

// validateFlags validates flag combinations
func validateFlags() error {
	// Count mutually exclusive operation flags
	operationCount := 0
	if *flagWrite || *flagWriteLong {
		operationCount++
	}
	if *flagCheck || *flagCheckLong {
		operationCount++
	}
	if *flagList || *flagListLong {
		operationCount++
	}
	if *flagDiff || *flagDiffLong {
		operationCount++
	}

	if operationCount > 1 {
		return fmt.Errorf("only one of -w/--write, -c/--check, -l/--list, -d/--diff can be specified")
	}

	if (*flagVerbose || *flagVerboseLong) && (*flagQuiet || *flagQuietLong) {
		return fmt.Errorf("-v/--verbose and -q/--quiet cannot be used together")
	}

	return nil
}

// printUsage prints the usage information
func printUsage() {
	fmt.Fprintf(os.Stderr, `mdfmt - Fast, reliable Markdown formatter

USAGE:
    mdfmt [OPTIONS] <files...>

DESCRIPTION:
    mdfmt formats Markdown files according to consistent style rules.
    By default, formatted output is written to stdout.

OPTIONS:
    Operation modes (mutually exclusive):
        -w, --write     Write formatted content back to files
        -c, --check     Check if files are formatted correctly (exit 1 if not)
        -l, --list      List files that need formatting
        -d, --diff      Show diff of changes without writing files

    Configuration:
        --config <file> Path to configuration file (.mdfmt.yaml)

    Output control:
        -v, --verbose   Verbose output (show processed files)
        -q, --quiet     Quiet mode (suppress non-error output)

    Information:
        -h, --help      Show this help message
        --version       Print version information

EXAMPLES:
    Format a single file to stdout:
        mdfmt README.md

    Format files in place:
        mdfmt --write *.md
        mdfmt -w docs/

    Check if files are properly formatted:
        mdfmt --check README.md docs/
        echo $?  # 0 if formatted, 1 if needs formatting

    Show what would change:
        mdfmt --diff README.md

    List files that need formatting:
        mdfmt --list docs/

    Use custom configuration:
        mdfmt --config .mdfmt.yaml --write docs/

    Verbose processing:
        mdfmt --verbose --write docs/

EXIT CODES:
    0   Success (no changes needed in check mode)
    1   Files need formatting (check mode only)
    2   Error occurred

CONFIGURATION:
    mdfmt looks for configuration in the following order:
    1. File specified by -config flag
    2. .mdfmt.yaml in current directory
    3. .mdfmt.yaml in parent directories (up to repository root)
    4. Built-in defaults

    Create example config: mdfmt -config example > .mdfmt.yaml

For more information: https://github.com/Gosayram/go-mdfmt
`)
}

// loadConfig loads the configuration from file or defaults
func loadConfig(configPath string) (*config.Config, error) {
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

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// createProcessingArgs creates processing arguments from flags
func createProcessingArgs() *ProcessingArgs {
	verbose := *flagVerbose || *flagVerboseLong
	quiet := *flagQuiet || *flagQuietLong

	return &ProcessingArgs{
		write:   *flagWrite || *flagWriteLong,
		check:   *flagCheck || *flagCheckLong,
		list:    *flagList || *flagListLong,
		diff:    *flagDiff || *flagDiffLong,
		verbose: verbose,
		quiet:   quiet,
	}
}

// processFiles processes the specified files
func processFiles(paths []string, cfg *config.Config) error {
	args := createProcessingArgs()
	fp := processor.NewFileProcessor(cfg, args.verbose)

	files, err := fp.FindFiles(paths)
	if err != nil {
		return fmt.Errorf("failed to find files: %w", err)
	}

	if len(files) == 0 {
		if args.verbose && !args.quiet {
			fmt.Println("No markdown files found")
		}
		return nil
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
		os.Exit(ExitCodeChangesNeeded)
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

	if args.verbose && !args.quiet && changed {
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
		if args.verbose && !args.quiet {
			fmt.Printf("Formatted: %s\n", filePath)
		}
	} else if args.verbose && !args.quiet {
		fmt.Printf("Already formatted: %s\n", filePath)
	}
	return nil
}

// handleCheckMode handles check mode output
func handleCheckMode(filePath string, changed bool, args *ProcessingArgs) error {
	if changed && args.verbose && !args.quiet {
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
