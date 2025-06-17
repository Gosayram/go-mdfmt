# CI/CD Documentation

## Overview

This project implements a comprehensive CI/CD pipeline using GitHub Actions for continuous integration, security scanning, automated dependency management, and automated releases. The pipeline ensures code quality, security, and reliable deployment through multiple automated workflows.

## Workflow Architecture

### Core Workflows

#### 1. CI - Lint and Test (`ci-lint-test.yml`)

**Purpose**: Primary continuous integration workflow for code quality assurance and testing.

**Triggers:**
- Push to main branch
- Pull requests to main branch
- Daily schedule at midnight UTC

**Job: lint_and_test**

The workflow performs comprehensive code validation including:

- **Security Hardening**: Uses Step Security Harden Runner with egress monitoring
- **Go Environment Setup**: Automatic Go version detection from `.go-version` file
- **Dependency Management**: Module download, verification, and caching
- **Build Validation**: Compilation check with `make build`
- **Testing Suite**:
  - Standard test execution with `make test`
  - Race condition detection with `make test-race`
  - Coverage analysis with `make test-coverage`
  - Integration testing with `make test-integration`
- **Code Quality Analysis**:
  - golangci-lint analysis with timeout configuration
  - Staticcheck static analysis
  - Import formatting verification
- **Coverage Reporting**: Upload to Codecov with error tolerance
- **Test Data Validation**: Verification of test fixtures

**Security Features:**
- Pinned action versions with SHA hashes
- Minimal required permissions
- Network egress monitoring
- Disabled sudo access

#### 2. Security Scanning (`security.yml`)

**Purpose**: Multi-layered security vulnerability detection and compliance monitoring.

**Triggers:**
- Push to main branch
- Pull requests to main branch
- Branch protection rule events
- Weekly schedule on Tuesdays at 07:20 UTC

**Jobs:**

**trivy_scan**
- Filesystem vulnerability scanning with Trivy
- SARIF format output for GitHub Security tab integration
- Critical and high severity focus
- Automated SARIF upload to GitHub Code Scanning

**nancy_check**
- Go module dependency vulnerability analysis
- Nancy vulnerability database integration
- Go.list generation for comprehensive scanning
- Security-events permission for vulnerability reporting

**ossf_scorecard**
- OpenSSF Scorecard security posture assessment
- Repository security practice evaluation
- Artifact retention with 5-day policy
- SARIF upload for centralized security reporting

#### 3. CodeQL Analysis (`codeql.yml`)

**Purpose**: Static application security testing (SAST) for code vulnerabilities.

**Triggers:**
- Push to main branch
- Pull requests to main branch
- Weekly schedule on Mondays

**Analysis Features:**
- Go language static analysis
- Automated build detection
- Security vulnerability identification
- Results integration with GitHub Security tab
- Categorized analysis results

#### 4. Dependency Review (`dependency-review.yml`)

**Purpose**: Pull request dependency change analysis and vulnerability prevention.

**Triggers:**
- Pull request events

**Security Controls:**
- Moderate severity threshold for blocking
- Approved license validation (MIT, Apache-2.0, BSD variants, ISC, MPL-2.0)
- Dependency manifest change detection
- Vulnerability blocking for PR merges

#### 5. Auto Tag (`auto-tag.yml`)

**Purpose**: Automated semantic versioning and tag management.

**Triggers:**
- Push to main branch with relevant file changes
- Manual workflow dispatch with force bump options

**Versioning Logic:**
- **Major Bump**: BREAKING, major, or breaking change keywords
- **Minor Bump**: feat, feature, or minor keywords
- **Patch Bump**: Default increment
- **Skip Tag**: [skip-tag] in commit message

**Version Management:**
- `.release-version` file as authoritative source
- Git tag synchronization and comparison
- Semantic version calculation
- Tag existence validation
- Manual force bump capability (patch/minor/major)

**File Updates:**
- Automatic `.release-version` file maintenance
- Git configuration for automated commits
- Tag creation with proper metadata

#### 6. Release (`release.yml`)

**Purpose**: Comprehensive multi-platform binary building and GitHub release creation.

**Triggers:**
- Git tags matching `v*` pattern
- Manual workflow dispatch with tag input

**Multi-Stage Pipeline:**

**build_and_test Job**
- Pre-release validation
- Test execution and linting
- Build verification

**build_binaries Job**
- Cross-platform compilation matrix:
  - Linux: amd64, arm64
  - macOS: amd64, arm64
  - Windows: amd64, arm64
- Binary signing with Cosign
- Checksum generation (SHA256, SHA512)
- Verification instructions creation
- Artifact upload with 1-day retention

**create_release Job**
- Artifact aggregation
- Changelog generation from git history
- GitHub release creation with:
  - Installation instructions
  - Verification procedures
  - Cosign public key distribution
  - Multi-platform download links

**build_docker Job**
- Multi-platform Docker image building
- GitHub Container Registry publishing
- Image tagging strategy:
  - Version-specific tags
  - Latest tag maintenance
- OpenContainer Image specification labels

## Security Implementation

### Step Security Harden Runner

All workflows implement security hardening through:
- Sudo access restriction
- Network egress monitoring and auditing
- Process execution tracking
- Security event logging

### Action Version Pinning

Security best practices implementation:
- SHA-based action pinning for immutability
- Version comments for maintainability
- Dependency update automation through Dependabot
- Supply chain attack prevention

### Permission Model

Principle of least privilege enforcement:
- Workflow-level read-all default permissions
- Job-specific permission grants
- Explicit permission documentation
- Security-sensitive operation isolation

### Code Signing

Binary integrity assurance:
- Cosign keyless signing implementation
- Private key password protection
- Public key distribution
- Signature verification workflows
- Checksum validation procedures

## Dependency Management

### Dependabot Configuration

Automated dependency updates across ecosystems:
- **GitHub Actions**: Daily security updates
- **Docker**: Daily base image updates
- **Go Modules**: Daily dependency updates

### Vulnerability Scanning

Multi-tool vulnerability detection:
- **Trivy**: Filesystem and container scanning
- **Nancy**: Go-specific vulnerability database
- **GitHub Advisory Database**: Integrated security advisories
- **OpenSSF Scorecard**: Security practice assessment

## Build Configuration

### Version Information Injection

Build-time metadata embedding:
```go
github.com/Gosayram/go-mdfmt/internal/version.Version
github.com/Gosayram/go-mdfmt/internal/version.Commit
github.com/Gosayram/go-mdfmt/internal/version.Date
github.com/Gosayram/go-mdfmt/internal/version.BuiltBy
```

### Cross-Platform Support

Platform matrix compilation:
- Linux: AMD64, ARM64
- macOS: AMD64, ARM64 (Intel and Apple Silicon)
- Windows: AMD64, ARM64

### Optimization Flags

Production build optimization:
- CGO disabled for static linking
- Dead code elimination (`-s -w`)
- Symbol table stripping
- Size optimization

## Container Strategy

### Multi-Stage Docker Builds

Efficient container image creation:
- Minimal runtime image
- Security-focused base images
- Multi-platform architecture support
- Layer optimization for caching

### Registry Management

GitHub Container Registry integration:
- Automated image publishing
- Tag lifecycle management
- Access control through GitHub permissions
- Image scanning integration

## Release Process

### Automated Release Workflow

1. **Tag Creation**: Push `v*` tag triggers release pipeline
2. **Validation**: Complete test suite execution
3. **Compilation**: Multi-platform binary generation
4. **Signing**: Cosign signature application
5. **Packaging**: Artifact collection and organization
6. **Publication**: GitHub release with changelog
7. **Distribution**: Docker image publishing

### Manual Release Process

Workflow dispatch capability:
1. Navigate to GitHub Actions
2. Select Release workflow
3. Specify target tag (e.g., `v1.0.0`)
4. Execute workflow manually

### Verification Procedures

Release integrity validation:
```bash
# Checksum verification
sha256sum -c mdfmt-*-linux-amd64.sha256

# Cosign signature verification
cosign verify-blob \
  --key cosign.pub \
  --signature mdfmt-*-linux-amd64.sig \
  mdfmt-*-linux-amd64
```

## Development Integration

### Local CI Commands

Makefile integration for local development:
```bash
make test           # Execute test suite
make lint           # Run linting analysis
make build          # Compile binary
make test-coverage  # Generate coverage report
make test-race      # Race condition detection
```

### Pre-commit Validation

Local validation matching CI pipeline:
- Code formatting verification
- Import organization
- Linting compliance
- Test execution
- Security scanning preparation

## Monitoring and Observability

### Security Dashboard Integration

Centralized security monitoring:
- GitHub Security tab consolidation
- SARIF result aggregation
- Vulnerability trend tracking
- Compliance reporting

### Build Status Monitoring

Pipeline health monitoring:
- GitHub Checks API integration
- Pull request status reporting
- Failure notification systems
- Performance metrics tracking

### Coverage Reporting

Code coverage analytics:
- Codecov integration
- Trend analysis
- Branch coverage requirements
- Coverage quality gates

## Configuration Management

### Environment Variables

Build-time configuration:
- `GOOS`, `GOARCH`: Target platform specification
- `CGO_ENABLED`: C library linking control
- `GOGC`: Garbage collector tuning for CI
- Version injection variables

### Secrets Management

Secure credential handling:
- GitHub token automatic provisioning
- Cosign private key protection
- Codecov token optional configuration
- Registry authentication management

## Best Practices Implementation

### Branch Protection

Recommended repository settings:
- Required status checks for all workflows
- Up-to-date branch requirements
- Code owner review requirements
- Push restriction enforcement

### Security Compliance

Industry standard adherence:
- OpenSSF Scorecard recommendations
- NIST Cybersecurity Framework alignment
- Supply chain security best practices
- Vulnerability disclosure procedures

### Performance Optimization

CI pipeline efficiency:
- Parallel job execution
- Dependency caching strategies
- Artifact reuse patterns
- Resource allocation optimization

## Troubleshooting

### Common Issues

**Build Failures:**
- Go version compatibility verification
- Dependency resolution conflicts
- Platform-specific compilation errors
- Resource constraint limitations

**Security Scan Failures:**
- Vulnerability threshold adjustments
- License compliance violations
- Dependency update requirements
- Configuration validation errors

**Release Pipeline Issues:**
- Tag format validation
- Permission configuration
- Artifact generation failures
- Registry connectivity problems

### Debug Procedures

**Workflow Debugging:**
1. Enable debug logging in Actions
2. Review job step outputs
3. Validate environment variables
4. Check permission configurations

**Local Reproduction:**
1. Use identical Go version
2. Replicate environment variables
3. Execute Makefile targets
4. Validate dependency versions 