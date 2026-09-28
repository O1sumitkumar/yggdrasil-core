#!/bin/sh
# Send one user message to /v1/chat/completions.
# Requires a running daemon and an installed model. See ../README.md.
set -eu

base="${YGG_BASE_URL:-http://127.0.0.1:7331}"
model="${YGG_MODEL:-profile:general-assistant}"

body=$(printf '%s' "{\"model\":\"$model\",\"messages\":[{\"role\":\"user\",\"content\":\"Hello\"}]}")

if [ -n "${YGG_API_KEY:-}" ]; then
  exec curl -fsS "$base/v1/chat/completions" \
    -H "Authorization: Bearer $YGG_API_KEY" \
    -H "Content-Type: application/json" \
    -d "$body"
fi

exec curl -fsS "$base/v1/chat/completions" \
  -H "Content-Type: application/json" \
  -d "$body"
