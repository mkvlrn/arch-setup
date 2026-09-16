#!/usr/bin/env bash
# STEP_NAME is read by install.sh after this file is sourced.
# shellcheck disable=SC2034

STEP_NAME='Stowing system files'

step_run() {
  run sudo rm -f /etc/pacman.conf /etc/makepkg.conf
  run sudo stow -R --no-folding -d "$SETUP_REPO_DIR/stow" -t / system
}
