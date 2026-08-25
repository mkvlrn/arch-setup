# initially for crontab, maybe others, idk
set -gx VISUAL "zed --wait"
set -gx EDITOR "zed --wait"
# paths
set -gx HOME_BIN "$HOME/.local/bin"
set -gx USR_LOCAL_BIN /usr/local/bin
set -gx GOPATH "$HOME/.go"
# PATH
fish_add_path --prepend $HOME_BIN $USR_LOCAL_BIN $GOPATH/bin
