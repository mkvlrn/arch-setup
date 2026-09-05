#!/usr/bin/env bash
# Globals are supplied by the main script and its native configuration.
# shellcheck disable=SC2154

install_yay() {
  run_command 'clone yay-bin' git clone https://aur.archlinux.org/yay-bin "$temp_dir/yay-bin" || return
  (
    cd -- "$temp_dir/yay-bin" || return
    run_command 'build yay' makepkg -si --noconfirm || return
  ) || return
  run_command 'track git packages' yay -Y --gendb || return
  run_command 'enable dev packages updates' yay -Y --devel --save || return
  run_command 'get best mirrors list' sudo reflector --latest 20 --protocol https --sort rate --save "$mirror_list_path" || return
  run_command 'update package data' yay -Syu --noconfirm || return
  # A Bash filter treats no debug packages as success without hiding query errors.
  run_command 'remove debug packages' bash -o pipefail -c "
    yay -Qq | {
      while IFS= read -r package; do
        if [[ \$package == *-debug ]]; then
          printf '%s\\n' \"\$package\" || exit
        fi
      done
    } | xargs -r yay -Rnsu
  " || return
}
