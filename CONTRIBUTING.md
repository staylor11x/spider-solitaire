# Contributing to Spider Solitaire

Thank you for your interest in contributing to Spider Solitaire! This document outlines our development workflow, branching strategy, and release process.

## Table of Contents

- [Development Workflow](#development-workflow)
- [Branching Strategy](#branching-strategy)
- [Versioning Strategy](#versioning-strategy)
- [Creating a Release](#creating-a-release)
- [Hotfix Process](#hotfix-process)
- [Testing](#testing)
- [Code Quality](#code-quality)

---

## Development Workflow

### Setting Up Your Development Environment

1. **Clone the repository:**
   ```bash
   git clone https://github.com/staylor11x/spider-solitaire.git
   cd spider-solitaire
   ```

2. **Install dependencies:**
   
   **Linux (Ubuntu):**
   ```bash
   sudo apt-get update
   sudo apt-get install -y libgl1-mesa-dev xorg-dev libasound2-dev
   go mod download
   ```
   
   **Windows:**
   ```bash
   go mod download
   ```

3. **Install development tools:**
   ```bash
   # golangci-lint for linting
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   
   # goimports for formatting
   go install golang.org/x/tools/cmd/goimports@latest
   ```

### Making Changes

1. Create a feature branch from `main`:
   ```bash
   git checkout main
   git pull
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and test locally:
   ```bash
   make test     # Run tests
   make lint     # Run linter
   make format   # Format code
   ```

3. Commit and push:
   ```bash
   git add .
   git commit -m "description of your changes"
   git push origin feature/your-feature-name
   ```

4. Open a Pull Request to `main` branch

---

## Branching Strategy

We use a **dual-versioning strategy** to support rapid development while maintaining stable releases.

### Branch Types

| Branch Pattern | Purpose | Auto-Versioning | Example |
|---------------|---------|-----------------|---------|
| `main` | Development branch | Patch versions (v1.2.X) | - |
| `release/v*` | Release branches | Minor/Major versions | `release/v1.3.0` |
| `feature/*` | Feature development | None | `feature/new-game-mode` |
| `fix/*` | Bug fixes | None | `fix/card-rendering` |

### Branch Lifecycle

```
feature/new-feature  →  PR  →  main  →  auto-tagged v1.2.1, v1.2.2, ...
                                  ↓
                         release/v1.3.0  →  tagged v1.3.0  →  builds & publishes
```

---

## Versioning Strategy

We follow [Semantic Versioning](https://semver.org/) (SemVer): `MAJOR.MINOR.PATCH`

### Main Branch (Development Versions)

- **Automatic patch versioning** on every merge to `main`
- Tags: `v1.2.0` → `v1.2.1` → `v1.2.2` → `v1.2.3`
- Purpose: Track development progress
- **Not published as releases**

### Release Branches (Production Versions)

- **Manual versioning** via branch name
- Tags: `v1.3.0`, `v2.0.0`, `v1.3.1` (hotfix)
- Purpose: Official releases
- **Published to GitHub Releases with binaries**

### Version Increment Guidelines

| Change Type | Version Bump | Example | When to Use |
|-------------|--------------|---------|-------------|
| Breaking changes | MAJOR | v1.x.x → v2.0.0 | API changes, major refactors |
| New features | MINOR | v1.2.x → v1.3.0 | New game modes, features |
| Bug fixes (hotfix) | PATCH | v1.3.0 → v1.3.1 | Critical bug in released version |

---

## Creating a Release

### Prerequisites

- All tests passing on `main` branch
- Code reviewed and merged
- [CHANGELOG.md](CHANGELOG.md) updated with release notes

### Step-by-Step Process

#### 1. Determine Version Number

Based on changes since last release:
- **Minor release** (new features): `v1.3.0`, `v1.4.0`
- **Major release** (breaking changes): `v2.0.0`, `v3.0.0`

#### 2. Create Release Branch

```bash
# Ensure you're on latest main
git checkout main
git pull

# Create release branch (replace with your version)
git checkout -b release/v1.3.0
```

#### 3. Push Release Branch

```bash
git push origin release/v1.3.0
```

#### 4. Automated Process

Once pushed, GitHub Actions will:
1. ✅ Validate version format
2. ✅ Ensure version > latest tag
3. ✅ Create Git tag `v1.3.0`
4. ✅ Push tag to repository
5. 🔄 (Future) Build and publish binaries

#### 5. Verify Release

1. Go to **Actions** tab in GitHub
2. Check "Release Branch Tagging" workflow succeeded
3. Verify tag created under **Releases** → **Tags**

### Supported Version Formats

| Format | Valid? | Type | Example |
|--------|--------|------|---------|
| `release/v1.3.0` | ✅ | Minor release | Standard |
| `release/v2.0.0` | ✅ | Major release | Breaking changes |
| `release/v1.3.1` | ✅ | Hotfix | Bug fix for v1.3.0 |
| `release/v1.3.0-beta` | ✅ | Pre-release | Testing version |
| `release/v1.3.0-rc1` | ✅ | Release candidate | Almost ready |
| `release/1.3.0` | ❌ | Missing 'v' | Invalid |
| `release/v1.3` | ❌ | Incomplete | Invalid |

---

## Hotfix Process

When a critical bug is found in a released version:

### Scenario: Bug in v1.3.0

#### Option 1: Quick Hotfix (Recommended for Critical Bugs)

```bash
# Create hotfix branch from main
git checkout main
git pull
git checkout -b release/v1.3.1

# Make your fix
# ... edit files ...

# Test thoroughly
make test
make lint

# Push hotfix branch
git push origin release/v1.3.1
```

This creates tag `v1.3.1` which will be built and released.

#### Option 2: Next Minor Version (Non-Critical Bugs)

For non-urgent fixes, include in next minor release:
```bash
git checkout main
# ... make fix ...
# ... wait for next release/v1.4.0 ...
```

### Hotfix Branch Lifecycle

```
v1.3.0 released  →  bug discovered
                      ↓
              fix merged to main
                      ↓
        create release/v1.3.1  →  tag v1.3.1  →  publish hotfix
```

---

## Testing

### Running Tests Locally

```bash
# Run all tests
make test

# Run tests with coverage
go test ./... -race -cover

# Run specific package tests
go test ./internal/game -v
```

### UI Tests (Headless)

On Linux, UI tests require a virtual display:
```bash
xvfb-run -a make test
```

### Continuous Integration

All PRs and commits to `main` automatically run:
- ✅ Full test suite with race detection
- ✅ Code linting (golangci-lint)
- ✅ Build verification (both CLI and game binaries)

Tests must pass before merging.

---

## Code Quality

### Formatting

Before committing, format your code:
```bash
make format
```

This runs:
- `gofmt -s -w .` - Standard Go formatting
- `goimports -w .` - Organize imports

### Linting

Fix linting issues:
```bash
make lint
```

Configuration: [`.golangci.yaml`](.golangci.yaml)

### Pre-commit Checklist

- [ ] Tests pass (`make test`)
- [ ] No linting errors (`make lint`)
- [ ] Code formatted (`make format`)
- [ ] Changes documented in [CHANGELOG.md](CHANGELOG.md)
- [ ] Commit message follows [conventional commits](https://www.conventionalcommits.org/)

---

## Commit Message Guidelines

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>: <description>

[optional body]

[optional footer]
```

### Types

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `chore:` - Maintenance tasks
- `refactor:` - Code restructuring
- `test:` - Test additions/changes
- `perf:` - Performance improvements

### Examples

```bash
feat: add undo move functionality
fix: correct card rendering issue in dark theme
docs: update installation instructions
chore: upgrade ebiten to v2.9.6
```

---

## Questions or Issues?

- **Bug reports:** [Open an issue](https://github.com/staylor11x/spider-solitaire/issues/new)
- **Feature requests:** [Open an issue](https://github.com/staylor11x/spider-solitaire/issues/new)
- **Questions:** [Discussions](https://github.com/staylor11x/spider-solitaire/discussions)

---

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.