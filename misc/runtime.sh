#!/usr/bin/env bash
# Commands retain their argument boundaries and stream diagnostics, including
# interactive prompts. Never eval a command assembled from configuration.
run_command() {
  local description=$1 status
  shift
  if "$@"; then
    return 0
  else
    status=$?
    printf 'Could not run %s (exit %d).\n' "$description" "$status" >&2
    return "$status"
  fi
}

sudo_pid=
start_sudo() {
  run_command 'validate sudo credentials' sudo -v || return
  (
    # Reap the sleeper too, so exiting setup leaves no refresh subprocesses.
    sleeper=
    trap 'if [[ -n $sleeper ]]; then kill "$sleeper" 2>/dev/null || :; wait "$sleeper" 2>/dev/null || :; fi; exit 0' TERM INT
    while :; do
      sleep 30 &
      sleeper=$!
      wait "$sleeper" || exit
      sleeper=
      sudo -n -v || exit
    done
  ) &
  sudo_pid=$!
}

stop_sudo() {
  if [[ -n $sudo_pid ]]; then
    kill "$sudo_pid" 2>/dev/null || :
    wait "$sudo_pid" 2>/dev/null || :
    sudo_pid=
  fi
}

# Explicit status handling matters: Bash disables errexit throughout functions
# invoked in conditionals. Each step also propagates individual command failures.
run_plan() {
  local mode=$1 index failures=0
  shift
  local -a plan=("$@")
  for ((index = 0; index < ${#plan[@]}; index++)); do
    printf '[%d/%d] %s\n' "$((index + 1))" "${#plan[@]}" "${plan[index]//_/ }"
    if "${plan[index]}"; then
      continue
    fi
    printf '%s failed: %s\n' "$mode" "${plan[index]}" >&2
    if [[ $mode == Setup ]]; then
      return 1
    fi
    failures=$((failures + 1))
  done
  if ((failures)); then
    printf '%d verification check(s) failed.\n' "$failures" >&2
    return 1
  fi
  if [[ $mode == Setup ]]; then
    printf 'Done.\n'
  else
    printf 'Verification passed.\n'
  fi
}
