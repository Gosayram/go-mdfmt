# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.4] - 2025-01-18

### Fixed
- Fixed auto-tag workflow to create lightweight tags instead of annotated tags for better GitHub UI compatibility
- Resolved issue where new tags were not visible in GitHub branch/tag selector
- Corrected tag creation logic that was blocking tag generation when version already existed in `.release-version` file
- Enhanced tag creation process with proper git configuration and improved logging

### Changed
- Simplified auto-tag workflow conditions to ensure reliable tag creation
- Improved tag creation debugging with commit SHA and message logging
- Updated tag creation to use lightweight tags for better GitHub integration

## [0.2.2] - 2025-01-18

### Fixed
- Corrected GitHub Actions workflow configurations for improved reliability
- Enhanced security scanning workflows with proper SARIF upload handling
- Fixed Docker build triggers and multi-platform compilation matrix

### Changed
- Completely rewrote CI/CD documentation with comprehensive workflow descriptions
- Updated security implementation documentation with detailed procedures
- Enhanced release process documentation with verification instructions
- Improved auto-tag workflow with better version management logic

### Documentation
- Updated `docs/CI_CD.md` with detailed workflow architecture descriptions
- Added comprehensive security implementation documentation
- Enhanced troubleshooting section with common issues and solutions
- Documented dependency management and container strategy

## [0.2.1] - 2025-01-16

### Fixed
- Corrected Dockerfile FROM syntax for proper Docker build functionality
- Updated version management to 0.2.1 with skip-tag functionality

### Changed
- Enhanced CI testing workflows with improved comments and documentation
- Minor workflow optimizations for better reliability

## [0.1.0] - 2025-01-15

### Added
- Initial release with markdown formatting capabilities
- Link preservation during text wrapping
- Nested list numbering fixes
- Comprehensive test suite with testdata management
- Linter compliance with zero magic numbers policy
- Professional coding standards documentation 
- Complete CI/CD pipeline with GitHub Actions
- Automated security scanning (Trivy, Nancy, OpenSSF Scorecard)
- CodeQL static analysis for security vulnerabilities
- Dependency review for pull requests
- Cross-platform binary builds (Linux, macOS, Windows for amd64/arm64)
- Multi-platform Docker images with GitHub Container Registry
- Automated release process with changelog generation
- Version management with build-time injection
- Comprehensive test coverage reporting
- Docker containerization with minimal scratch-based images
- Make commands for CI/CD operations
- Security hardening with Step Security Harden Runner
- Pinned GitHub Actions for reproducible builds

### Enhanced
- Build system with version information injection
- Makefile with CI/CD support commands
- Documentation with CI/CD guide and release instructions
- Security posture with vulnerability scanning and analysis

### Infrastructure
- `.github/workflows/ci-lint-test.yml` - Main CI pipeline
- `.github/workflows/security.yml` - Security scanning
- `.github/workflows/codeql.yml` - Static code analysis
- `.github/workflows/dependency-review.yml` - Dependency security
- `.github/workflows/release.yml` - Automated releases
- `Dockerfile` - Multi-stage container build
- `internal/version/` - Version management package
- `docs/CI_CD.md` - Comprehensive CI/CD documentation
- `RELEASE.md` - Release process guide