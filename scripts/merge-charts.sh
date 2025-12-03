#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly REPO_ROOT="$(dirname "$SCRIPT_DIR")"
readonly OUT_DIR="${REPO_ROOT}/helm-wrapper/upstream-charts"
readonly VALUES_FILE="${REPO_ROOT}/helm-wrapper/values.yaml"
readonly TEMP_DIR=$(mktemp -d)
readonly LOG_FILE="${REPO_ROOT}/.merge-charts.log"

trap 'rm -rf "$TEMP_DIR"' EXIT
trap 'echo "Script interrupted" >> "$LOG_FILE"' INT TERM

log() {
    local level="$1"
    shift
    local message="$*"
    local timestamp=$(date -u '+%Y-%m-%d %H:%M:%S')
    echo "[$timestamp] [$level] $message" | tee -a "$LOG_FILE"
}

error() {
    log "ERROR" "$@"
    exit 1
}

warn() {
    log "WARN" "$@"
}

info() {
    log "INFO" "$@"
}

validate_tools() {
    local required_tools=("curl" "jq" "tar" "gzip" "python3")
    for tool in "${required_tools[@]}"; do
        if ! command -v "$tool" &> /dev/null; then
            error "Required tool not found: $tool"
        fi
    done
    info "All required tools available"
}

validate_download() {
    local url="$1"
    local output_file="$2"
    
    if [[ ! -f "$output_file" ]]; then
        error "Download file not created: $output_file"
    fi
    
    local file_size=$(stat -f%z "$output_file" 2>/dev/null || stat -c%s "$output_file" 2>/dev/null || echo "0")
    if [[ "$file_size" -le 0 ]]; then
        error "Downloaded file is empty: $output_file"
    fi
    
    local mime_type=$(file -b --mime-type "$output_file" 2>/dev/null || echo "")
    if [[ "$mime_type" != "application/gzip" ]] && [[ "$mime_type" != "application/x-tar" ]] && [[ ! "$mime_type" =~ "compressed" ]]; then
        warn "Unexpected MIME type: $mime_type for $url"
    fi
    
    info "Validated download: $output_file (size: $file_size bytes)"
}

fetch_latest_release() {
    local owner="$1"
    local repo="$2"
    local api_url="https://api.github.com/repos/$owner/$repo/releases/latest"
    
    info "Fetching latest release from $owner/$repo"
    
    local response
    response=$(curl -s -w "\n%{http_code}" "$api_url" \
        -H "Accept: application/vnd.github.v3+json" \
        ${GITHUB_TOKEN:+-H "Authorization: token $GITHUB_TOKEN"})
    
    local http_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | head -n-1)
    
    if [[ "$http_code" != "200" ]]; then
        error "GitHub API returned HTTP $http_code for $owner/$repo"
    fi
    
    echo "$body"
}

find_helm_chart_asset() {
    local release_json="$1"
    local asset_keywords=("helm" "chart" ".tgz" "tar.gz")
    
    local url=$(echo "$release_json" | jq -r '.assets[]? | select(.name | test("helm|chart|tgz|tar\\.gz"; "i")) | .browser_download_url' | head -n1)
    
    if [[ -z "$url" || "$url" == "null" ]]; then
        warn "No chart asset found, attempting fallback to repository archive"
        local tag=$(echo "$release_json" | jq -r '.tag_name')
        if [[ -z "$tag" || "$tag" == "null" ]]; then
            error "No tag found in release"
        fi
        url="https://github.com/$1/$2/archive/refs/tags/${tag}.tar.gz"
    fi
    
    echo "$url"
}

download_file() {
    local url="$1"
    local output_file="$2"
    local max_retries=3
    local retry_count=0
    
    info "Downloading: $url"
    
    while [[ $retry_count -lt $max_retries ]]; do
        if curl -fsSL -o "$output_file" "$url"; then
            validate_download "$url" "$output_file"
            return 0
        fi
        
        retry_count=$((retry_count + 1))
        if [[ $retry_count -lt $max_retries ]]; then
            warn "Download failed, retrying ($retry_count/$max_retries)"
            sleep 2
        fi
    done
    
    error "Failed to download after $max_retries attempts: $url"
}

extract_version() {
    local url="$1"
    local basename=$(basename "$url")
    echo "$basename" | sed -E 's/.*-v?([0-9]+\.[0-9]+\.[0-9]+.*)\..*/\1/' | head -c 20
}

process_upstream() {
    local owner="$1"
    local repo="$2"
    
    local release_json
    release_json=$(fetch_latest_release "$owner" "$repo")
    
    local tag=$(echo "$release_json" | jq -r '.tag_name')
    local url=$(find_helm_chart_asset "$release_json" "$owner" "$repo")
    
    if [[ -z "$url" || "$url" == "null" ]]; then
        error "Could not determine download URL for $owner/$repo"
    fi
    
    local filename="${owner}-${repo}-${tag}.tgz"
    local output_path="${OUT_DIR}/${filename}"
    
    download_file "$url" "$output_path"
    
    info "Successfully processed $owner/$repo: $filename"
    echo "$output_path"
}

update_values_file() {
    if [[ ! -f "$VALUES_FILE" ]]; then
        error "Values file not found: $VALUES_FILE"
    fi
    
    info "Updating values file with generated timestamp"
    
    python3 - "$VALUES_FILE" <<'PYTHON_EOF'
import sys
import yaml
from datetime import datetime

vals_path = sys.argv[1]

try:
    with open(vals_path, 'r') as f:
        vals = yaml.safe_load(f) or {}
    
    vals['generated_at'] = datetime.utcnow().isoformat() + 'Z'
    
    with open(vals_path, 'w') as f:
        yaml.safe_dump(vals, f, default_flow_style=False, sort_keys=False)
    
    print(f"Updated {vals_path} with timestamp: {vals['generated_at']}")
except Exception as e:
    print(f"Error updating values file: {e}", file=sys.stderr)
    sys.exit(1)
PYTHON_EOF
}

main() {
    info "Starting chart merge process"
    
    validate_tools
    
    mkdir -p "$OUT_DIR"
    
    local upstreams=(
        "envoyproxy/gateway"
        "envoyproxy/ai-gateway-helm"
        "envoyproxy/ai-gateway-crds-helm"
        "envoyproxy/ai-gateway"
    )
    
    local processed_count=0
    local failed_count=0
    
    for upstream in "${upstreams[@]}"; do
        IFS='/' read -r owner repo <<< "$upstream"
        
        if process_upstream "$owner" "$repo" 2>&1 | tee -a "$LOG_FILE"; then
            ((processed_count++))
        else
            ((failed_count++))
            warn "Failed to process $upstream"
        fi
    done
    
    update_values_file
    
    info "Chart merge completed: $processed_count successful, $failed_count failed"
    
    if [[ $failed_count -gt 0 ]]; then
        error "Some upstream repos failed to process"
    fi
    
    info "All upstream charts successfully merged"
}

main "$@"
