#!/bin/sh

set -eu

cache_dir="${PACKAGE_CACHE_DIR:?PACKAGE_CACHE_DIR is required}"
pacman_cache_dir="$cache_dir/pacman"
yay_cache_dir="$cache_dir/yay"

ssh_vm() {
  sshpass -p arch ssh \
    -p 2222 \
    -o StrictHostKeyChecking=no \
    -o UserKnownHostsFile=/dev/null \
    arch@127.0.0.1 \
    "$@"
}

scp_vm() {
  sshpass -p arch scp \
    -P 2222 \
    -o StrictHostKeyChecking=no \
    -o UserKnownHostsFile=/dev/null \
    "$@"
}

mkdir -p \
  "$pacman_cache_dir" \
  "$yay_cache_dir"

# Copy the exact binary produced by this workflow.
scp_vm \
  ./bin/arch-setup \
  arch@127.0.0.1:/tmp/arch-setup

# Copy the post-install machine assertions.
scp_vm \
  ./.github/workflows/scripts/assert-vm.sh \
  arch@127.0.0.1:/tmp/assert-vm.sh

# Restore pacman's package cache into a directory owned by the VM user.
tar -C "$pacman_cache_dir" -cf - . |
  ssh_vm '
    mkdir -p "$HOME/.cache/pacman/pkg"
    tar -C "$HOME/.cache/pacman/pkg" -xf -
  '

# Make pacman transparently use the restored directory as its normal cache.
ssh_vm '
  printf "%s\n" arch |
    sudo -S mount --bind \
      "$HOME/.cache/pacman/pkg" \
      /var/cache/pacman/pkg
'

# Restore yay's normal build/cache directory.
tar -C "$yay_cache_dir" -cf - . |
  ssh_vm '
    mkdir -p "$HOME/.cache/yay"
    tar -C "$HOME/.cache/yay" -xf -
  '

# Prime sudo.
ssh_vm "printf '%s\n' arch | sudo -S -v"

# Run the installer.
ssh_vm \
  "chmod +x /tmp/arch-setup && MISE_GITHUB_TOKEN='$MISE_GITHUB_TOKEN' /tmp/arch-setup"

# Verify the resulting machine state.
ssh_vm \
  'chmod +x /tmp/assert-vm.sh && /tmp/assert-vm.sh'

# Keep the pacman cache from growing indefinitely. This retains packages
# useful to the resulting installed system while discarding obsolete ones.
ssh_vm \
  "printf '%s\n' arch | sudo -S pacman -Sc --noconfirm"

# Replace the host-side cache with the state produced by this VM run.
rm -rf \
  "$pacman_cache_dir" \
  "$yay_cache_dir"

mkdir -p \
  "$pacman_cache_dir" \
  "$yay_cache_dir"

ssh_vm \
  'tar -C "$HOME/.cache/pacman/pkg" -cf - .' |
  tar -C "$pacman_cache_dir" -xf -

ssh_vm \
  'tar -C "$HOME/.cache/yay" -cf - .' |
  tar -C "$yay_cache_dir" -xf -
