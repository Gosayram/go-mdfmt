# Automated Release System

## Overview

This project uses a fully automated release system that:
- ✅ **Auto-creates tags** when Go code changes in main branch
- ✅ **Builds signed binaries** for all platforms
- ✅ **Signs releases** with Cosign for security
- ✅ **Provides individual binaries** (not archives)
- ✅ **Generates checksums** for verification

## How It Works

### 1. Automatic Tagging (`.github/workflows/auto-tag.yml`)

**Triggers:**
- Push to `main` branch
- Changes in Go files (`**/*.go`, `go.mod`, `go.sum`)
- Changes in source directories (`cmd/`, `pkg/`, `internal/`)

**Version Bumping Logic:**
- **Patch** (default): `v1.0.0` → `v1.0.1`
- **Minor**: Commit message contains `feat`, `feature`, or `minor`
- **Major**: Commit message contains `BREAKING`, `major`, or `breaking change`

**Example:**
```bash
# This will create v1.0.1 (patch)
git commit -m "fix: resolve markdown link wrapping issue"

# This will create v1.1.0 (minor)
git commit -m "feat: add new formatting option"

# This will create v2.0.0 (major)
git commit -m "BREAKING: change CLI interface"
```

### 2. Release Build (`.github/workflows/release.yml`)

**Triggered by:** New tags (created automatically or manually)

**Build Matrix:**
- **Linux**: amd64, arm64
- **macOS**: amd64, arm64
- **Windows**: amd64, arm64

**Security Features:**
- ✅ **Cosign signatures** for all binaries
- ✅ **SHA256/SHA512 checksums**
- ✅ **SLSA provenance** metadata
- ✅ **Reproducible builds**

## Release Artifacts

Each release includes for every platform:

### Binaries
- `mdfmt-1.0.0-linux-amd64`
- `mdfmt-1.0.0-linux-arm64`
- `mdfmt-1.0.0-darwin-amd64`
- `mdfmt-1.0.0-darwin-arm64`
- `mdfmt-1.0.0-windows-amd64.exe`
- `mdfmt-1.0.0-windows-arm64.exe`

### Security Files
- `mdfmt-1.0.0-linux-amd64.sig` (Cosign signature)
- `cosign.pub` (Cosign public key for verification)
- `mdfmt-1.0.0-linux-amd64.sha256` (SHA256 checksum)
- `mdfmt-1.0.0-linux-amd64.sha512` (SHA512 checksum)
- `mdfmt-1.0.0-linux-amd64.verify` (Verification instructions)

### Docker Images
- `ghcr.io/gosayram/go-mdfmt:v1.0.0`
- `ghcr.io/gosayram/go-mdfmt:latest`

## Installation & Verification

### Quick Install (Linux/macOS)
```bash
# Download latest release
curl -L -o mdfmt https://github.com/Gosayram/go-mdfmt/releases/latest/download/mdfmt-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)

# Make executable
chmod +x mdfmt

# Move to PATH
sudo mv mdfmt /usr/local/bin/
```

### Secure Install with Verification
```bash
# Download binary and verification files
VERSION="1.0.0"
PLATFORM="linux-amd64"
BASE_URL="https://github.com/Gosayram/go-mdfmt/releases/download/v${VERSION}"

curl -L -o mdfmt "${BASE_URL}/mdfmt-${VERSION}-${PLATFORM}"
curl -L -o mdfmt.sha256 "${BASE_URL}/mdfmt-${VERSION}-${PLATFORM}.sha256"
curl -L -o mdfmt.sig "${BASE_URL}/mdfmt-${VERSION}-${PLATFORM}.sig"
curl -L -o cosign.pub "${BASE_URL}/cosign.pub"

# Verify checksum
sha256sum -c mdfmt.sha256

# Verify Cosign signature (requires cosign CLI)
cosign verify-blob \
  --key cosign.pub \
  --signature mdfmt.sig \
  mdfmt

# Install if verification passes
chmod +x mdfmt
sudo mv mdfmt /usr/local/bin/
```

### Install Cosign for Verification
```bash
# Linux
curl -O -L "https://github.com/sigstore/cosign/releases/latest/download/cosign-linux-amd64"
sudo mv cosign-linux-amd64 /usr/local/bin/cosign
chmod +x /usr/local/bin/cosign

# macOS
brew install cosign

# Windows
winget install sigstore.cosign
```

## Development Workflow

### Normal Development
```bash
# Make changes
git add .
git commit -m "fix: improve performance"
git push origin main

# Tag is created automatically (v1.0.1)
# Release is built automatically
# Binaries are available in ~5-10 minutes
```

### Feature Development
```bash
# For minor version bump
git commit -m "feat: add new export format"

# For major version bump
git commit -m "BREAKING: change configuration format"
```

### Manual Release
If you need to create a release manually:
```bash
git tag v1.2.3
git push origin v1.2.3
```

## Security Guarantees

### Cosign Signatures
- All binaries are signed with **private key Cosign**
- Public key (`cosign.pub`) is included with each release
- Private key is password-protected and stored in GitHub secrets
- Signatures are verifiable using the public key
- Provides **non-repudiation** and **integrity**

### SLSA Provenance
- **Level 3 SLSA** compliance
- Tracks build environment and materials
- Provides **supply chain security**
- Verifiable build metadata

### Checksums
- **SHA256** and **SHA512** for all binaries
- Detects **tampering** or **corruption**
- Enables **integrity verification**

## Troubleshooting

### Release Not Created
1. Check if Go files actually changed
2. Verify commit is on `main` branch
3. Check GitHub Actions logs
4. Ensure no existing tag with same version

### Verification Fails
```bash
# Check if you have the right files
ls -la mdfmt*

# Verify checksum format
cat mdfmt.sha256

# Check Cosign installation
cosign version

# Verify with verbose output
cosign verify-blob --key cosign.pub --signature mdfmt.sig mdfmt --verbose
```

### Manual Tag Creation
```bash
# If auto-tagging fails, create manually
git tag v1.0.1
git push origin v1.0.1

# Update .release-version file
echo "v1.0.1" > .release-version
git add .release-version
git commit -m "Update .release-version to v1.0.1"
git push origin main
```

## Configuration

### Customize Version Bumping
Edit `.github/workflows/auto-tag.yml` to change bump logic:

```yaml
# Check for different keywords
if echo "$COMMITS" | grep -i -E "(BREAKING|major|breaking change)" > /dev/null; then
  BUMP_TYPE="major"
elif echo "$COMMITS" | grep -i -E "(feat|feature|minor)" > /dev/null; then
  BUMP_TYPE="minor"
fi
```

### Skip Auto-Tagging
Add `[skip-tag]` to commit message:
```bash
git commit -m "docs: update README [skip-tag]"
```

## Monitoring

### GitHub Actions Status
- [Auto Tag Workflow](https://github.com/Gosayram/go-mdfmt/actions/workflows/auto-tag.yml)
- [Release Workflow](https://github.com/Gosayram/go-mdfmt/actions/workflows/release.yml)

### Release History
- [All Releases](https://github.com/Gosayram/go-mdfmt/releases)
- [Latest Release](https://github.com/Gosayram/go-mdfmt/releases/latest)

### Security Monitoring
- [Security Tab](https://github.com/Gosayram/go-mdfmt/security)
- [Dependency Graph](https://github.com/Gosayram/go-mdfmt/network/dependencies) 