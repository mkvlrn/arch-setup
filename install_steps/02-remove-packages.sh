#!/usr/bin/env bash
# STEP_NAME is read by install.sh after this file is sourced.
# shellcheck disable=SC2034

STEP_NAME='Removing unwanted packages'

step_run() {
  for package in "${REMOVE_PACKAGES[@]}"; do
    if pacman -Q "$package" >/dev/null 2>&1; then
      run sudo pacman -Rns --noconfirm "$package"
    fi
  done
}
