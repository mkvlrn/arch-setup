#!/bin/sh

set -eu

cache_dir="${PACKAGE_CACHE_DIR:?PACKAGE_CACHE_DIR is required}"
mise_github_token="${MISE_GITHUB_TOKEN:?MISE_GITHUB_TOKEN is required}"
pacman_cache_dir="$cache_dir/pacman"
yay_cache_dir="$cache_dir/yay"
mise_cache_dir="$cache_dir/mise"

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
  "$yay_cache_dir" \
  "$mise_cache_dir"

# Test the exact executable produced by the check-build job.
vm_step "Copying installer to VM"
scp_vm \
  ./bin/arch-setup \
  arch@127.0.0.1:/tmp/arch-setup

# Copy the exact repository contents checked out for this PR, including Git
# metadata required by the user Stow step.
vm_step "Copying candidate repository to VM"
{
  printf '%s\0' .git
  git ls-files -z
} |
  tar --null --files-from=- -cf - |
  ssh_vm '
    mkdir -p "$HOME/repos/arch-setup"
    tar -C "$HOME/repos/arch-setup" -xf -
  '

# Restore cached pacman packages into a directory owned by the VM user. Using a
# user-owned directory makes transferring the cache in and out of the VM easy.
vm_step "Restoring pacman package cache"
tar -C "$pacman_cache_dir" -cf - . |
  ssh_vm '
    mkdir -p "$HOME/.cache/pacman/pkg"
    tar -C "$HOME/.cache/pacman/pkg" -xf -
  '

# Bind the restored directory over pacman's normal cache location. Pacman can
# therefore reuse cached packages without any installer-specific configuration.
vm_step "Mounting pacman package cache"
ssh_vm '
  printf "%s\n" arch |
    sudo -S mount --bind \
      "$HOME/.cache/pacman/pkg" \
      /var/cache/pacman/pkg
'

# Restore yay's normal cache/build directory so downloaded sources and previous
# build artifacts are available to yay during installation.
vm_step "Restoring yay package cache"
tar -C "$yay_cache_dir" -cf - . |
  ssh_vm '
    mkdir -p "$HOME/.cache/yay"
    tar -C "$HOME/.cache/yay" -xf -
  '

# Restore downloads retained by mise so supported backends can reuse them while
# still performing a normal installation into a fresh machine.
vm_step "Restoring mise download cache"
tar -C "$mise_cache_dir" -cf - . |
  ssh_vm '
    mkdir -p "$HOME/.local/share/mise/downloads"
    tar -C "$HOME/.local/share/mise/downloads" -xf -
  '

# MISE_GITHUB_TOKEN belongs to the runner environment, so explicitly forward it
# to the installer process inside the VM.
vm_step "Running arch-setup"
ssh_vm \
  "chmod +x /tmp/arch-setup &&
   printf '%s\n' arch | sudo -S -v &&
   MISE_GITHUB_TOKEN='$mise_github_token' \
   MISE_ALWAYS_KEEP_DOWNLOAD=1 \
   /tmp/arch-setup --repo-ready"

# Run verification in a new login session so changes such as supplementary
# group membership are visible.
vm_step "Verifying machine state"
ssh_vm \
  '/tmp/arch-setup --verify'

# Pacman 7 may leave temporary download-* directories in its package cache.
# Remove only those temporary directories and preserve all actual package files.
vm_step "Cleaning pacman download leftovers"
ssh_vm \
  "printf '%s\n' arch | sudo -S find /var/cache/pacman/pkg \
    -mindepth 1 -maxdepth 1 \
    -type d -name 'download-*' \
    -exec rm -rf -- {} +"

# Replace the host-side directories restored by actions/cache with the state
# produced by this VM run. actions/cache will save them after the job succeeds.
vm_step "Preparing updated package caches"
rm -rf \
  "$pacman_cache_dir" \
  "$yay_cache_dir" \
  "$mise_cache_dir"
mkdir -p \
  "$pacman_cache_dir" \
  "$yay_cache_dir" \
  "$mise_cache_dir"

# Pacman creates package files as root, so archive the cache with sudo inside
# the VM and unpack it as the runner user on the host.
vm_step "Saving pacman package cache"
ssh_vm \
  "printf '%s\n' arch | sudo -S tar -C /var/cache/pacman/pkg -cf - ." |
  tar -C "$pacman_cache_dir" -xf -

# Yay runs as the normal Arch user, so its cache can be exported directly.
vm_step "Saving yay package cache"
ssh_vm \
  'tar -C "$HOME/.cache/yay" -cf - .' |
  tar -C "$yay_cache_dir" -xf -

# Mise runs as the normal Arch user, so its retained downloads can be exported
# directly without copying any installed tools.
vm_step "Saving mise download cache"
ssh_vm \
  'tar -C "$HOME/.local/share/mise/downloads" -cf - .' |
  tar -C "$mise_cache_dir" -xf -
