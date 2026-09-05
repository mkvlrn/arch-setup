#!/usr/bin/env bash

# These checks are sourced by the setup runner, which provides run_command and
# the configuration globals. Do not rely on errexit: callers use checks in if.
# shellcheck disable=SC2154

_verify_error() {
  printf '%s\n' "$*" >&2
}

_verify_trim() {
  REPLY=${1#"${1%%[![:space:]]*}"}
  REPLY=${REPLY%"${REPLY##*[![:space:]]}"}
}

verify_repository() {
  local remote status REPLY failed=0
  if [[ ! -d $repo_dir/.git ]]; then
    _verify_error "Repository metadata is missing or not a directory: $repo_dir/.git"
    return 1
  fi
  if remote=$(run_command 'Get repository origin' git --no-pager -C "$repo_dir" remote get-url origin); then
    _verify_trim "$remote"
    if [[ $REPLY != "$repo_ssh" ]]; then
      _verify_error "Expected origin '$repo_ssh', got '$REPLY'"
      failed=1
    fi
  else
    _verify_error "Cannot get repository origin: $repo_dir"
    failed=1
  fi
  if status=$(run_command 'Get repository status' git --no-pager --no-optional-locks -C "$repo_dir" status --porcelain); then
    _verify_trim "$status"
    if [[ -n $REPLY ]]; then
      _verify_error "Repository is dirty:
$status"
      failed=1
    fi
  else
    _verify_error "Cannot get repository status: $repo_dir"
    failed=1
  fi
  return "$failed"
}

_verify_stow_link() {
  local source=$1 destination=$2 source_real destination_real failed=0
  if [[ ! -L $destination ]]; then
    _verify_error "Stow destination is missing or not a symlink: $destination"
    return 1
  fi
  # A sentinel preserves trailing newlines in filenames during capture.
  if source_real=$(run_command "Resolve source: $source" realpath -e -- "$source" && printf '.'); then
    source_real=${source_real%$'\n.'}
  else
    _verify_error "Cannot resolve Stow source: $source"
    failed=1
  fi
  if destination_real=$(run_command "Resolve destination: $destination" realpath -e -- "$destination" && printf '.'); then
    destination_real=${destination_real%$'\n.'}
  else
    _verify_error "Cannot resolve Stow destination: $destination"
    failed=1
  fi
  if ((failed)); then
    return 1
  fi
  if [[ $source_real != "$destination_real" ]]; then
    _verify_error "Stow destination '$destination' points to '$destination_real' instead of '$source_real'"
    return 1
  fi
  return 0
}

_verify_stow_walk() {
  local source=$1 destination=$2 child name failed=0
  # Like WalkDir, do not descend through source symlinks, even to directories.
  if [[ -d $source && ! -L $source ]]; then
    if [[ ! -r $source || ! -x $source ]]; then
      _verify_error "Cannot read/traverse Stow directory: $source"
      return 1
    fi
    for child in "$source"/*; do
      name=${child##*/}
      if _verify_stow_walk "$child" "${destination%/}/$name"; then
        :
      else
        failed=1
      fi
    done
  elif [[ -e $source || -L $source ]]; then
    case ${source##*/} in
    .gitkeep | .stow-local-ignore) return 0 ;;
    esac
    if _verify_stow_link "$source" "$destination"; then
      :
    else
      failed=1
    fi
  else
    _verify_error "Cannot inspect Stow source: $source"
    failed=1
  fi
  return "$failed"
}

_verify_stow_tree() (
  # Isolate glob settings from the runner, and include hidden files regardless
  # of its GLOBIGNORE setting. nullglob makes empty directories valid.
  unset GLOBIGNORE
  set +f
  shopt -u failglob
  shopt -s dotglob nullglob
  if _verify_stow_walk "$1" "$2"; then
    return 0
  else
    return 1
  fi
)

verify_stow_system() {
  if _verify_stow_tree "$repo_dir/stow/system" /; then
    return 0
  else
    return 1
  fi
}

verify_stow_user() {
  if _verify_stow_tree "$repo_dir/stow/user" "$home_dir"; then
    return 0
  else
    return 1
  fi
}

verify_yay() {
  local packages package content failed=0
  if run_command 'Get Yay version' yay --version; then
    :
  else
    _verify_error 'Yay is not available'
    failed=1
  fi
  if packages=$(run_command 'List installed packages' yay -Qq); then
    while IFS= read -r package; do
      if [[ $package == *-debug ]]; then
        _verify_error "Debug package is installed: $package"
        failed=1
      fi
    done <<<"$packages"
  else
    _verify_error 'Cannot list installed packages'
    failed=1
  fi
  if content=$(run_command 'Read mirrorlist' cat -- "$mirror_list_path" && printf '.'); then
    content=${content%.}
    if [[ $content != *"$mirror_list_check"* ]]; then
      _verify_error "Mirrorlist '$mirror_list_path' lacks the expected Reflector command"
      failed=1
    fi
  else
    _verify_error "Cannot read mirrorlist: $mirror_list_path"
    failed=1
  fi
  return "$failed"
}

verify_packages() {
  local package failed=0
  for package in "${base_packages[@]}" "${main_packages[@]}"; do
    if run_command "Query package: $package" pacman -Q "$package"; then
      :
    else
      _verify_error "Package is not installed: $package"
      failed=1
    fi
  done
  return "$failed"
}

verify_removed_packages() {
  local package failed=0
  for package in "${remove_packages[@]}"; do
    # Match Go: a failed pacman query means the package is absent.
    if pacman -Q "$package" >/dev/null 2>&1; then
      _verify_error "Package is still installed: $package"
      failed=1
    fi
  done
  return "$failed"
}

verify_xdg() {
  local directory path detail failed=0
  for directory in "${xdg_mkdir[@]}"; do
    path=$home_dir/$directory
    if [[ ! -d $path ]]; then
      _verify_error "XDG path is missing, inaccessible, or not a directory: $path"
      failed=1
    fi
  done
  for directory in "${xdg_rmrf[@]}"; do
    path=$home_dir/$directory
    # Follow symlinks as os.Stat does, including accepting dangling links.
    # Preserve stat's diagnostic to distinguish absence from permission errors.
    if detail=$(LC_ALL=C stat -L -- "$path" 2>&1); then
      _verify_error "XDG path still exists: $path"
      failed=1
    elif [[ $detail != *': No such file or directory' ]]; then
      _verify_error "Cannot inspect removed XDG path '$path': $detail"
      failed=1
    fi
  done
  return "$failed"
}

verify_mise() {
  local missing REPLY failed=0
  if run_command 'Get mise version' "$home_dir/.local/bin/mise" --version; then
    :
  else
    _verify_error 'Mise is not available'
    failed=1
  fi
  if missing=$(run_command 'List missing mise tools' "$home_dir/.local/bin/mise" ls --global --missing --no-header); then
    _verify_trim "$missing"
    if [[ -n $REPLY ]]; then
      _verify_error "Mise tools are missing:
$REPLY"
      failed=1
    fi
  else
    _verify_error 'Cannot list missing mise tools'
    failed=1
  fi
  return "$failed"
}

_verify_passwd_field() {
  local account=$1 field=$2 expected=$3 entry REPLY remainder
  local -a fields=()
  if entry=$(run_command "Get passwd entry: $account" getent passwd "$account"); then
    _verify_trim "$entry"
    remainder=$REPLY
    # Unlike read -a, preserve an empty final field and reject extra fields.
    while [[ $remainder == *:* ]]; do
      fields+=("${remainder%%:*}")
      remainder=${remainder#*:}
    done
    fields+=("$remainder")
    if ((${#fields[@]} != 7)) || [[ $REPLY == *$'\n'* ]]; then
      _verify_error "Unexpected passwd entry for '$account'"
      return 1
    fi
    if [[ ${fields[field]} != "$expected" ]]; then
      _verify_error "Passwd field $field for '$account' is '${fields[field]}' instead of '$expected'"
      return 1
    fi
  else
    _verify_error "Cannot get passwd entry for '$account'"
    return 1
  fi
  return 0
}

verify_user() {
  local groups mode unit state completion failed=0
  if _verify_passwd_field "$username" 6 /usr/bin/fish; then
    :
  else
    failed=1
  fi
  if groups=$(run_command "Get groups: $username" id -nG "$username"); then
    if [[ ! $groups =~ (^|[[:space:]])docker($|[[:space:]]) ]]; then
      _verify_error "User '$username' is not in the docker group"
      failed=1
    fi
  else
    _verify_error "Cannot get groups for '$username'"
    failed=1
  fi
  if _verify_passwd_field ftp 5 "$home_dir/torrents"; then
    :
  else
    failed=1
  fi
  if mode=$(run_command 'Inspect home permissions' stat -L -c %a -- "$home_dir"); then
    if [[ ! $mode =~ ^[0-7]+$ ]]; then
      _verify_error "Invalid permission mode for '$home_dir': $mode"
      failed=1
    elif (((8#$mode & 1) == 0)); then
      _verify_error "Home is not traversable by other users: $home_dir"
      failed=1
    fi
  else
    _verify_error "Cannot inspect home directory: $home_dir"
    failed=1
  fi
  for unit in docker.socket pure-ftpd.service paccache.timer; do
    for state in enabled active; do
      if run_command "Check $unit is $state" systemctl "is-$state" --quiet "$unit"; then
        :
      else
        _verify_error "Systemd unit '$unit' is not $state"
        failed=1
      fi
    done
  done
  for completion in mise gh glab; do
    if [[ ! -e $home_dir/.config/fish/completions/$completion.fish ]]; then
      _verify_error "Completion file is missing or inaccessible: $home_dir/.config/fish/completions/$completion.fish"
      failed=1
    fi
  done
  return "$failed"
}
