# Contributing

Yggdrasil Core is the headless control plane. Changes should fit that role: runtimes, models, hardware detection, pairing, placement, health, and the HTTP API. Desktop and mobile clients live in other repositories.

## Contributor License Agreement

Contributors retain copyright ownership of their contributions.

By submitting a contribution to Yggdrasil Core, you agree to the [Yggdrasil Contributor License Agreement](CLA.md), which grants YEIXIO LLC the rights necessary to use, modify, distribute, sublicense, and relicense contributed code, including as part of commercial offerings.

This allows Yggdrasil Core to remain open source under the AGPL while preserving future commercial licensing options.

Pull requests include a checkbox for that agreement. Check it before requesting review.

## Development prerequisites

Go 1.26.3+, Node.js 22, and pnpm 9. Details, commands, and CI are in [docs/development.md](docs/development.md).

## Repository setup

```bash
git clone https://github.com/yeixio/yggdrasil-core.git
cd yggdrasil-core
make start
```

## Build

```bash
make daemon
```

`make package-headless` builds a headless package for the machine you are on. The release script is `scripts/build/package-core-release.sh`.

## Run

```bash
make start
```

Open `http://127.0.0.1:7331`. `make help` lists the other targets.

## Test

```bash
make test
cd web && pnpm test
```

## Lint

```bash
make lint
make vet
```

`make lint` runs golangci-lint v2.14.0 and `pnpm lint` in `web/`. CI runs the same checks. ESLint warnings are reported and do not fail the job.

## Formatting

```bash
make fmt
```

## Submitting an issue

Use the forms in `.github/ISSUE_TEMPLATE/`. Blank issues are disabled. Pick the form that matches the report: bug, hardware, performance, runtime, feature, or documentation. Strip secrets from logs. [SUPPORT.md](SUPPORT.md) says which form to use.

## Repository Labels

GitHub issue labels are defined in `.github/labels.yml`.

Maintainers can synchronize them with:

    ./scripts/sync-github-labels.sh

Running the script requires the GitHub CLI and permission to manage
labels in the repository.

Do not automatically run this workflow on contributor pull requests,
because repository label changes require elevated GitHub permissions.

## Submitting a pull request

Use the pull request template. Check "I agree to the Yggdrasil Contributor License Agreement." Link an issue when there is one. Describe how you tested and on which platform. Keep the change focused.

This repository does not include an `AGENTS.md`. Follow the style of the package you are editing. See [docs/development.md](docs/development.md).

## Runtime adapter contributions

Read [docs/runtimes.md](docs/runtimes.md). An adapter implements `pkg/pluginapi.Runtime` and needs tests for detection and lifecycle. Open a runtime request first if the backend is new, so scope and licensing can be discussed before a large patch.

## Hardware compatibility reports

Use the hardware issue form. Say whether it works, works with limitations, does not work, or you are not sure. Include OS, CPU, GPU, memory, Yggdrasil version, runtime, and models. This surface is large, and a careful report is a real contribution.

## Documentation contributions

User-facing behavior belongs in `docs/user-guide/guide.json` when it changes the public guide, and in `docs/` when it explains the daemon. Run the doc checks in [docs/user-guide/README.md](docs/user-guide/README.md) if you touch the guide.

## Security reports

Do not file a public issue for a vulnerability. Use [SECURITY.md](SECURITY.md).

## Good first contributions

- hardware testing on a machine you already have
- runtime support notes, including Windows GPU builds versus the CPU archive the installer selects
- documentation and examples
- platform packaging (Windows archives, confirming the Homebrew formula on the default branch)
- tests around runtime asset selection and API examples
- bug fixes with a reproduction

## Looking for something to work on?

Start with issues labeled [good first issue](https://github.com/yeixio/yggdrasil-core/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22), [help wanted](https://github.com/yeixio/yggdrasil-core/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22), [hardware](https://github.com/yeixio/yggdrasil-core/issues?q=is%3Aissue+is%3Aopen+label%3Ahardware), or [documentation](https://github.com/yeixio/yggdrasil-core/issues?q=is%3Aissue+is%3Aopen+label%3Adocumentation).

Hardware reports are genuine contributions. Yggdrasil has a large heterogeneous hardware surface, and CI does not generate tokens on those GPUs. A report that names the machine, the model, and the result is more useful than a guess.

Suggested starter issues are listed in [.github/SEED_ISSUES.md](.github/SEED_ISSUES.md). Those are drafts. They are not already filed on GitHub.
