#!/usr/bin/env bash
# STEP_NAME is read by install.sh after this file is sourced.
# shellcheck disable=SC2034

STEP_NAME='Installing main packages'

step_run() {
  run yay -S --noconfirm --needed "${MAIN_PACKAGES[@]}"
}
