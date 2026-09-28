#!/usr/bin/env bash
#MISE description="Sync Go dependencies and tools"

set -euo pipefail

mise install
mise prune -y
mise exec -- go mod tidy
