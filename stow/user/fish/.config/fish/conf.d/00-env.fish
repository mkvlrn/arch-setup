# initially for crontab, maybe others, idk
set -gx VISUAL "code --wait"
set -gx EDITOR "code --wait"
# paths
set -gx HOME_BIN "$HOME/.local/bin"
set -gx USR_LOCAL_BIN /usr/local/bin
# PATH
fish_add_path --prepend $HOME_BIN $USR_LOCAL_BIN
