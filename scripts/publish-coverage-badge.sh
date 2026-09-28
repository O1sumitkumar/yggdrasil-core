#!/usr/bin/env bash
# Turn a Go cover profile into a Shields endpoint and publish it on the coverage branch.
# Printing is the default. Set PUBLISH_COVERAGE_BADGE=1 to push coverage.json.
set -euo pipefail

profile="${1:?coverage profile is required}"
total="$(go tool cover -func="$profile" | awk '/^total:/ { print $3 }')"
if [[ -z "$total" ]]; then
  echo "No coverage total in $profile" >&2
  exit 1
fi

number="${total%\%}"
color="$(awk -v n="$number" 'BEGIN {
  if (n+0 >= 80) print "brightgreen"
  else if (n+0 >= 60) print "green"
  else if (n+0 >= 50) print "yellowgreen"
  else if (n+0 >= 30) print "yellow"
  else print "red"
}')"
json="$(printf '{"schemaVersion":1,"label":"coverage","message":"%s","color":"%s"}\n' "$total" "$color")"

if [[ "${PUBLISH_COVERAGE_BADGE:-}" != "1" ]]; then
  printf '%s' "$json"
  exit 0
fi

: "${GITHUB_TOKEN:?GITHUB_TOKEN is required to publish the coverage badge}"
: "${GITHUB_REPOSITORY:?GITHUB_REPOSITORY is required to publish the coverage badge}"

workdir="$(mktemp -d)"
trap 'rm -rf "$workdir"' EXIT
remote="https://x-access-token:${GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git"
if git ls-remote --exit-code --heads "$remote" coverage >/dev/null 2>&1; then
  git clone --quiet --depth 1 --branch coverage "$remote" "$workdir"
else
  git init -q -b coverage "$workdir"
fi

printf '%s' "$json" > "$workdir/coverage.json"
git -C "$workdir" add coverage.json
if git -C "$workdir" diff --cached --quiet; then
  exit 0
fi
git -C "$workdir" \
  -c user.name='github-actions[bot]' \
  -c user.email='41898282+github-actions[bot]@users.noreply.github.com' \
  commit -q -m "Update the coverage badge to ${total}."
git -C "$workdir" push -q "$remote" HEAD:coverage
