#!/usr/bin/env bash
# Globals are supplied by the main script and its native configuration.
# shellcheck disable=SC2154

install_base_packages() {
  run_command 'install packages' sudo pacman -Syu --noconfirm --needed "${base_packages[@]}" || return
}

install_main_packages() {
  run_command 'install packages' yay -S --noconfirm --needed "${main_packages[@]}" || return
}

remove_packages_step() {
  run_command 'remove unused packages' sudo pacman -Rns --noconfirm "${remove_packages[@]}" || return
}
