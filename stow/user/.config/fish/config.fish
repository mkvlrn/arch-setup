# no greeting
set -g fish_greeting

# mise
"$HOME/.local/bin/mise" activate fish | source

# interactive
if status is-interactive
    # oh-my-posh
    set omp_config ~/repos/ts-tools/packages/config/src/mkvlrn.omp.jsonc
    test -f $omp_config; or set omp_config https://raw.githubusercontent.com/mkvlrn/ts-tools/main/packages/config/src/mkvlrn.omp.jsonc
    oh-my-posh init fish --config $omp_config | source

    # keychain, .ssh
    keychain --eval --quiet | source
    ssh-add -q ~/.ssh/dev </dev/null

    # aliases
    # zed
    alias zed zeditor
    # eza to ls
    alias ls eza
    # repo eza to k
    alias k 'eza -al --git --git-repos --group-directories-first'
end
