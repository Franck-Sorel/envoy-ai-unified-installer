# Contributing to Envoy AI Unified Installer

Thank you for considering contributing to this project! This guide will help you get started.

## Getting Started

### Prerequisites

- Go 1.21 or later
- kubectl 1.21+
- helm 3.10+
- A Kubernetes cluster (local or remote)

### Development Setup

```bash
git clone https://github.com/<YOUR_USERNAME>/envoy-ai-unified-installer.git
cd envoy-ai-unified-installer/cli

go mod download
go mod tidy
```

### Building

```bash
cd cli
go build -o ../envoy-ai-installer
cd ..
```

### Testing Your Changes

1. **Run health check:**
   ```bash
   ./envoy-ai-installer doctor
   ```

2. **Test with dry-run:**
   ```bash
   ./envoy-ai-installer install --dry-run
   ```

3. **Local cluster testing:**
   ```bash
   kind create cluster --name test
   ./envoy-ai-installer install
   kubectl get pods -A
   kind delete cluster --name test
   ```

## Code Style

### Go Conventions

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `make fmt` to format code
- Use `make vet` to check for issues
- Use `make lint` for comprehensive linting

### File Organization

```
cli/
├── cmd/
│   ├── root.go        # Main command setup
│   ├── install.go     # Install subcommand
│   ├── version.go     # Version subcommand
│   └── doctor.go      # Health check subcommand
├── pkg/
│   ├── config/        # Configuration management
│   ├── helm/          # Helm operations
│   └── upstream/      # Upstream chart discovery
└── main.go            # Entry point
```

### Naming Conventions

- **Packages:** lowercase, single word when possible
- **Functions:** CamelCase, exported functions capitalized
- **Variables:** camelCase for local, CONSTANT_CASE for constants
- **Files:** lowercase with underscores (e.g., `config.go`, `helm.go`)

## Making Changes

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
```

Branch naming conventions:
- `feature/` — New features
- `fix/` — Bug fixes
- `docs/` — Documentation updates
- `refactor/` — Code refactoring
- `test/` — Test improvements

### 2. Make Your Changes

- Keep commits atomic and focused
- Write clear commit messages
- Reference issues when applicable

### 3. Test Locally

```bash
make clean
make all
./envoy-ai-installer doctor
./envoy-ai-installer install --dry-run
```

### 4. Before Submitting

```bash
make fmt
make vet
make lint
make test
```

## Submitting Changes

### 1. Push to Your Fork

```bash
git push origin feature/your-feature-name
```

### 2. Create a Pull Request

- Use a descriptive title
- Reference related issues (#123)
- Include a detailed description of changes
- Add screenshots/logs if applicable

### 3. PR Template

```markdown
## Description
Brief description of changes

## Changes Made
- [ ] Feature 1
- [ ] Feature 2

## Related Issues
Closes #123

## Testing
- [x] Dry-run tested
- [x] Linting passed
- [x] Tests passed

## Screenshots/Logs
(if applicable)
```

## Code Review Process

- At least one maintainer review required
- CI/CD checks must pass
- Conflicts must be resolved
- Comments should be addressed before merge

## Reporting Bugs

### Bug Report Template

```markdown
## Environment
- OS: (macOS/Linux/Windows)
- Go Version: 1.21+
- Kubernetes Version: 1.28+
- Helm Version: 3.12+

## Describe the Bug
Clear description of what went wrong

## Steps to Reproduce
1. Run `./envoy-ai-installer install --dry-run`
2. Observe error...

## Expected Behavior
What should happen

## Actual Behavior
What actually happened

## Error Messages/Logs
```
paste full error log
```

## Screenshots
(if applicable)
```

## Feature Requests

### Feature Request Template

```markdown
## Description
Clear description of requested feature

## Motivation
Why is this feature needed?

## Proposed Solution
How should it work?

## Alternatives Considered
Other approaches?

## Implementation Notes
(optional) Technical considerations
```

## Documentation

### Writing Documentation

- Use clear, concise language
- Include examples where possible
- Link to related documentation
- Update README.md if adding features
- Add inline code comments for complex logic

### Documentation Checklist

- [ ] README updated
- [ ] Comments added to complex functions
- [ ] Examples provided for new features
- [ ] API changes documented

## Commit Messages

### Convention

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Example

```
feat(install): add rate-limit values support

Add support for passing rate-limit configuration files
via the --values-extra flag to customize rate limiting.

Closes #42
```

### Types

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation
- `style:` Formatting
- `refactor:` Code restructuring
- `perf:` Performance improvement
- `test:` Test additions/updates
- `chore:` Build/dependency updates

## Development Workflow

### Example: Adding a New Command

1. Create new file in `cli/cmd/`:
   ```bash
   touch cli/cmd/newcommand.go
   ```

2. Implement command using Cobra:
   ```go
   package cmd

   import "github.com/spf13/cobra"

   var newCmd = &cobra.Command{
       Use:   "newcommand",
       Short: "Brief description",
       RunE:  runNew,
   }

   func runNew(cmd *cobra.Command, args []string) error {
       // Implementation
       return nil
   }

   func init() {
       rootCmd.AddCommand(newCmd)
   }
   ```

3. Add to root.go:
   ```go
   func init() {
       rootCmd.AddCommand(newCmd)
   }
   ```

4. Test:
   ```bash
   make build
   ./envoy-ai-installer newcommand --help
   ```

### Example: Adding a Helper Package

1. Create package directory:
   ```bash
   mkdir -p cli/pkg/mypackage
   touch cli/pkg/mypackage/mypackage.go
   ```

2. Implement package with exports:
   ```go
   package mypackage

   func MyExportedFunction() {
       // Implementation
   }
   ```

3. Use in commands:
   ```go
   package cmd

   import "github.com/franck-sorel/envoy-ai-unified-installer/pkg/mypackage"

   func runCmd(cmd *cobra.Command, args []string) error {
       mypackage.MyExportedFunction()
       return nil
   }
   ```

## Release Process

### Preparing a Release

1. Update version in `cli/main.go`:
   ```go
   var version = "0.2.0"
   ```

2. Update CHANGELOG (if exists):
   ```markdown
   ## [0.2.0] - 2024-01-15

   ### Added
   - Feature X
   - Feature Y

   ### Fixed
   - Bug fix A
   ```

3. Create release PR with description of all changes

4. After merge, create GitHub release:
   ```bash
   git tag -a v0.2.0 -m "Release v0.2.0"
   git push origin v0.2.0
   ```

## Getting Help

- **Documentation:** See [README.md](README.md) and [docs/](docs/)
- **Issues:** Search existing [GitHub issues](https://github.com/Franck-Sorel/envoy-ai-unified-installer/issues)
- **Discussions:** Start a [GitHub discussion](https://github.com/Franck-Sorel/envoy-ai-unified-installer/discussions)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing! 🎉
