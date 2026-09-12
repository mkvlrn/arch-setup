function upgrayedd --description "Update system, mise, tools, and completions"
    set -l old_pwd $PWD
    cd ~
    or return

    yay
    and mise self-update -y; or mise self-update -y
    and mise prune -y
    and mise cache clear
    and mise upgrade -b
    and mise completion fish >~/.config/fish/completions/mise.fish
    and gh completion -s fish >~/.config/fish/completions/gh.fish
    and glab completion -s fish >~/.config/fish/completions/glab.fish

    cd $old_pwd
end
