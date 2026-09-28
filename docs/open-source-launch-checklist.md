# Open-source launch checklist

Do this before the repository is presented as the public project. Do not treat a cloned working tree as a finished launch.

The decisions and gates to walk in order are in [launch-decisions.md](launch-decisions.md).

- [ ] AGPL-3.0-or-later confirmed. The text is [LICENSE](../LICENSE). The SPDX identifier is recorded in [NOTICE](../NOTICE).
- [x] Copyright owner confirmed. NOTICE names YEIXIO LLC.
- [x] Contributor licensing policy published in [CLA.md](../CLA.md). Pull requests include an agreement checkbox.
- [ ] Trademark policy reviewed. [TRADEMARKS.md](../TRADEMARKS.md) is a short policy, not a finished legal review.
- [x] Security reporting uses GitHub private vulnerability reporting. The steps are in [SECURITY.md](../SECURITY.md).
- [ ] GitHub Discussions enabled. Suggested categories: General, Ideas, Hardware & Runtimes, Show and Tell, Help. See [SUPPORT.md](../SUPPORT.md).
- [x] GitHub Sponsors account is `gopherstein` in [.github/FUNDING.yml](../.github/FUNDING.yml). The Sponsors profile still has to be enabled on that GitHub account.
- [ ] Repository description set. Suggested text: "Headless local-AI orchestration for the computers you already own."
- [ ] Repository topics set. Relevant topics: `local-ai`, `llm`, `ai`, `llama-cpp`, `distributed-computing`, `inference`, `golang`, `self-hosted`, `openai-api`, `apple-silicon`, `nvidia`, `amd`, `intel`.
- [ ] README complete and links checked.
- [ ] A current release available, with checksums, and the install path in the README tried on a clean machine.
- [ ] Clean install tested for macOS, Linux, and Windows.
- [ ] Issue forms working, and the labels in [.github/labels.yml](../.github/labels.yml) created so forms can apply them.
- [ ] 5–10 of the drafts in [.github/SEED_ISSUES.md](../.github/SEED_ISSUES.md) filed, including several `good first issue` items.
- [ ] Demo video or GIF prepared. `make screenshots` writes a demo-UI walkthrough to `docs/screenshots/demo.mp4` and `demo.gif`. A live recording of startup, discovery, a model, an API request, and Norn placement is still open.
- [ ] Announcement copy prepared.
- [ ] Git history reviewed for secrets. A scan of the current checkout did not find committed credentials. History can still contain something the current tree does not. Review it before publication if this repository was ever private.
- [ ] Private and proprietary code audit completed. Desktop and mobile source should stay out of this tree.

## AGPL network use

- [x] A running daemon offers corresponding source at `GET /about` and `GET /source`. The Settings page links to it. Release builds point at the tag. Forks override `version.SourceURL`. See [api.md](api.md).
