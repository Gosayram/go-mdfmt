# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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