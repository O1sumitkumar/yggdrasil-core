# Runtimes

A runtime adapter starts and stops a model and reports whether it is installed. Adapters register in-process. The interface is `pkg/pluginapi.Runtime`.

## llamacpp

| | |
| --- | --- |
| Id | `llamacpp` |
| Display name | llama.cpp (llama-server) |
| Model format | GGUF |
| Install | `POST /api/v1/runtimes/llamacpp/install`, or the web UI |

Install looks up recent GitHub releases of `ggml-org/llama.cpp` and downloads one archive whose name matches this host:

| Host | Asset name contains |
| --- | --- |
| macOS arm64 | `bin-macos-arm64` |
| macOS amd64 | `bin-macos-x64` |
| Linux amd64 | `bin-ubuntu-x64` |
| Linux arm64 | `bin-ubuntu-arm64` |
| Windows amd64 | `bin-win-cpu-x64` |
| anything else | the GOOS-GOARCH pair, which usually fails the lookup |

The binary is stored under the data directory `runtimes/llamacpp/` unless a macOS build finds a `llama-server` placed beside the daemon. Capability reporting adds `metal` on macOS and `vulkan` on Linux and Windows next to `cpu`. The Windows archive the installer selects is the CPU build named above. CUDA and ROCm are names the hardware inventory can attach to a GPU it sees. They are not a separate installer path in this code.

On a Mac App Store build, the sandbox cannot `fork` a binary downloaded into the container. Those builds are expected to ship a signed `llama-server` next to the daemon. That packaging is outside this repository.

## external-openai

| | |
| --- | --- |
| Id | `external-openai` |
| Display name | External OpenAI-compatible |

This adapter does not download a runtime. It needs a base URL in its configuration. When that URL is set, the adapter reports itself installed and points generation at that server. Prompts sent through it leave the machine. See [privacy.md](privacy.md).

## Adding an adapter

Implement the runtime interface, register it where the app constructs the runtime registry, and cover detection plus start/stop with tests. Open a runtime request or a pull request using the forms in `.github/`. Licensing of the upstream runtime should be compatible with shipping and downloading it. Code contributions are covered by the [Yggdrasil Contributor License Agreement](../CLA.md).
