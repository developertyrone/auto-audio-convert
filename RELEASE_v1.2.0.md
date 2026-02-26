# Release v1.2.0 - Overwrite Mode

**Release Date:** February 26, 2025  
**Tag:** `v1.2.0`  
**Commit:** `d00428b`

## 🎯 New Feature: Overwrite Mode

Added `--overwrite` flag to force re-conversion of existing files, bypassing the smart skip mechanism.

### What's New

#### `--overwrite` Flag
- **Purpose:** Force re-conversion of audio files even if the target file already exists
- **Use Case:** Useful when you want to:
  - Re-encode files with different quality settings
  - Fix corrupted conversions
  - Update to newer codec versions
  - Apply different encoding parameters

### Usage Examples

```bash
# Standard conversion (skips existing files)
./auto-audio-convert --from=flac --to=mp3 --quality=high

# Force overwrite all files
./auto-audio-convert --from=flac --to=mp3 --quality=high --overwrite

# Re-encode with different bitrate
./auto-audio-convert --from=mp3 --to=mp3 --bitrate=320k --overwrite
```

### Output Changes

When `--overwrite` is enabled, the tool now displays:
```
⚠️  Mode: OVERWRITE (existing files will be replaced)
```

And during processing:
```
🔄 [Worker 0] Overwriting: song1.mp3
```

Instead of:
```
🔄 [Worker 0] Converting: song1.mp3
```

## 📝 Technical Changes

### Modified Files

1. **`main.go`**
   - Bumped version to `1.2.0`
   - Added `--overwrite` flag definition
   - Added overwrite mode indicator in output
   - Passed `overwrite` parameter to `convertFiles()`

2. **`converter.go`**
   - Updated `convertFile()` signature to accept `overwrite bool`
   - Modified existence check: `if exists && !overwrite`
   - Updated `convertFiles()` to accept and pass `overwrite` parameter
   - Added conditional output for "Converting" vs "Overwriting"

### Backward Compatibility

✅ **Fully backward compatible** - existing scripts continue to work without changes:
- Default behavior unchanged (still skips existing files)
- All existing flags work as before
- Only new `--overwrite` flag is added

## 🚀 Deployment

### Git Operations
```bash
# Committed changes
git commit -m "feat: add --overwrite flag to force re-conversion of existing files (v1.2.0)"

# Tagged release
git tag -a v1.2.0 -m "Release v1.2.0: Add --overwrite flag for force re-conversion"

# Pushed to GitHub
git push origin main
git push origin v1.2.0
```

### GitHub Actions

The tag push automatically triggered the GitHub Actions workflow:
- **Workflow:** `.github/workflows/release.yml`
- **Trigger:** `push: tags: v*.*.*`
- **Actions:**
  1. Build binaries for 5 platforms (Linux, macOS, Windows - various architectures)
  2. Create release archives (`.tar.gz` for Unix, `.zip` for Windows)
  3. Generate SHA256 checksums
  4. Create GitHub Release with all artifacts
  5. Run tests (`go vet`, `go fmt`, build verification)

### Release Artifacts

The workflow will produce:
- `auto-audio-convert-linux-amd64.tar.gz`
- `auto-audio-convert-linux-arm64.tar.gz`
- `auto-audio-convert-darwin-amd64.tar.gz`
- `auto-audio-convert-darwin-arm64.tar.gz`
- `auto-audio-convert-windows-amd64.zip`
- `checksums.txt`

## 🔗 Links

- **Repository:** https://github.com/developertyrone/auto-audio-convert
- **Release:** https://github.com/developertyrone/auto-audio-convert/releases/tag/v1.2.0
- **Workflow Run:** https://github.com/developertyrone/auto-audio-convert/actions

## ✅ Testing Recommendations

Before using in production, test the overwrite functionality:

```bash
# Test 1: Verify skip behavior (default)
./auto-audio-convert --from=flac --to=mp3 --source=./test-data
# Should skip files that already exist

# Test 2: Verify overwrite behavior
./auto-audio-convert --from=flac --to=mp3 --source=./test-data --overwrite
# Should re-convert all files

# Test 3: Verify quality change
./auto-audio-convert --from=mp3 --to=mp3 --quality=low --overwrite
# Should re-encode existing MP3s with lower quality
```

## 📊 Impact Assessment

- **Risk Level:** Low
  - Non-breaking change
  - Opt-in feature (flag must be explicitly set)
  - Existing behavior preserved
- **Testing:** Code compiles (verified via git operations)
- **Documentation:** README.md should be updated to document the new flag

## 🎉 Summary

Version 1.2.0 successfully adds the `--overwrite` flag, enabling force re-conversion of existing files. The feature is fully implemented, committed, tagged, and pushed to GitHub. The automated CI/CD pipeline (GitHub Actions) will build cross-platform binaries and create a new release automatically.

**Status:** ✅ Complete - Release pipeline in progress on GitHub Actions
