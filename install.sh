#!/usr/bin/env bash

set -Eeuo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
SETUP_REPO_DIR=${SETUP_REPO_DIR:-"$HOME/repos/arch-setup"}
CI_MODE=false

usage() {
  printf 'Usage: %s [--ci]\n' "${BASH_SOURCE[0]}"
}

while (($#)); do
  case $1 in
  --ci) CI_MODE=true ;;
  -h | --help)
    usage
    exit 0
    ;;
  *)
    printf 'Unknown option: %s\n' "$1" >&2
    usage >&2
    exit 2
    ;;
  esac
  shift
done

export SETUP_REPO_DIR
# shellcheck source=config.sh
source "$SCRIPT_DIR/config.sh"

if [[ $CI_MODE == true ]]; then
  total_steps=9
else
  total_steps=10
fi
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
  fi
  status=$?
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

if [[ $CI_MODE == false ]]; then
  log 'Removing unwanted packages'
  run sudo pacman -Rns --noconfirm "${REMOVE_PACKAGES[@]}"

  log 'Cloning setup repository'
  if [[ -e "$SETUP_REPO_DIR" ]]; then
    printf 'Repository path already exists: %s\n' "$SETUP_REPO_DIR" >&2
    exit 1
  fi
  run git clone "$REPO_HTTP" "$SETUP_REPO_DIR"
  expected_revision=$(git -C "$SCRIPT_DIR" rev-parse HEAD)
  run git -C "$SETUP_REPO_DIR" checkout -B main "$expected_revision"
  run git -C "$SETUP_REPO_DIR" branch --set-upstream-to=origin/main main
  run git -C "$SETUP_REPO_DIR" remote set-url origin "$REPO_SSH"
else
  log 'Validating CI repository'
  [[ -d "$SETUP_REPO_DIR/.git" ]] || {
    printf 'Missing repository: %s\n' "$SETUP_REPO_DIR" >&2
    exit 1
  }
fi

log 'Stowing system files'
run sudo rm -f /etc/pacman.conf /etc/makepkg.conf
run sudo stow -R --no-folding -d "$SETUP_REPO_DIR/stow" -t / system

log 'Installing yay and updating mirrors'
yay_dir=${TMPDIR:-/tmp}/yay-bin
rm -rf "$yay_dir"
run git clone https://aur.archlinux.org/yay-bin "$yay_dir"
run makepkg -si --noconfirm -C -D "$yay_dir"
run yay -Y --gendb
run yay -Y --devel --save
run sudo reflector --latest 20 --protocol https --sort rate --save "$MIRROR_LIST"
run yay -Syu --noconfirm
printf 'Removing debug packages, if any\n'
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
run sudo usermod -d "$HOME/torrents" ftp
run chmod o+x "$HOME"
run sudo systemctl enable --now docker.socket
run sudo systemctl enable --now pure-ftpd.service
run sudo systemctl enable --now paccache.timer
completion_dir="$HOME/.config/fish/completions"
run mkdir -p "$completion_dir"
run "$HOME/.local/bin/mise" completion fish >"$completion_dir/mise.fish"
run "$HOME/.local/share/mise/shims/gh" completion -s fish >"$completion_dir/gh.fish"
run "$HOME/.local/share/mise/shims/glab" completion -s fish >"$completion_dir/glab.fish"

printf '\nInstallation complete. Run %q --verify to verify the machine.\n' "$SCRIPT_DIR/verify.sh"
