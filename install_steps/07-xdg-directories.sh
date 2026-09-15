#!/usr/bin/env bash
# STEP_NAME is read by install.sh after this file is sourced.
# shellcheck disable=SC2034

STEP_NAME='Updating XDG directories'

step_run() {
  run xdg-user-dirs-update
  run mkdir -p -- "${XDG_MKDIR[@]/#/$HOME/}"
  run rm -rf -- "${XDG_RMRF[@]/#/$HOME/}"
}
