#!/bin/sh
# List profiles exposed as OpenAI models.
# Requires a running daemon. See ../README.md.
set -eu

base="${YGG_BASE_URL:-http://127.0.0.1:7331}"

if [ -n "${YGG_API_KEY:-}" ]; then
  exec curl -fsS "$base/v1/models" -H "Authorization: Bearer $YGG_API_KEY"
fi

exec curl -fsS "$base/v1/models"
