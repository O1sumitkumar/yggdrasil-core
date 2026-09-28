# Suggested launch issues

The first public issues are filed. Do not file these drafts again.

- [#3](https://github.com/yeixio/yggdrasil-core/issues/3) Add shell completion for `yggctl`
- [#4](https://github.com/yeixio/yggdrasil-core/issues/4) Return an actionable error when llama-server is missing
- [#5](https://github.com/yeixio/yggdrasil-core/issues/5) Stream chat completions in the Python example
- [#6](https://github.com/yeixio/yggdrasil-core/issues/6) Stream chat completions in the JavaScript example
- [#7](https://github.com/yeixio/yggdrasil-core/issues/7) Report AMD GPU inference on Linux
- [#8](https://github.com/yeixio/yggdrasil-core/issues/8) Report Intel GPU detection and inference
- [#9](https://github.com/yeixio/yggdrasil-core/issues/9) Include stderr and probe errors on Heimdall health events
- [#10](https://github.com/yeixio/yggdrasil-core/issues/10) Document how to build a runtime adapter

`examples/python/chat.py` and `examples/javascript/chat.mjs` already call the API. Issues #5 and #6 add streaming. The helper binary is `yggctl`. There is no Ollama runtime.

These are drafts for later issues. Do not invent a failure that has not been observed.

Each item is an enhancement, a test gap, or a documentation gap that is visible in this repository.

1. **Record a hardware compatibility result**
   Labels: `hardware`, `good first issue`, `help wanted`
   Run detection and one GGUF completion on a machine you have. File it with the hardware form so [docs/compatibility.md](../docs/compatibility.md) can be updated from evidence.

2. **Exercise Windows diagnostics on a real Windows host**
   Labels: `windows`, `hardware`, `good first issue`
   `internal/hardware/windows.go` classifies GPU names. CI only cross-compiles Windows. Confirm CPU, memory, and GPU fields on Windows and note anything that comes back empty.

3. **Add `yggctl` shell completion**
   Labels: `enhancement`, `good first issue`
   `cmd/devctl` implements `version` and `paths` only. Completion for those two commands is a small, self-contained change.

4. **Add an examples test that the curl script is valid shell**
   Labels: `documentation`, `good first issue`, `api`
   [examples/](../examples/) calls `/v1/models` and `/v1/chat/completions`. A CI check that the shell example parses, without requiring a running model, would keep the sample from rotting.

5. **Surface llama.cpp install failures with the asset name that was requested**
   Labels: `runtime`, `enhancement`
   `platformAssetSuffix` fails with `no release asset matching`. Include the suffix and the HTTP status in the error the API returns so a Windows or unusual-arch install is diagnosable.

6. **Test llama.cpp asset selection**
   Labels: `runtime`, `good first issue`
   Table-test `platformAssetSuffix` for darwin/arm64, darwin/amd64, linux/amd64, linux/arm64, and windows/amd64. The Windows expectation is `bin-win-cpu-x64`.

7. **Document a clean-machine source install on each OS**
   Labels: `documentation`, `help wanted`
   The README quick start is a source build. Walk it on macOS, Windows, and Linux from a machine without a prior data directory and note missing steps.

8. **Report an AMD Linux machine**
   Labels: `hardware`, `linux`, `amd`, `help wanted`
   Detection uses `lspci` text and optional `rocminfo`. There is no recorded inference run. A hardware issue that says what happened is the contribution.

9. **Report an Intel GPU machine**
   Labels: `hardware`, `intel`, `help wanted`
   Linux Intel detection runs only when no other accelerator was added, and reports Vulkan and CPU. Windows uses a name heuristic. Both need a real machine before the compatibility table can say more.

10. **Add request identifiers to daemon logs**
    Labels: `enhancement`, `heimdall`
    Logs are JSON. HTTP logs do not attach a request id that also appears on task and placement events. That makes a single chat failure harder to follow across `daemon.log`.

11. **Forward prior messages, temperature, and max_tokens on `/v1/chat/completions`**
    Labels: `api`, `enhancement`
    `internal/api/openai/handler.go` keeps the last user message and ignores `temperature` and `max_tokens`. Documented in [docs/api.md](../docs/api.md). Decide whether to implement the pass-through or keep the limit and return a clear error when the request depends on it.

12. **Authenticate `/api/v1` when the API is not on loopback**
    Status: implemented. `/api/v1` and `/v1` require a bearer token whenever the API host is not loopback, and the daemon refuses to listen there until a key exists. Do not file this as an open issue.

13. **Select a GPU llama.cpp build on Windows when one is published**
    Labels: `runtime`, `windows`, `nvidia`
    The installer looks up `bin-win-cpu-x64` while capability reporting lists Vulkan. Track upstream asset names before changing the lookup.

14. **Offer corresponding source from a running daemon**
    Status: implemented. `GET /about`, `GET /source`, `yggdrasil-daemon -version`, and `yggctl version` already return the license and source URL. Do not file this as an open issue.
