#!/usr/bin/env bash

set -Eeuo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
SETUP_REPO_DIR=${SETUP_REPO_DIR:-"$HOME/repos/arch-setup"}

export SETUP_REPO_DIR
# shellcheck source=config.sh
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
  [[ -d "$SETUP_REPO_DIR/.git" ]] || {
    printf 'missing .git\n' >&2
    return 1
  }
  local remote head status expected
  remote=$(git -C "$SETUP_REPO_DIR" remote get-url origin)
  head=$(git -C "$SETUP_REPO_DIR" rev-parse HEAD)
  status=$(git -C "$SETUP_REPO_DIR" status --porcelain)
  expected=${ARCH_SETUP_EXPECTED_REVISION:-$head}
  [[ $remote == "$REPO_SSH" ]] || {
    printf 'expected remote %s, got %s\n' "$REPO_SSH" "$remote" >&2
    return 1
  }
  [[ $head == "$expected" ]] || {
    printf 'expected HEAD %s, got %s\n' "$expected" "$head" >&2
    return 1
  }
  [[ -z $status ]] || {
    printf 'repository is dirty:\n%s\n' "$status" >&2
    return 1
  }
}

check_stow_tree() {
  local package=$1 target=$2 source relative destination source_real destination_real
  while IFS= read -r -d '' source; do
    case ${source##*/} in .gitkeep | .stow-local-ignore) continue ;; esac
    relative=${source#"$SETUP_REPO_DIR/stow/$package/"}
    destination="$target/$relative"
    [[ -L "$destination" ]] || {
      printf '%s is not a symlink\n' "$destination" >&2
      return 1
    }
    source_real=$(realpath "$source")
    destination_real=$(realpath "$destination")
    [[ $source_real == "$destination_real" ]] || {
      printf '%s points to %s instead of %s\n' "$destination" "$destination_real" "$source_real" >&2
      return 1
    }
  done < <(find "$SETUP_REPO_DIR/stow/$package" -type f -print0)
}

check_packages() {
  local package
  for package in "${BASE_PACKAGES[@]}" "${MAIN_PACKAGES[@]}"; do
    pacman -Q "$package" >/dev/null || {
      printf 'package not installed: %s\n' "$package" >&2
      return 1
    }
  done
}

check_removed_packages() {
  local package
  for package in "${REMOVE_PACKAGES[@]}"; do
    if pacman -Q "$package" >/dev/null 2>&1; then
      printf 'package still installed: %s\n' "$package" >&2
      return 1
    fi
  done
}

check_yay() {
  yay --version >/dev/null
  local debug_packages
  debug_packages=$(yay -Qq | grep -- '-debug$' || true)
  [[ -z $debug_packages ]] || {
    printf 'debug packages installed:\n%s\n' "$debug_packages" >&2
    return 1
  }
  grep -Fq "$MIRROR_LIST_CHECK" "$MIRROR_LIST"
}

check_xdg() {
  local directory
  for directory in "${XDG_MKDIR[@]}"; do
    [[ -d "$HOME/$directory" ]] || {
      printf 'missing directory: %s\n' "$HOME/$directory" >&2
      return 1
    }
  done
  for directory in "${XDG_RMRF[@]}"; do
    [[ ! -e "$HOME/$directory" ]] || {
      printf 'path still exists: %s\n' "$HOME/$directory" >&2
      return 1
    }
  done
}

check_mise() {
  "$HOME/.local/bin/mise" --version >/dev/null
  [[ -z $("$HOME/.local/bin/mise" ls --global --missing --no-header) ]]
}

check_user() {
  local shell ftp_home groups completion unit
  shell=$(getent passwd "$USER" | cut -d: -f7)
  [[ $shell == /usr/bin/fish ]] || {
    printf 'shell is %s\n' "$shell" >&2
    return 1
  }
  groups=$(id -nG "$USER")
  [[ " $groups " == *' docker '* ]] || {
    printf 'user is not in docker group\n' >&2
    return 1
  }
  ftp_home=$(getent passwd ftp | cut -d: -f6)
  [[ $ftp_home == "$HOME/torrents" ]] || {
    printf 'ftp home is %s\n' "$ftp_home" >&2
    return 1
  }
  [[ $(stat -c '%a' "$HOME") =~ [13579]$ ]] || {
    printf 'home is not traversable\n' >&2
    return 1
  }
  for unit in docker.socket pure-ftpd.service paccache.timer; do
    systemctl is-enabled --quiet "$unit" || return 1
    systemctl is-active --quiet "$unit" || return 1
  done
  for completion in mise gh; do
    [[ -f "$HOME/.config/fish/completions/$completion.fish" ]] || {
      printf 'missing completion: %s\n' "$completion" >&2
      return 1
    }
  done
}

printf 'Verifying repository\n'
check 'repository state' check_repository
check 'system Stow links' check_stow_tree system /
check 'user Stow links' check_stow_tree user "$HOME"
check 'yay and mirror list' check_yay
check 'installed packages' check_packages
check 'removed packages' check_removed_packages
check 'XDG directories' check_xdg
check 'mise and managed tools' check_mise
check 'user settings and services' check_user

if ((failures)); then
  printf '\n%d verification check(s) failed.\n' "$failures" >&2
  exit 1
fi
printf '\nAll verification checks passed.\n'
