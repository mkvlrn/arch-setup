#!/bin/sh

# Run static checks for the shell installer and verifier.
set -eu

shellcheck \
  ./config.sh \
  ./install.sh \
  ./verify.sh \
  ./.github/workflows/scripts/*.sh

bash -n \
  ./config.sh \
  ./install.sh \
  ./verify.sh \
  ./.github/workflows/scripts/*.sh

git diff --check
