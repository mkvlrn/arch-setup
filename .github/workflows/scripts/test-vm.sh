#!/bin/sh

set -eu

# Copy the exact binary produced by this workflow.
sshpass -p arch scp \
  -P 2222 \
  -o StrictHostKeyChecking=no \
  -o UserKnownHostsFile=/dev/null \
  ./bin/arch-setup \
  arch@127.0.0.1:/tmp/arch-setup

# Copy the post-install machine assertions.
sshpass -p arch scp \
  -P 2222 \
  -o StrictHostKeyChecking=no \
  -o UserKnownHostsFile=/dev/null \
  ./.github/workflows/scripts/assert-vm.sh \
  arch@127.0.0.1:/tmp/assert-vm.sh

# Prime sudo.
sshpass -p arch ssh \
  -p 2222 \
  -o StrictHostKeyChecking=no \
  -o UserKnownHostsFile=/dev/null \
  arch@127.0.0.1 \
  "printf '%s\n' arch | sudo -S -v"

# Run the installer.
sshpass -p arch ssh \
  -p 2222 \
  -o StrictHostKeyChecking=no \
  -o UserKnownHostsFile=/dev/null \
  arch@127.0.0.1 \
  'chmod +x /tmp/arch-setup && /tmp/arch-setup'

# Verify the resulting machine state.
sshpass -p arch ssh \
  -p 2222 \
  -o StrictHostKeyChecking=no \
  -o UserKnownHostsFile=/dev/null \
  arch@127.0.0.1 \
  'chmod +x /tmp/assert-vm.sh && /tmp/assert-vm.sh'
