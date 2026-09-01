function upgrayedd --description "Update system, mise, tools, and completions"
    yay
    and flatpak update -y
    and mise self-update -y
    and mise cache clear
    and mise upgrade -b
    and mise completion fish >~/.config/fish/completions/mise.fish
    and gh completion -s fish >~/.config/fish/completions/gh.fish
    and glab completion -s fish >~/.config/fish/completions/glab.fish
end
