# Release Guide

This document provides comprehensive information about the go-mdfmt release process, automated workflows, and security features.

## Release Overview

go-mdfmt follows semantic versioning and provides automated release workflows with comprehensive security features including binary signing and supply chain verification.

## Release Process

### Automated Release Workflow

The release process is fully automated through GitHub Actions and consists of several coordinated workflows:

1. **Auto-tagging**: Automatic version detection and tag creation
2. **Security Scanning**: Comprehensive security analysis
3. **Binary Building**: Cross-platform binary compilation
4. **Binary Signing**: Cosign-based cryptographic signing
5. **Release Publishing**: GitHub release creation with complete changelog

### Version Management

Version information is managed through the `.release-version` file in the repository root:

```
0.2.6
```

The build system automatically injects version information at compile time:

```go
// Build-time variables injected via linker flags
var (
    Version     = "dev"      // From .release-version
    Commit      = "unknown"  // Git commit hash
    Date        = "unknown"  // Build timestamp
    BuiltBy     = "unknown"  // Builder identity
    BuildNumber = "0"        // CI build number
)
```

## Release Artifacts

Each release includes comprehensive artifacts for verification and deployment:

### Binary Artifacts

**Supported Platforms**:
- Linux AMD64: `mdfmt-{version}-linux-amd64`
- Linux ARM64: `mdfmt-{version}-linux-arm64`
- macOS AMD64: `mdfmt-{version}-darwin-amd64`
- macOS ARM64: `mdfmt-{version}-darwin-arm64`
- Windows AMD64: `mdfmt-{version}-windows-amd64.exe`
- Windows ARM64: `mdfmt-{version}-windows-arm64.exe`

### Security Artifacts

**For each binary**:
- **Signature**: `.sig` file containing Cosign cryptographic signature
- **SHA256 Checksum**: `.sha256` file for integrity verification
- **SHA512 Checksum**: `.sha512` file for additional integrity verification
- **Verification Instructions**: `.verify` file with usage instructions

**Release-wide**:
- **Public Key**: `cosign.pub` for signature verification
- **Source Code**: Automated source code archives

## Security and Verification

### Binary Signing with Cosign

All binaries are cryptographically signed using Cosign with a private key for supply chain security.

**Signing Process**:
1. Cosign private key securely stored in GitHub repository secrets
2. Each binary signed during build process
3. Signatures include timestamping for temporal verification
4. Public key published with each release for verification

### Verification Instructions

**Step 1: Install Cosign**
```bash
# macOS (Homebrew)
brew install cosign

# Linux (direct download)
curl -L -o cosign https://github.com/sigstore/cosign/releases/latest/download/cosign-linux-amd64
chmod +x cosign
sudo mv cosign /usr/local/bin/
```

**Step 2: Download Release Files**
```bash
# Download binary (example for Linux AMD64)
curl -L -o mdfmt https://github.com/Gosayram/go-mdfmt/releases/download/v0.2.6/mdfmt-0.2.6-linux-amd64

# Download signature
curl -L -o mdfmt.sig https://github.com/Gosayram/go-mdfmt/releases/download/v0.2.6/mdfmt-0.2.6-linux-amd64.sig

# Download public key
curl -L -o cosign.pub https://github.com/Gosayram/go-mdfmt/releases/download/v0.2.6/cosign.pub

# Download checksums (optional)
curl -L -o mdfmt.sha256 https://github.com/Gosayram/go-mdfmt/releases/download/v0.2.6/mdfmt-0.2.6-linux-amd64.sha256
```

**Step 3: Verify Signature**
```bash
# Verify Cosign signature
cosign verify-blob \
  --key cosign.pub \
  --signature mdfmt.sig \
  mdfmt

# Verify checksum
sha256sum -c mdfmt.sha256
```

### Security Scanning

Each release undergoes comprehensive security analysis:

**Automated Security Checks**:
- **Dependency Review**: Vulnerability scanning of all dependencies
- **CodeQL Analysis**: Static security analysis of source code
- **OpenSSF Scorecard**: Supply chain security assessment
- **SLSA Compliance**: Software supply chain integrity verification

**Security Scorecard Results** available at: [OpenSSF Scorecard](https://securityscorecards.dev/viewer/?uri=github.com/Gosayram/go-mdfmt)

## Automated Workflows

### Auto-tag Workflow (`.github/workflows/auto-tag.yml`)

**Triggers**:
- Push to main branch
- Manual workflow dispatch with version override

**Functionality**:
- Detects version changes in `.release-version`
- Creates Git tags automatically
- Prevents duplicate tags
- Supports manual version bumping

**Configuration Options**:
```yaml
# Manual trigger with version override
workflow_dispatch:
  inputs:
    force_bump:
      description: 'Force version bump (patch/minor/major)'
      required: false
      type: choice
      options: ['', 'patch', 'minor', 'major']
```

### Release Workflow (`.github/workflows/release.yml`)

**Triggers**:
- Tag creation (from auto-tag workflow)
- Manual workflow dispatch with tag specification

**Build Matrix**:
```yaml
strategy:
  matrix:
    include:
      - goos: linux
        goarch: amd64
      - goos: linux
        goarch: arm64
      - goos: darwin
        goarch: amd64
      - goos: darwin
        goarch: arm64
      - goos: windows
        goarch: amd64
      - goos: windows
        goarch: arm64
```

**Security Features**:
- Binary signing with Cosign private key
- Checksum generation (SHA256 and SHA512)
- Verification instruction generation
- Public key publication

### Security Workflow (`.github/workflows/security.yml`)

**Scheduled Security Scans**:
- Weekly OpenSSF Scorecard assessment
- Dependency vulnerability monitoring
- Security advisory integration

**Manual Security Analysis**:
- On-demand security assessment
- Supply chain verification
- Compliance reporting

## Release Notes Generation

Release notes are automatically generated with comprehensive information:

**Changelog Content**:
- Commit history since last release
- Structured release notes from CHANGELOG.md (when available)
- Installation instructions for all platforms
- Security verification instructions
- Asset download links

**Enhanced Changelog Format**:
```markdown
## Changes
### Commits since v0.2.5:
- fix: correct Cosign verification instructions (296ba73)
- feat: add new formatting option (abc1234)

## Release Notes
### Fixed
- Corrected semantic versioning logic in auto-tag workflow
- Fixed incorrect version jumps by implementing proper increment rules

### Changed
- Implemented custom versioning scheme with 9-limit logic
- Enhanced auto-tag workflow with detailed logging

### Added
- Custom versioning logic with examples
- Automatic CHANGELOG.md extraction for release notes

## Installation
### Linux/macOS:
[Installation commands]

### Windows:
[Installation commands]

### Verification:
[Security verification instructions]
```

**CHANGELOG.md Integration**:
The release workflow now automatically extracts structured information from CHANGELOG.md:
- Searches for section matching current version (e.g., `## [0.2.6]`)
- Includes Fixed, Changed, Added, and other standard sections
- Falls back to commit-based changelog if no CHANGELOG.md entry exists
- Provides both technical commit history and user-friendly release notes

## Development Release Process

### Creating a Release

**Method 1: Automatic (Recommended)**
1. Update `.release-version` file with new version
2. Commit and push to main branch
3. Auto-tag workflow creates tag automatically
4. Release workflow builds and publishes release

**Method 2: Manual**
1. Update `.release-version` file
2. Use manual workflow dispatch to force tag creation
3. Specify version bump type if needed

### Version Bumping

**Custom Semantic Versioning with 9-Limit Logic**:

go-mdfmt implements a custom versioning scheme where each component has a limit of 9:
- **Patch** (0.0.X): Bug fixes, documentation updates (0.2.0 → 0.2.1 → ... → 0.2.9)
- **Minor** (0.X.0): New features, backward-compatible changes (0.2.9 → 0.3.0)
- **Major** (X.0.0): Breaking changes, API modifications (0.9.9 → 1.0.0)

**Version Increment Examples**:
```
# Patch increments
0.2.4 → 0.2.5 → 0.2.6 → ... → 0.2.9

# When patch reaches 9, bump minor and reset patch
0.2.9 → 0.3.0

# When minor reaches 9 and patch is at 9, bump major
0.9.9 → 1.0.0

# Manual minor bump (skips patch limit check)
0.2.4 → 0.3.0 (via workflow_dispatch with minor)

# Manual major bump
0.2.4 → 1.0.0 (via workflow_dispatch with major)
```

**Automated Version Logic**:
The auto-tag workflow implements this logic automatically:
- Checks current patch/minor values before incrementing
- Automatically cascades version bumps when limits are reached
- Logs version bump decisions for debugging
- Prevents incorrect version jumps (e.g., 0.2.4 → 0.3.0 for patch bump)

**Manual Version Bumping**:
```bash
# Update version file (note: no 'v' prefix in file)
echo "0.2.6" > .release-version

# Commit and push
git add .release-version
git commit -m "chore: bump version to 0.2.6"
git push origin main
```

### Pre-release Testing

Before creating releases, ensure comprehensive testing:

```bash
# Run all tests
make test

# Run security checks
make lint staticcheck

# Test cross-platform builds
make build-cross

# Verify configuration
mdfmt --version
```

## Release Artifacts Structure

Each release provides the following file structure:

```
Release v0.2.6/
├── mdfmt-0.2.6-linux-amd64           # Linux AMD64 binary
├── mdfmt-0.2.6-linux-amd64.sig       # Cosign signature
├── mdfmt-0.2.6-linux-amd64.sha256    # SHA256 checksum
├── mdfmt-0.2.6-linux-amd64.sha512    # SHA512 checksum
├── mdfmt-0.2.6-linux-amd64.verify    # Verification instructions
├── mdfmt-0.2.6-linux-arm64           # Linux ARM64 binary
├── mdfmt-0.2.6-linux-arm64.sig       # Cosign signature
├── mdfmt-0.2.6-linux-arm64.sha256    # SHA256 checksum
├── mdfmt-0.2.6-linux-arm64.sha512    # SHA512 checksum
├── mdfmt-0.2.6-linux-arm64.verify    # Verification instructions
├── mdfmt-0.2.6-darwin-amd64          # macOS AMD64 binary
├── mdfmt-0.2.6-darwin-amd64.sig      # Cosign signature
├── mdfmt-0.2.6-darwin-amd64.sha256   # SHA256 checksum
├── mdfmt-0.2.6-darwin-amd64.sha512   # SHA512 checksum
├── mdfmt-0.2.6-darwin-amd64.verify   # Verification instructions
├── mdfmt-0.2.6-darwin-arm64          # macOS ARM64 binary
├── mdfmt-0.2.6-darwin-arm64.sig      # Cosign signature
├── mdfmt-0.2.6-darwin-arm64.sha256   # SHA256 checksum
├── mdfmt-0.2.6-darwin-arm64.sha512   # SHA512 checksum
├── mdfmt-0.2.6-darwin-arm64.verify   # Verification instructions
├── mdfmt-0.2.6-windows-amd64.exe     # Windows AMD64 binary
├── mdfmt-0.2.6-windows-amd64.exe.sig # Cosign signature
├── mdfmt-0.2.6-windows-amd64.exe.sha256 # SHA256 checksum
├── mdfmt-0.2.6-windows-amd64.exe.sha512 # SHA512 checksum
├── mdfmt-0.2.6-windows-amd64.exe.verify # Verification instructions
├── mdfmt-0.2.6-windows-arm64.exe     # Windows ARM64 binary
├── mdfmt-0.2.6-windows-arm64.exe.sig # Cosign signature
├── mdfmt-0.2.6-windows-arm64.exe.sha256 # SHA256 checksum
├── mdfmt-0.2.6-windows-arm64.exe.sha512 # SHA512 checksum
├── mdfmt-0.2.6-windows-arm64.exe.verify # Verification instructions
├── cosign.pub                         # Public key for verification
├── Source code (zip)                  # Automated source archive
└── Source code (tar.gz)               # Automated source archive
```

## Docker Integration

The release workflow includes Docker image building and publishing:

**Docker Images**:
- Registry: `ghcr.io/gosayram/go-mdfmt`
- Tags: `latest`, `v0.2.6`, `0.2.6`
- Multi-architecture support (AMD64, ARM64)

**Docker Usage**:
```bash
# Pull and run latest version
docker pull ghcr.io/gosayram/go-mdfmt:latest
docker run --rm -v $(pwd):/workspace ghcr.io/gosayram/go-mdfmt:latest --write .
```

## Troubleshooting Releases

### Common Release Issues

**Auto-tag Workflow Fails**:
- Verify `.release-version` file format
- Check for existing tags with same version
- Ensure proper Git permissions

**Release Workflow Fails**:
- Verify Cosign secrets are configured
- Check cross-platform build compatibility
- Ensure artifact upload permissions

**Security Verification Fails**:
- Verify public key matches private key
- Check signature file integrity
- Ensure Cosign CLI version compatibility

### Security Secrets Configuration

**Required Repository Secrets**:
- `COSIGN_PRIVATE_KEY`: Private key for binary signing
- `COSIGN_PASSWORD`: Password for private key access
- `GITHUB_TOKEN`: Automatic GitHub Actions token

**Secret Generation**:
```bash
# Generate Cosign key pair
cosign generate-key-pair

# Add cosign.key content to COSIGN_PRIVATE_KEY secret
# Add key password to COSIGN_PASSWORD secret
# Commit cosign.pub to repository
```

## Release Quality Assurance

### Pre-release Checklist

- [ ] All tests passing (`make test`)
- [ ] Code quality checks passing (`make lint staticcheck`)
- [ ] Cross-platform builds successful (`make build-cross`)
- [ ] Configuration validation working
- [ ] Documentation updated
- [ ] Version number updated in `.release-version`

### Post-release Verification

- [ ] All platform binaries downloadable
- [ ] Cosign signature verification working
- [ ] Checksum verification successful
- [ ] Docker images published correctly
- [ ] Release notes accurate and complete
- [ ] Security scanning completed without issues

## Release Schedule

**Regular Releases**:
- **Patch releases**: As needed for bug fixes
- **Minor releases**: Monthly for new features
- **Major releases**: Quarterly for significant changes

**Security Releases**:
- **Critical vulnerabilities**: Immediate release
- **Security updates**: Within 48 hours
- **Dependency updates**: Weekly automated scanning 