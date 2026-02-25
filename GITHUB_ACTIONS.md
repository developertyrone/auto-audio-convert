# GitHub Actions - Quick Reference

## ✅ **What Was Set Up**

### 1. **Automated Release Workflow** (`.github/workflows/release.yml`)
**Triggers on:** Version tags (`v1.0.0`, `v2.1.3`, etc.)

**What it does:**
- ✅ Builds binaries for 5 platforms (Linux x64/ARM, macOS Intel/ARM, Windows)
- ✅ Creates compressed archives (.tar.gz, .zip)
- ✅ Generates SHA256 checksums
- ✅ Creates GitHub Release with auto-generated notes
- ✅ Uploads all assets automatically

### 2. **CI Test Workflow** (`.github/workflows/ci.yml`)
**Triggers on:** Every push to `main` branch, all pull requests

**What it does:**
- ✅ Tests on Linux, macOS, Windows
- ✅ Runs `go vet` (code analysis)
- ✅ Runs `gofmt` check (formatting)
- ✅ Builds binary on each OS
- ✅ Tests `--version` flag

---

## 🚀 **How to Create a Release**

### **Method 1: Automated (Recommended)**

```bash
# 1. Make sure all changes are committed
git add .
git commit -m "feat: add new feature"
git push origin main

# 2. Create and push a version tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

**That's it!** GitHub Actions will:
1. Build all binaries (takes ~3-5 minutes)
2. Create release at: `https://github.com/developertyrone/auto-audio-convert/releases`
3. Upload binaries, archives, and checksums

### **Method 2: Manual Trigger**

1. Go to: https://github.com/developertyrone/auto-audio-convert/actions
2. Click "Build and Release" workflow
3. Click "Run workflow" button
4. Select branch
5. Click green "Run workflow"

*(This builds artifacts but doesn't create a release)*

---

## 📋 **Release Checklist**

When you're ready to release `v1.0.0`:

```bash
# ✅ Step 1: Ensure code is ready
git status  # Should be clean

# ✅ Step 2: Create tag
git tag -a v1.0.0 -m "Initial release"

# ✅ Step 3: Push tag
git push origin v1.0.0

# ✅ Step 4: Wait for GitHub Actions
# Go to: https://github.com/developertyrone/auto-audio-convert/actions
# Watch the workflow run (usually 3-5 minutes)

# ✅ Step 5: Verify release
# Go to: https://github.com/developertyrone/auto-audio-convert/releases
# Check that v1.0.0 appears with all assets
```

---

## 📦 **What Gets Released**

For each version tag, the release will include:

```
auto-audio-convert-linux-amd64.tar.gz       (~5MB)
auto-audio-convert-linux-arm64.tar.gz       (~5MB)
auto-audio-convert-darwin-amd64.tar.gz      (~5MB)
auto-audio-convert-darwin-arm64.tar.gz      (~5MB)
auto-audio-convert-windows-amd64.zip        (~5MB)
checksums.txt                                (SHA256 hashes)
```

Plus auto-generated release notes with:
- Features list
- Installation instructions
- Usage examples
- Checksums

---

## 🔍 **How to Check Workflow Status**

### **View Workflows:**
https://github.com/developertyrone/auto-audio-convert/actions

### **View Releases:**
https://github.com/developertyrone/auto-audio-convert/releases

### **Check Latest CI:**
Look for green ✅ or red ❌ badges in the Actions tab

---

## 🐛 **Troubleshooting**

### **Workflow Failed?**
1. Click on the failed workflow in Actions tab
2. Expand the failed step
3. Read the error message
4. Fix the issue and push again

### **Want to Delete a Bad Release?**
```bash
# Delete remote tag
git push --delete origin v1.0.0

# Delete local tag
git tag -d v1.0.0

# Go to GitHub releases page and delete the release manually
```

Then fix the issue and create a new tag (e.g., `v1.0.1`).

---

## 🎯 **Next Steps**

### **Create Your First Release:**

```bash
cd /home/openclaw/.openclaw/coding/auto-audio-convert/repo

# Tag and release
git tag -a v1.0.0 -m "Initial release - auto audio converter with FFmpeg auto-install"
git push origin v1.0.0
```

Then watch the magic happen at:
- **Actions:** https://github.com/developertyrone/auto-audio-convert/actions
- **Releases:** https://github.com/developertyrone/auto-audio-convert/releases

---

## 📚 **Documentation Files**

- **RELEASE.md** - Complete release process guide
- **.github/workflows/release.yml** - Release automation
- **.github/workflows/ci.yml** - Continuous integration tests

---

**Your repository is now fully automated. Every tag you push becomes a release! 🎉**
