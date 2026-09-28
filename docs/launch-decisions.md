# Launch decisions

Work this list before the repository is presented as public. Check an item only after the decision is made and the matching repo or GitHub change exists.

The longer inventory is [open-source-launch-checklist.md](open-source-launch-checklist.md). This file is the one to walk in order.

## Decisions

- [x] **1. Copyright owner.** YEIXIO LLC. Recorded in `NOTICE`.
- [x] **2. Contributor licensing.** Published in [CLA.md](../CLA.md). Pull requests include an agreement checkbox.
- [ ] **3. Trademark review.** `TRADEMARKS.md` is a short policy. Have it reviewed before a commercial launch. The repository can stay private or public while that review is pending, as long as the TODO stays visible.
- [x] **4. Vulnerability reporting.** GitHub private vulnerability reporting is enabled. The process is in `SECURITY.md`.
- [x] **5. Code of conduct contact.** Reports go to conduct@yeix.io. Recorded in `CODE_OF_CONDUCT.md`.
- [x] **6. Sponsorship.** `.github/FUNDING.yml` lists GitHub Sponsors user `gopherstein`.
- [x] **7. AGPL source offer.** `GET /about` and `GET /source` return name, version, commit, `AGPL-3.0-or-later`, and the source URL for that build. Release builds point at `tree/v<version>`. `yggdrasil-daemon -version` and `yggctl version` print the same notice. Forks set `version.SourceURL`.
- [x] **8. Control API authentication.** Decision: `/api/v1` requires bearer-token authentication whenever Yggdrasil binds beyond loopback. Loopback access may remain unauthenticated by default. Yggdrasil must refuse non-loopback control API exposure unless authentication is configured. `/v1` uses the same bind rule. Keys are stored as hashes. A bearer token on plain HTTP does not encrypt traffic.
- [x] **9. Release footer.** `docs/releases/_footer.md` describes Core release artifacts and checksum verification. Desktop builds are not published from this repository.

## Actions before public

- [ ] **10. GitHub Discussions.** Enable Discussions. Suggested categories: General, Ideas, Hardware & Runtimes, Show and Tell, Help.
- [x] **11. Labels.** `.github/labels.yml` is the definition. `./scripts/sync-github-labels.sh` created or updated those labels on GitHub. CI does not run that script.
- [x] **12. Seed issues.** Issues [#3](https://github.com/yeixio/yggdrasil-core/issues/3) through [#10](https://github.com/yeixio/yggdrasil-core/issues/10) are open. Four are `good first issue`: shell completion, the missing llama-server error, and streaming in the Python and JavaScript examples.
- [x] **13. Git history.** Reviewed all 15 commits, from the 2026-09-25 initial commit through `629b0af`. No private keys, token prefixes, credential files, or assigned secrets. The GitHub repository is public. History was not rewritten.
- [ ] **14. Homebrew formula.** Do not document `brew install` until `Formula/yggdrasil.rb` is on `main` and a clean install works. The release workflow builds the macOS archives, fills the formula from `SHA256SUMS.txt`, and opens a pull request from the `formula` branch. This repository is the tap. A separate `homebrew-yggdrasil` repository, and a formula in Homebrew core, are not set up. The installed commands are `yggdrasil-daemon` and `yggctl`.
  - [ ] Release created
  - [ ] Release workflow produced the macOS headless archives
  - [ ] Formula pull request opened
  - [ ] Formula pull request merged into `main`
  - [ ] Formula URL and SHA256 match that release
  - [ ] `brew install` tested on a machine without a local dev build
  - [ ] `yggdrasil-daemon -version` runs
  - [ ] `yggctl version` runs
  - [ ] Uninstall and reinstall work
  - [ ] README matches the command that was tested
- [ ] **15. Clean install.** Try a clean install on macOS, Linux, and Windows. The Windows amd64 archive is produced by the release script and is not on a published release yet.
- [ ] **16. Demo.** `make screenshots` records a walkthrough of the demo UI into `docs/screenshots/demo.mp4`. A live recording is still open: start the daemon, discover another computer, install and run a model, send an API request, and show Norn placing a role.

## Suggested order

1. Items 1, 2, and 5 are text you can settle in this repository.
2. Item 7 is a product decision. Items 8 and 9 are implemented.
3. Items 4, 6, 10, and 14 are GitHub settings or a release. They are not done by editing files alone. Items 11 and 12 are done.
4. Items 15 and 16 are verification. Item 13 is done. Item 3 can stay open until a lawyer reviews it, as long as the TODO remains.
