# Workflow Logic Explanation

## Release System Workflow Sequence

### 1. Development and Commit
```bash
# Developer makes changes to Go code
git add .
git commit -m "fix: improve markdown parsing"
git push origin main
```

### 2. Auto-Tag Workflow (`.github/workflows/auto-tag.yml`)

**Trigger:** Push to `main` with changes to Go files

**Logic:**
1. **Gets latest tag:** `git describe --tags --abbrev=0`
2. **Analyzes commits** from last tag to HEAD
3. **Determines version type:**
   - `BREAKING` / `major` → major bump (1.0.0 → 2.0.0)
   - `feat` / `feature` / `minor` → minor bump (1.0.0 → 1.1.0)
   - Default → patch bump (1.0.0 → 1.0.1)
4. **Updates `.release-version`** file with new version
5. **Commits and pushes** with `[skip-tag]` to avoid cycle
6. **Creates new tag:** `git tag v1.0.1`
7. **Pushes tag:** `git push origin v1.0.1`

**Important:** Updates `.release-version` BEFORE creating tag!

### 3. Release Workflow (`.github/workflows/release.yml`)

**Trigger:** New tag (created by auto-tag or manually)

**Sequence:**

#### 3.1 Build and Test
- Checks out code
- Runs tests
- Runs linter

#### 3.2 Build Binaries (parallel for all platforms)
- Builds binary with version in name: `mdfmt-1.0.1-linux-amd64`
- **Signs with Cosign** using private key from secrets
- Generates **checksums** (SHA256, SHA512)
- Creates **verification instructions**
- Uploads as artifacts

#### 3.3 Create Release
- Generates **public key** from private Cosign key
- Downloads all artifacts
- Generates **changelog** from commits
- Creates **GitHub Release** with all files + public key

#### 3.4 Build Docker
- Builds Docker image
- Pushes to GitHub Container Registry

### 4. Result

**GitHub Release contains:**
```
mdfmt-1.0.1-linux-amd64          # Binary
mdfmt-1.0.1-linux-amd64.sig     # Cosign signature
cosign.pub                       # Cosign public key
mdfmt-1.0.1-linux-amd64.sha256  # SHA256 checksum
mdfmt-1.0.1-linux-amd64.sha512  # SHA512 checksum
mdfmt-1.0.1-linux-amd64.verify  # Verification instructions
```

**Docker images:**
```
ghcr.io/gosayram/go-mdfmt:v1.0.1
ghcr.io/gosayram/go-mdfmt:latest
```

## Key Features

### Private Key Cosign Signing
```bash
# Signing with private key
cosign sign-blob --yes \
  --key cosign.key \
  --output-signature binary.sig \
  binary

# Verification with public key
cosign verify-blob \
  --key cosign.pub \
  --signature binary.sig \
  binary
```

### Cycle Prevention
- Auto-tag commits `.release-version` update with `[skip-tag]` → doesn't create new tag
- Release workflow does NOT update `.release-version` (auto-tag does this)

### Security
- **id-token: write** for Cosign OIDC
- **contents: write** only for release workflow
- All actions pinned to SHA

## Example Scenarios

### Regular Development
```bash
git commit -m "fix: resolve issue with links"
# → Creates v1.0.1
# → Release in 5-10 minutes
```

### New Feature
```bash
git commit -m "feat: add new export format"
# → Creates v1.1.0
# → Release in 5-10 minutes
```

### Breaking Change
```bash
git commit -m "BREAKING: change CLI interface"
# → Creates v2.0.0
# → Release in 5-10 minutes
```

### Skip Auto-tagging
```bash
git commit -m "docs: update README [skip-tag]"
# → Tag NOT created
# → Release does NOT happen
```

### Manual Release
```bash
git tag v1.2.3
git push origin v1.2.3
# → Release triggered manually
```

## Monitoring

### Status Check
- [Auto Tag Actions](https://github.com/Gosayram/go-mdfmt/actions/workflows/auto-tag.yml)
- [Release Actions](https://github.com/Gosayram/go-mdfmt/actions/workflows/release.yml)
- [Releases](https://github.com/Gosayram/go-mdfmt/releases)

### Debugging
```bash
# Check latest tag
git describe --tags --abbrev=0

# Check commits since last tag
git log $(git describe --tags --abbrev=0)..HEAD --oneline

# Check .release-version
cat .release-version
``` 