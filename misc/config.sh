#!/usr/bin/env bash
# Configuration stays in config.json; jq is required before any system changes.
# Globals are consumed by the sourced setup and verification modules.
# shellcheck disable=SC2034,SC2154
base_packages=() remove_packages=() main_packages=() xdg_mkdir=() xdg_rmrf=()
repo_http='' repo_ssh='' mirror_list_path='' mirror_list_check=''

_config_array() {
  local json=$1 key=$2
  local -n destination=$3
  # NUL delimiters preserve whitespace, glob characters, and empty array entries.
  # Wait explicitly: process-substitution failures do not propagate via mapfile.
  mapfile -d '' -t destination < <(jq -j --arg key "$key" '.[$key][] | ., "\u0000"' <<<"$json")
  wait "$!"
}

_config_string() {
  local json=$1 key=$2 value
  # A sentinel prevents command substitution from stripping trailing newlines.
  value=$(jq -j --arg key "$key" '.[$key], "."' <<<"$json") || return
  printf -v "$3" '%s' "${value%.}"
}

load_config() {
  local path=$1 json
  if ! command -v jq >/dev/null 2>&1; then
    printf 'jq is required to read config.json. Install it with: sudo pacman -S jq\n' >&2
    return 1
  fi
  # Validate one complete document before assigning any configuration globals.
  # Bash cannot represent NUL bytes, so reject them instead of silently truncating.
  if ! json=$(jq -cse '
    def text: type == "string" and (contains("\u0000") | not);
    if length == 1 and (.[0] | type == "object") then .[0]
    else error("expected one configuration object") end
    | if (
        all(.basePackages, .removePackages, .mainPackages, .xdgMkDir, .xdgRmRf;
          type == "array" and all(.[]; text))
        and all(.repoHttp, .repoSsh, .mirrorListPath, .mirrorListCheck; text)
      ) then . else error("expected string arrays and string configuration fields") end
  ' -- "$path"); then
    printf 'Could not load configuration: %s\n' "$path" >&2
    return 1
  fi
  _config_array "$json" basePackages base_packages || return
  _config_array "$json" removePackages remove_packages || return
  _config_array "$json" mainPackages main_packages || return
  _config_array "$json" xdgMkDir xdg_mkdir || return
  _config_array "$json" xdgRmRf xdg_rmrf || return
  _config_string "$json" repoHttp repo_http || return
  _config_string "$json" repoSsh repo_ssh || return
  _config_string "$json" mirrorListPath mirror_list_path || return
  _config_string "$json" mirrorListCheck mirror_list_check || return
}

bootstrap() {
  load_config "$script_dir/config.json" || return
  if [[ -z ${HOME:-} || $HOME != /* ]]; then
    printf 'HOME must be a nonempty absolute path.\n' >&2
    return 1
  fi
  home_dir=${HOME%/}
  [[ -n $home_dir ]] || home_dir=/
  username=$(id -un) || return
  repo_dir=$(realpath -e -- "$home_dir/repos/arch-setup") || return
  if [[ $script_dir != "$repo_dir" || ! -d $repo_dir/.git ]]; then
    printf 'Run the installer from the checkout at %s/repos/arch-setup.\n' "$home_dir" >&2
    return 1
  fi
  temp_dir=${TMPDIR:-/tmp}
  # Resolve relative TMPDIR before steps change working directories.
  temp_dir=$(realpath -e -- "$temp_dir") || return
}
