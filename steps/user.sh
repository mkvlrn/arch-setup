#!/usr/bin/env bash
# Globals are supplied by the main script and its native configuration.
# shellcheck disable=SC2154

configure_user() {
  local completions="$home_dir/.config/fish/completions"
  local mise_path="$home_dir/.local/bin/mise"
  local gh_path="$home_dir/.local/share/mise/shims/gh"
  local glab_path="$home_dir/.local/share/mise/shims/glab"

  run_command 'set user shell' sudo chsh -s /usr/bin/fish "$username" || return
  run_command 'add user to docker group' sudo usermod -aG docker "$username" || return
  run_command 'set anonymous ftp user root dir' sudo usermod -d "$home_dir/torrents" ftp || return
  run_command 'allow ftp user to traverse to download dir' chmod o+x "$home_dir" || return
  run_command 'start docker service' sudo systemctl enable --now docker.socket || return
  run_command 'start pure-ftpd service' sudo systemctl enable --now pure-ftpd.service || return
  run_command 'start paccache service' sudo systemctl enable --now paccache.timer || return
  run_command 'create completions directory' mkdir -p "$completions" || return
  # Keep redirections inside the reported command so open/write failures are reported too.
  run_command 'generate mise completions' bash -c "\"\$1\" completion fish > \"\$2\" || exit" bash "$mise_path" "$completions/mise.fish" || return
  run_command 'generate gh completions' bash -c "\"\$1\" completion -s fish > \"\$2\" || exit" bash "$gh_path" "$completions/gh.fish" || return
  run_command 'generate glab completions' bash -c "\"\$1\" completion -s fish > \"\$2\" || exit" bash "$glab_path" "$completions/glab.fish" || return
}
