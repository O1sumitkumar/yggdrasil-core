# Troubleshooting

## The web UI does not load

The daemon only serves a UI when it finds a directory containing `index.html`. From a clone, build it first:

```bash
make frontend
make daemon
./bin/yggdrasil-daemon
```

Then open `http://127.0.0.1:7331`. `GET /api/v1/health` still works if the UI was not built.

## Port already in use

The API defaults to `7331` and Bifrost to `7332`. Stop the other process or start with another data directory only after you change `api_port` and `internal_port` in that directory's `config.json`. The daemon reads the file from the data directory on startup.

## Where are the logs?

`yggctl paths` prints the log directory. The daemon appends JSON to `logs/daemon.log` and to standard output.

Remove API keys, tokens, personal data, private URLs, and credentials before you paste a log into an issue.

## A model will not start

- The file must be a GGUF. `POST /api/v1/models/install-from-url` rejects other URLs.
- llama.cpp must be installed (`POST /api/v1/runtimes/llamacpp/install` or the web UI). Install needs a network path to GitHub.
- On Windows the installer looks for a CPU llama.cpp archive (`bin-win-cpu-x64`). A GPU listed in hardware detection does not by itself select a GPU build.
- A model marked as a tight fit can be refused until you confirm. That check is about memory, not about the file being corrupt.

## Another computer does not appear

- Discovery defaults to on. Both daemons need to be running.
- Bifrost must be reachable on port 7332. On macOS, allow Local Network for the process.
- mDNS does not cross every Docker or VPN boundary. Set `YGGDRASIL_STATIC_PEERS` or `static_peers` to `host:7332`.
- After you enable discovery on a daemon that already bound Bifrost to loopback, restart it.

## Pairing fails

Approve the offer on the second computer. Protected Bifrost routes answer `401` until that pairing is stored. Health and pairing routes are the ones that answer before trust exists.

## OpenAI calls return 401

That happens when the daemon is not bound to loopback and the request has no `Authorization: Bearer YOUR_API_KEY` header. `/api/v1` and `/v1` both use that check. On the default loopback bind, neither requires a key. Create a key in the web UI. Do not paste a real key into a bug report, and do not put one in a URL.

## Chat returns an error about no model

`/v1/chat/completions` does not download a model. Install and start one, or use a profile whose role has a model id. The handler uses the last user message only.

## Reset

`POST /api/v1/settings/reset` clears application state. Pass `delete_models=true` only when you also want model files removed. This is local and not reversible.
