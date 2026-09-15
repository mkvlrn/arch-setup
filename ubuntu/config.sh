#!/usr/bin/env bash

# Kubuntu package and application sources. Every item in these lists is
# required for this personal workstation setup.
# shellcheck disable=SC2034

APT_PACKAGES=(
  age
  ark
  curl
  docker-buildx
  docker-compose-v2
  docker.io
  dolphin
  fish
  flatpak
  ghostty
  kate
  less
  okular
  openssh-server
  power-profiles-daemon
  qalculate-qt
  stow
  unzip
  xdg-user-dirs
)

FLATPAKS=(
  app.zen_browser.zen
  com.usebruno.Bruno
  dev.zed.Zed
  org.ferdium.Ferdium
)

GETNF_FONTS='Hack IosevkaTerm Noto ZedMono'

XDG_MKDIR=(repos
  work
  documents
  downloads
  media
  torrents
)
XDG_RMRF=(Desktop
  Documents
  Downloads
  Music
  Pictures
  Projects
  Public
  Templates
  Videos)

REPO_SSH='git@github.com:mkvlrn/arch-setup'
