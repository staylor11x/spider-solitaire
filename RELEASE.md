# Release Guide

For architecture decisions and rationale, see [DescisionLog.md](DescisionLog.md).

---

## Quick Reference

### Create a Release

```bash
git checkout main
git pull
git checkout -b release/v1.3.0
git push origin release/v1.3.0
```

Workflows will:
1. Validate version format
2. Build binaries (Windows + Linux)
3. Generate SHA256 checksums
4. Create GitHub Release with all files

**Result:** Public release on [GitHub Releases](https://github.com/staylor11x/spider-solitaire/releases)

### Create a Hotfix

```bash
git checkout main
git pull
git checkout -b release/v1.3.1
git push origin release/v1.3.1
```

Same process as standard release. Use patch version (v1.3.**1**).

### Create a Pre-release

```bash
git checkout main
git pull
git checkout -b release/v1.3.0-beta
git push origin release/v1.3.0-beta
```

Automatically marked as pre-release on GitHub (won't appear in latest).

---

## How It Works

```
Push release/vX.Y.Z branch
         ↓
Release Branch Tagging workflow
  ├─ Validate version format
  ├─ Verify version > current tag
  └─ Create & push tag vX.Y.Z
         ↓
Build Release Binaries workflow (triggered by tag)
  ├─ Build Windows executable
  ├─ Build Linux executable
  ├─ Generate SHA256 checksums
  └─ Upload artifacts
         ↓
Create GitHub Release workflow
  ├─ Download all artifacts
  ├─ Generate release notes
  └─ Publish to GitHub Releases
         ↓
✅ Release live at github.com/staylor11x/spider-solitaire/releases
```

---

## Troubleshooting

### "Invalid branch name format"

**Problem:** Release branch name doesn't match `release/vX.Y.Z`

**Solutions:**
- ✅ `release/v1.3.0` - Correct
- ✅ `release/v1.3.0-beta` - Correct
- ❌ `release/1.3.0` - Missing 'v'
- ❌ `release/v1.3` - Incomplete version
- ❌ `release/v1.3.0.1` - Too many segments

**Fix:** Delete branch, recreate with correct name

```bash
git push origin --delete release/v1.3.0
git checkout -b release/v1.3.0
git push origin release/v1.3.0
```

### "Version must be greater than current latest tag"

**Problem:** Trying to release v1.2.0 when v1.3.0 already exists

**Solution:** Use a version greater than the current latest

```bash
# Check current tags
git tag -l --sort=-version:refname | head -5

# Create release with higher version
git checkout -b release/v1.4.0
git push origin release/v1.4.0
```

### "Tag already exists"

**Problem:** Trying to release v1.3.0 twice

**Solution:** Delete the tag locally and remotely, then retry

```bash
# Delete remote tag
git push origin --delete v1.3.0

# Delete local tag
git tag -d v1.3.0

# Retry release
git checkout -b release/v1.3.0
git push origin release/v1.3.0
```

### "Build workflow didn't trigger"

**Problem:** Release tag was created but build didn't start

**Solution:** 
1. Check the **Actions** tab for any failures in "Release Branch Tagging"
2. If that passed, check "Build Release Binaries" started (may take 1-2 minutes)
3. If neither appeared, check your branch was pushed correctly: `git push origin release/vX.Y.Z`

### "GitHub Release not created"

**Problem:** Build succeeded but no release appeared

**Possible causes:**
- Check Actions tab for "Create GitHub Release" job
- Verify all artifacts were uploaded (look for checksums.txt)
- Check GITHUB_TOKEN permissions (should be automatic)

**Manual fix:** Create release manually from GitHub UI using tag you created.

---

## Versioning Strategy

See [CONTRIBUTING.md - Versioning Strategy](CONTRIBUTING.md#versioning-strategy) for details.

**Quick guide:**
- **Breaking changes** → v1.0.0 → v2.0.0 (MAJOR)
- **New features** → v1.2.0 → v1.3.0 (MINOR)
- **Bug fixes** → v1.2.5 → v1.2.6 (PATCH)
- **Main branch** → auto-increments patch (v1.2.0 → v1.2.1 → v1.2.2)
- **Release branches** → manual minor/major (release/v1.3.0)

---

## Tips

✅ **Before releasing:**
- Ensure all tests pass on `main`
- Update user-facing docs if needed
- Test the game locally

✅ **After releasing:**
- Verify release appears on [GitHub Releases](https://github.com/staylor11x/spider-solitaire/releases)
- Download and test binary if possible
- Announce release if you have users

✅ **For hotfixes:**
- Create from `main` (same as standard release)
- Use patch version (v1.3.1 for v1.3.0 hotfix)
- No need to merge release branch back to main
