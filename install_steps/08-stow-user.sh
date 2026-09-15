#!/usr/bin/env bash
# STEP_NAME is read by install.sh after this file is sourced.
# shellcheck disable=SC2034

STEP_NAME='Stowing user files'

step_run() {
  run stow -R --no-folding --adopt -d "$SETUP_REPO_DIR/stow" -t "$HOME" user
  run git -C "$SETUP_REPO_DIR" restore .
  run git -C "$SETUP_REPO_DIR" clean -fd
}
