# CI/CD Documentation

## Overview

This project uses GitHub Actions for continuous integration and deployment. The CI/CD pipeline includes automated testing, security scanning, building, and releasing.

## Workflows

### 1. CI - Lint and Test (`.github/workflows/ci-lint-test.yml`)

**Triggers:**
- Push to any branch
- Pull requests to any branch
- Daily at midnight (scheduled)

**Jobs:**
- **Lint and Test**: Runs linting, testing, and validation
  - Go version detection from `.go-version`
  - Dependency caching
  - Code formatting check
  - Unit tests with race detection
  - Coverage reporting to Codecov
  - golangci-lint analysis
  - staticcheck analysis
  - Integration tests
  - Test data validation

### 2. Security Scanning (`.github/workflows/security.yml`)

**Triggers:**
- Push to main branch
- Pull requests to main branch
- Daily at midnight (scheduled)

**Jobs:**
- **Trivy Security Scan**: Filesystem vulnerability scanning
- **Nancy Vulnerability Check**: Go dependency vulnerability scanning
- **OpenSSF Scorecard**: Security posture assessment

### 3. CodeQL Analysis (`.github/workflows/codeql.yml`)

**Triggers:**
- Push to main branch
- Pull requests to main branch
- Weekly on Mondays (scheduled)

**Jobs:**
- **Analyze**: Static code analysis for security vulnerabilities

### 4. Dependency Review (`.github/workflows/dependency-review.yml`)

**Triggers:**
- Pull requests to main branch

**Jobs:**
- **Dependency Review**: Reviews dependency changes for security issues

### 5. Release (`.github/workflows/release.yml`)

**Triggers:**
- Push tags matching `v*` pattern
- Manual workflow dispatch

**Jobs:**
- **Build and Test**: Pre-release validation
- **Build Binaries**: Cross-platform binary compilation
  - Linux (amd64, arm64)
  - macOS (amd64, arm64)
  - Windows (amd64, arm64)
- **Create Release**: GitHub release creation with changelog
- **Build Docker**: Multi-platform Docker image build and push

## Security Features

### Step Security Harden Runner
All workflows use `step-security/harden-runner` to:
- Disable sudo access
- Monitor network egress
- Provide security audit trails

### Pinned Action Versions
All GitHub Actions are pinned to specific SHA hashes for security and reproducibility.

### Minimal Permissions
Each workflow uses the principle of least privilege with minimal required permissions.

## Build Configuration

### Version Information
Build-time version information is injected via ldflags:
- Version from git tags or `.release-version`
- Commit SHA
- Build date
- Built by information

### Cross-Platform Builds
The release workflow builds binaries for:
- `linux/amd64`
- `linux/arm64`
- `darwin/amd64`
- `darwin/arm64`
- `windows/amd64`
- `windows/arm64`

### Docker Images
Multi-platform Docker images are built and pushed to GitHub Container Registry:
- `ghcr.io/gosayram/go-mdfmt:latest`
- `ghcr.io/gosayram/go-mdfmt:v<version>`

## Release Process

### Automatic Releases
1. Create and push a git tag: `git tag v1.0.0 && git push origin v1.0.0`
2. GitHub Actions automatically:
   - Runs all tests and security checks
   - Builds cross-platform binaries
   - Creates GitHub release with changelog
   - Builds and pushes Docker images

### Manual Releases
1. Go to GitHub Actions
2. Select "Release" workflow
3. Click "Run workflow"
4. Enter the desired tag (e.g., `v1.0.0`)

## Local Development

### Make Commands for CI
```bash
# Run CI linting checks
make ci-lint

# Run CI tests
make ci-test

# Run CI build
make ci-build

# Run complete CI pipeline
make ci-release
```

### Docker Development
```bash
# Build Docker image
make docker-build

# Run Docker image
make docker-run
```

## Configuration Files

### `.go-version`
Specifies the Go version used in CI/CD pipelines.

### `.release-version`
Contains the current release version (managed automatically).

### `Dockerfile`
Multi-stage Docker build configuration for minimal production images.

### `.dockerignore`
Excludes unnecessary files from Docker build context.

## Secrets Required

The following GitHub secrets are required for full CI/CD functionality:

- `GITHUB_TOKEN`: Automatically provided by GitHub
- `CODECOV_TOKEN`: For coverage reporting (optional)

## Monitoring and Notifications

### Security Alerts
- Trivy scan results are uploaded to GitHub Security tab
- Nancy vulnerability checks report to Security tab
- OpenSSF Scorecard results are published
- CodeQL analysis results appear in Security tab

### Build Status
- All workflows report status to GitHub checks
- Failed builds block pull request merging
- Release failures are visible in Actions tab

## Best Practices

### Branch Protection
Recommended branch protection rules for `main`:
- Require status checks to pass
- Require branches to be up to date
- Require review from code owners
- Restrict pushes to specific people/teams

### Security
- All dependencies are scanned for vulnerabilities
- Code is analyzed for security issues
- Build artifacts are signed and verified
- Minimal container images reduce attack surface

### Performance
- Dependency caching reduces build times
- Parallel job execution where possible
- Efficient Docker layer caching
- Minimal artifact retention periods 