# mise
mise activate fish | source

# oh-my-posh
oh-my-posh init fish --config ~/repos/ts-tools/packages/config/src/mkvlrn.omp.jsonc | source

# keychain, .ssh
keychain --eval --quiet | source
ssh-add -q ~/.ssh/id_ed25519 </dev/null