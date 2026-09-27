#!/usr/bin/env bash
# Build the web UI, serve demo data, and capture the README screenshots.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
created_workspaces=()
cleanup_workspaces() {
  if ((${#created_workspaces[@]})); then
    rm -f "${created_workspaces[@]}"
  fi
}
trap cleanup_workspaces EXIT

install_pnpm() {
  local dir="$1"
  local major
  major="$(cd /tmp && pnpm --version | cut -d. -f1)"
  if [[ "$major" -ge 10 && ! -f "$dir/pnpm-workspace.yaml" ]]; then
    printf '%s\n' 'packages:' '  - "."' 'allowBuilds:' '  esbuild: true' '  playwright: true' > "$dir/pnpm-workspace.yaml"
    created_workspaces+=("$dir/pnpm-workspace.yaml")
  fi
  (cd "$dir" && pnpm install)
}

if [[ ! -f web/dist/index.html ]]; then
  install_pnpm web
  (cd web && pnpm build)
fi

install_pnpm scripts/screenshots
if [[ "$(uname -s)" == "Linux" ]]; then
  (cd scripts/screenshots && pnpm exec playwright install --with-deps chromium)
else
  (cd scripts/screenshots && pnpm exec playwright install chromium)
fi

node scripts/screenshots/capture.mjs

if command -v magick >/dev/null 2>&1 || command -v convert >/dev/null 2>&1; then
  chmod +x scripts/process-screenshots.sh scripts/validate-screenshots.sh
  ./scripts/process-screenshots.sh
  ./scripts/validate-screenshots.sh
fi

echo "README screenshots: docs/screenshots"
