# Configuration Guide

This document provides comprehensive information about configuring go-mdfmt for your project needs.

## Configuration Overview

go-mdfmt uses YAML-based configuration files with intelligent discovery and validation. The configuration system provides sensible defaults while allowing complete customization of formatting behavior.

## Configuration Discovery

Configuration files are discovered automatically in the following order:

1. **Explicit Path**: File specified with `--config` flag
2. **Current Directory**: `.mdfmt.yaml` in the working directory
3. **Parent Directories**: Walking up the directory tree to find `.mdfmt.yaml`
4. **Built-in Defaults**: Comprehensive default configuration

### Supported Configuration File Names

The following file names are recognized (in order of precedence):

- `.mdfmt.yaml`
- `.mdfmt.yml`  
- `.mdfmt.json`
- `mdfmt.yaml`
- `mdfmt.yml`
- `mdfmt.json`

## Complete Configuration Reference

### Basic Configuration Structure

```yaml
# Maximum line width for text reflow
line_width: 80

# Heading formatting configuration
heading:
  style: "atx"
  normalize_levels: true

# List formatting configuration  
list:
  bullet_style: "-"
  number_style: "."
  consistent_indentation: true

# Code block formatting configuration
code:
  fence_style: "```"
  language_detection: true

# Whitespace handling configuration
whitespace:
  max_blank_lines: 2
  trim_trailing_spaces: true
  ensure_final_newline: true

# File processing configuration
files:
  extensions: [".md", ".markdown", ".mdown"]
  ignore_patterns: ["node_modules/**", ".git/**", "vendor/**"]
```

## Configuration Options

### Line Width (`line_width`)

**Type**: Integer  
**Default**: `80`  
**Valid Range**: Greater than 0

Controls the maximum line width for paragraph text reflow. Longer paragraphs will be wrapped to fit within this limit.

**Examples**:
```yaml
line_width: 80   # Standard terminal width
line_width: 100  # Wider format for modern displays
line_width: 120  # Maximum readable width
```

### Heading Configuration (`heading`)

Controls heading formatting and normalization behavior.

#### Style (`heading.style`)

**Type**: String  
**Default**: `"atx"`  
**Valid Values**: `"atx"`, `"setext"`

Defines the heading style to use:
- `"atx"`: Use hash symbols (`# Heading`)
- `"setext"`: Use underline style for H1/H2 (`Heading\n=======`)

```yaml
heading:
  style: "atx"     # Use # ## ### style
  # OR
  style: "setext"  # Use underline style
```

#### Normalize Levels (`heading.normalize_levels`)

**Type**: Boolean  
**Default**: `true`

When enabled, fixes heading level jumps (e.g., H1 directly to H3 becomes H1 to H2).

```yaml
heading:
  normalize_levels: true   # Fix heading level jumps
  normalize_levels: false  # Preserve original levels
```

### List Configuration (`list`)

Controls formatting of bulleted and numbered lists.

#### Bullet Style (`list.bullet_style`)

**Type**: String  
**Default**: `"-"`  
**Valid Values**: `"-"`, `"*"`, `"+"`

Defines the character used for bullet list items.

```yaml
list:
  bullet_style: "-"  # Use hyphens
  bullet_style: "*"  # Use asterisks  
  bullet_style: "+"  # Use plus signs
```

#### Number Style (`list.number_style`)

**Type**: String  
**Default**: `"."`  
**Valid Values**: `"."`, `")"`

Defines the punctuation used for numbered list items.

```yaml
list:
  number_style: "."  # 1. 2. 3.
  number_style: ")"  # 1) 2) 3)
```

#### Consistent Indentation (`list.consistent_indentation`)

**Type**: Boolean  
**Default**: `true`

When enabled, ensures consistent indentation for nested list items.

```yaml
list:
  consistent_indentation: true   # Standardize indentation
  consistent_indentation: false  # Preserve original indentation
```

### Code Block Configuration (`code`)

Controls formatting of code blocks and inline code.

#### Fence Style (`code.fence_style`)

**Type**: String  
**Default**: `"```"`  
**Valid Values**: `"```"`, `"~~~"`

Defines the fence character sequence for code blocks.

```yaml
code:
  fence_style: "```"  # Use backticks
  fence_style: "~~~"  # Use tildes
```

#### Language Detection (`code.language_detection`)

**Type**: Boolean  
**Default**: `true`

When enabled, attempts to detect and add language identifiers to code blocks.

```yaml
code:
  language_detection: true   # Auto-detect languages
  language_detection: false  # Preserve existing language tags only
```

### Whitespace Configuration (`whitespace`)

Controls whitespace handling and cleanup behavior.

#### Maximum Blank Lines (`whitespace.max_blank_lines`)

**Type**: Integer  
**Default**: `2`  
**Valid Range**: Greater than or equal to 0

Limits the number of consecutive blank lines allowed in the document.

```yaml
whitespace:
  max_blank_lines: 0  # No blank lines allowed
  max_blank_lines: 1  # Maximum one blank line
  max_blank_lines: 2  # Default: maximum two blank lines
```

#### Trim Trailing Spaces (`whitespace.trim_trailing_spaces`)

**Type**: Boolean  
**Default**: `true`

When enabled, removes trailing whitespace from all lines.

```yaml
whitespace:
  trim_trailing_spaces: true   # Remove trailing spaces
  trim_trailing_spaces: false  # Preserve trailing spaces
```

#### Ensure Final Newline (`whitespace.ensure_final_newline`)

**Type**: Boolean  
**Default**: `true`

When enabled, ensures files end with a newline character.

```yaml
whitespace:
  ensure_final_newline: true   # Add final newline
  ensure_final_newline: false  # Preserve original ending
```

### File Processing Configuration (`files`)

Controls which files are processed and which are ignored.

#### Extensions (`files.extensions`)

**Type**: Array of Strings  
**Default**: `[".md", ".markdown", ".mdown"]`

Defines file extensions that should be processed as Markdown files.

```yaml
files:
  extensions: [".md", ".markdown", ".mdown", ".mkd"]
```

#### Ignore Patterns (`files.ignore_patterns`)

**Type**: Array of Strings  
**Default**: `["node_modules/**", ".git/**", "vendor/**"]`

Defines glob patterns for files and directories to ignore during processing.

```yaml
files:
  ignore_patterns:
    - "node_modules/**"
    - ".git/**"
    - "vendor/**"
    - "build/**"
    - "dist/**"
    - "*.tmp"
```

## Configuration Examples

### Minimal Configuration

```yaml
line_width: 100
```

### Documentation Project Configuration

```yaml
line_width: 100
heading:
  style: "atx"
  normalize_levels: true
list:
  bullet_style: "-"
  number_style: "."
  consistent_indentation: true
code:
  fence_style: "```"
  language_detection: true
whitespace:
  max_blank_lines: 1
  trim_trailing_spaces: true
  ensure_final_newline: true
files:
  extensions: [".md", ".markdown"]
  ignore_patterns: 
    - "node_modules/**"
    - ".git/**"
    - "build/**"
```

### Blog/Content Configuration

```yaml
line_width: 80
heading:
  style: "atx"
  normalize_levels: false  # Preserve author's heading structure
list:
  bullet_style: "*"
  number_style: "."
  consistent_indentation: true
code:
  fence_style: "```"
  language_detection: true
whitespace:
  max_blank_lines: 3  # Allow more spacing for readability
  trim_trailing_spaces: true
  ensure_final_newline: true
files:
  extensions: [".md", ".markdown", ".mdown"]
  ignore_patterns:
    - "node_modules/**"
    - ".git/**"
    - "drafts/**"
    - "_site/**"
```

### Technical Documentation Configuration

```yaml
line_width: 120
heading:
  style: "atx"
  normalize_levels: true
list:
  bullet_style: "-"
  number_style: "."
  consistent_indentation: true
code:
  fence_style: "```"
  language_detection: true
whitespace:
  max_blank_lines: 2
  trim_trailing_spaces: true
  ensure_final_newline: true
files:
  extensions: [".md", ".markdown", ".mdown"]
  ignore_patterns:
    - "node_modules/**"
    - ".git/**"
    - "vendor/**"
    - "target/**"
    - "*.tmp"
    - "test-output/**"
```

## Configuration Validation

All configuration values are validated when the configuration is loaded. Invalid configurations will result in descriptive error messages.

### Common Validation Errors

**Invalid line width**:
```
Error: line_width must be greater than 0
```

**Invalid heading style**:
```
Error: heading.style must be 'atx' or 'setext'
```

**Invalid bullet style**:
```
Error: list.bullet_style must be '-', '*', or '+'
```

**Invalid number style**:
```
Error: list.number_style must be '.' or ')'
```

**Invalid fence style**:
```
Error: code.fence_style must be '```' or '~~~'
```

**Invalid max blank lines**:
```
Error: whitespace.max_blank_lines must be >= 0
```

## Using Configuration Files

### Command Line Usage

```bash
# Use specific configuration file
mdfmt --config .mdfmt-custom.yaml --write docs/

# Use default configuration discovery
mdfmt --write docs/

# Check configuration with verbose output
mdfmt --config .mdfmt.yaml --verbose --check docs/
```

### Generating Example Configuration

Create a complete example configuration file:

```bash
# Create example configuration with all options
cat > .mdfmt.yaml << 'EOF'
line_width: 80
heading:
  style: "atx"
  normalize_levels: true
list:
  bullet_style: "-"
  number_style: "."
  consistent_indentation: true
code:
  fence_style: "```"
  language_detection: true
whitespace:
  max_blank_lines: 2
  trim_trailing_spaces: true
  ensure_final_newline: true
files:
  extensions: [".md", ".markdown", ".mdown"]
  ignore_patterns: ["node_modules/**", ".git/**", "vendor/**"]
EOF
```

### Project-Specific Configuration

For projects with specific formatting requirements, place `.mdfmt.yaml` in the repository root:

```yaml
# Project-specific formatting rules
line_width: 100

# Use consistent bullet style across team
list:
  bullet_style: "-"
  number_style: "."

# Enforce strict whitespace rules
whitespace:
  max_blank_lines: 1
  trim_trailing_spaces: true
  ensure_final_newline: true

# Ignore build artifacts and dependencies
files:
  ignore_patterns:
    - "node_modules/**"
    - ".git/**"
    - "build/**"
    - "dist/**"
    - "coverage/**"
```

## Default Configuration

When no configuration file is found, go-mdfmt uses these default values:

```yaml
line_width: 80
heading:
  style: "atx"
  normalize_levels: true
list:
  bullet_style: "-"
  number_style: "."
  consistent_indentation: true
code:
  fence_style: "```"
  language_detection: true
whitespace:
  max_blank_lines: 2
  trim_trailing_spaces: true
  ensure_final_newline: true
files:
  extensions: [".md", ".markdown", ".mdown"]
  ignore_patterns: ["node_modules/**", ".git/**", "vendor/**"]
```

## Troubleshooting Configuration

### Configuration Not Found

If configuration files are not being discovered:

1. Verify file name matches supported patterns
2. Check file permissions (must be readable)
3. Ensure YAML syntax is valid
4. Use `--config` flag to specify exact path

### Configuration Validation Errors

If configuration validation fails:

1. Check all required fields are present
2. Verify data types match expected values
3. Ensure values are within valid ranges
4. Review error messages for specific issues

### Performance Considerations

Configuration choices can affect performance:

- **Line Width**: Shorter line widths may require more processing for text reflow
- **Language Detection**: Disabling can improve performance for large codebases
- **File Extensions**: Limiting extensions reduces file scanning overhead
- **Ignore Patterns**: Effective patterns prevent processing unnecessary files 