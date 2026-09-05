#!/usr/bin/env bash
# Ordinary commands are silent on success. Checks that need stdout explicitly
# use capture_command inside command substitution. Never eval command arguments.
run_command() {
  capture_command "$@" >/dev/null
}

capture_command() (
  local description=$1 status logs
  shift
  logs=$(mktemp -d) || return
  # Isolate traps from the runner and remove logs on success, failure, or signals.
  trap 'rm -rf -- "$logs"' EXIT
  trap 'exit 129' HUP
  trap 'exit 130' INT
  trap 'exit 143' TERM
  if "$@" >"$logs/stdout" 2>"$logs/stderr"; then
    cat -- "$logs/stdout"
  else
    status=$?
    printf 'Could not run %s (exit %d).\n' "$description" "$status" >&2
    if [[ -s $logs/stdout ]]; then
      printf '\nstdout:\n' >&2
      cat -- "$logs/stdout" >&2
      printf '\n' >&2
    fi
    if [[ -s $logs/stderr ]]; then
      printf '\nstderr:\n' >&2
      cat -- "$logs/stderr" >&2
      printf '\n' >&2
    fi
    return "$status"
  fi
)

sudo_pid=
start_sudo() {
  # Authentication must remain interactive; subsequent commands are captured.
  sudo -v || return
  (
    # Reap the sleeper too, so exiting setup leaves no refresh subprocesses.
    sleeper=
    trap 'if [[ -n $sleeper ]]; then kill "$sleeper" 2>/dev/null || :; wait "$sleeper" 2>/dev/null || :; fi; exit 0' TERM INT
    while :; do
      sleep 30 &
      sleeper=$!
      wait "$sleeper" || exit
      sleeper=
      run_command 'refresh sudo credentials' sudo -n -v || exit
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

animate_progress() (
  local label=$1 frame sleeper=''
  trap 'if [[ -n $sleeper ]]; then kill "$sleeper" 2>/dev/null || :; wait "$sleeper" 2>/dev/null || :; fi; exit 0' TERM INT HUP
  while :; do
    for frame in '|' '/' '-' $'\\'; do
      printf '\r\033[K%s %s' "$label" "$frame"
      sleep 0.1 &
      sleeper=$!
      wait "$sleeper" || exit
      sleeper=''
    done
  done
)

# Keep the step in the foreground (including its stdin). Only the animation is
# backgrounded. Buffer diagnostics so the spinner cannot overwrite an error.
run_animated_step() (
  local label=$1 step=$2 spinner='' errors status=0
  errors=$(mktemp) || return
  finish_animation() {
    if [[ -n $spinner ]]; then
      kill "$spinner" 2>/dev/null || :
      wait "$spinner" 2>/dev/null || :
      spinner=''
    fi
  }
  # Called by the EXIT trap, including when a signal interrupts the step.
  # shellcheck disable=SC2329
  cleanup_progress() {
    finish_animation
    rm -f -- "$errors"
  }
  trap cleanup_progress EXIT
  trap 'finish_animation; printf "\r\033[K%s interrupted\n" "$label"; exit 129' HUP
  trap 'finish_animation; printf "\r\033[K%s interrupted\n" "$label"; exit 130' INT
  trap 'finish_animation; printf "\r\033[K%s interrupted\n" "$label"; exit 143' TERM
  animate_progress "$label" &
  spinner=$!
  "$step" 2>"$errors" || status=$?
  finish_animation
  if ((status == 0)); then
    printf '\r\033[K%s done\n' "$label"
  else
    printf '\r\033[K%s failed\n' "$label"
  fi
  cat -- "$errors" >&2
  return "$status"
)

# Explicit status handling matters: Bash disables errexit throughout functions
# invoked in conditionals. Each step also propagates individual command failures.
run_plan() {
  local mode=$1 index label status failures=0 animated=false
  if [[ -t 1 && -t 2 && ${TERM:-dumb} != dumb &&
    ${CI:-false} != true && ${GITHUB_ACTIONS:-false} != true ]]; then
    animated=true
  fi
  shift
  local -a plan=("$@")
  for ((index = 0; index < ${#plan[@]}; index++)); do
    printf -v label '[%d/%d] %s' "$((index + 1))" "${#plan[@]}" "${plan[index]//_/ }"
    status=0
    if [[ $animated == true ]]; then
      run_animated_step "$label" "${plan[index]}" || status=$?
    else
      printf '%s\n' "$label"
      "${plan[index]}" || status=$?
    fi
    if ((status == 0)); then
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
