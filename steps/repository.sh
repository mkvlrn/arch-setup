#!/usr/bin/env bash
# Globals are supplied by the main script and its native configuration.
# shellcheck disable=SC2154

clone_repository() (
  cd -- "$repo_dir" || return
  run_command 'set ssh upstream' git remote set-url origin "$repo_ssh" || return
)
