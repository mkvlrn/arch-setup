#!/bin/sh

# Copy the exact executable produced by the build job into the fresh Arch
# VM and run it directly.
#
# This deliberately does not use GitHub Pages or a GitHub Release, so the
# test always exercises the exact artifact produced by this workflow run.
set -eu

# Copy the built executable into the VM.
sshpass -p arch scp \
  -P 2222 \
  -o StrictHostKeyChecking=no \
  -o UserKnownHostsFile=/dev/null \
  ./bin/arch-setup \
  arch@127.0.0.1:/tmp/arch-setup

# Prime sudo for the VM's `arch` user so the non-interactive installer can
# perform privileged operations without waiting for a password prompt.
sshpass -p arch ssh \
  -p 2222 \
  -o StrictHostKeyChecking=no \
  -o UserKnownHostsFile=/dev/null \
  arch@127.0.0.1 \
  "printf '%s\n' arch | sudo -S -v"

# Run the exact binary produced by the build job.
sshpass -p arch ssh \
  -p 2222 \
  -o StrictHostKeyChecking=no \
  -o UserKnownHostsFile=/dev/null \
  arch@127.0.0.1 \
  'chmod +x /tmp/arch-setup && /tmp/arch-setup'
