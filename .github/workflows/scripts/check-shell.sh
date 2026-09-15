#!/bin/sh

# Run static checks for the shell installer and verifier.
set -eu

shellcheck \
  ./index.html \
  ./arch/*.sh \
  ./ubuntu/*.sh \
  ./secrets/*.sh \
  ./.github/workflows/scripts/*.sh

bash -n \
  ./index.html \
  ./arch/*.sh \
  ./ubuntu/*.sh \
  ./secrets/*.sh \
  ./.github/workflows/scripts/*.sh

git diff --check
