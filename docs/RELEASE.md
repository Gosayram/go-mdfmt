# Release Guide

## Quick Release Process

### 1. Automatic Release (Recommended)
```bash
# Create and push a new tag
git tag v1.0.0
git push origin v1.0.0
```

This will automatically:
- ✅ Run all tests and security checks
- ✅ Build cross-platform binaries (Linux, macOS, Windows)
- ✅ Create GitHub release with changelog
- ✅ Build and push Docker images to GHCR

### 2. Manual Release
1. Go to [GitHub Actions](https://github.com/Gosayram/go-mdfmt/actions)
2. Select "Release" workflow
3. Click "Run workflow"
4. Enter tag (e.g., `v1.0.0`)
5. Click "Run workflow"

## Version Management

### Bump Version Locally
```bash
# Patch version (1.0.0 → 1.0.1)
make bump-patch

# Minor version (1.0.0 → 1.1.0)
make bump-minor

# Major version (1.0.0 → 2.0.0)
make bump-major
```

### Check Current Version
```bash
make version
```

## Pre-Release Checklist

### Local Testing
```bash
# Run complete CI pipeline locally
make ci-release

# Test Docker build
make docker-build
make docker-run

# Test cross-platform builds
make build-cross
```

### Code Quality
```bash
# Format and lint code
make fmt
make lint

# Run security checks
make staticcheck

# Run all tests
make test
make test-race
```

## Release Artifacts

Each release creates:

### Binaries
- `mdfmt-<version>-linux-amd64.tar.gz`
- `mdfmt-<version>-linux-arm64.tar.gz`
- `mdfmt-<version>-darwin-amd64.tar.gz`
- `mdfmt-<version>-darwin-arm64.tar.gz`
- `mdfmt-<version>-windows-amd64.zip`
- `mdfmt-<version>-windows-arm64.zip`

### Docker Images
- `ghcr.io/gosayram/go-mdfmt:latest`
- `ghcr.io/gosayram/go-mdfmt:v<version>`

## Installation Instructions

### From GitHub Releases
```bash
# Linux/macOS
curl -L https://github.com/Gosayram/go-mdfmt/releases/latest/download/mdfmt-*-linux-amd64.tar.gz | tar xz
chmod +x mdfmt
sudo mv mdfmt /usr/local/bin/
```

### Using Docker
```bash
# Pull and run
docker pull ghcr.io/gosayram/go-mdfmt:latest
docker run --rm -v $(pwd):/workspace ghcr.io/gosayram/go-mdfmt:latest /workspace/README.md
```

### From Source
```bash
git clone https://github.com/Gosayram/go-mdfmt.git
cd go-mdfmt
make install
```

## Troubleshooting

### Release Failed
1. Check [GitHub Actions](https://github.com/Gosayram/go-mdfmt/actions) for error details
2. Common issues:
   - Tests failing → Fix tests and re-tag
   - Linting errors → Run `make lint-fix`
   - Security issues → Check security tab

### Docker Build Failed
```bash
# Test locally
make docker-build

# Check Dockerfile syntax
docker build --no-cache .
```

### Version Issues
```bash
# Reset version file
echo "v0.1.0" > .release-version

# Check version info
./bin/mdfmt --version
```

## CI/CD Status

Check the status of all workflows:
- [CI - Lint and Test](https://github.com/Gosayram/go-mdfmt/actions/workflows/ci-lint-test.yml)
- [Security Scanning](https://github.com/Gosayram/go-mdfmt/actions/workflows/security.yml)
- [CodeQL Analysis](https://github.com/Gosayram/go-mdfmt/actions/workflows/codeql.yml)
- [Release](https://github.com/Gosayram/go-mdfmt/actions/workflows/release.yml)

## Security

All releases are:
- ✅ Scanned for vulnerabilities (Trivy, Nancy)
- ✅ Analyzed for security issues (CodeQL)
- ✅ Built with hardened runners
- ✅ Signed with GitHub's signing key
- ✅ Reproducible builds with pinned dependencies 