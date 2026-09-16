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
  bruno-bin # export in distrobox
  docker-buildx
  docker-compose
  docker
  dolphin
  ferdium-bin # export in distrobox
  fish
  ghostty # export in distrobox
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
  zed             # export in distrobox
  zen-browser-bin # export in distrobox
)

GETNF_FONTS=(
  Hack
  IosevkaTerm
  Noto
  ZedMono
)

REPO_SSH='git@github.com:mkvlrn/arch-setup'

XDG_MKDIR=(
  repos
  work
  documents
  downloads
  media
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
