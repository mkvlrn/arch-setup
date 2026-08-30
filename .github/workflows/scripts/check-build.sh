#!/bin/sh

# Run all static checks, unit tests, and build the executable for downstream jobs.
set -eu

make lint
make test
make build
