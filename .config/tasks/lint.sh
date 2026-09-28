#!/usr/bin/env bash
#MISE description="Run the Go linter"

set -euo pipefail

mise exec -- golangci-lint run ./... "$@"
