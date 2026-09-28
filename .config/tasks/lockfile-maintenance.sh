#!/usr/bin/env bash
#MISE description="Regenerate Go dependency files"

set -euo pipefail

mise exec -- go mod tidy

git add go.mod go.sum

if git diff --cached --quiet -- go.mod go.sum; then
  echo "Dependency files unchanged; nothing to commit."
  exit 0
fi

git commit -m "chore(deps): update Go dependency files"
