# Pre-Commit Hooks Setup Guide

Automated checks to ensure code quality, security, and consistency before commits.

## 🎯 Overview

This project uses **pre-commit** hooks to automatically validate:

✅ **Go Code**
- Formatting (gofmt)
- Linting (golangci-lint)
- Static analysis (go vet)

✅ **Shell Scripts**
- Syntax validation (shellcheck)
- Formatting (shfmt)

✅ **YAML/JSON**
- YAML validation
- JSON validation

✅ **Documentation**
- Markdown linting
- Line length validation

✅ **Security**
- Secret detection
- Private key detection

✅ **Files**
- Large file detection
- Trailing whitespace removal
- CRLF line endings check

✅ **Commit Messages**
- Conventional Commits format validation

---

## 📦 Installation

### Prerequisites

- Python 3.7+
- pip or conda
- Git 2.9+

### Step 1: Install pre-commit

```bash
pip install pre-commit
```

Or with conda:

```bash
conda install -c conda-forge pre-commit
```

Or on macOS:

```bash
brew install pre-commit
```

### Step 2: Install Git Hooks

In the repository root:

```bash
pre-commit install
pre-commit install --hook-type commit-msg
```

This creates files in `.git/hooks/`:
- `pre-commit` — Runs before staging commits
- `commit-msg` — Validates commit messages

### Step 3: Verify Installation

```bash
pre-commit --version
```

Expected output:
```
pre-commit 3.5.0
```

---

## 🚀 Quick Start

### First Run (Install Hook Dependencies)

The first time you commit, pre-commit will download and install all tools:

```bash
git add cli/cmd/install.go
git commit -m "feat(install): add new feature"
```

Output:
```
Trim trailing whitespace.................................................Passed
Fix End of File Fixer.................................................Passed
Check Yaml...........................................................Passed
golangci-lint.........................................................Failed
go fmt...............................................................Passed
...
```

### Hook Failures

If a hook fails, pre-commit will:

1. **Show errors** — Display what failed
2. **Auto-fix** — Many hooks auto-fix issues
3. **Require action** — Re-stage and commit fixed files

**Example: gofmt auto-fix**

```bash
$ git commit -m "feat(install): add feature"

go fmt...................................................Failed

# Pre-commit automatically fixed formatting
# You need to re-stage and commit

$ git add cli/cmd/install.go
$ git commit -m "feat(install): add feature"

go fmt...................................................Passed
```

---

## 🛠️ Common Tasks

### Run All Hooks Manually

```bash
pre-commit run --all-files
```

Output:
```
Trim trailing whitespace.................................................Passed
Fix End of File Fixer.................................................Passed
Check Yaml...........................................................Passed
Check Json...........................................................Passed
golangci-lint.........................................................Passed
go fmt...............................................................Passed
go vet...............................................................Passed
ShellCheck......................................................Passed
shfmt.................................................................Passed
markdownlint..........................................................Passed
yamllint...............................................................Passed
Detect secrets........................................................Passed
commitlint............................................................Passed
```

### Run Specific Hook

```bash
pre-commit run golangci-lint --all-files
pre-commit run go-fmt --all-files
pre-commit run shellcheck --all-files
```

### Skip Hooks (Emergency Only)

```bash
git commit --no-verify
```

> ⚠️ **Only use `--no-verify` in emergencies!**

### Update Hooks

```bash
pre-commit autoupdate
```

This updates all hooks to their latest versions.

---

## 🔍 Hook Descriptions

### Go Formatting & Linting

**golangci-lint** — Comprehensive Go linting
- Checks code quality, security, performance
- Auto-fixable issues: `golangci-lint run ./... --fix`

**go fmt** — Format Go code
- Enforces standard Go formatting
- Auto-fixes formatting issues

**go vet** — Static analysis
- Detects suspicious constructs
- Cannot auto-fix; requires manual review

### Shell Scripts

**shellcheck** — Shell script linting
- Detects syntax errors and bad practices
- Requires manual fixes

**shfmt** — Shell script formatting
- Auto-fixes formatting
- Uses 4-space indentation

### File Checks

**Trim trailing whitespace** — Remove spaces at line ends
- Auto-fixes

**Fix End of File Fixer** — Ensure newline at EOF
- Auto-fixes

**Check Yaml** — Validate YAML syntax
- Auto-fixes invalid structure
- Requires manual review of warnings

**Check Json** — Validate JSON syntax
- Requires manual fixes

### Documentation

**markdownlint** — Lint markdown files
- Line length: max 120 characters (configurable)
- Heading format, list format, etc.
- Requires manual fixes

**yamllint** — Lint YAML files
- Proper indentation, spacing, quotes
- Auto-fixable: `yamllint -d relaxed .`

### Security

**Detect secrets** — Find potential secrets
- Scans for AWS keys, GitHub tokens, etc.
- **False positives possible** — Review carefully
- Baseline file: `.secrets.baseline`

**Forbid CRLF** — Prevent Windows line endings
- Auto-fixes for non-Makefile files

**Forbid tabs** — Prevent tabs (except Makefile)
- Requires manual fixes

### Commit Messages

**commitlint** — Validate commit message format
- Requires: `type(scope): subject`
- Examples:
  - ✅ `feat(install): add gateway step`
  - ✅ `fix(doctor): handle missing kubectl`
  - ❌ `updated code` (no type/scope)
  - ❌ `FEAT(INSTALL): ADD STEP` (wrong case)

---

## 🔧 Configuration

### Main Config: `.pre-commit-config.yaml`

Defines which hooks to run:

```yaml
repos:
  - repo: https://github.com/golangci/golangci-lint
    rev: v1.55.2
    hooks:
      - id: golangci-lint
```

### Commitlint Config: `.commitlintrc.json`

Validates commit message format:

```json
{
  "extends": ["@commitlint/config-conventional"],
  "rules": {
    "type-enum": ["feat", "fix", "docs", ...]
  }
}
```

### Secrets Baseline: `.secrets.baseline`

Allows false positives:

```json
{
  "version": "1.4.0",
  "results": {}
}
```

---

## 📋 Troubleshooting

### Hook Not Running

**Problem:** Pre-commit hook doesn't run on commit

**Solution:**
```bash
pre-commit install
pre-commit install --hook-type commit-msg
```

Verify hooks installed:
```bash
ls -la .git/hooks/
```

Should show:
- `pre-commit`
- `commit-msg`

### "Hook X Not Found"

**Problem:** `golangci-lint: command not found`

**Solution:**
Pre-commit manages its own environment. Rebuild:

```bash
pre-commit clean
pre-commit run --all-files
```

### "Python Version Mismatch"

**Problem:** `python: version conflict`

**Solution:**
```bash
pre-commit autoupdate
```

### False Positive in Secret Detection

**Problem:** Code flagged as secret but it's not

**Solution:**
```bash
# Review the finding
pre-commit run detect-secrets --all-files

# If it's a false positive, add to baseline
detect-secrets scan --baseline .secrets.baseline
git add .secrets.baseline
git commit -m "chore: update secrets baseline"
```

### markdownlint Too Strict

**Problem:** Markdown formatting errors you don't want

**Solution:**
Add to file (inline ignore):

```markdown
<!-- markdownlint-disable MD013 -->
Very long line that exceeds 120 characters...
<!-- markdownlint-enable MD013 -->
```

### commitlint Rejects Message

**Problem:** Commit message fails validation

**Solution:**
Use correct format:

```bash
# ❌ Wrong
git commit -m "updated code"

# ✅ Correct
git commit -m "feat(install): add new feature"
git commit -m "fix(doctor): handle error gracefully"
git commit -m "docs: add troubleshooting guide"
```

---

## 🎯 Best Practices

### 1. **Let Hooks Auto-Fix**

Many issues are auto-fixed. Just re-stage and commit:

```bash
git add <file>
git commit -m "message"  # Hook fails, auto-fixes

git add <file>           # Re-stage fixed file
git commit -m "message"  # Now passes
```

### 2. **Fix Lint Errors Locally**

Before staging:

```bash
make lint --fix
make fmt
git add -A
git commit -m "fix: lint issues"
```

### 3. **Use Make Targets**

```bash
make fmt    # Format code
make vet    # Run vet
make lint   # Run linter
make test   # Run tests
```

### 4. **Bypass Only in Emergencies**

```bash
git commit --no-verify  # Only when absolutely necessary!
```

### 5. **Keep Commits Small**

Large commits take longer to lint. Keep changes focused:

```bash
# Good: One feature per commit
git commit -m "feat(install): add validation"
git commit -m "test(install): add validation tests"
git commit -m "docs: document validation"

# Bad: Everything at once
git commit -m "updated various things"
```

---

## 📚 Additional Resources

- [pre-commit Framework](https://pre-commit.com/)
- [golangci-lint](https://golangci-lint.run/)
- [commitlint](https://commitlint.js.org/)
- [shellcheck](https://www.shellcheck.net/)
- [Conventional Commits](https://www.conventionalcommits.org/)

---

## ✅ Verification

After setup, verify everything works:

```bash
# 1. Hooks installed
pre-commit --version
ls -la .git/hooks/

# 2. Test a commit (should pass)
touch test.txt
git add test.txt
git commit -m "test: verify hooks work"

# 3. Cleanup
git reset HEAD~1
rm test.txt
```

Expected output:
```
All checks passed! ✅
```

---

## 🚀 Next Steps

1. ✅ Install pre-commit hooks
2. ✅ Run `pre-commit run --all-files` to validate existing code
3. ✅ Make a test commit
4. ✅ Read [COMMIT_RULES.md](../COMMIT_RULES.md) for commit practices
5. ✅ Start developing!

---

**Happy committing with automated quality checks! 🎉**
