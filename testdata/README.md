# Test Data

This directory contains test files for mdfmt development and testing.

## Files

- `test_complex.md` - Complex markdown with nested lists, inline formatting, and code blocks
- `test_links.md` - Test file with various link formats
- `test_simple_link.md` - Simple test with a long link
- `test_debug.md` - Debug test file

## Usage

Use the Makefile commands to safely test with these files:

```bash
# Run all tests on copies (safe)
make test-data

# Create copies for manual testing
make test-data-copy

# Format test files (copies only)
make test-data-format

# Check if files need formatting
make test-data-check

# Show differences
make test-data-diff

# Clean up test copies and results
make test-data-clean
```

## Safety

The Makefile commands always work with copies in `testdata/copies/` to preserve the original test files. Results are saved in `testdata/results/` for inspection.

## Adding New Tests

1. Add new `.md` files to this directory
2. Update the Makefile `test-data` target if needed
3. Run `make test-data` to test the new files 