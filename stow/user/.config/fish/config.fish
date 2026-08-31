# no greeting
set -g fish_greeting

# for ghostty, via cosmic
set -gx GTK_IM_MODULE simple

# initially for crontab, maybe others, idk
set -gx VISUAL "code --wait"
set -gx EDITOR "code --wait"

# paths
set -g HOME_BIN "$HOME/.local/bin"
set -g USR_LOCAL_BIN /usr/local/bin
set -gx GOPATH "$HOME/.go"

# PATH
fish_add_path --prepend $HOME_BIN $USR_LOCAL_BIN $GOPATH/bin

# mise
mise activate fish | source

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
