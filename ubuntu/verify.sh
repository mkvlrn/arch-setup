#!/usr/bin/env bash

set -Eeuo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
SETUP_REPO_DIR=${SETUP_REPO_DIR:-"$HOME/repos/arch-setup"}

export SETUP_REPO_DIR
# shellcheck source=ubuntu/config.sh
source "$SCRIPT_DIR/config.sh"

failures=0
check() {
  local name=$1
  shift
  printf ' %-55s' "$name"
  if "$@"; then
    printf ' OK\n'
  else
    printf ' FAIL\n'
    failures=$((failures + 1))
  fi
}

check_repository() {
  [[ -d "$SETUP_REPO_DIR/.git" ]] || return 1
  local remote status
  remote=$(git -C "$SETUP_REPO_DIR" remote get-url origin)
  status=$(git -C "$SETUP_REPO_DIR" status --porcelain)
  [[ $remote == "$REPO_SSH" && -z $status ]]
}

check_stow_tree() {
  local package=$1 target=$2
  local source_root=$3 source relative destination source_real destination_real
  while IFS= read -r -d '' source; do
    case ${source##*/} in .gitkeep | .stow-local-ignore) continue ;; esac
    relative=${source#"$source_root/$package/"}
    destination="$target/$relative"
    [[ -L "$destination" ]] || return 1
    source_real=$(realpath "$source")
    destination_real=$(realpath "$destination")
    [[ $source_real == "$destination_real" ]] || return 1
  done < <(find "$source_root/$package" -type f -print0)
}

check_apt_packages() {
  local package
  for package in "${APT_PACKAGES[@]}"; do
    dpkg-query -W -f='${db:Status-Abbrev}' "$package" 2>/dev/null | grep -q '^ii ' || return 1
  done
}

check_flatpaks() {
  local app
  for app in "${FLATPAKS[@]}"; do
    flatpak info --system "$app" >/dev/null 2>&1 || return 1
  done
}

check_fonts() {
  local font installed_fonts
  installed_fonts=$("$HOME/.local/bin/getnf" -l)
  for font in Hack IosevkaTerm Noto ZedMono; do
    grep -Fq "$font - " <<<"$installed_fonts" || return 1
  done
}

check_xdg() {
  local directory
  for directory in "${XDG_MKDIR[@]}"; do
    [[ -d "$HOME/$directory" ]] || return 1
  done
  for directory in "${XDG_RMRF[@]}"; do
    [[ ! -e "$HOME/$directory" ]] || return 1
  done
}

check_mise() {
  "$HOME/.local/bin/mise" --version >/dev/null
  [[ -z $("$HOME/.local/bin/mise" ls --global --missing --no-header) ]]
}

check_user_settings() {
  local shell groups completion
  shell=$(getent passwd "$USER" | cut -d: -f7)
  [[ $shell == /usr/bin/fish ]] || return 1
  groups=$(id -nG "$USER")
  [[ " $groups " == *' docker '* ]] || return 1
  systemctl is-enabled --quiet docker.service || return 1
  systemctl is-active --quiet docker.service || return 1
  for completion in mise gh; do
    [[ -f "$HOME/.config/fish/completions/$completion.fish" ]] || return 1
  done
}

printf 'Verifying Ubuntu setup\n'
check 'repository state' check_repository
check 'system Stow links' check_stow_tree system_ubuntu / "$SETUP_REPO_DIR/stow"
check 'user Stow links' check_stow_tree user "$HOME" "$SETUP_REPO_DIR/stow"
check 'apt packages' check_apt_packages
check 'Flatpak applications' check_flatpaks
check 'Nerd Fonts' check_fonts
check 'XDG directories' check_xdg
check 'mise and managed tools' check_mise
check 'user settings and Docker' check_user_settings

if ((failures)); then
  printf '\n%d verification check(s) failed.\n' "$failures" >&2
  exit 1
fi
printf '\nAll Ubuntu verification checks passed.\n'
