#!/bin/sh
set -e
if ! getent passwd yggdrasil >/dev/null 2>&1; then
  useradd --system --home-dir /var/lib/yggdrasil --create-home --shell /usr/sbin/nologin yggdrasil
fi
if command -v systemctl >/dev/null 2>&1; then
  systemctl daemon-reload || true
  systemctl enable yggdrasil.service || true
  systemctl restart yggdrasil.service || true
fi
