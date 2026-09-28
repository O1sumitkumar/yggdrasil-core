# Development

## Prerequisites

- Go 1.26.3 or newer (`go.mod`)
- Node.js 22
- pnpm 9 (the root `package.json` pins `pnpm@9.15.0`)
- Git

Linux package builds also need `nfpm` 2.41.3, as used by the release workflow. Day-to-day daemon work does not.

## Setup

```bash
git clone https://github.com/yeixio/yggdrasil-core.git
cd yggdrasil-core
make frontend
make daemon
```

`make daemon` writes `bin/yggdrasil-daemon` and `bin/yggctl`, and stamps the current git commit into both. `yggctl version` and `yggdrasil-daemon -version` print the license and the corresponding-source URL. Release packaging sets the version as well, so a tagged build points at `tree/v<version>`.

A fork that serves a modified daemon over the network sets its own source URL at build time:

```text
-X github.com/yeixio/yggdrasil-core/internal/version.SourceURL=<url-of-your-corresponding-source>
```

## Run

```bash
./bin/yggdrasil-daemon
```

Optional data directory:

```bash
./bin/yggdrasil-daemon -data-dir "$PWD/.ygg-dev-data"
```

The web UI is served from `web/dist` when that build exists. The Vite dev server is separate:

```bash
make run-web
```

That listens on `http://127.0.0.1:5173`. The daemon API stays on port 7331.

## Test

```bash
make test
make vet
```

`make test` is `go test ./...`. `make ci` runs format, vet, tests, and the frontend job locally.

Web tests alone:

```bash
cd web && pnpm test
```

Cluster check (Docker, stub inference, not a GPU test):

```bash
make test-cluster
```

Documentation snapshots:

```bash
python3 scripts/test_publish_docs.py
python3 scripts/publish-docs.py --destination /tmp/ygg-docs-check \
  --version 0.0.0-test --commit "$(git rev-parse HEAD)"
```

## Format and lint

Go formatting is `gofmt` via `make fmt`. CI fails if a Go file outside `web/`, `vendor/`, and `node_modules/` is not `gofmt`-clean.

Go lint in CI is `go vet ./...`. There is no golangci-lint config.

The web check in CI is `pnpm exec tsc -b --pretty false`, `pnpm test`, and `pnpm build`. There is no ESLint script in `web/package.json`.

## CI

[`.github/workflows/ci.yml`](../.github/workflows/ci.yml) on Ubuntu:

- user-guide publish checks
- `gofmt`, `go vet`, `go test ./...`
- cross-compiles of `yggdrasil-daemon` and `yggctl` for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, and windows/amd64 (`CGO_ENABLED=0`)
- web typecheck, test, and build

[`.github/workflows/security.yml`](../.github/workflows/security.yml) runs `govulncheck ./...` and `pnpm audit --prod` in `web/`. The audit step does not fail the job (`|| true`).

[`.github/workflows/release.yml`](../.github/workflows/release.yml) runs on tags matching `v*`. It builds packages, writes `SHA256SUMS.txt`, and publishes a GitHub Release. It does not sign binaries.

[`.github/workflows/screenshots.yml`](../.github/workflows/screenshots.yml) recaptures `docs/screenshots` on a tag or when started by hand, then holds those stills into `demo.mp4` and `demo.gif`.

Dependabot is configured for Go modules, the web and screenshot npm trees, and GitHub Actions.

## Labels

Issue labels are defined in [`.github/labels.yml`](../.github/labels.yml). GitHub does not create them from that file. Maintainers synchronize them with `./scripts/sync-github-labels.sh`. CI does not run that script.

## Release flow

Tag `v*` → release workflow → Linux `.deb` and `.rpm` (amd64 and arm64), macOS headless archives (arm64 and amd64), Windows amd64 headless archive, `SHA256SUMS.txt` → GitHub Release.

The same job opens a Homebrew formula pull request and updates the `apt` branch. Signing is not part of this workflow. The checklist is [release-checklist.md](release-checklist.md).

## Conventions

This repository does not include an `AGENTS.md`. Match the package you are editing. Keep handlers thin and put behavior in `internal/`. Do not add license headers file by file. The project license is the root `LICENSE`.
