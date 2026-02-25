# Release Guide

## How to Create a New Release

### 1. **Update Version** (Optional)
If you want a specific version number in the binary:

```bash
# Edit main.go and update the version constant
const version = "1.1.0"
```

### 2. **Commit Changes**
```bash
git add .
git commit -m "chore: bump version to v1.1.0"
git push origin main
```

### 3. **Create and Push Tag**
```bash
# Create annotated tag
git tag -a v1.1.0 -m "Release v1.1.0"

# Push tag to GitHub
git push origin v1.1.0
```

### 4. **Automated Build & Release**
GitHub Actions will automatically:
1. ✅ Build binaries for all platforms
2. ✅ Create compressed archives (.tar.gz, .zip)
3. ✅ Generate SHA256 checksums
4. ✅ Create GitHub Release
5. ✅ Upload all assets

**Release will be available at:**
```
https://github.com/developertyrone/auto-audio-convert/releases
```

---

## Manual Release (Alternative)

If you prefer to create releases manually:

### 1. **Build Locally**
```bash
./build.sh
```

### 2. **Create Archives**
```bash
cd dist
tar czf auto-audio-convert-linux-amd64.tar.gz auto-audio-convert-linux-amd64
tar czf auto-audio-convert-linux-arm64.tar.gz auto-audio-convert-linux-arm64
tar czf auto-audio-convert-darwin-amd64.tar.gz auto-audio-convert-darwin-amd64
tar czf auto-audio-convert-darwin-arm64.tar.gz auto-audio-convert-darwin-arm64
zip auto-audio-convert-windows-amd64.zip auto-audio-convert-windows-amd64.exe

# Generate checksums
sha256sum *.tar.gz *.zip > checksums.txt
```

### 3. **Create GitHub Release**
1. Go to https://github.com/developertyrone/auto-audio-convert/releases
2. Click "Draft a new release"
3. Choose tag: `v1.1.0` (create new)
4. Release title: `v1.1.0`
5. Upload files from `dist/`:
   - All `.tar.gz` files
   - All `.zip` files
   - `checksums.txt`
6. Click "Publish release"

---

## Versioning Convention

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR** (v2.0.0): Breaking changes
- **MINOR** (v1.1.0): New features, backward-compatible
- **PATCH** (v1.0.1): Bug fixes, backward-compatible

### Examples:
- `v1.0.0` - Initial release
- `v1.0.1` - Bug fix (ffmpeg detection)
- `v1.1.0` - New feature (progress bar)
- `v2.0.0` - Breaking change (new CLI flags)

---

## Automated Workflow Triggers

### Push Tag (Recommended)
```bash
git tag v1.1.0
git push origin v1.1.0
```
→ Builds + Creates Release

### Manual Dispatch
1. Go to Actions tab
2. Select "Build and Release"
3. Click "Run workflow"
4. Choose branch
5. Click "Run workflow"

→ Builds artifacts only (no release)

---

## Testing Before Release

Run CI tests locally:

```bash
# Format check
gofmt -l .

# Vet
go vet ./...

# Build all platforms
./build.sh

# Test binary
./dist/auto-audio-convert-linux-amd64 --version
```

---

## Release Checklist

- [ ] Update CHANGELOG.md (if exists)
- [ ] Update version in main.go (optional)
- [ ] Commit and push changes
- [ ] Create and push tag
- [ ] Wait for GitHub Actions to complete
- [ ] Verify release assets on GitHub
- [ ] Test download links
- [ ] Announce release (optional)

---

## Rollback a Release

If you need to delete a bad release:

```bash
# Delete remote tag
git push --delete origin v1.1.0

# Delete local tag
git tag -d v1.1.0

# Delete GitHub release manually via web UI
```

Then fix the issue and re-release with a new patch version (e.g., `v1.1.1`).

---

## First Release

For the initial `v1.0.0` release:

```bash
# Ensure all code is committed
git add .
git commit -m "feat: initial release v1.0.0"
git push origin main

# Create tag
git tag -a v1.0.0 -m "Initial release"
git push origin v1.0.0
```

GitHub Actions will handle the rest! 🚀
