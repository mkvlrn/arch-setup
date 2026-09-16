#!/usr/bin/env bash
# STEP_NAME is read by install.sh after this file is sourced.
# shellcheck disable=SC2034

STEP_NAME='Installing mise and managed tools'

step_run() {
  run sh -c 'curl https://mise.run | sh'
  GOPATH="$HOME/.go" run "$HOME/.local/bin/mise" install
}
