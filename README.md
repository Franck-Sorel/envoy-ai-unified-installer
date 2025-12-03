# Envoy AI Unified Installer

A **production-grade** Go CLI + GitHub Actions pipeline for installing Envoy AI Gateway with automatic upstream synchronization.

## 🎯 Features

✅ **Automated upstream sync** — GitHub Actions periodically fetches latest Envoy Gateway & AI Gateway releases  
✅ **Production-ready CLI** — Cobra/Viper-based with full configuration support  
✅ **Official 4-step install** — Implements exact Envoy AI Gateway installation process  
✅ **Dry-run mode** — Preview all changes before applying  
✅ **Health checks** — `doctor` command validates prerequisites  
✅ **Version tracking** — Show CLI and upstream component versions  
✅ **GitHub Pages Helm repo** — Optional pre-built chart repository  
✅ **Comprehensive logging** — Detailed output and error messages  

## 📋 Table of Contents

- [Quick Start](#quick-start)
- [CLI Commands](#cli-commands)
- [Project Structure](#project-structure)
- [GitHub Actions Setup](#github-actions-setup)
- [Configuration](#configuration)
- [Development](#development)
- [Security](#security)

---

## 🚀 Quick Start

### Prerequisites

- **kubectl** (1.21+) — [Install](https://kubernetes.io/docs/tasks/tools/)
- **helm** (3.10+) — [Install](https://helm.sh/docs/intro/install/)
- **Go** (1.21+) — [Install](https://golang.org/doc/install) (for building from source)
- Kubernetes cluster access

### Installation Steps

#### 1. Clone & Build

```bash
git clone https://github.com/Franck-Sorel/envoy-ai-unified-installer.git
cd envoy-ai-unified-installer/cli

go build -o ../envoy-ai-installer
cd ..
```

#### 2. Run Health Check

```bash
./envoy-ai-installer doctor
```

Expected output:
```
🏥 System Health Check

🔍 kubectl:            ✅ v1.28.0
🔍 Helm:               ✅ v3.12.0
🔍 Kubernetes cluster: ✅ CONNECTED
🔍 Namespace 'envoy-gateway-system':    ⚠️ NOT FOUND (will be created)
🔍 Namespace 'envoy-ai-gateway-system': ⚠️ NOT FOUND (will be created)

✅ All checks passed! You're ready to install Envoy AI Gateway.
```

#### 3. Install

```bash
./envoy-ai-installer install
```

With optional Redis for rate limiting:

```bash
./envoy-ai-installer install --with-redis
```

#### 4. Verify

```bash
kubectl get pods -n envoy-gateway-system
kubectl get pods -n envoy-ai-gateway-system
```

---

## 🛠️ CLI Commands

### `install` — Install Envoy AI Gateway

Implements the official 4-step installation process:

1. Clean previous installations (optional)
2. Install Envoy Gateway with official values
3. Install Envoy AI Gateway CRDs
4. Install Envoy AI Gateway controller

**Flags:**

```bash
--namespace-gateway string          Kubernetes namespace for Envoy Gateway (default: envoy-gateway-system)
--namespace-ai string                Kubernetes namespace for Envoy AI (default: envoy-ai-gateway-system)
--values-extra string                Comma-separated list of additional values files
--with-redis                         Install Redis (bitnami) for rate limiting
--skip-clean                         Skip cleaning up previous installations
--dry-run                            Preview changes without applying
--config string                      Config file path
```

**Examples:**

```bash
./envoy-ai-installer install

./envoy-ai-installer install --namespace-gateway prod-gw --namespace-ai prod-ai

./envoy-ai-installer install --values-extra rate-limit.yaml,inference-pool.yaml

./envoy-ai-installer install --dry-run
```

### `version` — Show Version Information

Display CLI version and upstream component versions.

```bash
./envoy-ai-installer version
```

Output:
```
📦 envoy-ai-installer Version Information

  CLI Version:    0.1.0
  Git Commit:     a1b2c3d
  Build Time:     2024-01-10T15:30:00Z

  Helm Version:   v3.12.0

📋 Upstream Component Versions

  envoyproxy/gateway:              v0.6.0
  envoyproxy/ai-gateway-helm:      v0.2.1
  envoyproxy/ai-gateway-crds-helm: v0.2.1
  envoyproxy/ai-gateway:           v0.2.1
```

### `doctor` — Health Check

Validate system prerequisites and cluster connectivity.

```bash
./envoy-ai-installer doctor
```

Checks:
- kubectl availability and version
- Helm availability and version
- Kubernetes cluster connectivity
- Required namespaces
- Optional Redis installation

---

## 📂 Project Structure

```
envoy-ai-unified-installer/
├── cli/                           # Go CLI source
│   ├── main.go                    # Entry point
│   ├── go.mod                     # Module definition
│   ├── go.sum                     # Dependency checksums
│   ├── cmd/                       # Cobra commands
│   │   ├── root.go                # Root command & config
│   │   ├── install.go             # Install command
│   │   ├── version.go             # Version command
│   │   └── doctor.go              # Doctor command
│   └── pkg/                       # Internal packages
│       ├── config/                # Configuration management (Viper)
│       │   └── config.go
│       ├── helm/                  # Helm operations
│       │   └── helm.go
│       └── upstream/              # Upstream chart discovery
│           └── upstream.go
├── helm-wrapper/                  # Helm chart for unified installation
│   ├── Chart.yaml                 # Chart metadata
│   ├── values.yaml                # Default values
│   └── upstream-charts/           # Generated by CI
├── .github/
│   └── workflows/                 # GitHub Actions workflows
│       ├── sync-upstream.yml      # Sync upstream releases (6h schedule)
│       └── release-chart.yml      # Package & publish chart
├── scripts/
│   └── merge-charts.sh            # Download & validate upstream charts
├── docs/
│   └── github-actions-setup.md    # Complete setup guide
├── README.md
└── LICENSE
```

---

## 🔄 GitHub Actions Setup

For detailed setup instructions, see **[docs/github-actions-setup.md](docs/github-actions-setup.md)**.

### Quick Setup

1. **Create secrets** in repository Settings → Secrets → Actions:
   - `GH_PAGES_DEPLOY_PAT` (Personal Access Token with `public_repo`, `workflow` scopes)
   - `ACTIONS_DEPLOY_KEY` (SSH private key for GitHub Pages) — optional but recommended

2. **Enable GitHub Pages**:
   - Settings → Pages
   - Source: `gh-pages` branch, `/` (root) folder
   - Save

3. **Run workflows manually** to test:
   - Actions → "Sync Upstream Releases" → Run workflow
   - Actions → "Build & Publish Helm Chart" → Run workflow

### How It Works

**sync-upstream.yml** (Every 6 hours):
- Fetches latest releases from upstream repos
- Validates downloads (HTTP 200, correct MIME type, non-empty files)
- Updates `helm-wrapper/values.yaml`
- Commits & pushes changes to default branch

**release-chart.yml** (On helm-wrapper changes):
- Packages Helm chart
- Publishes to GitHub Pages as Helm repository

### Use Published Chart

After workflows complete:

```bash
helm repo add envoy-ai https://<USERNAME>.github.io/<REPO>
helm repo update
helm upgrade --install envoy-ai-unified envoy-ai/unified \
  -n envoy-ai-gateway-system --create-namespace
```

---

## ⚙️ Configuration

### Config File

Create `~/.envoy-ai-installer/config.yaml`:

```yaml
namespace_gateway: envoy-gateway-system
namespace_ai: envoy-ai-gateway-system
skip_clean: false
dry_run: false
values_extra:
  - /path/to/rate-limit.yaml
  - /path/to/inference-pool.yaml
```

### Environment Variables

Override config with `EAIG_*` prefix:

```bash
export EAIG_NAMESPACE_GATEWAY=prod-gateway
export EAIG_NAMESPACE_AI=prod-ai
export EAIG_DRY_RUN=true

./envoy-ai-installer install
```

### Command-Line Flags

Flags override both config and environment variables:

```bash
./envoy-ai-installer install \
  --config ~/.envoy-ai-installer/config.yaml \
  --namespace-gateway prod-gw \
  --namespace-ai prod-ai \
  --values-extra custom-values.yaml \
  --dry-run
```

---

## 🔧 Development

### Building from Source

```bash
cd cli
go mod download
go mod tidy
go build -o ../envoy-ai-installer
```

### Building with Version Info

```bash
go build \
  -ldflags="-X main.version=0.1.0 \
            -X main.gitCommit=$(git rev-parse --short HEAD) \
            -X main.buildTime=$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
  -o ../envoy-ai-installer
```

### Local Testing

Use `kind`, `minikube`, or `k3s`:

```bash
kind create cluster --name envoy-test
./envoy-ai-installer install --dry-run
./envoy-ai-installer install
kubectl get pods -A
kind delete cluster --name envoy-test
```

### Testing Workflows Locally

Install [act](https://github.com/nektos/act):

```bash
act -j sync
act -j build
```

---

## 🔒 Security

### Principles

- **Zero trust upstream:** All artifacts downloaded from official upstream sources
- **Validation:** All downloads validated (HTTP 200, file size, MIME type)
- **Dry-run mode:** Preview all changes before applying
- **Least privilege:** CLI only performs required Helm operations
- **No secrets in code:** All credentials managed via GitHub Secrets

### Best Practices

1. **Rotate PAT tokens** every 90 days
2. **Use SSH keys** for GitHub Pages deployment (see setup guide)
3. **Enable branch protection** on main branch
4. **Review workflows** in pull requests
5. **Audit Actions logs** regularly
6. **Monitor cluster resources** post-installation

### Secret Management

- `GITHUB_TOKEN` — Auto-provided by GitHub Actions (7-hour expiration)
- `GH_PAGES_DEPLOY_PAT` — Personal Access Token with limited scopes
- `ACTIONS_DEPLOY_KEY` — SSH key for secure Git operations

---

## 📚 Additional Resources

- [Official Envoy AI Gateway Docs](https://github.com/envoyproxy/ai-gateway)
- [Envoy Gateway Docs](https://gateway.envoyproxy.io/)
- [Helm Documentation](https://helm.sh/docs/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [GitHub Actions Guide](docs/github-actions-setup.md)

---

## 🐛 Troubleshooting

### `doctor` shows warnings

Run the doctor command for diagnostics:

```bash
./envoy-ai-installer doctor
```

### Helm charts not found

Clear Helm cache and update:

```bash
helm repo update --force-update
helm search repo envoy
```

### Workflow failed

1. Check **Actions** → workflow run logs
2. Verify secrets: Settings → Secrets and variables → Actions
3. Check `.merge-charts.log` in repository
4. Run `merge-charts.sh` locally with debugging:
   ```bash
   bash -x scripts/merge-charts.sh
   ```

### Installation fails

Use `--dry-run` to preview:

```bash
./envoy-ai-installer install --dry-run --verbose
```

Check Kubernetes events:

```bash
kubectl get events -n envoy-gateway-system
kubectl get events -n envoy-ai-gateway-system
```

---

## 📄 License

See [LICENSE](LICENSE) file for details.

---

## 🤝 Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make changes with clear commit messages
4. Test locally with `--dry-run`
5. Submit pull request with detailed description

---

## ⭐ Support

If you found this helpful, please star the repository and share feedback!
