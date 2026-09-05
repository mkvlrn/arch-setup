#!/bin/sh

# Run static checks and tests before testing the installer in the VM.
set -eu

make lint
make test
