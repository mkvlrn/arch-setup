#!/usr/bin/env bash

mkdir -p bin
cd bin
cp ../cb/PKGBUILD .
makepkg -si
cd ..
sudo systemctl enable --now systemd-resolved gpaservice.service gpawatchdog.service wapptunnel.service
sudo ln -sf /run/systemd/resolve/stub-resolv.conf /etc/resolv.conf
