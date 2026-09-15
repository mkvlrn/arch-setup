#!/bin/sh

set -u

if [ -f ubuntu-qemu.pid ]; then
  kill "$(cat ubuntu-qemu.pid)" 2>/dev/null || true
fi
