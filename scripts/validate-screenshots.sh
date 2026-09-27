#!/usr/bin/env bash
# Check that screenshots/appstore contains a PNG of each configured screen and size.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

CONFIG="${SCREENSHOT_CONFIG:-$ROOT/screenshots/config.json}"
OUT_DIR="$ROOT/screenshots/appstore"

if command -v magick >/dev/null 2>&1; then
  MAGICK=magick
elif command -v identify >/dev/null 2>&1; then
  MAGICK=identify
else
  echo "ImageMagick is required to validate dimensions." >&2
  exit 1
fi

python3 - "$CONFIG" "$OUT_DIR" "$MAGICK" <<'PY'
import json, os, subprocess, sys
config_path, out_dir, magick = sys.argv[1:4]
cfg = json.load(open(config_path))
errors = []
previous = []

def identify(path):
    cmd = [magick, "identify", "-format", "%w %h", path] if magick == "magick" else ["identify", "-format", "%w %h", path]
    return subprocess.check_output(cmd, text=True).split()

for size in cfg["dimensions"]:
    want_w, want_h = (int(part) for part in size.split("x"))
    names = []
    for screen in cfg["screens"]:
        filename = screen["filename"]
        names.append(filename)
        path = os.path.join(out_dir, size, filename)
        if not os.path.isfile(path):
            errors.append(f"Missing screenshot: {path}")
            continue
        with open(path, "rb") as handle:
            header = handle.read(8)
        if header != b"\x89PNG\r\n\x1a\n":
            errors.append(f"Not a PNG: {path}")
        if os.path.getsize(path) < 10000:
            errors.append(f"Screenshot looks empty: {path}")
            continue
        got_w, got_h = (int(part) for part in identify(path))
        if got_w != want_w or got_h != want_h:
            errors.append(f"{path} is {got_w}x{got_h}, expected {want_w}x{want_h}")
    if names != sorted(names):
        errors.append(f"Filenames for {size} are not in numeric order: {', '.join(names)}")
    if previous and previous != names:
        errors.append(f"Filenames for {size} do not match the previous canvas")
    previous = names
if errors:
    print("\n".join(errors), file=sys.stderr)
    raise SystemExit("Screenshot validation failed")
print(f"Screenshot validation passed: {out_dir}")
PY
