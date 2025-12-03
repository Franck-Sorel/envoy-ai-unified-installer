# Getting Started

Quick setup guide to start developing or using the Envoy AI Unified Installer.

## 🚀 For Users (Install Envoy AI Gateway)

### Prerequisites

- `kubectl` (1.21+)
- `helm` (3.10+)
- Kubernetes cluster access

### Installation

```bash
git clone https://github.com/<YOUR_USERNAME>/envoy-ai-unified-installer.git
cd envoy-ai-unified-installer/cli

go build -o ../envoy-ai-installer
cd ..
```

### Run

```bash
# Check system readiness
./envoy-ai-installer doctor

# Install Envoy AI Gateway
./envoy-ai-installer install

# With optional Redis
./envoy-ai-installer install --with-redis

# Dry-run to preview
./envoy-ai-installer install --dry-run
```

For full documentation, see [README.md](README.md).

---

## 🔧 For Developers

### Initial Setup

```bash
# 1. Clone repository
git clone https://github.com/<YOUR_USERNAME>/envoy-ai-unified-installer.git
cd envoy-ai-unified-installer

# 2. Install pre-commit hooks
pip install pre-commit
pre-commit install
pre-commit install --hook-type commit-msg

# 3. Verify installation
make build
./envoy-ai-installer doctor
```

### Daily Development Workflow

```bash
# 1. Create feature branch
git checkout -b feature/your-feature

# 2. Make changes
vim cli/cmd/install.go

# 3. Test locally
make fmt
make vet
make lint
make test

# 4. Build and verify
make build
./envoy-ai-installer doctor

# 5. Commit (pre-commit hooks run automatically)
git add -A
git commit -m "feat(install): add new feature"

# 6. Push
git push origin feature/your-feature

# 7. Create pull request
```

### Make Targets

```bash
make help          # Show all available targets
make build         # Build binary
make clean         # Remove build artifacts
make fmt           # Format Go code
make vet           # Run go vet
make lint          # Run golangci-lint
make test          # Run tests
make doctor        # Run health check
```

### Project Structure

```
cli/
├── main.go         # Entry point
├── cmd/            # Cobra commands
│   ├── root.go     # Root command setup
│   ├── install.go  # Install command
│   ├── version.go  # Version command
│   └── doctor.go   # Health check
└── pkg/            # Packages
    ├── config/     # Configuration (Viper)
    ├── helm/       # Helm operations
    └── upstream/   # Upstream discovery
```

### Key Files

| File | Purpose |
|------|---------|
| `README.md` | User guide & feature overview |
| `CONTRIBUTING.md` | Development guidelines |
| `COMMIT_RULES.md` | Commit practices & guidelines |
| `.pre-commit-config.yaml` | Pre-commit hooks configuration |
| `.commitlintrc.json` | Commit message validation |
| `Makefile` | Build automation |
| `docs/github-actions-setup.md` | CI/CD setup guide |
| `docs/pre-commit-setup.md` | Pre-commit hooks guide |

---

## 📝 Commit Workflow

**Golden Rule:** Commit frequently after every meaningful advancement!

### Commit Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Examples

```bash
git commit -m "feat(install): add gateway installation step"
git commit -m "fix(doctor): handle missing kubectl gracefully"
git commit -m "docs: add rate-limiting configuration guide"
git commit -m "test(helm): add chart validation tests"
```

### Valid Types & Scopes

**Types:** `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`, `ci`

**Scopes:** `cli`, `install`, `config`, `helm`, `upstream`, `doctor`, `version`, `ci`, `docs`, `scripts`

See [COMMIT_RULES.md](COMMIT_RULES.md) for full guidelines.

---

## 🔍 Pre-Commit Hooks

Automatic checks before each commit:

✅ Go formatting & linting  
✅ Shell script validation  
✅ YAML/JSON validation  
✅ Markdown linting  
✅ Secret detection  
✅ Commit message validation  

### Setup

```bash
pre-commit install
pre-commit install --hook-type commit-msg
```

### Run Manually

```bash
pre-commit run --all-files
```

See [docs/pre-commit-setup.md](docs/pre-commit-setup.md) for full guide.

---

## 🧪 Testing

### Local Testing

```bash
# Unit tests
make test

# Build verification
make build

# Health check
./envoy-ai-installer doctor

# Dry-run installation
./envoy-ai-installer install --dry-run
```

### Local Kubernetes Cluster

```bash
# Create test cluster
kind create cluster --name envoy-test

# Test installation
./envoy-ai-installer install

# Verify
kubectl get pods -A

# Cleanup
kind delete cluster --name envoy-test
```

### Testing Workflows Locally

```bash
# Install act (https://github.com/nektos/act)
brew install act

# Run workflows
act -j sync
act -j build
```

---

## 📚 Documentation

- **[README.md](README.md)** — Feature overview & user guide
- **[CONTRIBUTING.md](CONTRIBUTING.md)** — Developer guidelines
- **[COMMIT_RULES.md](COMMIT_RULES.md)** — Commit practices
- **[docs/github-actions-setup.md](docs/github-actions-setup.md)** — CI/CD setup
- **[docs/pre-commit-setup.md](docs/pre-commit-setup.md)** — Pre-commit hooks

---

## 🐛 Troubleshooting

### "go: command not found"

Install Go 1.21+: https://golang.org/doc/install

### "pre-commit: command not found"

```bash
pip install pre-commit
pre-commit install
```

### "golangci-lint not found"

```bash
make lint
```

It installs automatically if needed.

### Build fails

```bash
cd cli
go mod download
go mod tidy
go build -o ../envoy-ai-installer
```

### Pre-commit hook fails

```bash
# Fix automatically
make fmt

# Re-stage and commit
git add -A
git commit -m "message"
```

---

## 📞 Getting Help

1. **Check documentation** — See files above
2. **Run health check** — `./envoy-ai-installer doctor`
3. **Check existing issues** — GitHub Issues
4. **Open new issue** — Include logs and error details

---

## ✅ Verification Checklist

After setup, verify:

- [ ] Repository cloned
- [ ] Pre-commit hooks installed
- [ ] `make build` succeeds
- [ ] `./envoy-ai-installer doctor` shows no errors
- [ ] `make test` passes
- [ ] `make lint` passes
- [ ] Sample commit works

---

## 🎯 Next Steps

1. ✅ Clone and setup (see above)
2. ✅ Read [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines
3. ✅ Read [COMMIT_RULES.md](COMMIT_RULES.md) for commit practices
4. ✅ Create feature branch: `git checkout -b feature/your-feature`
5. ✅ Make changes and commit frequently
6. ✅ Submit pull request

---

## 🔗 Quick Links

- **GitHub:** https://github.com/Franck-Sorel/envoy-ai-unified-installer
- **Envoy AI Gateway:** https://github.com/envoyproxy/ai-gateway
- **Envoy Gateway:** https://gateway.envoyproxy.io/
- **Helm:** https://helm.sh/docs/

---

**Happy developing! 🚀**
