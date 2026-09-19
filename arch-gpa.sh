#!/usr/bin/env bash

mkdir -p bin
cd bin
makepkg -p ../cb/PKGBUILD -si
cd ..
sudo systemctl enable --now systemd-resolved
sudo ln -sf /run/systemd/resolve/stub-resolv.conf /etc/resolv.conf
