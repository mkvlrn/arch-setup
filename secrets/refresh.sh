#!/bin/sh

set -eu

dir="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
output="$dir/secrets.tar.age"

printf '\nEncrypting machine secrets\n\n'

tar -C "$HOME" -cf - \
  --exclude='.ssh/agent' \
  .ssh \
  .aws \
  .config/fish/conf.d/secrets.fish |
  age -p -o "$output"

printf '\nSecrets written to %s\n' "$output"
