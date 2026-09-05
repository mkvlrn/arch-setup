#!/usr/bin/env bash
# Globals are supplied by the main script and its native configuration.
# shellcheck disable=SC2154

install_mise() {
  run_command 'install mise' bash -o pipefail -c 'curl -f https://mise.run | sh' || return
  run_command 'install tools managed by mise' env GOPATH="$home_dir/.go" "$home_dir/.local/bin/mise" install || return
}
