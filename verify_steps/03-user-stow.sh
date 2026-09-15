#!/usr/bin/env bash
# CHECK_NAME is read by verify.sh after this file is sourced.
# shellcheck disable=SC2034

CHECK_NAME='user Stow links'

check_run() {
  check_stow_tree user "$HOME"
}
