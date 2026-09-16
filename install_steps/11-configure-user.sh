#!/usr/bin/env bash
# STEP_NAME is read by install.sh after this file is sourced.
# shellcheck disable=SC2034

STEP_NAME='Configuring user settings'

step_run() {
  run sudo chsh -s /usr/bin/fish "$USER"
  run sudo usermod -aG docker "$USER"
  run chmod o+x "$HOME"

  run sudo systemctl enable --now docker.socket
  run sudo systemctl enable --now paccache.timer
  run systemctl --user enable --now zen-notification-sound
  run systemctl --user enable --now ssh-agent.socket

  run ssh-add "$HOME/.ssh/dev"
  run ssh-add "$HOME/.ssh/cb"

  local completion_dir="$HOME/.config/fish/completions"
  run mkdir -p "$completion_dir"
  run "$HOME/.local/bin/mise" completion fish >"$completion_dir/mise.fish"
  run "$HOME/.local/share/mise/shims/gh" completion -s fish >"$completion_dir/gh.fish"
}
