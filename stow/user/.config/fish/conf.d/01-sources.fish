# mise
mise activate fish | source

# oh-my-posh
set omp_config ~/repos/ts-tools/packages/config/src/mkvlrn.omp.jsonc; test -f $omp_config; or set omp_config https://raw.githubusercontent.com/mkvlrn/ts-tools/main/packages/config/src/mkvlrn.omp.jsonc
oh-my-posh init fish --config $omp_config | source

# keychain, .ssh
if status is-interactive
    keychain --eval --quiet | source
    ssh-add -q ~/.ssh/dev </dev/null
end
