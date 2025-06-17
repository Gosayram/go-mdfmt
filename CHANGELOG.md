# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.6] - 2025-01-18

### Fixed
- Corrected semantic versioning logic in auto-tag workflow to implement custom 9-limit versioning
- Fixed incorrect version jumps (0.2.4 → 0.3.0) by implementing proper version increment rules
- Restored correct version sequence with custom logic: patch limit of 9 before minor bump
- Enhanced release workflow to extract changelog information from CHANGELOG.md for better release notes

### Changed
- Implemented custom versioning scheme: when patch reaches 9, bump minor (0.2.9 → 0.3.0)
- Added intelligent changelog integration in release workflow to include structured release notes
- Updated version file to correct sequence (0.2.6) from incorrect jump to 0.3.0
- Enhanced auto-tag workflow with detailed logging for version bump decisions

### Added
- Custom versioning logic with examples: 0.1.9 → 0.2.0, 0.9.9 → 1.0.0
- Automatic CHANGELOG.md extraction for release notes generation
- Better debugging and logging in version calculation process

## [0.2.5] - 2025-01-18

### Fixed
- Corrected changelog generation in release workflow to properly identify previous tag
- Fixed release notes showing incorrect "commits since" information
- Improved tag sorting logic to use semantic versioning for accurate previous tag detection
- Enhanced changelog formatting with better commit message filtering

### Changed
- Updated release workflow to exclude merge commits from changelog
- Improved changelog generation with proper tag comparison logic
- Added current tag logging for better debugging in release process

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