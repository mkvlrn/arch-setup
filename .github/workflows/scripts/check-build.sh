#!/bin/sh

# Run all static checks, unit tests, and build the executable that will
# later be tested inside a real Arch VM.
set -eu

bun install --frozen-lockfile

mise typecheck
mise lint
mise test-ci
mise build
