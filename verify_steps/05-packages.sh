#!/usr/bin/env bash
# CHECK_NAME is read by verify.sh after this file is sourced.
# shellcheck disable=SC2034

CHECK_NAME='installed packages'

check_run() {
  check_packages
}
