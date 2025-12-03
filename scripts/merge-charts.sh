#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly REPO_ROOT="$(dirname "$SCRIPT_DIR")"
readonly OUTPUT_DIR="${REPO_ROOT}/helm-wrapper/upstream-charts"
readonly VALUES_FILE="${REPO_ROOT}/helm-wrapper/values.yaml"
readonly TEMP_DIR=$(mktemp -d)
readonly LOG_FILE="${REPO_ROOT}/.merge-charts.log"

readonly GITHUB_TOKEN="${GITHUB_TOKEN:-}"
readonly GITHUB_API_VERSION="2022-11-28"
readonly MAX_RETRIES=3

declare -a REPOS=(
    "envoyproxy/gateway"
    "envoyproxy/ai-gateway-helm"
    "envoyproxy/ai-gateway-crds-helm"
)

trap 'rm -rf "$TEMP_DIR"' EXIT
trap 'log ERROR "Script interrupted"' INT TERM

log() {
    local level="$1"
    shift
    local message="$*"
    local timestamp=$(date -u '+%Y-%m-%d %H:%M:%S UTC')
    printf "%s\n" "[$timestamp] [$level] $message" | tee -a "$LOG_FILE"
}

error() {
    log ERROR "$@"
}

warn() {
    log WARN "$@"
}

info() {
    log INFO "$@"
}

validate_tools() {
    local required_tools=("curl" "jq" "tar" "gzip" "python3")
    for tool in "${required_tools[@]}"; do
        if ! command -v "$tool" &>/dev/null; then
            log ERROR "Required tool not found: $tool"
            return 1
        fi
    done
    info "All required tools available"
    return 0
}

get_auth_headers() {
    if [[ -n "$GITHUB_TOKEN" ]]; then
        printf "%s\n" '-H'
        printf "%s\n" "Authorization: Bearer $GITHUB_TOKEN"
        printf "%s\n" '-H'
        printf "%s\n" "X-GitHub-Api-Version: $GITHUB_API_VERSION"
    else
        warn "GITHUB_TOKEN not set - using unauthenticated requests (rate limits apply)"
    fi
}

fetch_latest_release() {
    local owner_repo="$1"
    local api_url="https://api.github.com/repos/${owner_repo}/releases/latest"

    info "Fetching latest release from: $api_url"

    local http_code
    local api_response

    for attempt in 1 2 3; do
        http_code=$(curl -sL -w "%{http_code}" -o /tmp/api_response.json \
            $(get_auth_headers) \
            "$api_url" 2>/dev/null || printf "%s" "000")

        if [[ "$http_code" == "200" ]]; then
            api_response=$(<"/tmp/api_response.json")
            break
        elif [[ "$http_code" == "403" ]]; then
            warn "Rate limited or forbidden (HTTP $http_code) for $owner_repo - attempt $attempt/$MAX_RETRIES"
            if [[ $attempt -lt $MAX_RETRIES ]]; then
                sleep $((2 ** attempt))
                continue
            else
                error "Failed to fetch $owner_repo: HTTP $http_code (rate limit exceeded)"
                return 1
            fi
        elif [[ "$http_code" == "404" ]]; then
            error "Repository not found: $owner_repo (HTTP 404)"
            return 1
        else
            warn "HTTP $http_code from $api_url for $owner_repo - attempt $attempt/$MAX_RETRIES"
            if [[ $attempt -lt $MAX_RETRIES ]]; then
                sleep $((2 ** attempt))
                continue
            else
                error "Failed to fetch $owner_repo after $MAX_RETRIES attempts"
                return 1
            fi
        fi
    done

    if ! jq empty 2>/dev/null <<< "$api_response"; then
        error "Invalid JSON response from GitHub API for $owner_repo"
        warn "Raw response: ${api_response:0:200}"
        return 1
    fi

    printf "%s" "$api_response"
}

get_download_url() {
    local api_json="$1"
    local repo_key="$2"

    local tag_name
    tag_name=$(jq -r '.tag_name // empty' 2>/dev/null <<< "$api_json")

    if [[ -z "$tag_name" ]] || [[ "$tag_name" == "null" ]]; then
        error "No tag_name in API response for $repo_key - skipping"
        return 1
    fi

    info "Extracted tag_name for $repo_key: $tag_name"

    local tgz_url
    tgz_url=$( jq -r '.assets[] | select(.browser_download_url | endswith(".tgz")) | .browser_download_url' 2>/dev/null <<< "$api_json" | head -n 1 || true)

    if [[ -n "$tgz_url" ]] && [[ "$tgz_url" != "null" ]]; then
        printf "%s" "$tgz_url"
        return 0
    fi

    local owner_repo
    owner_repo=$(jq -r '.repository.full_name // empty' 2>/dev/null <<< "$api_json")

    if [[ -z "$owner_repo" ]] || [[ "$owner_repo" == "null" ]]; then
        error "Cannot determine repository name from API response for $repo_key - skipping"
        return 1
    fi

    printf "%s" "https://github.com/${owner_repo}/archive/refs/tags/${tag_name}.tar.gz"
}

download_file() {
    local url="$1"
    local output_file="$2"

    info "Downloading: $url"

    for attempt in 1 2 3; do
        if curl -fL --retry 3 --retry-all-errors -o "$output_file" "$url" 2>/dev/null; then
            if [[ -f "$output_file" ]] && [[ -s "$output_file" ]]; then
                info "Successfully downloaded: $(basename "$output_file")"
                return 0
            fi
        fi

        if [[ $attempt -lt $MAX_RETRIES ]]; then
            warn "Download attempt $attempt/$MAX_RETRIES failed for $url - retrying..."
            sleep 2
        else
            error "Failed to download $url after $MAX_RETRIES attempts"
            return 1
        fi
    done

    return 1
}

process_repo() {
    local owner_repo="$1"

    info "Processing repository: $owner_repo"

    local api_json
    api_json=$(fetch_latest_release "$owner_repo") || return 1

    local tag_name
    tag_name=$(jq -r '.tag_name // empty' 2>/dev/null <<< "$api_json")

    if [[ -z "$tag_name" ]] || [[ "$tag_name" == "null" ]]; then
        error "No valid tag_name for $owner_repo - skipping"
        return 1
    fi

    local download_url
    download_url=$(get_download_url "$api_json" "$owner_repo") || return 1

    if [[ -z "$download_url" ]]; then
        error "No valid download URL for $owner_repo - skipping"
        return 1
    fi

    local repo_name
    repo_name=$(printf "%s" "$owner_repo" | sed 's/.*\///')
    local filename="${repo_name}-${tag_name}.tgz"

    if [[ "$download_url" == *.tar.gz ]]; then
        filename="${repo_name}-${tag_name}.tar.gz"
    fi

    local output_path="${OUTPUT_DIR}/${filename}"

    mkdir -p "$OUTPUT_DIR"

    download_file "$download_url" "$output_path" || return 1

    info "Saved: $filename"
    return 0
}

update_values_yaml() {
    if [[ ! -f "$VALUES_FILE" ]]; then
        warn "values.yaml not found at $VALUES_FILE - skipping update"
        return 0
    fi

    local utc_timestamp
    utc_timestamp=$(date -u '+%Y-%m-%dT%H:%M:%SZ')

    info "Updating values.yaml with generated_at: $utc_timestamp"

    python3 << 'PYTHON_EOF'
import yaml
import sys

try:
    with open(sys.argv[1], 'r') as f:
        values = yaml.safe_load(f) or {}

    values['generated_at'] = sys.argv[2]

    with open(sys.argv[1], 'w') as f:
        yaml.dump(values, f, default_flow_style=False, sort_keys=False)

    print("values.yaml updated successfully", file=sys.stderr)
except Exception as e:
    print(f"Failed to update values.yaml: {e}", file=sys.stderr)
    sys.exit(1)
PYTHON_EOF
    "$VALUES_FILE" "$utc_timestamp"
}

main() {
    info "=== Envoy AI Unified Installer - Merge Charts Script ==="
    info "Output directory: $OUTPUT_DIR"
    info "Values file: $VALUES_FILE"

    if ! validate_tools; then
        error "Tool validation failed"
        return 1
    fi

    mkdir -p "$OUTPUT_DIR"

    local processed_count=0
    local failed_count=0
    local skipped_count=0

    for owner_repo in "${REPOS[@]}"; do
        if process_repo "$owner_repo"; then
            ((processed_count++))
        else
            ((failed_count++))
            ((skipped_count++))
        fi
    done

    update_values_yaml || warn "Failed to update values.yaml"

    info "=== Summary ==="
    info "Processed: $processed_count successful, $skipped_count skipped"

    if [[ -d "$OUTPUT_DIR" ]]; then
        local count=0
        count=$(find "$OUTPUT_DIR" -type f | wc -l)
        info "Total artifacts downloaded: $count"
        if [[ $count -gt 0 ]]; then
            find "$OUTPUT_DIR" -type f -exec ls -lh {} \; | while read -r line; do
                info "  $line"
            done
        fi
    fi

    if [[ $skipped_count -gt 0 ]]; then
        warn "Some repositories were skipped - check logs above"
    fi

    info "Script completed"
}

main "$@"
