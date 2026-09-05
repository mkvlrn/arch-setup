#!/usr/bin/env bash
# Isolated regression tests: never invoke package managers, sudo, or services.
set -euo pipefail
root=$(CDPATH='' cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd -P)
# shellcheck source=main.sh
source "$root/main.sh"
sandbox=$(mktemp -d)
trap 'rm -rf -- "$sandbox"' EXIT

fail() {
  printf 'FAIL: %s\n' "$*" >&2
  exit 1
}

# Load the real configuration from outside the checkout, then exercise JSON
# validation and argument preservation without executing any setup operations.
mkdir -p "$sandbox/bootstrap/repos/arch-setup/.git"
cp "$root/config.json" "$sandbox/bootstrap/repos/arch-setup/config.json"
(
  HOME=$sandbox/bootstrap
  script_dir=$HOME/repos/arch-setup
  cd "$sandbox"
  bootstrap || fail 'bootstrap depends on working directory'
  [[ $repo_dir == "$script_dir" ]] || fail 'bootstrap chose a different checkout'
  script_dir=$root
  if bootstrap >/dev/null 2>&1; then fail 'mismatched checkout accepted'; fi
) || fail 'checkout bootstrap validation'
load_config "$root/config.json" || fail 'real configuration rejected'
[[ ${base_packages[0]} == git && $repo_http == https://github.com/mkvlrn/arch-setup ]] ||
  fail 'configuration fields not loaded'
jq '.basePackages = ["space name", "wild*card", "line\nbreak", "", "$(false)"]
    | .removePackages = [] | .repoSsh = "remote\n\n"' \
  "$root/config.json" >"$sandbox/config.json"
load_config "$sandbox/config.json" || fail 'valid special characters rejected'
# Literal shell syntax in JSON must remain data, not execute as a command.
# shellcheck disable=SC2016
[[ ${#base_packages[@]} == 5 && ${base_packages[0]} == 'space name' &&
  ${base_packages[1]} == 'wild*card' && ${base_packages[2]} == $'line\nbreak' &&
  ${base_packages[3]} == '' && ${base_packages[4]} == '$(false)' &&
  ${#remove_packages[@]} == 0 && $repo_ssh == $'remote\n\n' ]] ||
  fail 'JSON string boundaries changed'
for filter in '.basePackages = 123' '.mainPackages = [false]' \
  'del(.repoHttp)' '.repoSsh = "nul\u0000byte"' '.xdgMkDir = null'; do
  jq "$filter" "$root/config.json" >"$sandbox/invalid.json"
  if load_config "$sandbox/invalid.json" >/dev/null 2>&1; then
    fail "invalid configuration accepted: $filter"
  fi
done
printf '{' >"$sandbox/invalid.json"
if load_config "$sandbox/invalid.json" >/dev/null 2>&1; then fail 'malformed JSON accepted'; fi
cat "$root/config.json" "$root/config.json" >"$sandbox/invalid.json"
if load_config "$sandbox/invalid.json" >/dev/null 2>&1; then fail 'multiple documents accepted'; fi
if load_config "$sandbox/missing.json" >/dev/null 2>&1; then fail 'missing JSON accepted'; fi
if (PATH=/nonexistent load_config "$root/config.json") >"$sandbox/error" 2>&1; then
  fail 'missing jq accepted'
fi
[[ $(<"$sandbox/error") == *'sudo pacman -S jq'* ]] || fail 'missing jq guidance absent'

# Capture plans without executing any of their functions.
bootstrap() { :; }
start_sudo() { :; }
run_plan() { printf '%s\n' "$@"; }
expected=$'Setup\ninstall_base_packages\nremove_packages_step\nclone_repository\nstow_system\ninstall_yay\ninstall_main_packages\nconfigure_xdg\nstow_user\ninstall_mise_and_tools\nconfigure_user'
actual=$(GITHUB_ACTIONS=false main)
[[ $actual == "$expected" ]] || fail 'normal setup plan'
actual=$(GITHUB_ACTIONS=true main)
[[ $actual == *clone_repository* && $actual != *remove_packages_step* ]] || fail 'CI setup plan'
actual=$(GITHUB_ACTIONS=false main --verify)
[[ $actual == Verification* && $actual == *verify_removed_packages* ]] || fail 'verification plan'
actual=$(GITHUB_ACTIONS=true main -verify)
[[ $actual != *verify_removed_packages* && $actual == *verify_user ]] || fail 'CI verification plan'
if main --invalid >/dev/null 2>&1; then fail 'unknown option accepted'; fi

# Restore the real runner and prove stop-on-error versus aggregate-check behavior.
# shellcheck source=misc/runtime.sh
source "$root/misc/runtime.sh"
first() {
  printf 'first\n' >>"$sandbox/calls"
  return 1
}
second() { printf 'second\n' >>"$sandbox/calls"; }
if run_plan Setup first second >/dev/null 2>&1; then fail 'setup failure lost'; fi
[[ $(<"$sandbox/calls") == first ]] || fail 'setup continued after failure'
: >"$sandbox/calls"
if run_plan Verification first second >/dev/null 2>&1; then fail 'verification failure lost'; fi
[[ $(<"$sandbox/calls") == $'first\nsecond' ]] || fail 'verification stopped early'
if run_command failure bash -c 'exit 23' >/dev/null 2>&1; then
  fail 'command failure lost'
else
  [[ $? == 23 ]] || fail 'command status changed'
fi

# Successful commands hide both streams; failed commands report both streams
# on stderr without contaminating values captured by verification checks.
mkdir "$sandbox/logs"
noisy_command() {
  printf 'command output\n\n'
  printf 'command diagnostic\n' >&2
  return "${1:-0}"
}
TMPDIR=$sandbox/logs run_command quiet noisy_command >"$sandbox/out" 2>"$sandbox/err"
[[ ! -s $sandbox/out && ! -s $sandbox/err ]] || fail 'successful command was noisy'
TMPDIR=$sandbox/logs capture_command query noisy_command >"$sandbox/out" 2>"$sandbox/err"
printf 'command output\n\n' >"$sandbox/expected"
cmp "$sandbox/out" "$sandbox/expected" || fail 'captured stdout changed'
[[ ! -s $sandbox/err ]] || fail 'successful query leaked stderr'
if TMPDIR=$sandbox/logs capture_command broken noisy_command 23 >"$sandbox/out" 2>"$sandbox/err"; then
  fail 'noisy failure accepted'
else
  [[ $? == 23 ]] || fail 'noisy failure status changed'
fi
[[ ! -s $sandbox/out ]] || fail 'failure contaminated captured stdout'
error=$(<"$sandbox/err")
[[ $error == *'broken (exit 23)'* && $error == *$'stdout:\ncommand output'* &&
  $error == *$'stderr:\ncommand diagnostic'* ]] || fail 'failure output missing'
[[ -z $(find "$sandbox/logs" -mindepth 1 -print) ]] || fail 'command logs leaked'
quiet_step() { run_command quiet noisy_command; }
actual=$(run_plan Setup quiet_step)
[[ $actual == $'[1/1] quiet step\nDone.' ]] || fail 'main progress changed or command output leaked'

# Exercise animation separately from TTY detection, with no real setup commands.
animated_success() { sleep 0.2; }
animated_failure() {
  printf 'step diagnostic\n' >&2
  return 17
}
TMPDIR=$sandbox/logs run_animated_step '[1/1] test' animated_success \
  >"$sandbox/out" 2>"$sandbox/err"
[[ $(<"$sandbox/out") == *'[1/1] test done'* && ! -s $sandbox/err ]] ||
  fail 'animation did not finish cleanly'
if TMPDIR=$sandbox/logs run_animated_step '[1/1] test' animated_failure \
  >"$sandbox/out" 2>"$sandbox/err"; then
  fail 'animated failure accepted'
else
  [[ $? == 17 ]] || fail 'animation changed failure status'
fi
[[ $(<"$sandbox/out") == *'[1/1] test failed'* &&
$(<"$sandbox/err") == 'step diagnostic' ]] || fail 'animation hid diagnostics'
[[ -z $(find "$sandbox/logs" -mindepth 1 -print) ]] || fail 'animation logs leaked'
[[ -z $(jobs -pr) ]] || fail 'animation process leaked'
actual=$(GITHUB_ACTIONS=true run_plan Setup quiet_step)
[[ $actual == $'[1/1] quiet step\nDone.' ]] || fail 'CI output contains animation'

home_dir=$sandbox/home
repo_dir=$sandbox/repo
mkdir -p "$home_dir" "$repo_dir/stow/user/.config"
printf 'fixture\n' >"$repo_dir/stow/user/.config/file with spaces"
mkdir -p "$home_dir/.config"
ln -s "$repo_dir/stow/user/.config/file with spaces" "$home_dir/.config/file with spaces"
verify_stow_user || fail 'valid Stow link rejected'
rm -- "$home_dir/.config/file with spaces"
printf 'not a link\n' >"$home_dir/.config/file with spaces"
if verify_stow_user >/dev/null 2>&1; then fail 'non-symlink accepted'; fi
xdg_mkdir=(.config)
xdg_rmrf=(Desktop)
verify_xdg || fail 'absent XDG directory rejected'
mkdir "$home_dir/Desktop"
if verify_xdg >/dev/null 2>&1; then fail 'existing removed directory accepted'; fi

# Package checks must inspect every package even after an earlier failure.
base_packages=(missing)
main_packages=(present)
remove_packages=(absent)
pacman() {
  printf '%s\n' "$2" >>"$sandbox/packages"
  [[ $2 == present ]]
}
if verify_packages >/dev/null 2>&1; then fail 'missing package accepted'; fi
[[ $(<"$sandbox/packages") == $'missing\npresent' ]] || fail 'package checks stopped early'
verify_removed_packages || fail 'absent package rejected'
# Assert the destructive Git operations are restricted to the adopted package,
# and that a failed restore prevents cleanup. No real Git writes are performed.
run_command() {
  local description=$1
  shift
  printf '%s\n' "$*" >>"$sandbox/stow-commands"
  [[ $description != "${fail_command:-}" ]]
}
stow_user || fail 'Stow user step failed'
expected="stow -R --no-folding --adopt -d $repo_dir/stow -t $home_dir user"
expected+=$'\ngit restore -- stow/user\ngit clean -fd -- stow/user'
[[ $(<"$sandbox/stow-commands") == "$expected" ]] || fail 'Stow cleanup escaped its package'
: >"$sandbox/stow-commands"
fail_command='restore adopted files'
if stow_user; then fail 'Stow restore failure lost'; fi
[[ $(<"$sandbox/stow-commands") != *'git clean'* ]] || fail 'cleanup ran after restore failure'
printf 'All Bash regression tests passed.\n'
