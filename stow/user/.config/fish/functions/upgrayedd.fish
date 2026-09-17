function upgrayedd --description "Update system, mise, tools, completions, and fonts"
    set -l old_pwd $PWD
    cd ~
    or return

    # yay
    sudo apt update
    and sudo apt upgrade
    and mise self-update -y; or mise self-update -y
    and mise prune -y
    and mise cache clear
    and mise upgrade -b
    and flatpak update -u
    and mise completion fish >~/.config/fish/completions/mise.fish
    and gh completion -s fish >~/.config/fish/completions/gh.fish
    and curl -fsSL https://raw.githubusercontent.com/getnf/getnf/main/install.sh | bash
    and getnf -U

    cd $old_pwd
end
