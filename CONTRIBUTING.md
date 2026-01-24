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

**For detailed release procedures and troubleshooting, see [RELEASE.md](RELEASE.md).**

### Quick Summary

```bash
git checkout main
git pull
git checkout -b release/vX.Y.Z
git push origin release/vX.Y.Z
```

GitHub Actions will automatically:
1. ✅ Validate version format
2. ✅ Build Windows and Linux binaries
3. ✅ Generate SHA256 checksums
4. ✅ Create GitHub Release with files attached

See [RELEASE.md](RELEASE.md) for hotfixes, pre-releases, and troubleshooting.

---

## Hotfix Process

See [RELEASE.md - Create a Hotfix](RELEASE.md#create-a-hotfix) for detailed instructions.

**Quick summary:** Create `release/v1.3.1` branch to patch v1.3.0. Same process as standard release using patch version number.

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