#!/usr/bin/env bash
# STEP_NAME is read by install.sh after this file is sourced.
# shellcheck disable=SC2034

STEP_NAME='Installing base packages'

step_run() {
  run sudo pacman -Syu --noconfirm --needed "${BASE_PACKAGES[@]}"
}
