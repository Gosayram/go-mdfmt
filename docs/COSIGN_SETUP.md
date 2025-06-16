# Cosign Setup for Release Signing

## Overview

This project uses Cosign to sign all release binaries with a private key. To make the system work, you need to configure GitHub secrets.

## Required Secrets

### 1. COSIGN_PRIVATE_KEY
Cosign private key in PEM format.

### 2. COSIGN_PASSWORD  
Password to protect the private key.

## Key Generation

### Install Cosign

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

### Create Key Pair

```bash
# Generate private and public keys
cosign generate-key-pair

# You will be asked to enter a password to protect the private key
# Remember this password - you'll need it for the COSIGN_PASSWORD secret

# Files will be created:
# - cosign.key (private key) 
# - cosign.pub (public key)
```

## GitHub Secrets Configuration

### 1. Go to repository settings
```
https://github.com/YOUR_USERNAME/YOUR_REPO/settings/secrets/actions
```

### 2. Add COSIGN_PRIVATE_KEY secret
- Click "New repository secret"
- Name: `COSIGN_PRIVATE_KEY`
- Value: Contents of `cosign.key` file (including `-----BEGIN ENCRYPTED COSIGN PRIVATE KEY-----` and `-----END ENCRYPTED COSIGN PRIVATE KEY-----`)

### 3. Add COSIGN_PASSWORD secret
- Click "New repository secret"  
- Name: `COSIGN_PASSWORD`
- Value: Password you used when generating the keys

## Setup Verification

### Local Verification
```bash
# Create test file
echo "test content" > test.txt

# Sign file
cosign sign-blob --key cosign.key --output-signature test.txt.sig test.txt

# Verify signature
cosign verify-blob --key cosign.pub --signature test.txt.sig test.txt

# If command completed successfully, keys work correctly
```

### GitHub Actions Verification
After configuring secrets:

1. Make a commit with changes to Go code
2. Wait for tag creation
3. Check that release workflow completed successfully
4. Release should contain `.sig` files and `cosign.pub`

## Security

### Recommendations
- **DO NOT** commit private key to repository
- **DO NOT** share private key
- Use strong password to protect key
- Regularly rotate keys (recommended once per year)

### Key Rotation
1. Generate new key pair
2. Update secrets in GitHub
3. Create new release to test
4. Securely delete old keys

## Release Verification

### For Users
```bash
# Download binary, signature and public key
curl -L -o mdfmt https://github.com/YOUR_USERNAME/YOUR_REPO/releases/latest/download/mdfmt-VERSION-PLATFORM
curl -L -o mdfmt.sig https://github.com/YOUR_USERNAME/YOUR_REPO/releases/latest/download/mdfmt-VERSION-PLATFORM.sig  
curl -L -o cosign.pub https://github.com/YOUR_USERNAME/YOUR_REPO/releases/latest/download/cosign.pub

# Verify signature
cosign verify-blob --key cosign.pub --signature mdfmt.sig mdfmt

# If verification passed, binary is authentic
```

## Troubleshooting

### "invalid password" Error
- Check that password in `COSIGN_PASSWORD` secret matches key password
- Make sure there are no extra spaces or characters in password

### "invalid key format" Error  
- Check that private key is copied completely
- Make sure `-----BEGIN` and `-----END` lines are included
- Check for extra characters or line breaks

### Workflow Fails at Signing Step
- Check GitHub Actions logs
- Make sure both secrets are configured correctly
- Check workflow permissions (should have `id-token: write`)

## Additional Information

- [Official Cosign Documentation](https://docs.sigstore.dev/cosign/overview/)
- [GitHub Secrets Documentation](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
- [Sigstore Security Model](https://docs.sigstore.dev/security/) 