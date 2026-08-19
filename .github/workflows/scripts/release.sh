#!/bin/sh

# Release only after the exact executable has passed the Arch VM test.
#
# Only one release is retained. The tag itself is disposable; consumers
# always use GitHub's /releases/latest/download/arch-setup endpoint.
set -eu

TAG="build-${GITHUB_SHA}"

# Find the currently published latest release, if one exists.
OLD_TAG="$(
  gh release view \
    --json tagName \
    --jq '.tagName' \
    2>/dev/null || true
)"

# Remove the old release and its tag.
if [ -n "$OLD_TAG" ]; then
  gh release delete "$OLD_TAG" \
    --cleanup-tag \
    --yes
fi

# Publish the already-tested binary as the new latest release.
gh release create "$TAG" \
  ./bin/arch-setup \
  --title latest \
  --latest
