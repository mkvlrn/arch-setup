#!/usr/bin/env bash
#MISE description="Format Go code"

set -euo pipefail

mise exec -- golangci-lint fmt ./... "$@"
