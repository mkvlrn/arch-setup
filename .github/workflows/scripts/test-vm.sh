#!/bin/sh

set -eu

mise_github_token="${MISE_GITHUB_TOKEN:?MISE_GITHUB_TOKEN is required}"

# Run a command inside the ephemeral Arch VM.
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

# Make VM preparation work visible in the Actions log.
vm_step() {
  printf '\n🖥️  %s\n' "$1"
}

# Test the exact executable produced by the check-build job.
vm_step "Copying installer to VM"
scp_vm \
  ./bin/arch-setup \
  arch@127.0.0.1:/tmp/arch-setup

# Copy the exact repository contents checked out for this PR, including Git
# metadata required by the user Stow step.
vm_step "Copying candidate repository to VM"
# shellcheck disable=SC2016
{
  printf '%s\0' .git
  git ls-files -z
} |
  tar --null --files-from=- -cf - |
  ssh_vm '
    mkdir -p "$HOME/repos/arch-setup"
    tar -C "$HOME/repos/arch-setup" -xf -
  '

# MISE_GITHUB_TOKEN belongs to the runner environment, so explicitly forward it
# to the installer process inside the VM.
vm_step "Running arch-setup"
ssh_vm \
  "chmod +x /tmp/arch-setup &&
   printf '%s\n' arch | sudo -S -v &&
   MISE_GITHUB_TOKEN='$mise_github_token' \
   MISE_ALWAYS_KEEP_DOWNLOAD=1 \
   CI=true \
   /tmp/arch-setup"

# Run verification in a new login session so changes such as supplementary
# group membership are visible.
vm_step "Verifying machine state"
ssh_vm \
  'CI=true /tmp/arch-setup --verify'
