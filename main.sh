#!/usr/bin/env bash
# Arch Linux setup entrypoint. Run as the target user, not through sudo.
# Setup intentionally removes configured XDG directories and resets adopted
# repository files. Use --verify for read-only checks instead.
set -euo pipefail

script_dir=$(CDPATH='' cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)
# shellcheck source=misc/config.sh
source "$script_dir/misc/config.sh"
# shellcheck source=misc/runtime.sh
source "$script_dir/misc/runtime.sh"
# shellcheck source=steps/packages.sh
source "$script_dir/steps/packages.sh"
# shellcheck source=steps/repository.sh
source "$script_dir/steps/repository.sh"
# shellcheck source=steps/stow.sh
source "$script_dir/steps/stow.sh"
# shellcheck source=steps/yay.sh
source "$script_dir/steps/yay.sh"
# shellcheck source=steps/xdg.sh
source "$script_dir/steps/xdg.sh"
# shellcheck source=steps/mise.sh
source "$script_dir/steps/mise.sh"
# shellcheck source=steps/user.sh
source "$script_dir/steps/user.sh"
# shellcheck source=steps/verify.sh
source "$script_dir/steps/verify.sh"

usage() {
  printf 'Usage: %s [--verify | -verify] [--help]\n' "${0##*/}"
  printf 'Run from ~/repos/arch-setup; GITHUB_ACTIONS=true skips package removal.\n'
}

main() {
  local verify_only=false arg
  for arg in "$@"; do
    case $arg in
    --verify | -verify | --verify=true | -verify=true) verify_only=true ;;
    --verify=false | -verify=false) verify_only=false ;;
    --help | -help | -h)
      usage
      return 0
      ;;
    *)
      printf 'Unknown argument: %s\n' "$arg" >&2
      usage >&2
      return 2
      ;;
    esac
  done

  bootstrap || return
  local -a plan
  if [[ $verify_only == true ]]; then
    plan=(verify_repository verify_stow_system verify_yay verify_packages)
    if [[ ${GITHUB_ACTIONS:-} != true ]]; then
      plan+=(verify_removed_packages)
    fi
    plan+=(verify_xdg verify_stow_user verify_mise verify_user)
    run_plan Verification "${plan[@]}"
  else
    trap stop_sudo EXIT
    trap 'exit 130' INT
    trap 'exit 143' TERM
    start_sudo || return
    plan=(install_base_packages)
    if [[ ${GITHUB_ACTIONS:-} != true ]]; then
      plan+=(remove_packages_step)
    fi
    plan+=(clone_repository)
    plan+=(stow_system install_yay install_main_packages configure_xdg stow_user)
    plan+=(install_mise_and_tools configure_user)
    run_plan Setup "${plan[@]}"
  fi
}

# Sourcing the entrypoint lets the isolated test harness replace system commands.
if [[ ${BASH_SOURCE[0]} == "$0" ]]; then
  main "$@"
fi
