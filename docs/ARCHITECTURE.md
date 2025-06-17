# Architecture Overview

This document provides a comprehensive overview of the go-mdfmt architecture, components, and design decisions.

## System Architecture

go-mdfmt follows a modular, pipeline-based architecture that separates concerns and enables maintainable, testable code.

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   CLI       │───►│   Config    │───►│  Processor  │───►│   Output    │
│  Interface  │    │ Management  │    │  Pipeline   │    │  Handler    │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
                                              │
                                              ▼
                                    ┌─────────────────┐
                                    │                 │
                                    │   Core Engine   │
                                    │                 │
                                    └─────────────────┘
                                              │
                           ┌──────────────────┼──────────────────┐
                           ▼                  ▼                  ▼
                   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
                   │   Parser    │    │  Formatter  │    │  Renderer   │
                   │  (goldmark) │    │   Engine    │    │   Engine    │
                   └─────────────┘    └─────────────┘    └─────────────┘
```

## Core Components

### CLI Interface (`cmd/mdfmt`)

**Responsibility**: Command-line argument parsing, user interaction, and application entry point.

**Key Features**:
- Flag validation and mutual exclusivity checking
- Help and version information display
- Error handling and exit code management
- User input sanitization

**Constants Defined**:
```go
const (
    ExitCodeError         = 2
    ExitCodeChangesNeeded = 1
    OutputFilePermissions = 0o600
)
```

### Configuration Management (`pkg/config`)

**Responsibility**: Configuration file discovery, parsing, validation, and default value management.

**Key Features**:
- Hierarchical configuration file discovery
- YAML and JSON configuration support
- Comprehensive validation with descriptive error messages
- Default configuration generation

**Constants Defined**:
```go
const (
    DefaultLineWidth         = 80
    DefaultMaxBlankLines     = 2
    ConfigFilePermissions    = 0o600
)
```

**Configuration Structure**:
```go
type Config struct {
    LineWidth  int             `yaml:"line_width"`
    Heading    HeadingConfig   `yaml:"heading"`
    List       ListConfig      `yaml:"list"`
    Code       CodeConfig      `yaml:"code"`
    Whitespace WhitespaceConfig `yaml:"whitespace"`
    Files      FilesConfig     `yaml:"files"`
}
```

### File Processor (`pkg/processor`)

**Responsibility**: File system operations, batch processing, and file pattern matching.

**Key Features**:
- Recursive directory traversal
- Glob pattern matching for file inclusion/exclusion
- File metadata preservation
- Concurrent processing support

### Parser (`pkg/parser`)

**Responsibility**: Markdown document parsing using goldmark library.

**Key Features**:
- CommonMark specification compliance
- AST (Abstract Syntax Tree) generation
- Extension support for enhanced markdown features
- Error-tolerant parsing with recovery

**Design Decision**: Uses goldmark parser for reliable, specification-compliant parsing instead of regex-based approaches.

### Formatter (`pkg/formatter`)

**Responsibility**: Core formatting logic and rule application.

**Key Features**:
- Rule-based formatting engine
- Configurable formatting rules
- AST manipulation and normalization
- Content transformation without data loss

**Formatting Rules Applied**:
- Text reflow with configurable line width
- Heading level normalization
- List style consistency
- Code block fence standardization
- Whitespace cleanup

### Renderer (`pkg/renderer`)

**Responsibility**: Converting formatted AST back to Markdown text.

**Key Features**:
- AST to Markdown conversion
- Output formatting consistency
- Character encoding preservation
- Line ending normalization

### Version Management (`internal/version`)

**Responsibility**: Build information, version tracking, and release metadata.

**Constants Defined**:
```go
const (
    ShortCommitHashLength = 7
    UnknownValue         = "unknown"
)
```

**Build Variables** (set at compile time):
```go
var (
    Version     = "dev"
    Commit      = "unknown"
    Date        = "unknown"
    BuiltBy     = "unknown"
    BuildNumber = "0"
)
```

## Data Flow

### Input Processing Flow

1. **CLI Parsing**: Command-line arguments parsed and validated
2. **Configuration Loading**: Configuration files discovered and loaded
3. **File Discovery**: Input files/directories resolved to file list
4. **File Filtering**: Files filtered based on extensions and ignore patterns

### Formatting Pipeline

1. **Parse**: Markdown content parsed into AST using goldmark
2. **Format**: AST modified according to formatting rules
3. **Render**: Modified AST converted back to Markdown text
4. **Output**: Formatted content written to destination

### Error Handling

Each component implements comprehensive error handling:
- **Validation Errors**: Configuration and input validation with descriptive messages
- **File System Errors**: Graceful handling of permission and I/O errors
- **Parse Errors**: Error recovery during markdown parsing
- **Processing Errors**: Context-aware error reporting with file/line information

## Design Principles

### Modularity

Each package has a single, well-defined responsibility with minimal dependencies between components.

### Testability

Components are designed for unit testing with dependency injection and interface-based design where appropriate.

### Performance

- **Memory Efficiency**: Files processed individually to minimize memory usage
- **CPU Efficiency**: AST-based processing avoids expensive regex operations
- **I/O Efficiency**: Batch file operations where possible

### Security

- **Input Validation**: All user inputs validated before processing
- **File System Safety**: Proper permission checking and path sanitization
- **Memory Safety**: No unsafe operations or manual memory management

## Configuration System

### Configuration Discovery

Configuration files are searched in the following order:

1. Path specified by `--config` flag
2. `.mdfmt.yaml` in current directory
3. `.mdfmt.yaml` in parent directories (walking up to repository root)
4. Built-in default configuration

### Configuration Validation

All configuration options are validated at startup with comprehensive error messages:

```go
func (c *Config) Validate() error {
    if c.LineWidth < 1 {
        return fmt.Errorf("line_width must be greater than 0")
    }
    // Additional validations...
}
```

## Build System

### Make Targets

The build system provides comprehensive targets for all development activities:

**Development Targets**:
- `make build` - Build for current platform
- `make build-cross` - Cross-platform builds
- `make test` - Run test suite
- `make benchmark` - Performance benchmarks

**Quality Assurance**:
- `make lint` - Code linting with golangci-lint
- `make staticcheck` - Static analysis
- `make fmt` - Code formatting
- `make vet` - Go vet analysis

### Version Management

Build-time version information is injected using linker flags:

```makefile
LDFLAGS=-ldflags "-s -w \
    -X github.com/Gosayram/go-mdfmt/internal/version.Version=$(VERSION) \
    -X github.com/Gosayram/go-mdfmt/internal/version.Commit=$(COMMIT) \
    -X github.com/Gosayram/go-mdfmt/internal/version.Date=$(DATE) \
    -X github.com/Gosayram/go-mdfmt/internal/version.BuiltBy=$(BUILT_BY)"
```

## Testing Strategy

### Unit Tests

Each package includes comprehensive unit tests with table-driven test patterns:

```go
func TestFormatMarkdown(t *testing.T) {
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

### Integration Tests

Integration tests verify end-to-end functionality with real file processing scenarios.

### Benchmark Tests

Performance benchmarks measure processing speed for different file sizes and content types.

## Security Considerations

### Input Validation

All user inputs undergo validation:
- File paths sanitized to prevent directory traversal
- Configuration values validated against acceptable ranges
- Command-line arguments validated for type and format

### File System Access

File operations use appropriate permission checks and error handling:
- Read-only access for input files
- Controlled write access for output files
- Proper cleanup of temporary files

### Memory Safety

Go's memory safety features provide protection against buffer overflows and memory corruption, while the application avoids unsafe operations.

## Future Architecture Considerations

### Plugin System

Future versions may include a plugin architecture for custom formatting rules:

```go
type FormatterPlugin interface {
    Name() string
    Apply(node ast.Node, config *Config) error
}
```

### Language Server Protocol

Architecture supports future LSP implementation for editor integration:

```go
type LSPServer struct {
    formatter *Formatter
    config    *Config
}
```

### Web Interface

Component separation enables future web interface development with the same core formatting engine.

## Performance Characteristics

### Benchmarks

Typical performance on modern hardware:
- Small files (< 10KB): 1-2ms per file
- Medium files (10-100KB): 5-15ms per file
- Large files (> 100KB): 50-200ms per file

### Memory Usage

Memory usage scales linearly with file size, with typical usage:
- Base application: ~5MB
- Per file processing: ~2-3x file size in memory
- No memory leaks in long-running scenarios

### Concurrency

The architecture supports concurrent file processing with goroutines while maintaining thread safety for shared configuration and formatting logic. 