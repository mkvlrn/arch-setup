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
# shellcheck source=ubuntu/config.sh
source "$SCRIPT_DIR/config.sh"

total_steps=9
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

log 'Installing Ubuntu packages'
run sudo apt-get update
run sudo apt-get install -y --no-install-recommends "${APT_PACKAGES[@]}"

log 'Validating existing repository'
[[ -d "$SETUP_REPO_DIR/.git" ]] || {
  printf 'Missing repository: %s\n' "$SETUP_REPO_DIR" >&2
  exit 1
}
run git -C "$SETUP_REPO_DIR" remote set-url origin "$REPO_SSH"

log 'Deploying Ubuntu system files'
run sudo stow -R --no-folding -d "$SETUP_REPO_DIR/stow" -t / system_ubuntu

log 'Configuring Flatpak applications'
run sudo flatpak remote-add --if-not-exists --system flathub https://dl.flathub.org/repo/flathub.flatpakrepo
run sudo flatpak install --system -y flathub "${FLATPAKS[@]}"

log 'Updating XDG directories'
run xdg-user-dirs-update
run mkdir -p -- "${XDG_MKDIR[@]/#/$HOME/}"
run rm -rf -- "${XDG_RMRF[@]/#/$HOME/}"

log 'Stowing user files'
run stow -R --no-folding --adopt -d "$SETUP_REPO_DIR/stow" -t "$HOME" user
run git -C "$SETUP_REPO_DIR" restore .
run git -C "$SETUP_REPO_DIR" clean -fd

log 'Installing Nerd Fonts'
run sh -c 'curl -fsSL https://raw.githubusercontent.com/getnf/getnf/main/install.sh | bash'
run "$HOME/.local/bin/getnf" -i "${GETNF_FONTS// /,}"
run fc-cache -f

log 'Installing mise and managed tools'
if [[ -z ${MISE_GITHUB_TOKEN:-} ]]; then
  printf 'MISE_GITHUB_TOKEN is required to install mise-managed tools.\n' >&2
  exit 1
fi
run sh -c 'curl https://mise.run | sh'
GOPATH="$HOME/.go" run "$HOME/.local/bin/mise" install

log 'Configuring user settings'
run sudo chsh -s /usr/bin/fish "$USER"
run sudo usermod -aG docker "$USER"
run sudo systemctl enable --now docker.service
completion_dir="$HOME/.config/fish/completions"
run mkdir -p "$completion_dir"
run "$HOME/.local/bin/mise" completion fish >"$completion_dir/mise.fish"
run "$HOME/.local/share/mise/shims/gh" completion -s fish >"$completion_dir/gh.fish"
