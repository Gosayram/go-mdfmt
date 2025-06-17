# Development Guide

This document provides comprehensive information for developers working on go-mdfmt, including setup, development workflows, testing, and contribution guidelines.

## Development Environment Setup

### Prerequisites

**Required Software**:
- Go 1.24.4 or later
- Git 2.30.0 or later
- Make (for build automation)

**Optional Tools**:
- Docker (for containerized testing)
- golangci-lint (for advanced linting)
- staticcheck (for static analysis)

### Repository Setup

```bash
# Clone the repository
git clone https://github.com/Gosayram/go-mdfmt.git
cd go-mdfmt

# Install dependencies and development tools
make deps
make install-tools

# Verify setup
make build
./bin/mdfmt --version
```

### Development Tools Installation

The project provides automated tool installation:

```bash
# Install all development tools
make install-tools

# Individual tool installation
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install golang.org/x/tools/cmd/goimports@latest
```

## Project Structure

```
go-mdfmt/
├── cmd/
│   └── mdfmt/              # CLI application entry point
├── pkg/
│   ├── config/             # Configuration management
│   ├── formatter/          # Core formatting logic
│   ├── parser/             # Markdown parsing (goldmark)
│   ├── processor/          # File processing pipeline
│   └── renderer/           # AST to Markdown conversion
├── internal/
│   └── version/            # Version information management
├── testdata/               # Test fixtures and examples
├── docs/                   # Project documentation
├── .github/
│   └── workflows/          # CI/CD automation
├── Makefile               # Build automation
├── go.mod                 # Go module definition
├── .mdfmt.yaml           # Project formatting configuration
├── .golangci.yml         # Linter configuration
└── .cursorrules          # Development standards
```

## Development Workflow

### Building the Project

**Basic Build**:
```bash
# Build for current platform
make build

# Build with debug symbols
make build-debug

# Cross-platform builds
make build-cross
```

**Build Output**:
- Binaries created in `bin/` directory
- Cross-platform binaries include platform suffix
- Debug builds include symbol information

### Running the Application

**Development Mode**:
```bash
# Run from source (development)
make run ARGS="README.md"
make dev ARGS="--write docs/"

# Run built binary
make run-built ARGS="--check ."
```

**Testing Configuration**:
```bash
# Test with current project configuration
make run ARGS="--config .mdfmt.yaml --verbose --diff ."

# Test with example configuration
make example-config
make run ARGS="--config .mdfmt.example.yaml --check ."
```

## Code Quality Standards

### Formatting and Style

Following .cursorrules requirements, all code must meet strict quality standards:

**Code Formatting**:
```bash
# Format Go code
make fmt

# Format imports
make imports

# Check formatting
go fmt ./...
goimports -w .
```

**Linting**:
```bash
# Run linters
make lint

# Run linters with auto-fix
make lint-fix

# Static analysis
make staticcheck
```

### Code Quality Checks

**Complete Quality Pipeline**:
```bash
# Run all quality checks
make check-all

# Individual checks
make vet          # Go vet analysis
make lint         # golangci-lint
make staticcheck  # Static analysis
make fmt          # Code formatting
make imports      # Import formatting
```

### Constants and Magic Numbers

Per .cursorrules requirements, all magic numbers must be replaced with named constants:

**Example Constant Definitions**:
```go
// pkg/config/config.go
const (
    DefaultLineWidth         = 80
    DefaultMaxBlankLines     = 2
    ConfigFilePermissions    = 0o600
)

// cmd/mdfmt/main.go
const (
    ExitCodeError         = 2
    ExitCodeChangesNeeded = 1
    OutputFilePermissions = 0o600
)
```

## Testing

### Test Categories

**Unit Tests**:
```bash
# Run all unit tests
make test

# Run tests with race detection
make test-race

# Run tests with coverage
make test-coverage
```

**Integration Tests**:
```bash
# Run integration tests
make test-integration

# Run all tests including benchmarks
make test-all
```

**Test Data Management**:
```bash
# Create test data copies
make test-data-copy

# Format test data (safe copies)
make test-data-format

# Check test data formatting
make test-data-check

# Show test data differences
make test-data-diff

# Clean test data artifacts
make test-data-clean
```

### Test Structure

Following Go testing best practices and .cursorrules requirements:

```go
// Example test structure
func TestFormatMarkdown(t *testing.T) {
    const (
        // Test constants
        ExpectedStatusCode = 200
        TestTimeout        = 5 * time.Second
        MaxTestRetries     = 3
    )
    
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "simple paragraph",
            input:    "Hello world",
            expected: "Hello world\n\n",
        },
        {
            name:     "long paragraph reflow",
            input:    "This is a very long paragraph that should be wrapped",
            expected: "This is a very long paragraph that should be\nwrapped\n\n",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := FormatMarkdown(tt.input)
            if result != tt.expected {
                t.Errorf("got %q, want %q", result, tt.expected)
            }
        })
    }
}
```

### Benchmarking

**Performance Testing**:
```bash
# Run basic benchmarks
make benchmark

# Run comprehensive benchmarks
make benchmark-long

# Run formatting-specific benchmarks
make benchmark-format

# Generate benchmark report
make benchmark-report
```

**Benchmark Structure**:
```go
// Example benchmark
func BenchmarkFormatMarkdown(b *testing.B) {
    const (
        SmallFileSize  = 1024    // 1KB
        MediumFileSize = 10240   // 10KB
        LargeFileSize  = 102400  // 100KB
    )
    
    testCases := []struct {
        name string
        size int
    }{
        {"small", SmallFileSize},
        {"medium", MediumFileSize},
        {"large", LargeFileSize},
    }
    
    for _, tc := range testCases {
        b.Run(tc.name, func(b *testing.B) {
            content := generateTestContent(tc.size)
            b.ResetTimer()
            
            for i := 0; i < b.N; i++ {
                FormatMarkdown(content)
            }
        })
    }
}
```

## Debugging

### Debug Builds

```bash
# Build with debug symbols
make build-debug

# Run with debugger
dlv exec ./bin/mdfmt-debug -- --write README.md
```

### Verbose Output

```bash
# Enable verbose logging
make run ARGS="--verbose --diff ."

# Debug configuration loading
make run ARGS="--config debug.yaml --verbose ."
```

### Profiling

```bash
# CPU profiling
go run -cpuprofile=cpu.prof ./cmd/mdfmt --write large-file.md

# Memory profiling
go run -memprofile=mem.prof ./cmd/mdfmt --write large-file.md

# Analyze profiles
go tool pprof cpu.prof
go tool pprof mem.prof
```

## Dependency Management

### Managing Dependencies

```bash
# Install dependencies
make deps

# Upgrade dependencies
make upgrade-deps

# Clean dependencies
make clean-deps

# Check for vulnerabilities
go list -json -m all | nancy sleuth
```

### Dependency Guidelines

**Dependency Selection Criteria**:
- Well-maintained with active development
- Minimal transitive dependencies
- Compatible licensing
- Security track record

**Current Core Dependencies**:
- `github.com/yuin/goldmark` - Markdown parsing
- `gopkg.in/yaml.v3` - YAML configuration

## Configuration Development

### Configuration Schema

**Adding New Configuration Options**:

1. **Update Config Struct** (`pkg/config/config.go`):
```go
type Config struct {
    // ... existing fields ...
    NewOption NewOptionConfig `yaml:"new_option" json:"new_option"`
}

type NewOptionConfig struct {
    Enabled bool   `yaml:"enabled" json:"enabled"`
    Value   string `yaml:"value" json:"value"`
}
```

2. **Update Default Configuration**:
```go
func Default() *Config {
    return &Config{
        // ... existing defaults ...
        NewOption: NewOptionConfig{
            Enabled: true,
            Value:   "default_value",
        },
    }
}
```

3. **Add Validation**:
```go
func (c *Config) Validate() error {
    // ... existing validation ...
    
    if c.NewOption.Enabled && c.NewOption.Value == "" {
        return fmt.Errorf("new_option.value cannot be empty when enabled")
    }
    
    return nil
}
```

### Configuration Testing

```bash
# Create test configuration
make example-config

# Validate configuration
make validate-config

# Test with custom configuration
make run ARGS="--config test.yaml --check ."
```

## Build System

### Makefile Targets

The build system provides comprehensive automation:

**Development Targets**:
- `make build` - Build for current platform
- `make build-cross` - Cross-platform builds  
- `make run` - Run from source
- `make dev` - Development mode

**Quality Assurance**:
- `make fmt` - Code formatting
- `make lint` - Linting
- `make staticcheck` - Static analysis
- `make test` - Run tests

**Utilities**:
- `make clean` - Clean artifacts
- `make deps` - Install dependencies
- `make help` - Show all targets

### Build Configuration

**Build Variables**:
```makefile
# Version management
TAG_NAME ?= $(shell head -n 1 .release-version 2>/dev/null || echo "v0.1.0")
VERSION ?= $(shell head -n 1 .release-version 2>/dev/null | sed 's/^v//' || echo "dev")

# Build information
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
BUILT_BY ?= $(shell git config user.name 2>/dev/null | tr ' ' '_' || echo "unknown")

# Linker flags
LDFLAGS=-ldflags "-s -w \
  -X github.com/Gosayram/go-mdfmt/internal/version.Version=$(VERSION) \
  -X github.com/Gosayram/go-mdfmt/internal/version.Commit=$(COMMIT) \
  -X github.com/Gosayram/go-mdfmt/internal/version.Date=$(DATE) \
  -X github.com/Gosayram/go-mdfmt/internal/version.BuiltBy=$(BUILT_BY)"
```

## Error Handling

### Error Handling Patterns

Following .cursorrules requirements for proper error handling:

**Error Wrapping**:
```go
func processFile(filename string) error {
    content, err := os.ReadFile(filename)
    if err != nil {
        return fmt.Errorf("failed to read file %s: %w", filename, err)
    }
    
    formatted, err := formatContent(content)
    if err != nil {
        return fmt.Errorf("failed to format file %s: %w", filename, err)
    }
    
    return nil
}
```

**Custom Error Types**:
```go
// ConfigValidationError represents configuration validation errors
type ConfigValidationError struct {
    Field   string
    Value   interface{}
    Message string
}

func (e *ConfigValidationError) Error() string {
    return fmt.Sprintf("configuration validation failed for field %s: %s", e.Field, e.Message)
}
```

### Logging

**Structured Logging**:
```go
const (
    LogLevelError = "ERROR"
    LogLevelWarn  = "WARN"
    LogLevelInfo  = "INFO"
    LogLevelDebug = "DEBUG"
)

func logError(msg string, err error) {
    fmt.Fprintf(os.Stderr, "[%s] %s: %v\n", LogLevelError, msg, err)
}
```

## Documentation

### Code Documentation

**Function Documentation** (following GoDoc standards):
```go
// FormatMarkdown formats the given markdown content according to the provided configuration.
// It returns the formatted content as a string and any error encountered during processing.
//
// The function performs the following operations:
// 1. Parses the markdown content into an AST
// 2. Applies formatting rules based on configuration
// 3. Renders the AST back to markdown text
//
// Example usage:
//   cfg := config.Default()
//   formatted, err := FormatMarkdown(content, cfg)
//   if err != nil {
//       return fmt.Errorf("formatting failed: %w", err)
//   }
func FormatMarkdown(content string, cfg *Config) (string, error) {
    // Implementation...
}
```

**Package Documentation**:
```go
// Package formatter provides core markdown formatting functionality.
// It implements a rule-based formatting engine that processes markdown
// Abstract Syntax Trees (AST) to apply consistent formatting rules.
package formatter
```

### Documentation Generation

```bash
# Generate API documentation
make docs-api

# Generate documentation
make docs

# View documentation locally
godoc -http=:6060
```

## Git Workflow

### Commit Message Format

Following conventional commit format:

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix  
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test additions/modifications
- `chore`: Maintenance tasks

**Examples**:
```
feat(formatter): add heading level normalization
fix(config): resolve validation error for empty files
docs(readme): update installation instructions
test(parser): add test cases for edge cases
```

### Branch Strategy

**Branch Types**:
- `main` - Production ready code
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring

**Workflow**:
1. Create feature branch from `main`
2. Implement changes with tests
3. Run quality checks locally
4. Submit pull request
5. Address review feedback
6. Merge after approval

## Performance Guidelines

### Performance Considerations

**Memory Management**:
- Process files individually to limit memory usage
- Use streaming where possible for large files
- Avoid loading entire directory trees into memory

**CPU Optimization**:
- Use goldmark parser for efficient AST operations
- Minimize regex usage in favor of AST traversal
- Cache expensive operations where appropriate

**I/O Optimization**:
- Batch file operations when possible
- Use buffered I/O for large files
- Minimize filesystem traversal

### Performance Testing

```bash
# Memory usage testing
go run -memprofile=mem.prof ./cmd/mdfmt large-directory/

# CPU profiling
go run -cpuprofile=cpu.prof ./cmd/mdfmt --write large-file.md

# Benchmark comparison
make benchmark > before.txt
# Make changes
make benchmark > after.txt
benchcmp before.txt after.txt
```

## Security Considerations

### Secure Development Practices

**Input Validation**:
- Validate all user inputs
- Sanitize file paths to prevent directory traversal
- Validate configuration values within acceptable ranges

**File Operations**:
- Use appropriate file permissions
- Validate file extensions before processing
- Handle symlinks securely

**Error Handling**:
- Avoid exposing sensitive information in error messages
- Log security events appropriately
- Handle resource exhaustion gracefully

### Security Testing

```bash
# Vulnerability scanning
make security-scan

# Dependency vulnerability check
go list -json -m all | nancy sleuth

# Static security analysis
gosec ./...
```

## Troubleshooting

### Common Development Issues

**Build Failures**:
```bash
# Clean and rebuild
make clean
make build

# Check Go version
go version

# Verify dependencies
go mod verify
```

**Test Failures**:
```bash
# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestSpecificFunction ./pkg/formatter

# Debug test with race detection
go test -race -v ./...
```

**Linting Issues**:
```bash
# Run linters with detailed output
golangci-lint run --verbose

# Auto-fix issues where possible
make lint-fix

# Check specific linter
golangci-lint run --enable=gofmt
```

### Debug Environment Setup

```bash
# Install debugging tools
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug application
dlv debug ./cmd/mdfmt -- --write README.md

# Debug tests
dlv test ./pkg/formatter
``` 