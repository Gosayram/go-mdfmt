# go-mdfmt Makefile

# Project-specific variables
BINARY_NAME := mdfmt
OUTPUT_DIR := bin
CMD_DIR := cmd/mdfmt

TAG_NAME ?= $(shell head -n 1 .release-version 2>/dev/null || echo "v0.1.0")
VERSION ?= $(shell head -n 1 .release-version 2>/dev/null | sed 's/^v//' || echo "dev")
BUILD_INFO ?= $(shell date +%s)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GO_VERSION := $(shell cat .go-version 2>/dev/null || echo "1.24.2")
GO_FILES := $(wildcard $(CMD_DIR)/*.go internal/**/*.go pkg/**/*.go)
GOPATH ?= $(shell go env GOPATH)
GOLANGCI_LINT = $(GOPATH)/bin/golangci-lint
STATICCHECK = $(GOPATH)/bin/staticcheck
GOIMPORTS = $(GOPATH)/bin/goimports

# Build flags
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
BUILT_BY ?= $(shell git remote get-url origin 2>/dev/null | sed -n 's/.*[:/]\([^/]*\)\/[^/]*\.git.*/\1/p' || git config user.name 2>/dev/null | tr ' ' '_' || unknown)

# Linker flags for version information
LDFLAGS=-ldflags "-X github.com/Gosayram/go-mdfmt/internal/version.Version=$(VERSION) \
				  -X github.com/Gosayram/go-mdfmt/internal/version.Commit=$(COMMIT) \
				  -X github.com/Gosayram/go-mdfmt/internal/version.Date=$(DATE) \
				  -X github.com/Gosayram/go-mdfmt/internal/version.BuiltBy=$(BUILT_BY) \
				  -X github.com/Gosayram/go-mdfmt/internal/version.BuildNumber=$(BUILD_INFO)"

# Ensure the output directory exists
$(OUTPUT_DIR):
	@mkdir -p $(OUTPUT_DIR)

# Default target
.PHONY: default
default: fmt vet imports lint staticcheck build quicktest

# Display help information
.PHONY: help
help:
	@echo "go-mdfmt - Fast, Reliable Markdown Formatter"
	@echo ""
	@echo "Available targets:"
	@echo "  Building and Running:"
	@echo "  ===================="
	@echo "  default         - Run formatting, vetting, linting, staticcheck, build, and quick tests"
	@echo "  run             - Run the application locally"
	@echo "  dev             - Run in development mode"
	@echo "  build           - Build the application for the current OS/architecture"
	@echo "  build-debug     - Build debug version with debug symbols"
	@echo "  build-cross     - Build binaries for multiple platforms (Linux, macOS, Windows)"
	@echo "  install         - Install binary to /usr/local/bin"
	@echo "  uninstall       - Remove binary from /usr/local/bin"
	@echo ""
	@echo "  Testing and Validation:"
	@echo "  ======================"
	@echo "  test            - Run all tests with standard coverage"
	@echo "  test-with-race  - Run all tests with race detection and coverage"
	@echo "  quicktest       - Run quick tests without additional checks"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  test-race       - Run tests with race detection"
	@echo "  test-integration- Run integration tests"
	@echo "  test-all        - Run all tests and benchmarks"
	@echo ""
	@echo "  Benchmarking:"
	@echo "  ============="
	@echo "  benchmark       - Run basic benchmarks"
	@echo "  benchmark-long  - Run comprehensive benchmarks with longer duration"
	@echo "  benchmark-format- Run markdown formatting benchmarks"
	@echo "  benchmark-report- Generate a markdown report of all benchmarks"
	@echo ""
	@echo "  Code Quality:"
	@echo "  ============"
	@echo "  fmt             - Check and format Go code"
	@echo "  vet             - Analyze code with go vet"
	@echo "  imports         - Format imports with goimports"
	@echo "  lint            - Run golangci-lint"
	@echo "  lint-fix        - Run linters with auto-fix"
	@echo "  staticcheck     - Run staticcheck static analyzer"
	@echo "  check-all       - Run all code quality checks"
	@echo ""
	@echo "  Dependencies:"
	@echo "  ============="
	@echo "  deps            - Install project dependencies"
	@echo "  install-deps    - Install project dependencies (alias for deps)"
	@echo "  upgrade-deps    - Upgrade all dependencies to latest versions"
	@echo "  clean-deps      - Clean up dependencies"
	@echo "  install-tools   - Install development tools"
	@echo ""
	@echo "  Configuration:"
	@echo "  =============="
	@echo "  example-config  - Create example configuration file"
	@echo "  validate-config - Validate configuration file syntax"
	@echo ""
	@echo "  Version Management:"
	@echo "  =================="
	@echo "  version         - Show current version information"
	@echo "  bump-patch      - Bump patch version"
	@echo "  bump-minor      - Bump minor version"
	@echo "  bump-major      - Bump major version"
	@echo "  release         - Build release version with all optimizations"
	@echo ""
	@echo "  Cleanup:"
	@echo "  ========"
	@echo "  clean           - Clean build artifacts"
	@echo "  clean-coverage  - Clean coverage and benchmark files"
	@echo "  clean-all       - Clean everything including dependencies"
	@echo ""
	@echo "  Documentation:"
	@echo "  =============="
	@echo "  docs            - Generate documentation"
	@echo "  docs-api        - Generate API documentation"
	@echo ""
	@echo "Examples:"
	@echo "  make build                    - Build the binary"
	@echo "  make test                     - Run all tests"
	@echo "  make build-cross              - Build for multiple platforms"
	@echo "  make run ARGS=\"README.md\"     - Run with arguments"
	@echo "  make example-config           - Create .mdfmt.example.yaml"
	@echo ""
	@echo "For CLI usage instructions, run: ./bin/mdfmt --help"

# Build and run the application locally
.PHONY: run
run:
	@echo "Running $(BINARY_NAME)..."
	go run ./$(CMD_DIR) $(ARGS)

# Dependencies
.PHONY: deps install-deps upgrade-deps clean-deps install-tools
deps: install-deps

install-deps:
	@echo "Installing Go dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed successfully"

upgrade-deps:
	@echo "Upgrading all dependencies to latest versions..."
	go get -u ./...
	go mod tidy
	@echo "Dependencies upgraded. Please test thoroughly before committing!"

clean-deps:
	@echo "Cleaning up dependencies..."
	rm -rf vendor

install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "Development tools installed successfully"

# Build targets
.PHONY: build build-debug build-cross

build: $(OUTPUT_DIR)
	@echo "Building $(BINARY_NAME) with version $(VERSION)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build \
		$(LDFLAGS) \
		-o $(OUTPUT_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

build-debug: $(OUTPUT_DIR)
	@echo "Building debug version..."
	CGO_ENABLED=0 go build \
		-gcflags="all=-N -l" \
		$(LDFLAGS) \
		-o $(OUTPUT_DIR)/$(BINARY_NAME)-debug ./$(CMD_DIR)

build-cross: $(OUTPUT_DIR)
	@echo "Building cross-platform binaries..."
	GOOS=linux   GOARCH=amd64   CGO_ENABLED=0 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)
	GOOS=darwin  GOARCH=arm64   CGO_ENABLED=0 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(CMD_DIR)
	GOOS=darwin  GOARCH=amd64   CGO_ENABLED=0 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(CMD_DIR)
	GOOS=windows GOARCH=amd64   CGO_ENABLED=0 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)
	@echo "Cross-platform binaries are available in $(OUTPUT_DIR):"
	@ls -1 $(OUTPUT_DIR)

# Development targets
.PHONY: dev run-built

dev:
	@echo "Running in development mode..."
	go run ./$(CMD_DIR) $(ARGS)

run-built: build
	./$(OUTPUT_DIR)/$(BINARY_NAME) $(ARGS)

# Testing
.PHONY: test test-with-race quicktest test-coverage test-race test-integration test-all

test:
	@echo "Running Go tests..."
	go test -v ./... -cover

test-with-race:
	@echo "Running all tests with race detection and coverage..."
	go test -v -race -cover ./...

quicktest:
	@echo "Running quick tests..."
	go test ./...

test-coverage:
	@echo "Running tests with coverage report..."
	go test -v -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race:
	@echo "Running tests with race detection..."
	go test -v -race ./...

test-integration: build
	@echo "Running integration tests..."
	# Create test markdown files
	@mkdir -p testdata/integration
	@echo "# Test Document" > testdata/integration/test.md
	@echo "" >> testdata/integration/test.md
	@echo "This is a test document with inconsistent formatting." >> testdata/integration/test.md
	@echo "" >> testdata/integration/test.md
	@echo "" >> testdata/integration/test.md
	@echo "" >> testdata/integration/test.md
	@echo "* Item 1" >> testdata/integration/test.md
	@echo "+ Item 2" >> testdata/integration/test.md
	@echo "- Item 3" >> testdata/integration/test.md
	# Test basic formatting
	./$(OUTPUT_DIR)/$(BINARY_NAME) --check testdata/integration/test.md || echo "Expected formatting differences"
	# Test with write mode
	./$(OUTPUT_DIR)/$(BINARY_NAME) --write testdata/integration/test.md
	@echo "Integration tests completed"

test-all: test-coverage test-race benchmark
	@echo "All tests and benchmarks completed"

# Benchmark targets
.PHONY: benchmark benchmark-long benchmark-format benchmark-report

benchmark:
	@echo "Running benchmarks..."
	go test -v -bench=. -benchmem ./...

benchmark-long:
	@echo "Running comprehensive benchmarks (longer duration)..."
	go test -v -bench=. -benchmem -benchtime=5s ./...

benchmark-format: build
	@echo "Running markdown formatting benchmarks..."
	@mkdir -p testdata/benchmark
	@echo "# Large Document" > testdata/benchmark/large.md
	@echo "" >> testdata/benchmark/large.md
	@echo "This is a large document for benchmarking." >> testdata/benchmark/large.md
	@for i in $$(seq 1 1000); do echo "This is paragraph $$i with some text that needs to be wrapped at 80 characters to test the formatting performance." >> testdata/benchmark/large.md; done
	time ./$(OUTPUT_DIR)/$(BINARY_NAME) testdata/benchmark/large.md > /dev/null
	@echo "Format benchmarks completed"

benchmark-report:
	@echo "Generating benchmark report..."
	@echo "# Benchmark Results" > benchmark-report.md
	@echo "\nGenerated on \`$$(date)\`\n" >> benchmark-report.md
	@echo "## Go Benchmarks" >> benchmark-report.md
	@echo "| Test | Iterations | Time/op | Memory/op | Allocs/op |" >> benchmark-report.md
	@echo "|------|------------|---------|-----------|-----------|" >> benchmark-report.md
	@go test -bench=. -benchmem ./... 2>/dev/null | grep "Benchmark" | awk '{print "| " $$1 " | " $$2 " | " $$3 " | " $$5 " | " $$7 " |"}' >> benchmark-report.md
	@echo "Benchmark report generated: benchmark-report.md"

# Code quality
.PHONY: fmt vet imports lint staticcheck check-all

fmt:
	@echo "Checking and formatting code..."
	@go fmt ./...
	@echo "Code formatting completed"

vet:
	@echo "Running go vet..."
	go vet ./...

# Run goimports
.PHONY: imports
imports:
	@if command -v $(GOIMPORTS) >/dev/null 2>&1; then \
		echo "Running goimports..."; \
		$(GOIMPORTS) -local github.com/Gosayram/go-mdfmt -w $(GO_FILES); \
		echo "Imports formatting completed!"; \
	else \
		echo "goimports is not installed. Installing..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
		echo "Running goimports..."; \
		$(GOIMPORTS) -local github.com/Gosayram/go-mdfmt -w $(GO_FILES); \
		echo "Imports formatting completed!"; \
	fi

# Run linter
.PHONY: lint
lint:
	@if command -v $(GOLANGCI_LINT) >/dev/null 2>&1; then \
		echo "Running linter..."; \
		$(GOLANGCI_LINT) run; \
		echo "Linter completed!"; \
	else \
		echo "golangci-lint is not installed. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		echo "Running linter..."; \
		$(GOLANGCI_LINT) run; \
		echo "Linter completed!"; \
	fi

# Run staticcheck tool
.PHONY: staticcheck
staticcheck:
	@if command -v $(STATICCHECK) >/dev/null 2>&1; then \
		echo "Running staticcheck..."; \
		$(STATICCHECK) ./...; \
		echo "Staticcheck completed!"; \
	else \
		echo "staticcheck is not installed. Installing..."; \
		go install honnef.co/go/tools/cmd/staticcheck@latest; \
		echo "Running staticcheck..."; \
		$(STATICCHECK) ./...; \
		echo "Staticcheck completed!"; \
	fi

.PHONY: lint-fix
lint-fix:
	@echo "Running linters with auto-fix..."
	@$(GOLANGCI_LINT) run --fix
	@echo "Auto-fix completed"

check-all: fmt vet imports lint staticcheck
	@echo "All code quality checks completed"

# Configuration targets
.PHONY: example-config validate-config

example-config:
	@echo "Creating example configuration file..."
	@echo "# go-mdfmt configuration file" > .mdfmt.example.yaml
	@echo "line_width: 80" >> .mdfmt.example.yaml
	@echo "" >> .mdfmt.example.yaml
	@echo "heading:" >> .mdfmt.example.yaml
	@echo "  style: \"atx\"              # atx (#) or setext (===)" >> .mdfmt.example.yaml
	@echo "  normalize_levels: true    # Fix heading level jumps" >> .mdfmt.example.yaml
	@echo "" >> .mdfmt.example.yaml
	@echo "list:" >> .mdfmt.example.yaml
	@echo "  bullet_style: \"-\"         # -, *, or +" >> .mdfmt.example.yaml
	@echo "  number_style: \".\"         # . or )" >> .mdfmt.example.yaml
	@echo "  consistent_indentation: true" >> .mdfmt.example.yaml
	@echo "" >> .mdfmt.example.yaml
	@echo "code:" >> .mdfmt.example.yaml
	@echo "  fence_style: \"\`\`\`\"        # \`\`\` or ~~~" >> .mdfmt.example.yaml
	@echo "  language_detection: true  # Auto-detect and add language labels" >> .mdfmt.example.yaml
	@echo "" >> .mdfmt.example.yaml
	@echo "whitespace:" >> .mdfmt.example.yaml
	@echo "  max_blank_lines: 2        # Maximum consecutive blank lines" >> .mdfmt.example.yaml
	@echo "  trim_trailing_spaces: true" >> .mdfmt.example.yaml
	@echo "  ensure_final_newline: true" >> .mdfmt.example.yaml
	@echo "" >> .mdfmt.example.yaml
	@echo "files:" >> .mdfmt.example.yaml
	@echo "  extensions: [\".md\", \".markdown\", \".mdown\"]" >> .mdfmt.example.yaml
	@echo "  ignore_patterns: [\"node_modules/**\", \".git/**\", \"vendor/**\"]" >> .mdfmt.example.yaml
	@echo "Example configuration created as .mdfmt.example.yaml"

validate-config: build
	@echo "Validating configuration file..."
	@if [ -f .mdfmt.yaml ]; then \
		./$(OUTPUT_DIR)/$(BINARY_NAME) --config .mdfmt.yaml --help > /dev/null && echo "✓ .mdfmt.yaml is valid"; \
	elif [ -f .mdfmt.example.yaml ]; then \
		./$(OUTPUT_DIR)/$(BINARY_NAME) --config .mdfmt.example.yaml --help > /dev/null && echo "✓ .mdfmt.example.yaml is valid"; \
	else \
		echo "No configuration file found to validate"; \
	fi

# Release and installation
.PHONY: release install uninstall

release: test lint staticcheck
	@echo "Building release version $(VERSION)..."
	@mkdir -p $(OUTPUT_DIR)
	CGO_ENABLED=0 go build \
		$(LDFLAGS) \
		-ldflags="-s -w" \
		-o $(OUTPUT_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Release build completed: $(OUTPUT_DIR)/$(BINARY_NAME)"

install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(OUTPUT_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installation completed"

uninstall:
	@echo "Removing $(BINARY_NAME) from /usr/local/bin..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstallation completed"

# Cleanup
.PHONY: clean clean-coverage clean-all

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(OUTPUT_DIR)
	rm -f coverage.out coverage.html benchmark-report.md
	rm -rf testdata/integration testdata/benchmark
	go clean -cache
	@echo "Cleanup completed"

clean-coverage:
	@echo "Cleaning coverage and benchmark files..."
	rm -f coverage.out coverage.html benchmark-report.md
	@echo "Coverage files cleaned"

clean-all: clean clean-deps
	@echo "Deep cleaning everything including dependencies..."
	go clean -modcache
	@echo "Deep cleanup completed"

# Version management
.PHONY: version bump-patch bump-minor bump-major

version:
	@echo "Project: go-mdfmt"
	@echo "Go version: $(GO_VERSION)"
	@echo "Release version: $(VERSION)"
	@echo "Tag name: $(TAG_NAME)"
	@echo "Build target: $(GOOS)/$(GOARCH)"
	@echo "Commit: $(COMMIT)"
	@echo "Built by: $(BUILT_BY)"
	@echo "Build info: $(BUILD_INFO)"

bump-patch:
	@if [ ! -f .release-version ]; then echo "0.1.0" > .release-version; fi
	@current=$$(cat .release-version); \
	new=$$(echo $$current | awk -F. '{$$3=$$3+1; print $$1"."$$2"."$$3}'); \
	echo $$new > .release-version; \
	echo "Version bumped from $$current to $$new"

bump-minor:
	@if [ ! -f .release-version ]; then echo "0.1.0" > .release-version; fi
	@current=$$(cat .release-version); \
	new=$$(echo $$current | awk -F. '{$$2=$$2+1; $$3=0; print $$1"."$$2"."$$3}'); \
	echo $$new > .release-version; \
	echo "Version bumped from $$current to $$new"

bump-major:
	@if [ ! -f .release-version ]; then echo "0.1.0" > .release-version; fi
	@current=$$(cat .release-version); \
	new=$$(echo $$current | awk -F. '{$$1=$$1+1; $$2=0; $$3=0; print $$1"."$$2"."$$3}'); \
	echo $$new > .release-version; \
	echo "Version bumped from $$current to $$new"

# Docker support
.PHONY: docker-build docker-run

docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(TAG_NAME) .
	@echo "Docker image built: $(BINARY_NAME):$(TAG_NAME)"

docker-run:
	@echo "Running Docker image..."
	docker run -it --rm $(BINARY_NAME):$(TAG_NAME) --version

# Documentation
.PHONY: docs docs-api

docs:
	@echo "Generating documentation..."
	@mkdir -p docs
	@echo "# go-mdfmt API Documentation" > docs/api.md
	@echo "\nGenerated on \`$$(date)\`\n" >> docs/api.md
	@go doc -all ./... >> docs/api.md
	@echo "Documentation generated in docs/api.md"

docs-api:
	@echo "Generating API documentation..."
	@mkdir -p docs
	go doc -all ./... > docs/api.md
	@echo "API documentation generated" 