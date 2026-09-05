#!/usr/bin/env bash
# Isolated POSIX launcher tests; no package manager, network, or real TTY needed.
set -euo pipefail
# Test normal user execution even when the harness itself runs in CI.
unset GITHUB_ACTIONS
root=$(CDPATH='' cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd -P)
sandbox=$(mktemp -d)
trap 'rm -rf -- "$sandbox"' EXIT

fail() {
  printf 'FAIL: %s\n' "$*" >&2
  exit 1
}

# Strip only the final invocation; exercise the real definitions with sh.
sed '$d' "$root/index.html" >"$sandbox/definitions.sh"
cat >"$sandbox/runner.sh" <<'SH'
#!/bin/sh
set -eu
. "$TEST_ROOT/definitions.sh"
id() { printf '%s\n' "$TEST_UID"; }
stat() { printf '%s\n' "$TEST_OWNER"; }
sudo() {
  printf 'sudo %s\n' "$*" >>"$TEST_LOG"
  [ "$*" = 'pacman -Syu --needed --noconfirm git jq' ] || return 99
  return "$SUDO_STATUS"
}
git() {
  printf 'git %s\n' "$*" >>"$TEST_LOG"
  [ "$#" = 3 ] && [ "$1" = clone ] &&
    [ "$2" = https://github.com/mkvlrn/arch-setup ] &&
    [ "$3" = "$HOME/repos/arch-setup" ] || return 99
  return "$GIT_STATUS"
}
bash() {
  printf 'bash %s\n' "$*" >>"$TEST_LOG"
  [ "$#" = 1 ] && [ "$1" = "$HOME/repos/arch-setup/main.sh" ] &&
    [ "${GITHUB_ACTIONS+x}" != x ] || return 99
  IFS= read -r input
  [ "$input" = installer-input ] || return 99
  return "$SETUP_STATUS"
}
age() {
  printf 'age %s\n' "$*" >>"$TEST_LOG"
  [ "$#" = 2 ] && [ "$1" = -d ] &&
    [ "$2" = "$HOME/repos/arch-setup/secrets/secrets.tar.age" ] || return 99
  IFS= read -r input
  [ "$input" = secret-input ] || return 99
  printf 'decrypted archive\n'
  return "$AGE_STATUS"
}
tar() {
  printf 'tar\n' >>"$TEST_LOG"
  [ "$#" = 4 ] && [ "$1" = -C ] && [ "$2" = "$HOME" ] &&
    [ "$3" = -xf ] && [ -f "$4" ] || return 99
  [ "$(cat "$4")" = 'decrypted archive' ] || return 99
  return "$TAR_STATUS"
}
# Guard against accidentally adding direct privileged/network calls.
pacman() { return 99; }
curl() { return 99; }
jq() { return 99; }
with_terminal() { setup <"$TEST_ROOT/terminal-input"; }
main
SH

export TEST_ROOT=$sandbox TEST_LOG=$sandbox/calls
export TEST_UID=1000 TEST_OWNER=1000
export SUDO_STATUS=0 GIT_STATUS=0 SETUP_STATUS=0 AGE_STATUS=0 TAR_STATUS=0
export TMPDIR=$sandbox/tmp
mkdir "$TMPDIR"
printf 'installer-input\nsecret-input\n' >"$sandbox/terminal-input"

reset_case() {
  export HOME=$sandbox/"home $1"
  mkdir -p "$HOME"
  : >"$TEST_LOG"
}
run_case() {
  local expected=$1 status=0
  # Model curl | sh: stdin contains shell source, never interactive answers.
  cat "$sandbox/runner.sh" | sh >"$sandbox/output" 2>&1 || status=$?
  [[ $status == "$expected" ]] || fail "expected status $expected, got $status: $(<"$sandbox/output")"
  [[ -z $(find "$TMPDIR" -mindepth 1 -print) ]] || fail 'secret temporary files leaked'
}
no_effects() {
  [[ ! -s $TEST_LOG ]] || fail 'precondition failure ran external commands'
}

reset_case success
run_case 0
expected=$(printf 'sudo pacman -Syu --needed --noconfirm git jq\ngit clone https://github.com/mkvlrn/arch-setup %s/repos/arch-setup\nbash %s/repos/arch-setup/main.sh\nage -d %s/repos/arch-setup/secrets/secrets.tar.age\ntar' "$HOME" "$HOME" "$HOME")
[[ $(<"$TEST_LOG") == "$expected" ]] || fail 'incorrect success order'

for kind in directory file dangling; do
  reset_case "$kind"
  mkdir "$HOME/repos"
  case $kind in
  directory) mkdir "$HOME/repos/arch-setup" ;;
  file) printf 'keep me\n' >"$HOME/repos/arch-setup" ;;
  dangling) ln -s "$HOME/missing" "$HOME/repos/arch-setup" ;;
  esac
  run_case 1
  no_effects
  [[ -e $HOME/repos/arch-setup || -L $HOME/repos/arch-setup ]] || fail 'existing path removed'
done

reset_case root
TEST_UID=0 run_case 1
no_effects
reset_case owner
TEST_OWNER=2000 run_case 1
no_effects
[[ ! -e $HOME/repos ]] || fail 'invalid owner created repos'
for invalid in '' relative / "$sandbox/missing" "$sandbox/terminal-input"; do
  HOME=$invalid run_case 1
  no_effects
done

for stage in SUDO GIT SETUP AGE TAR; do
  reset_case "$stage"
  export "${stage}_STATUS=23"
  run_case 23
  export "${stage}_STATUS=0"
  case $stage in
  SUDO) [[ $(wc -l <"$TEST_LOG") == 1 ]] || fail 'sudo failure continued' ;;
  GIT) [[ $(wc -l <"$TEST_LOG") == 2 ]] || fail 'clone failure continued' ;;
  SETUP) [[ $(wc -l <"$TEST_LOG") == 3 ]] || fail 'setup failure restored secrets' ;;
  AGE) [[ $(wc -l <"$TEST_LOG") == 4 ]] || fail 'age failure extracted archive' ;;
  TAR) [[ $(wc -l <"$TEST_LOG") == 5 ]] || fail 'tar not reached' ;;
  esac
done
printf 'All launcher tests passed.\n'
