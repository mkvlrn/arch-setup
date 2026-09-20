#!/bin/sh

set -eu

dir="$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)"
output="$dir/cmd/arch-setup/secrets.tar.age"

printf '\nEncrypting machine secrets\n\n'

tar -C "$HOME" -cf - \
  --exclude='.ssh/agent' \
  .ssh \
  .aws \
  .config/rclone/proton-secrets \
  .config/environment.d/10-secrets.conf |
  age -p -o "$output"

printf '\nSecrets written to %s\n' "$output"
