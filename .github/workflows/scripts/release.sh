#!/bin/sh

# Replace the current release with the executable built after merging a PR.
set -eu

TAG="build-${GITHUB_SHA}"

OLD_TAG="$(
  gh release view \
    --json tagName \
    --jq '.tagName' \
    2>/dev/null || true
)"

if [ -n "$OLD_TAG" ]; then
  gh release delete "$OLD_TAG" \
    --cleanup-tag \
    --yes
fi

gh release create "$TAG" \
  ./bin/arch-setup \
  --title latest \
  --latest
