#!/usr/bin/env bash
#MISE description="Run Go tests with coverage"

set -euo pipefail

mise exec -- go test ./... -cover "$@"
