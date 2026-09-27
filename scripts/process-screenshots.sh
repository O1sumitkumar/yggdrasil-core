#!/usr/bin/env bash
# Letterbox README screenshots onto the canvases listed in screenshots/config.json.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

CONFIG="${SCREENSHOT_CONFIG:-$ROOT/screenshots/config.json}"
RAW_DIR="$ROOT/docs/screenshots"
OUT_DIR="$ROOT/screenshots/appstore"

if [[ ! -f "$CONFIG" ]]; then
  echo "Missing screenshot config: $CONFIG" >&2
  exit 1
fi

if command -v magick >/dev/null 2>&1; then
  MAGICK=magick
elif command -v convert >/dev/null 2>&1; then
  MAGICK=convert
else
  echo "ImageMagick is required." >&2
  exit 1
fi

python3 - "$CONFIG" "$RAW_DIR" "$OUT_DIR" "$MAGICK" <<'PY'
import json, os, subprocess, sys
config_path, raw_dir, out_dir, magick = sys.argv[1:5]
cfg = json.load(open(config_path))
background = cfg.get("background", "#121C26")

def convert(source, dest, width, height):
    if magick == "magick":
        cmd = [
            "magick", source,
            "-resize", f"{width}x{height}",
            "-background", background,
            "-gravity", "center",
            "-extent", f"{width}x{height}",
            dest,
        ]
    else:
        cmd = [
            "convert", source,
            "-resize", f"{width}x{height}",
            "-background", background,
            "-gravity", "center",
            "-extent", f"{width}x{height}",
            dest,
        ]
    subprocess.check_call(cmd)

for size in cfg["dimensions"]:
    width, height = (int(part) for part in size.split("x"))
    dest_dir = os.path.join(out_dir, size)
    os.makedirs(dest_dir, exist_ok=True)
    for screen in cfg["screens"]:
        source = os.path.join(raw_dir, screen["filename"])
        if not os.path.isfile(source):
            raise SystemExit(f"Missing screenshot: {source}")
        dest = os.path.join(dest_dir, screen["filename"])
        convert(source, dest, width, height)
        print(dest)
PY
