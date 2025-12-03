# Commit Rules & Guidelines

> **⚡ Golden Rule: Commit frequently after every meaningful advancement!**

This document establishes commit practices and rules for the Envoy AI Unified Installer project.

---

## 🎯 Core Principle

**Commit early, commit often.** After each logical, testable advancement:

✅ **DO COMMIT** if you've:
- ✓ Fixed a bug
- ✓ Added a feature
- ✓ Modified a function
- ✓ Updated documentation
- ✓ Refactored code
- ✓ Added tests
- ✓ Fixed formatting issues
- ✓ Updated dependencies

---

## 📋 Commit Workflow

### Step 1: Make Changes

Write focused code that addresses a single concern.

### Step 2: Test Before Committing

Run pre-commit checks before staging:

```bash
make fmt
make vet
make lint
make test
```

Or let pre-commit hooks run automatically (see Section 3).

### Step 3: Verify Hooks Pass

Pre-commit hooks will validate:
- ✓ Go formatting (gofmt)
- ✓ Go linting (golangci-lint)
- ✓ Go vet checks
- ✓ Shell script validation
- ✓ YAML validation
- ✓ Markdown linting
- ✓ Secret detection
- ✓ Commit message format

### Step 4: Commit

```bash
git add <files>
git commit -m "type(scope): description"
```

### Step 5: Push

```bash
git push origin branch-name
```

---

## 🔧 Setting Up Pre-Commit Hooks

### Installation

```bash
pip install pre-commit
pre-commit install
pre-commit install --hook-type commit-msg
```

### Run Manually

```bash
pre-commit run --all-files
```

### Bypass Hooks (Emergency Only)

```bash
git commit --no-verify
```

> ⚠️ **Only use `--no-verify` in emergencies!**

---

## 📝 Commit Message Format

Follow **Conventional Commits** format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type

| Type | Usage |
|------|-------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation |
| `style` | Code formatting (no logic change) |
| `refactor` | Code restructuring |
| `perf` | Performance improvement |
| `test` | Test additions/updates |
| `chore` | Build, dependencies, tooling |
| `ci` | CI/CD configuration |

### Scope

| Scope | Component |
|-------|-----------|
| `cli` | CLI commands |
| `install` | Install command |
| `config` | Configuration management |
| `helm` | Helm operations |
| `upstream` | Upstream discovery |
| `doctor` | Doctor command |
| `version` | Version command |
| `ci` | GitHub Actions workflows |
| `docs` | Documentation files |
| `scripts` | Helper scripts |

### Examples

#### ✅ Good Commits

```
feat(install): add --values-extra flag support

Allow users to pass multiple additional values files
for Helm chart customization (rate limiting, pools, etc).

Closes #42
```

```
fix(doctor): improve namespace detection logic

Previously failed silently if namespace creation needed.
Now provides clearer feedback.
```

```
docs: add rate-limiting configuration guide
```

```
test(helm): add unit tests for chart validation
```

```
chore(deps): upgrade go-github from v55 to v56
```

```
style(code): fix formatting issues detected by golangci-lint
```

---

## 🚦 Commit Frequency Guidelines

### Commit After Each...

#### Feature Development
```bash
feat(install): implement step 1 - cleanup
git add -A && git commit -m "feat(install): implement step 1 - cleanup"

feat(install): implement step 2 - gateway install
git add -A && git commit -m "feat(install): implement step 2 - gateway install"

feat(install): implement step 3 - CRDs install
git add -A && git commit -m "feat(install): implement step 3 - CRDs install"
```

#### Bug Fix
```bash
fix(doctor): handle missing kubectl gracefully
git add -A && git commit -m "fix(doctor): handle missing kubectl gracefully"
```

#### Documentation Update
```bash
docs: add rate-limiting examples
git add -A && git commit -m "docs: add rate-limiting examples"
```

#### Test Addition
```bash
test(helm): add validation tests
git add -A && git commit -m "test(helm): add validation tests"
```

#### Refactoring
```bash
refactor(config): simplify Viper initialization
git add -A && git commit -m "refactor(config): simplify Viper initialization"
```

---

## ✅ Pre-Commit Checklist

Before running `git commit`, verify:

- [ ] Code builds: `make build`
- [ ] Tests pass: `make test`
- [ ] Formatting: `make fmt`
- [ ] Linting: `make lint`
- [ ] Vet checks: `make vet`
- [ ] No debug code left in
- [ ] Comments added for complex logic
- [ ] Commit message follows convention
- [ ] Related issues referenced (#123)
- [ ] No secrets or credentials in files

---

## 🔍 Pre-Commit Hooks Reference

### Hooks Applied

#### Go Quality
- **gofmt** — Automatic formatting
- **golangci-lint** — Comprehensive linting
- **go vet** — Built-in static analysis

#### Shell Scripts
- **shellcheck** — Shell syntax validation
- **shfmt** — Format shell scripts

#### Files
- **trailing-whitespace** — Remove trailing spaces
- **end-of-file-fixer** — Ensure newline at EOF
- **check-yaml** — Validate YAML syntax
- **check-json** — Validate JSON syntax
- **detect-private-key** — Detect secrets
- **check-large-files** — Prevent large files (>5MB)

#### Documentation
- **markdownlint** — Markdown validation
- **yamllint** — YAML linting

#### Commit Messages
- **commitlint** — Validate message format

### Hook Failures

If a hook fails:

1. **Read the error** — Understand what failed
2. **Fix the issue** — Address the problem (most auto-fix)
3. **Re-stage** — `git add` fixed files
4. **Re-commit** — `git commit` again

#### Common Fixes

**gofmt errors:**
```bash
make fmt
git add cli/
git commit
```

**shellcheck errors:**
```bash
bash -n scripts/merge-charts.sh
# Fix issues manually
git add scripts/
git commit
```

**Trailing whitespace:**
```bash
# Pre-commit auto-fixes this
git add -A
git commit
```

---

## 📊 Commit History Best Practices

### ✅ Good History

```
* a1b2c3d - chore: release v0.2.0
* d4e5f6g - docs: add troubleshooting guide
* h7i8j9k - test(doctor): add comprehensive checks
* k0l1m2n - fix(doctor): handle kubectl errors gracefully
* n3o4p5q - feat(doctor): implement health check command
* q6r7s8t - refactor(config): simplify Viper setup
```

Small, focused commits with clear messages.

### ❌ Bad History

```
* a1b2c3d - wip - many changes
* d4e5f6g - fixed stuff
* h7i8j9k - updates
* k0l1m2n - asdfghjkl
```

Large, unclear commits with poor messages.

---

## 🚀 Commit Before Pull Request

Before submitting a PR, ensure:

```bash
git log --oneline origin/main..HEAD
```

Shows:
- ✓ Clear commit messages
- ✓ Logical progression
- ✓ No `wip` commits
- ✓ No `fix previous commit` commits
- ✓ Related issues referenced

---

## 🔄 Interactive Rebase (Advanced)

Clean up commit history before PR:

```bash
git rebase -i origin/main
```

Actions:
- `pick` — Keep commit as-is
- `reword` — Change commit message
- `squash` — Combine with previous commit
- `drop` — Remove commit

Example:
```
pick a1b2c3d feat(doctor): implement health check
reword d4e5f6g fix: typo in help text
squash h7i8j9k test: add doctor tests
```

---

## ⚡ Quick Commit Aliases

Add to `~/.gitconfig` for faster commits:

```bash
git config --global alias.cm 'commit -m'
git config --global alias.add-commit '!git add -A && git commit -m'
git config --global alias.amend 'commit --amend --no-edit'
```

Usage:
```bash
git add-commit "feat(install): add new flag"
git amend  # Add forgotten changes
git cm "fix: typo"
```

---

## 📚 Examples by Scenario

### Scenario: Bug Fix in Helm Package

```bash
# 1. Create branch
git checkout -b fix/helm-validation

# 2. Fix the bug
vim cli/pkg/helm/helm.go

# 3. Add tests
vim cli/pkg/helm/helm_test.go

# 4. Test
make test

# 5. Commit bug fix
git add cli/pkg/helm/helm.go
git commit -m "fix(helm): validate release name format"

# 6. Commit tests
git add cli/pkg/helm/helm_test.go
git commit -m "test(helm): add release name validation tests"

# 7. Push
git push origin fix/helm-validation
```

### Scenario: Feature Development (4-step install)

```bash
# Step 1
git add -A && git commit -m "feat(install): implement cleanup step"

# Step 2
git add -A && git commit -m "feat(install): implement gateway installation"

# Step 3
git add -A && git commit -m "feat(install): implement CRD installation"

# Step 4
git add -A && git commit -m "feat(install): implement controller installation"

# Tests
git add -A && git commit -m "test(install): add e2e tests"

# Documentation
git add -A && git commit -m "docs: add installation guide"
```

### Scenario: Documentation Update

```bash
git checkout -b docs/add-troubleshooting

vim docs/github-actions-setup.md

git add docs/github-actions-setup.md
git commit -m "docs: add troubleshooting section"

git push origin docs/add-troubleshooting
```

---

## 🛡️ Protection Rules

### Protected Branches

The `main` branch is protected:
- ✓ Requires PR review (1+ approval)
- ✓ Requires all checks to pass
- ✓ Requires branches to be up-to-date
- ✓ Automatically deletes head branches after merge

### Merge Strategy

Always use **Squash and merge** for clean history:
```
✓ Squash and merge
```

---

## 📖 Reading Further

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Git Commit Best Practices](https://www.git-scm.com/book/en/v2/Git-Basics-Recording-Changes-to-the-Repository)
- [Pre-commit Framework](https://pre-commit.com/)
- [commitlint](https://commitlint.js.org/)

---

## 🎯 Summary

| What | When | Example |
|------|------|---------|
| **Commit** | After each logical unit | `feat(install): add gateway step` |
| **Push** | When feature complete | After PR-ready commits |
| **Test** | Before each commit | `make test` |
| **Format** | Before committing | `make fmt` |
| **Lint** | Before committing | `make lint` |

**Remember:** Commit frequently, keep commits small and focused, write clear messages. This makes:
- ✓ History readable
- ✓ Debugging easier
- ✓ Reviews faster
- ✓ Collaboration better

---

**Happy committing! 🚀**
