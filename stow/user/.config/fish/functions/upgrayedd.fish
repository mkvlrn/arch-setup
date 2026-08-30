function upgrayedd --description "Update system, mise, tools, and completions"
    yay
    and mise self-update -y
    and mise cache clear
    and mise upgrade -b
    and glab completion -s fish >~/.config/fish/completions/glab.fish
    and mise completion fish >~/.config/fish/completions/mise.fish
    and gh completion -s fish >~/.config/fish/completions/gh.fish
end
