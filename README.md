# go-mdfmt

A fast, reliable, and opinionated Markdown formatter written in Go. Provides consistent, configurable formatting for `.md` files across projects, making your documentation readable, lintable, and style-consistent.

## Features

**Core Formatting Capabilities**
- **Text Reflow**: Wraps long paragraphs at configurable line width (default 80 characters)
- **Heading Normalization**: Ensures consistent heading levels and spacing using ATX style
- **List Standardization**: Consistent bullet and numbered list formatting with proper indentation
- **Code Block Formatting**: Auto-corrects indentation and applies consistent fence styles
- **Whitespace Management**: Removes excessive blank lines and trailing spaces
- **Link Preservation**: Maintains markdown link structure and formatting

**Operation Modes**
- **Write Mode**: Format files in-place (`--write`)
- **Check Mode**: Verify formatting without changes (`--check`) - ideal for CI/CD
- **Diff Mode**: Preview changes before applying (`--diff`)
- **List Mode**: Show files that need formatting (`--list`)
- **Standard Output**: Display formatted content to stdout (default)

**Configuration and Integration**
- YAML-based configuration with sensible defaults
- File pattern matching with ignore capabilities
- Verbose and quiet output modes
- CI/CD ready with meaningful exit codes
- Cross-platform binary releases

## Installation

### Pre-built Binaries

Download the latest binary from the [releases page](https://github.com/Gosayram/go-mdfmt/releases):

**Linux:**
```bash
curl -L -o mdfmt https://github.com/Gosayram/go-mdfmt/releases/latest/download/mdfmt-0.2.6-linux-amd64
chmod +x mdfmt
sudo mv mdfmt /usr/local/bin/
```

**macOS:**
```bash
curl -L -o mdfmt https://github.com/Gosayram/go-mdfmt/releases/latest/download/mdfmt-0.2.6-darwin-arm64
chmod +x mdfmt
sudo mv mdfmt /usr/local/bin/
```

**Windows:**
```powershell
curl -L -o mdfmt.exe https://github.com/Gosayram/go-mdfmt/releases/latest/download/mdfmt-0.2.6-windows-amd64.exe
```

### From Source

**Prerequisites:**
- Go 1.24.4 or later
- Git

**Build from source:**
```bash
git clone https://github.com/Gosayram/go-mdfmt.git
cd go-mdfmt
make build
# Binary will be available in bin/mdfmt
```

**Install from Go modules:**
```bash
go install github.com/Gosayram/go-mdfmt/cmd/mdfmt@latest
```

## Quick Start

### Basic Usage

```bash
# Format file to stdout
mdfmt README.md

# Format multiple files to stdout
mdfmt docs/*.md

# Format directory (finds all .md files)
mdfmt docs/
```

### Write Changes to Files

```bash
# Format single file in-place
mdfmt --write README.md

# Format multiple files in-place
mdfmt --write docs/*.md

# Format all markdown files in project
mdfmt --write .
```

### Check Formatting (CI/CD)

```bash
# Check if files are properly formatted
mdfmt --check docs/
echo $? # Exit code: 0 = formatted, 1 = needs formatting, 2 = error

# Show what would change
mdfmt --diff README.md

# List files that need formatting
mdfmt --list docs/
```

## Command Line Interface

```
USAGE:
    mdfmt [OPTIONS] <files...>

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

EXIT CODES:
    0   Success (no changes needed in check mode)
    1   Files need formatting (check mode only)
    2   Error occurred
```

## Configuration

mdfmt uses YAML configuration files with automatic discovery. Configuration files are searched in this order:

1. File specified by `--config` flag
2. `.mdfmt.yaml` in current directory
3. `.mdfmt.yaml` in parent directories (up to repository root)
4. Built-in defaults

### Configuration File Structure

Create `.mdfmt.yaml` in your project root:

```yaml
# Line width for paragraph reflow
line_width: 80

# Heading configuration
heading:
  style: "atx"              # Use # headings instead of === underline
  normalize_levels: true    # Fix heading level jumps

# List formatting
list:
  bullet_style: "-"         # Use - for bullets (options: -, *, +)
  number_style: "."         # Use 1. for numbered lists (options: ., ))
  consistent_indentation: true

# Code block formatting
code:
  fence_style: "```"        # Use ``` for code blocks (options: ```, ~~~)
  language_detection: true  # Auto-detect and add language labels

# Whitespace handling
whitespace:
  max_blank_lines: 2        # Maximum consecutive blank lines
  trim_trailing_spaces: true
  ensure_final_newline: true

# File processing
files:
  extensions: [".md", ".markdown", ".mdown"]
  ignore_patterns: ["node_modules/**", ".git/**", "vendor/**"]
```

### Configuration Validation

All configuration options are validated on startup. Invalid configurations will result in descriptive error messages.

## Development

### Building from Source

```bash
# Install dependencies
make deps

# Build for current platform
make build

# Build for all platforms
make build-cross

# Run tests
make test

# Run linters and static analysis
make lint staticcheck

# Format code
make fmt

# Install to /usr/local/bin
make install
```

### Available Make Targets

**Building and Running:**
- `make build` - Build for current OS/architecture
- `make build-cross` - Build binaries for multiple platforms
- `make install` - Install binary to /usr/local/bin
- `make run ARGS="README.md"` - Run application locally

**Testing and Validation:**
- `make test` - Run all tests with coverage
- `make test-race` - Run tests with race detection
- `make benchmark` - Run performance benchmarks
- `make check-all` - Run all code quality checks

**Code Quality:**
- `make fmt` - Format Go code
- `make lint` - Run golangci-lint
- `make staticcheck` - Run static analysis
- `make vet` - Run go vet

**Development:**
- `make deps` - Install dependencies
- `make clean` - Clean build artifacts
- `make help` - Show all available targets

## CI/CD Integration

### GitHub Actions

```yaml
name: Markdown Format Check
on: [push, pull_request]

jobs:
  markdown-format:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Download mdfmt
        run: |
          curl -L -o mdfmt https://github.com/Gosayram/go-mdfmt/releases/latest/download/mdfmt-0.2.6-linux-amd64
          chmod +x mdfmt
          sudo mv mdfmt /usr/local/bin/
      
      - name: Check markdown formatting
        run: mdfmt --check --verbose .
```

### Pre-commit Hook

```bash
#!/bin/sh
# .git/hooks/pre-commit

markdown_files=$(git diff --cached --name-only --diff-filter=ACM | grep '\.md$')

if [ -n "$markdown_files" ]; then
    echo "Checking markdown formatting..."
    if ! mdfmt --check $markdown_files; then
        echo "Markdown files need formatting. Run: mdfmt --write $markdown_files"
        exit 1
    fi
fi
```

### Makefile Integration

```makefile
.PHONY: fmt-markdown check-markdown

fmt-markdown:
	mdfmt --write .

check-markdown:
	mdfmt --check --diff .

# Include in existing targets
fmt: fmt-go fmt-markdown
check: check-go check-markdown
```

## Architecture

go-mdfmt follows a modular architecture with clear separation of concerns:

**Core Components:**
- **Parser** (`pkg/parser`) - Goldmark-based Markdown AST parsing
- **Formatter** (`pkg/formatter`) - Rule-based formatting engine
- **Renderer** (`pkg/renderer`) - AST to Markdown conversion
- **Config** (`pkg/config`) - YAML-based configuration management
- **Processor** (`pkg/processor`) - File handling and batch processing

**Internal Modules:**
- **Version** (`internal/version`) - Build information and versioning

**CLI Interface:**
- **Main** (`cmd/mdfmt`) - Command-line interface and argument parsing

## Performance

mdfmt is designed for performance with large codebases:

- **Memory Efficient**: Processes files individually, not batch-loaded
- **Fast Parsing**: Uses goldmark parser for reliable AST generation
- **Concurrent Safe**: Supports concurrent file processing
- **Minimal Dependencies**: Small binary size with fast startup

**Benchmarks** (on test machine):
- Small files (< 10KB): ~1ms per file
- Medium files (10-100KB): ~5-15ms per file
- Large files (> 100KB): ~50-200ms per file

## Security

### Release Verification

All binaries are signed with Cosign for supply chain security:

```bash
# Download release files
curl -L -o mdfmt https://github.com/Gosayram/go-mdfmt/releases/download/v0.2.5/mdfmt-0.2.6-linux-amd64
curl -L -o mdfmt.sig https://github.com/Gosayram/go-mdfmt/releases/download/v0.2.5/mdfmt-0.2.6-linux-amd64.sig
curl -L -o cosign.pub https://github.com/Gosayram/go-mdfmt/releases/download/v0.2.5/cosign.pub

# Verify signature
cosign verify-blob --key cosign.pub --signature mdfmt.sig mdfmt

# Verify checksum
curl -L -o mdfmt.sha256 https://github.com/Gosayram/go-mdfmt/releases/download/v0.2.5/mdfmt-0.2.6-linux-amd64.sha256
sha256sum -c mdfmt.sha256
```

### Security Scanning

The project includes automated security scanning:
- **Dependency Review**: Automated dependency vulnerability scanning
- **CodeQL Analysis**: Static security analysis
- **OpenSSF Scorecard**: Supply chain security assessment
- **SLSA Compliance**: Software supply chain integrity

## Comparison with Alternatives

**vs Prettier:**
- **Specialization**: Purpose-built for Markdown vs universal formatter
- **Performance**: Native binary vs Node.js runtime overhead
- **Configuration**: Markdown-specific options vs generic formatting rules
- **Integration**: Better suited for Go/CLI environments

**vs Other Markdown Tools:**
- **Reliability**: Consistent AST-based parsing vs regex-based approaches
- **Speed**: Optimized for large codebases and CI/CD pipelines
- **Standards**: Follows CommonMark specification
- **Maintenance**: Active development with regular releases

## Contributing

Contributions are welcome. Please ensure all code follows the project standards:

1. **Code Quality**: All Go code must pass `make check-all`
2. **Testing**: Include tests for new functionality
3. **Documentation**: Update relevant documentation
4. **Commit Messages**: Use conventional commit format

**Development Setup:**
```bash
git clone https://github.com/Gosayram/go-mdfmt.git
cd go-mdfmt
make deps
make build
make test
```

## License

Licensed under the MIT License. See [LICENSE](LICENSE) for details.

## Support

- **Issues**: Report bugs and feature requests on [GitHub Issues](https://github.com/Gosayram/go-mdfmt/issues)
- **Documentation**: Additional documentation available in [docs/](docs/)
- **Security**: Report security issues privately to project maintainers 