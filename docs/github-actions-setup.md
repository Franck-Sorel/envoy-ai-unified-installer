# GitHub Actions Setup Guide

This guide explains how to set up GitHub Actions workflows for the Envoy AI Unified Installer to automatically sync upstream releases and publish Helm charts to GitHub Pages.

## Prerequisites

- A GitHub repository (either created from this template or forked)
- Write access to the repository
- GitHub Pages enabled on your repository
- Administrator access to configure secrets

---

## 1. Required Secrets

The GitHub Actions workflows require the following secrets to be configured in your repository:

| Secret Name | Purpose | Required | Default |
|-------------|---------|----------|---------|
| `GITHUB_TOKEN` | GitHub API access (auto-provided by Actions) | ✅ Yes | Auto-provided |
| `GH_PAGES_DEPLOY_PAT` | Deploy Helm charts to GitHub Pages | ✅ Yes | Configure manually |

### Additional Optional Secrets

| Secret Name | Purpose | Optional |
|-------------|---------|----------|
| `PAT_RELEASES` | Higher rate limits for GitHub API | ❌ No |
| `ACTIONS_DEPLOY_KEY` | SSH key for pushing to gh-pages | ❌ No |

---

## 2. Creating Secrets: Step-by-Step

### Step 1: Generate a Personal Access Token (PAT)

1. Go to https://github.com/settings/tokens
2. Click **"Generate new token"** → **"Generate new token (classic)"**
3. Fill in the form:
   - **Token name:** `GH_PAGES_DEPLOY_PAT`
   - **Expiration:** 90 days (recommended) or No expiration
   - **Scopes:** Select:
     - ✅ `public_repo` (push to public repos)
     - ✅ `workflow` (modify workflows)
4. Click **"Generate token"**
5. **Copy the token immediately** (you won't see it again)

### Step 2: Add the Secret to Your Repository

1. Go to your repository on GitHub
2. Click **Settings** → **Secrets and variables** → **Actions**
3. Click **"New repository secret"**
4. Fill in:
   - **Name:** `GH_PAGES_DEPLOY_PAT`
   - **Secret:** Paste the token from Step 1
5. Click **"Add secret"**

### Step 3: (Optional) Create Additional PAT for Higher API Rate Limits

If you want higher GitHub API rate limits (useful for CI jobs that run frequently):

1. Repeat Steps 1-2 above
2. Name this token `PAT_RELEASES`
3. Same scopes as above

---

## 3. SSH Key Setup (Optional but Recommended)

For secure, key-based authentication to GitHub Pages deployment:

### Generate SSH Key

Run this command on your local machine:

```bash
ssh-keygen -t ed25519 -C "envoy-ai-installer-gh-pages" -N "" -f ~/.ssh/envoy_ai_gh_pages
```

This creates two files:
- `~/.ssh/envoy_ai_gh_pages` (private key)
- `~/.ssh/envoy_ai_gh_pages.pub` (public key)

### Add Public Key to GitHub Deploy Keys

1. Go to your repository on GitHub
2. Click **Settings** → **Deploy keys**
3. Click **"Add deploy key"**
4. Fill in:
   - **Title:** `GitHub Pages Deploy Key`
   - **Key:** Paste the contents of `~/.ssh/envoy_ai_gh_pages.pub`
   - **Allow write access:** ✅ YES (required for publishing)
5. Click **"Add key"**

### Add Private Key to Repository Secrets

1. Go to **Settings** → **Secrets and variables** → **Actions**
2. Click **"New repository secret"**
3. Fill in:
   - **Name:** `ACTIONS_DEPLOY_KEY`
   - **Secret:** Paste the contents of `~/.ssh/envoy_ai_gh_pages`
4. Click **"Add secret"**

### Secure the Local Files

```bash
chmod 600 ~/.ssh/envoy_ai_gh_pages
chmod 644 ~/.ssh/envoy_ai_gh_pages.pub
```

---

## 4. Enable GitHub Pages

### Enable GitHub Pages Hosting

1. Go to your repository on GitHub
2. Click **Settings** → **Pages**
3. Under **Source**, select:
   - **Deploy from a branch**
   - **Branch:** `gh-pages`
   - **Folder:** `/ (root)`
4. Click **"Save"**

GitHub will automatically create the `gh-pages` branch and enable Pages. Wait 1-2 minutes for it to be available.

### Verify GitHub Pages URL

After enabling, you should see:

```
Your site is live at https://<YOUR_USERNAME>.github.io/<REPO_NAME>/
```

---

## 5. How CI Sync Works

### Workflow: `sync-upstream.yml`

**Trigger:** Runs on a schedule (every 6 hours) or manually via `workflow_dispatch`

**Steps:**

1. **Checkout:** Clone the repository
2. **Setup tools:** Install Node.js, Helm, and required CLI tools
3. **Run merge-charts.sh:** Fetch latest upstream releases from:
   - `envoyproxy/gateway`
   - `envoyproxy/ai-gateway-helm`
   - `envoyproxy/ai-gateway-crds-helm`
   - `envoyproxy/ai-gateway`
4. **Validate:** Check that all downloads are valid (HTTP 200, correct MIME type, non-empty)
5. **Commit & Push:** If new charts were downloaded:
   - Create a commit with message: `ci: sync upstream charts <timestamp>`
   - Push changes to the default branch

**Environment Variables:**
- `GITHUB_TOKEN`: Used for GitHub API requests (provided by Actions)

### Workflow: `release-chart.yml`

**Trigger:** Runs when `helm-wrapper/` directory changes or via `workflow_dispatch`

**Steps:**

1. **Checkout:** Clone the repository
2. **Setup Helm:** Install Helm CLI
3. **Package Chart:** Run `helm package helm-wrapper -d packaged`
4. **Deploy to Pages:** Publish packaged chart to `gh-pages` branch using `peaceiris/actions-gh-pages@v4`

**Result:**
- Helm chart is packaged as `.tgz`
- Published to GitHub Pages at:
  ```
  https://<USERNAME>.github.io/<REPO>/
  ```

---

## 6. Testing the Workflows

### Manually Trigger Sync Workflow

1. Go to your repository
2. Click **Actions**
3. Select **"Sync Upstream Releases"** workflow
4. Click **"Run workflow"** → **"Run workflow"**

Check logs to verify:
```
✅ All upstream charts successfully merged
```

### Manually Trigger Release Workflow

1. Click **Actions** → **"Build & Publish Helm Chart"**
2. Click **"Run workflow"** → **"Run workflow"**

After ~2 minutes, verify:
```bash
helm repo add envoy-ai https://<USERNAME>.github.io/<REPO>
helm repo update
helm search repo envoy-ai
```

You should see the `envoy-ai-unified` chart in the output.

---

## 7. Troubleshooting

### Issue: "403 Forbidden" on Push

**Cause:** The PAT doesn't have sufficient permissions.

**Solution:**
1. Verify PAT has these scopes: `public_repo`, `workflow`
2. Try regenerating with "No expiration"
3. Use SSH keys instead (see Section 3)

### Issue: "Cannot find module" in Go build

**Cause:** Go dependencies not downloaded.

**Solution:**
```bash
cd cli
go mod download
go mod tidy
```

### Issue: Workflow shows "missing chart repo metadata"

**Cause:** The `helm-wrapper/values.yaml` doesn't have proper structure.

**Solution:**
```bash
helm lint helm-wrapper/
```

Check output and ensure `Chart.yaml` and `values.yaml` are valid YAML.

### Issue: GitHub Pages shows 404

**Cause:** GitHub Pages branch or folder configuration is incorrect.

**Solution:**
1. Verify **Settings** → **Pages** shows:
   - Source: `gh-pages` branch, `/` (root) folder
2. Push a test commit to `gh-pages`:
   ```bash
   git push origin gh-pages --force
   ```
3. Wait 2-5 minutes and refresh

### Issue: Helm charts not syncing

**Cause:** GitHub API rate limits or network issues.

**Solution:**
1. Check workflow logs: **Actions** → select workflow run → view logs
2. Look for HTTP 403 errors (rate limit)
3. If rate-limited, add `PAT_RELEASES` secret and update workflow to use it:
   ```yaml
   -H "Authorization: token ${{ secrets.PAT_RELEASES }}"
   ```
4. Check network connectivity: run `curl -I https://api.github.com/`

### Issue: Caching problems

**Cause:** Helm repo cache is stale.

**Solution:**
```bash
helm repo remove envoy-ai
helm repo add envoy-ai https://<USERNAME>.github.io/<REPO>
helm repo update --force-update
```

---

## 8. Next Steps

### Local Installation Using Published Chart

After the workflows complete successfully:

```bash
helm repo add envoy-ai https://<USERNAME>.github.io/<REPO>
helm repo update
helm upgrade --install envoy-ai-unified envoy-ai/unified -n envoy-ai-gateway --create-namespace
```

Or use the CLI:

```bash
envoy-ai-installer install
```

### Monitoring Updates

- **Manual check:** Run `envoy-ai-installer version` to see latest upstream versions
- **Automatic notifications:** Enable GitHub Actions notifications in your account settings
- **Workflow runs:** Visit **Actions** tab to see sync history

---

## 9. Security Considerations

### Token Security

- **Never commit secrets** to the repository
- **Rotate PAT regularly** (every 90 days recommended)
- Use **GitHub's secret scanning** to detect accidental commits
- **Scope tokens minimally** (only required permissions)

### SSH Key Security

- **Encrypt** the private key when storing locally:
  ```bash
  ssh-keygen -p -f ~/.ssh/envoy_ai_gh_pages
  ```
- **Rotate keys** annually or after any security incident
- **Never share** the private key

### Repository Security

- **Protect main branch:** Require PR reviews before merge
- **Enable signed commits:** Configure GPG signing
- **Audit Actions:** Review workflow permissions and secrets regularly
- **Monitor CI/CD logs** for suspicious activity

---

## Reference

### Useful Links

- [GitHub Secrets Documentation](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
- [GitHub Pages Documentation](https://docs.github.com/en/pages)
- [Helm Chart Publishing](https://helm.sh/docs/chart_repository/)
- [GitHub Actions Best Practices](https://docs.github.com/en/actions/guides)

### Support

For issues or questions:
1. Check workflow logs in **Actions** tab
2. Review `.merge-charts.log` in repository root
3. Run `envoy-ai-installer doctor` for local environment checks
4. Open an issue on GitHub with detailed logs and error messages
