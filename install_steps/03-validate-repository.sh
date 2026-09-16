#!/usr/bin/env bash
# STEP_NAME is read by install.sh after this file is sourced.
# shellcheck disable=SC2034

STEP_NAME='Validating existing repository'

step_run() {
  [[ -d "$SETUP_REPO_DIR/.git" ]] || {
    printf 'Missing repository: %s\n' "$SETUP_REPO_DIR" >&2
    return 1
  }
  run git -C "$SETUP_REPO_DIR" remote set-url origin "$REPO_SSH"
}
