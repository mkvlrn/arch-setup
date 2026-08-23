#!/bin/sh

set -eu

cache_dir="${PACKAGE_CACHE_DIR:?PACKAGE_CACHE_DIR is required}"
mise_github_token="${MISE_GITHUB_TOKEN:?MISE_GITHUB_TOKEN is required}"
pacman_cache_dir="$cache_dir/pacman"
yay_cache_dir="$cache_dir/yay"

# Run a command inside the Arch VM.
#
# Host key persistence is intentionally disabled because the VM is ephemeral.
# LogLevel=ERROR suppresses the expected host-key warning without hiding
# command failures.
ssh_vm() {
  sshpass -p arch ssh \
    -p 2222 \
    -o StrictHostKeyChecking=no \
    -o UserKnownHostsFile=/dev/null \
    -o LogLevel=ERROR \
    arch@127.0.0.1 \
    "$@"
}

# Copy files between the runner and the Arch VM using the same ephemeral SSH
# connection settings as ssh_vm.
scp_vm() {
  sshpass -p arch scp \
    -P 2222 \
    -o StrictHostKeyChecking=no \
    -o UserKnownHostsFile=/dev/null \
    -o LogLevel=ERROR \
    "$@"
}

# Make VM preparation work visible in the Actions log while keeping the actual
# installer and assertion output unchanged.
vm_step() {
  printf '\n🖥️  %s\n' "$1"
}

# actions/cache restores these directories before this script runs. On the
# first run they will simply be empty.
mkdir -p \
  "$pacman_cache_dir" \
  "$yay_cache_dir"
vm_step "Copying installer to VM"

# Test the exact executable produced by the check-build job.
scp_vm \
  ./bin/arch-setup \
  arch@127.0.0.1:/tmp/arch-setup
vm_step "Copying machine assertions to VM"

# Assertions run after installation to verify the resulting machine state.
scp_vm \
  ./.github/workflows/scripts/assert-vm.sh \
  arch@127.0.0.1:/tmp/assert-vm.sh
vm_step "Restoring pacman package cache"

# Restore cached pacman packages into a directory owned by the VM user. Using a
# user-owned directory makes transferring the cache in and out of the VM easy.
tar -C "$pacman_cache_dir" -cf - . |
  ssh_vm '
    mkdir -p "$HOME/.cache/pacman/pkg"
    tar -C "$HOME/.cache/pacman/pkg" -xf -
  '
vm_step "Mounting pacman package cache"

# Bind the restored directory over pacman's normal cache location. Pacman can
# therefore reuse cached packages without any installer-specific configuration.
ssh_vm '
  printf "%s\n" arch |
    sudo -S mount --bind \
      "$HOME/.cache/pacman/pkg" \
      /var/cache/pacman/pkg
'
vm_step "Restoring yay package cache"

# Restore yay's normal cache/build directory so downloaded sources and previous
# build artifacts are available to yay during installation.
tar -C "$yay_cache_dir" -cf - . |
  ssh_vm '
    mkdir -p "$HOME/.cache/yay"
    tar -C "$HOME/.cache/yay" -xf -
  '
vm_step "Preparing sudo"

# Prime sudo before running the installer so later privileged commands can use
# the credentials already configured for the test VM.
ssh_vm "printf '%s\n' arch | sudo -S -v"
vm_step "Running arch-setup"

# MISE_GITHUB_TOKEN belongs to the runner environment, so explicitly forward it
# to the installer process inside the VM. Bun and mise inherit it from there.
ssh_vm \
  "chmod +x /tmp/arch-setup && MISE_GITHUB_TOKEN='$mise_github_token' /tmp/arch-setup"
vm_step "Running machine state checks"

# Keep assertion output untouched so the CI log remains the same as a normal
# installation test.
ssh_vm \
  'chmod +x /tmp/assert-vm.sh && /tmp/assert-vm.sh'
vm_step "Cleaning pacman download leftovers"

# Pacman 7 may leave temporary download-* directories in its package cache.
# Remove only those temporary directories and preserve all actual package files.
ssh_vm \
  "printf '%s\n' arch | sudo -S find /var/cache/pacman/pkg \
    -mindepth 1 -maxdepth 1 \
    -type d -name 'download-*' \
    -exec rm -rf -- {} +"
vm_step "Preparing updated package caches"

# Replace the host-side directories restored by actions/cache with the state
# produced by this VM run. actions/cache will save them after the job succeeds.
rm -rf \
  "$pacman_cache_dir" \
  "$yay_cache_dir"
mkdir -p \
  "$pacman_cache_dir" \
  "$yay_cache_dir"
vm_step "Saving pacman package cache"

# Pacman creates package files as root, so archive the cache with sudo inside
# the VM and unpack it as the runner user on the host.
ssh_vm \
  "printf '%s\n' arch | sudo -S tar -C /var/cache/pacman/pkg -cf - ." |
  tar -C "$pacman_cache_dir" -xf -
vm_step "Saving yay package cache"

# Yay runs as the normal Arch user, so its cache can be exported directly.
ssh_vm \
  'tar -C "$HOME/.cache/yay" -cf - .' |
  tar -C "$yay_cache_dir" -xf -
