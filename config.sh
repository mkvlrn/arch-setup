#!/usr/bin/env bash

# This file is sourced by install.sh and verify.sh; its variables are consumed there.
# shellcheck disable=SC2034

BASE_PACKAGES=(
  git
  base-devel
  reflector
  stow
)

REMOVE_PACKAGES=(
  vim
)

MAIN_PACKAGES=(
  age
  ark
  bruno-bin
  deluge-gtk
  deluge
  docker-buildx
  docker-compose
  docker
  dolphin
  ferdium-bin
  fish
  ghostty
  kate
  less
  okular
  openssh
  pacman-contrib
  power-profiles-daemon
  pure-ftpd
  qalculate-qt
  ttf-hack-nerd
  ttf-iosevkaterm-nerd
  ttf-noto-nerd
  ttf-zed-mono-nerd
  unzip
  visual-studio-code-bin
  xdg-user-dirs
  zed
  zen-browser-bin
)

REPO_HTTP='https://github.com/mkvlrn/arch-setup'
REPO_SSH='git@github.com:mkvlrn/arch-setup'
MIRROR_LIST='/etc/pacman.d/mirrorlist'
MIRROR_LIST_CHECK='# With:       reflector --latest 20 --protocol https --sort rate --save /etc/pacman.d/mirrorlist'

XDG_MKDIR=(repos work documents downloads media torrents)
XDG_RMRF=(Desktop Documents Downloads Music Pictures Projects Public Templates Videos)
