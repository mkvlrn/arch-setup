#!/bin/sh

# Prepare and boot a disposable Arch Linux VM using QEMU/KVM.
#
# The VM runs in the background and remains alive for the following
# workflow steps in this job.
set -eu

# Install the host-side tools required to run and access the VM.
sudo apt-get update
sudo apt-get install -y \
  qemu-system-x86 \
  qemu-utils \
  sshpass

# Download Arch's current official basic VM image.
curl -fsSL \
  https://geo.mirror.pkgbuild.com/images/latest/Arch-Linux-x86_64-basic.qcow2 \
  -o arch.qcow2

# Create a disposable copy-on-write disk backed by the pristine image.
# Changes made by arch-setup go into this overlay only.
qemu-img create \
  -f qcow2 \
  -F qcow2 \
  -b "$PWD/arch.qcow2" \
  arch-test.qcow2

# Allow the runner user to access hardware virtualization.
sudo chmod 666 /dev/kvm

# Boot Arch in the background.
#
# Host port 2222 is forwarded to the VM's SSH port 22.
qemu-system-x86_64 \
  -enable-kvm \
  -cpu host \
  -m 2G \
  -smp 2 \
  -drive file=arch-test.qcow2,format=qcow2 \
  -nic user,hostfwd=tcp::2222-:22 \
  -nographic \
  >qemu.log 2>&1 &

echo "$!" >qemu.pid

# Wait until an actual SSH login succeeds.
#
# Testing the TCP port alone is insufficient because QEMU can expose the
# forwarded port before sshd inside the guest is ready.
for _ in $(seq 1 60); do
  if sshpass -p arch ssh \
    -p 2222 \
    -o StrictHostKeyChecking=no \
    -o UserKnownHostsFile=/dev/null \
    -o ConnectTimeout=2 \
    arch@127.0.0.1 \
    'true' \
    >/dev/null 2>&1; then
    exit 0
  fi

  sleep 2
done

# Surface the VM console if boot/SSH readiness failed.
cat qemu.log
exit 1
