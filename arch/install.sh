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
# shellcheck source=arch/config.sh
source "$SCRIPT_DIR/config.sh"

total_steps=10
step_number=0

log() {
  step_number=$((step_number + 1))
  printf '[%d/%d] %s\n' "$step_number" "$total_steps" "$1"
}

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

log 'Installing base packages'
run sudo pacman -Syu --noconfirm --needed "${BASE_PACKAGES[@]}"

log 'Removing unwanted packages'
for package in "${REMOVE_PACKAGES[@]}"; do
  if pacman -Q "$package" >/dev/null 2>&1; then
    run sudo pacman -Rns --noconfirm "$package"
  fi
done

log 'Validating existing repository'
[[ -d "$SETUP_REPO_DIR/.git" ]] || {
  printf 'Missing repository: %s\n' "$SETUP_REPO_DIR" >&2
  exit 1
}
run git -C "$SETUP_REPO_DIR" remote set-url origin "$REPO_SSH"

log 'Stowing system files'
run sudo rm -f /etc/pacman.conf /etc/makepkg.conf
run sudo stow -R --no-folding -d "$SETUP_REPO_DIR/stow" -t / system_arch

log 'Installing yay and updating mirrors'
yay_dir=${TMPDIR:-/tmp}/yay-bin
rm -rf "$yay_dir"
run git clone https://aur.archlinux.org/yay-bin "$yay_dir"
run makepkg -si --noconfirm -C -D "$yay_dir"
run yay -Y --gendb
run yay -Y --devel --save
run sudo reflector --latest 20 --protocol https --sort rate --save "$MIRROR_LIST"
run yay -Syu --noconfirm
debug_packages=$(yay -Qq | grep -- '-debug$' || true)
if [[ -n $debug_packages ]]; then
  printf '%s\n' "$debug_packages" | xargs -r yay -Rnsu
fi

log 'Installing main packages'
run yay -S --noconfirm --needed "${MAIN_PACKAGES[@]}"

log 'Updating XDG directories'
run xdg-user-dirs-update
run mkdir -p -- "${XDG_MKDIR[@]/#/$HOME/}"
run rm -rf -- "${XDG_RMRF[@]/#/$HOME/}"

log 'Stowing user files'
run stow -R --no-folding --adopt -d "$SETUP_REPO_DIR/stow" -t "$HOME" user
run git -C "$SETUP_REPO_DIR" restore .
run git -C "$SETUP_REPO_DIR" clean -fd

log 'Installing mise and managed tools'
run sh -c 'curl https://mise.run | sh'
GOPATH="$HOME/.go" run "$HOME/.local/bin/mise" install

log 'Configuring user settings'
run sudo chsh -s /usr/bin/fish "$USER"
run sudo usermod -aG docker "$USER"

run chmod o+x "$HOME"
run sudo systemctl enable --now docker.socket

run sudo systemctl enable --now paccache.timer
completion_dir="$HOME/.config/fish/completions"
run mkdir -p "$completion_dir"
run "$HOME/.local/bin/mise" completion fish >"$completion_dir/mise.fish"
run "$HOME/.local/share/mise/shims/gh" completion -s fish >"$completion_dir/gh.fish"
