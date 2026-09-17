#!/usr/bin/env bash

# This file is sourced by install.sh and verify.sh; its variables are consumed there.
# shellcheck disable=SC2034

BASE_PACKAGES=(
  base-devel
  git
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
  docker
  docker-buildx
  docker-compose
  dolphin
  ferdium-bin
  fish
  ghostty
  kate
  kcalc
  less
  okular
  openssh
  pacman-contrib
  power-profiles-daemon
  unzip
  wl-clipboard
  xdg-user-dirs
  zed
  zen-browser-bin
)

GETNF_FONTS=(
  Hack
  IosevkaTerm
  Noto
  ZedMono
)

REPO_SSH='git@github.com:mkvlrn/arch-setup'

XDG_MKDIR=(
  documents
  downloads
  media
  proton
  repos
  work
)
XDG_RMRF=(
  Desktop
  Documents
  Downloads
  Music
  Pictures
  Projects
  Public
  Templates
  Videos
)
