#!/usr/bin/env bash
# STEP_NAME is read by install.sh after this file is sourced.
# shellcheck disable=SC2034

STEP_NAME='Installing Nerd Fonts'

step_run() {
  run sh -c 'curl -fsSL https://raw.githubusercontent.com/getnf/getnf/main/install.sh | bash'
  run "$HOME/.local/bin/getnf" -i "$(
    IFS=,
    printf '%s' "${GETNF_FONTS[*]}"
  )"
  run fc-cache -f
}
