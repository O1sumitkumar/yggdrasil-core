## Downloads

| Package | Description |
|---------|-------------|
| `*-darwin-arm64-desktop.zip` | macOS desktop app (Apple Silicon) |
| `*-windows-amd64-desktop.zip` | Windows desktop app |
| `*-linux-amd64-desktop.tar.gz` | Linux desktop app (GTK/WebKit) |
| `*-linux-amd64-headless.tar.gz` | Daemon + web UI (servers / no GUI) |
| `*-linux-arm64-headless.tar.gz` | Daemon + web UI for ARM64 servers |

### Headless

```bash
tar -xzf yggdrasil-*-linux-amd64-headless.tar.gz
cd yggdrasil-*-linux-amd64-headless
./yggdrasil-daemon
# open http://127.0.0.1:7331
```

### Builds

- macOS desktop builds are Developer ID signed and notarized.
- Windows builds are **unsigned**.
- Desktop packages include `yggdrasil-daemon` beside the UI shell.
