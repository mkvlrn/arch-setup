#!/usr/bin/env bash

set -Eeuo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
SETUP_REPO_DIR=${SETUP_REPO_DIR:-"$HOME/repos/arch-setup"}

if (($#)); then
  case $1 in
  -h | --help)
    printf 'Usage: %s\n' "${BASH_SOURCE[0]}"
    exit 0
    ;;
  *)
    printf 'This script does not accept arguments.\n' >&2
    exit 2
    ;;
  esac
fi

export SETUP_REPO_DIR
# shellcheck source=config.sh
source "$SCRIPT_DIR/config.sh"

export_distrobox_apps() {
  command -v distrobox-export >/dev/null 2>&1 || return 0

  local package
  while IFS= read -r package; do
    [[ -n $package ]] || continue
    run distrobox-export --app "$package"
  done < <(
    awk '
      /^MAIN_PACKAGES=\(/ { in_packages = 1; next }
      in_packages && /^\)/ { exit }
      in_packages && /# export in distrobox/ { print $1 }
    ' "$SCRIPT_DIR/config.sh"
  )
}

shopt -s nullglob
step_files=("$SCRIPT_DIR/install_steps"/[0-9][0-9]-*.sh)
if ((${#step_files[@]} == 0)); then
  printf 'No installer steps found in %s/install_steps\n' "$SCRIPT_DIR" >&2
  exit 1
fi

total_steps=${#step_files[@]}
step_number=0

run() {
  local output status
  output=$(mktemp)
  if "$@" >"$output" 2>&1; then
    rm -f "$output"
    return 0
  else
    status=$?
  fi
  printf 'Command failed: '
  printf '%q ' "$@"
  printf '\n'
  cat "$output" >&2
  rm -f "$output"
  return "$status"
}

if [[ $EUID -eq 0 ]]; then
  printf 'Do not run this script as root.\n' >&2
  exit 1
fi

for step_file in "${step_files[@]}"; do
  unset STEP_NAME
  unset -f step_run 2>/dev/null || true

  # shellcheck source=/dev/null
  source "$step_file"

  [[ -n ${STEP_NAME:-} ]] || {
    printf 'Missing STEP_NAME in %s\n' "$step_file" >&2
    exit 1
  }
  declare -F step_run >/dev/null || {
    printf 'Missing step_run function in %s\n' "$step_file" >&2
    exit 1
  }

  step_number=$((step_number + 1))
  printf '[%d/%d] %s\n' "$step_number" "$total_steps" "$STEP_NAME"
  step_run

  if [[ $step_file == *'/06-main-packages.sh' ]]; then
    export_distrobox_apps
  fi
done
