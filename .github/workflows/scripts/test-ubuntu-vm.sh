#!/bin/sh

# shellcheck disable=SC2016
set -eu

mise_github_token="${MISE_GITHUB_TOKEN:?MISE_GITHUB_TOKEN is required}"

ssh_vm() {
  sshpass -p ubuntu ssh \
    -p 2223 \
    -o StrictHostKeyChecking=no \
    -o UserKnownHostsFile=/dev/null \
    -o LogLevel=ERROR \
    mkvlrn@127.0.0.1 \
    "$@"
}

# Confirm the downloaded image is the intended Ubuntu LTS before changing it.
ssh_vm '. /etc/os-release && test "$VERSION_ID" = 26.04'

# Copy the exact repository contents, including Git metadata required by the
# repository verification step.
{
  printf '%s\0' .git
  git ls-files -z
} |
  tar --null --files-from=- -cf - |
  ssh_vm 'mkdir -p "$HOME/repos/arch-setup" && tar -C "$HOME/repos/arch-setup" -xf -'

ssh_vm \
  "MISE_GITHUB_TOKEN='$mise_github_token' \
   SETUP_REPO_DIR=\"\$HOME/repos/arch-setup\" \
   \"\$HOME/repos/arch-setup/ubuntu/install.sh\""

ssh_vm \
  'SETUP_REPO_DIR="$HOME/repos/arch-setup" \
   "$HOME/repos/arch-setup/ubuntu/verify.sh"'
