#!/usr/bin/env bash
# STEP_NAME is read by install.sh after this file is sourced.
# shellcheck disable=SC2034

STEP_NAME='Installing yay and updating mirrors'

step_run() {
  local yay_dir=${TMPDIR:-/tmp}/yay-bin
  rm -rf "$yay_dir"
  run git clone https://aur.archlinux.org/yay-bin "$yay_dir"
  run makepkg -si --noconfirm -C -D "$yay_dir"
  run yay -Y --gendb
  run yay -Y --devel --save
  run sudo reflector --latest 20 --protocol https --sort rate --save /etc/pacman.d/mirrorlist
  run yay -Syu --noconfirm

  local debug_packages
  debug_packages=$(yay -Qq | grep -- '-debug$' || true)
  if [[ -n $debug_packages ]]; then
    printf '%s\n' "$debug_packages" | xargs -r yay -Rnsu
  fi
}
