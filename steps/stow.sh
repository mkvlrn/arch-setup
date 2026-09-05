#!/usr/bin/env bash
# Globals are supplied by the main script and its native configuration.
# shellcheck disable=SC2154

stow_system() {
  run_command 'remove existing confs' sudo rm -f /etc/pacman.conf /etc/makepkg.conf || return
  run_command 'stow system files' sudo stow -R --no-folding -d "$repo_dir/stow" -t / system || return
}

stow_user() {
  run_command 'stow user files' stow -R --no-folding --adopt -d "$repo_dir/stow" -t "$home_dir" user || return
  (
    cd -- "$repo_dir" || return
    # Adoption only changes this package; leave unrelated checkout work alone.
    run_command 'restore adopted files' git restore -- stow/user || return
    run_command 'clean adopted files' git clean -fd -- stow/user || return
  ) || return
}
