#!/usr/bin/env bash
# Globals are supplied by the main script and its native configuration.
# shellcheck disable=SC2154

configure_xdg() {
  run_command 'update xdg dirs' xdg-user-dirs-update || return
  (
    cd -- "$home_dir" || return
    run_command 'create new xdg set' mkdir -p "${xdg_mkdir[@]}" || return
    run_command 'remove old xdg set' rm -rf "${xdg_rmrf[@]}" || return
  ) || return
}
