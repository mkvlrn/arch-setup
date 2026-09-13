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

# Make VM preparation work visible in the Actions log.
vm_step() {
  printf '\n🖥️  %s\n' "$1"
}

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
vm_step "Running install.sh"
ssh_vm \
  "chmod +x \"\$HOME/repos/arch-setup/install.sh\" \"\$HOME/repos/arch-setup/verify.sh\" &&
   printf '%s\n' arch | sudo -S -v &&
   MISE_GITHUB_TOKEN='$mise_github_token' \
   GITHUB_ACTIONS='$GITHUB_ACTIONS' \
   SETUP_REPO_DIR=\"\$HOME/repos/arch-setup\" \
   \"\$HOME/repos/arch-setup/install.sh\" --ci"

# Run verification in a new login session so changes such as supplementary
# group membership are visible.
vm_step "Verifying machine state"
# shellcheck disable=SC2016
ssh_vm \
  'SETUP_REPO_DIR="$HOME/repos/arch-setup" \
   ARCH_SETUP_EXPECTED_REVISION="$(git -C "$HOME/repos/arch-setup" rev-parse HEAD)" \
   "$HOME/repos/arch-setup/verify.sh" --ci'
