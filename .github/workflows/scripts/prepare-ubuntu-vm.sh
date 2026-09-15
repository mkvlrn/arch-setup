#!/bin/sh

# Boot a disposable Ubuntu 26.04 cloud-image VM for installer tests.
set -eu

sudo apt-get update
sudo apt-get install -y qemu-system-x86 qemu-utils cloud-image-utils sshpass

curl -fsSL \
  https://cloud-images.ubuntu.com/resolute/current/resolute-server-cloudimg-amd64.img \
  -o ubuntu.qcow2

cat >ubuntu-user-data <<'EOF'
#cloud-config
users:
  - name: mkvlrn
    gecos: CI User
    groups: [adm, sudo]
    sudo: ALL=(ALL) NOPASSWD:ALL
    shell: /bin/bash
    lock_passwd: false
chpasswd:
  list: |
    mkvlrn: ubuntu
  expire: false
ssh_pwauth: true
growpart:
  mode: auto
resize_rootfs: true
package_update: true
packages:
  - openssh-server
EOF

cat >ubuntu-meta-data <<'EOF'
instance-id: os-setup-ubuntu
local-hostname: os-setup-ubuntu
EOF

cloud-localds ubuntu-seed.img ubuntu-user-data ubuntu-meta-data
qemu-img create -f qcow2 -F qcow2 -b ubuntu.qcow2 ubuntu-test.qcow2 32G

qemu-system-x86_64 \
  -enable-kvm \
  -cpu host \
  -m 4G \
  -smp 2 \
  -drive file=ubuntu-test.qcow2,format=qcow2 \
  -drive file=ubuntu-seed.img,format=raw \
  -nic user,hostfwd=tcp::2223-:22 \
  -nographic \
  >ubuntu-vm.log 2>&1 &
printf '%s\n' "$!" >ubuntu-qemu.pid

for _ in $(seq 1 90); do
  if sshpass -p ubuntu ssh \
    -p 2223 \
    -o StrictHostKeyChecking=no \
    -o UserKnownHostsFile=/dev/null \
    -o ConnectTimeout=2 \
    mkvlrn@127.0.0.1 \
    true 2>/dev/null; then
    exit 0
  fi
  sleep 2
done

printf 'Ubuntu VM did not become ready.\n' >&2
exit 1
