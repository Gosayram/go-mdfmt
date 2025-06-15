package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Gosayram/go-mdfmt/internal/version"
	"github.com/Gosayram/go-mdfmt/pkg/config"
)

// CLI flags
var (
	flagWrite     = flag.Bool("w", false, "write result to (source) file instead of stdout")
	flagWrite2    = flag.Bool("write", false, "write result to (source) file instead of stdout")
	flagDiff      = flag.Bool("d", false, "display diffs instead of rewriting files")
	flagDiff2     = flag.Bool("diff", false, "display diffs instead of rewriting files")
	flagCheck     = flag.Bool("c", false, "don't write the files back, just return the status. Return code 0 if nothing would change, 1 if some files would be reformatted")
	flagCheck2    = flag.Bool("check", false, "don't write the files back, just return the status. Return code 0 if nothing would change, 1 if some files would be reformatted")
	flagLineWidth = flag.Int("line-width", 0, "maximum line width for text reflow (0 = use config default)")
	flagConfig    = flag.String("config", "", "path to configuration file")
	flagVerbose   = flag.Bool("v", false, "verbose output")
	flagVerbose2  = flag.Bool("verbose", false, "verbose output")
	flagVersion   = flag.Bool("version", false, "print version information")
	flagHelp      = flag.Bool("h", false, "print help information")
	flagHelp2     = flag.Bool("help", false, "print help information")
)

func main() {
	flag.Parse()

	// Handle version flag
	if *flagVersion {
		fmt.Println(version.GetFullVersionInfo())
		return
	}

	// Handle help flag
	if *flagHelp || *flagHelp2 {
		printUsage()
		return
	}

	// Combine short and long flag variants
	write := *flagWrite || *flagWrite2
	diff := *flagDiff || *flagDiff2
	check := *flagCheck || *flagCheck2
	verbose := *flagVerbose || *flagVerbose2

	// Validate flag combinations
	if err := validateFlags(write, diff, check); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	// Load configuration
	cfg, err := loadConfig(*flagConfig, *flagLineWidth)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(2)
	}

	// Get files to process
	files := flag.Args()
	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no files specified\n")
		printUsage()
		os.Exit(2)
	}

	// Create application
	app := &App{
		config:  cfg,
		verbose: verbose,
		write:   write,
		diff:    diff,
		check:   check,
	}

	// Run the application
	exitCode := app.Run(files)
	os.Exit(exitCode)
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

// validateFlags validates flag combinations
func validateFlags(write, diff, check bool) error {
	count := 0
	if write {
		count++
	}
	if diff {
		count++
	}
	if check {
		count++
	}

	if count > 1 {
		return fmt.Errorf("only one of --write, --diff, or --check can be specified")
	}

	return nil
}

// loadConfig loads the configuration from file or defaults
func loadConfig(configPath string, lineWidth int) (*config.Config, error) {
	var cfg *config.Config
	var err error

	if configPath != "" {
		// Load from specified config file
		cfg, err = config.LoadFromFile(configPath)
		if err != nil {
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
			cfg, err = config.LoadFromFile(configFile)
			if err != nil {
				return nil, fmt.Errorf("failed to load config from %s: %w", configFile, err)
			}
		} else {
			// Use default configuration
			cfg = config.Default()
		}
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

// App represents the main application
type App struct {
	config  *config.Config
	verbose bool
	write   bool
	diff    bool
	check   bool
}

// Run runs the application with the given files
func (a *App) Run(files []string) int {
	if a.verbose {
		fmt.Fprintf(os.Stderr, "Processing %d files...\n", len(files))
	}

	// TODO: Implement actual file processing
	// For now, just print what we would do
	for _, file := range files {
		if a.verbose {
			fmt.Fprintf(os.Stderr, "Processing: %s\n", file)
		}

		if a.write {
			fmt.Printf("Would write formatted content to: %s\n", file)
		} else if a.diff {
			fmt.Printf("Would show diff for: %s\n", file)
		} else if a.check {
			fmt.Printf("Would check formatting for: %s\n", file)
		} else {
			fmt.Printf("Would format to stdout: %s\n", file)
		}
	}

	if a.verbose {
		fmt.Fprintf(os.Stderr, "Configuration:\n")
		fmt.Fprintf(os.Stderr, "  Line width: %d\n", a.config.LineWidth)
		fmt.Fprintf(os.Stderr, "  Heading style: %s\n", a.config.Heading.Style)
		fmt.Fprintf(os.Stderr, "  Bullet style: %s\n", a.config.List.BulletStyle)
	}

	return 0
}
