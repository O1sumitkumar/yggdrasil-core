# Linux packaging

## Headless (recommended for servers)

Built by CI as `yggdrasil-<version>-linux-<arch>-headless.tar.gz`.

```bash
VERSION=0.1.0 GOOS=linux GOARCH=amd64 ./scripts/build/package-headless.sh
```

Contains:

- `yggdrasil-daemon`
- `yggctl`
- `web/` (served by the daemon)

```bash
./yggdrasil-daemon
# open http://127.0.0.1:7331
```

Optional LAN bind:

```bash
YGGDRASIL_API_HOST=0.0.0.0 ./yggdrasil-daemon
```

## Desktop (GTK / WebKit)

Built by CI on Ubuntu with Wails:

```bash
VERSION=0.1.0 ./scripts/build/package-desktop.sh
```

Requires `libgtk-3-dev` and `libwebkit2gtk-4.1-dev`.

On Ubuntu 24.04+, the script passes `-tags webkit2_41` so Wails links
against webkit2gtk 4.1 (4.0 is unavailable).

## Icons (hicolor)

Canonical icons live under `assets/brand/generated/linux/<SIZE>x<SIZE>/apps/yggdrasil.png`
(regenerate with `pnpm icons`). Desktop entry: `packaging/linux/yggdrasil.desktop`
with `Icon=yggdrasil` (no absolute path).

Install layout for packages / AppImage:

```text
/usr/share/icons/hicolor/<SIZE>x<SIZE>/apps/yggdrasil.png
/usr/share/applications/yggdrasil.desktop
```

## Future

- `.deb` / `.rpm` / AppImage wrappers around the headless or desktop layouts
- systemd user unit for headless daemon mode
