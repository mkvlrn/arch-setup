#!/usr/bin/env bash
#MISE description="Build the Go application binary"

set -euo pipefail

revision=$(git rev-parse --verify HEAD)
if [ -n "$(git --no-optional-locks status --porcelain --untracked-files=normal)" ]; then
  printf 'Cannot build a revision-aligned installer from a dirty checkout.\n' >&2
  exit 1
fi
mkdir -p ./bin
go build -ldflags "-X github.com/mkvlrn/arch-setup/internal/revision.Commit=$revision" -o ./bin/arch-setup ./cmd/arch-setup
