#!/bin/sh

# Stop QEMU if it was successfully started.
#
# This is mostly hygiene on GitHub-hosted runners, but also ensures QEMU
# is terminated before the job finishes when a later test step fails.
set -u

if [ -f qemu.pid ]; then
  kill "$(cat qemu.pid)" 2>/dev/null || true
fi
