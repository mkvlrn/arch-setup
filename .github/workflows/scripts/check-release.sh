#!/bin/sh

# Determine whether the current main push was produced by merging a pull request.
set -eu

if [ "$GITHUB_EVENT_NAME" != "push" ]; then
  echo "release=false" >>"$GITHUB_OUTPUT"
  exit 0
fi

if gh api \
  "repos/${GITHUB_REPOSITORY}/commits/${GITHUB_SHA}/pulls" |
  jq -e \
    --arg sha "$GITHUB_SHA" \
    'any(.[]; .merged_at != null and .base.ref == "main" and .merge_commit_sha == $sha)' \
    >/dev/null; then
  echo "release=true" >>"$GITHUB_OUTPUT"
else
  echo "release=false" >>"$GITHUB_OUTPUT"
fi
