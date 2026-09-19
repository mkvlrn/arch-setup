# no greeting
set -g fish_greeting

# mise
$HOME/.local/bin/mise activate fish | source

# interactive
if status is-interactive
    # oh-my-posh
    set omp_config ~/repos/ts-tools/packages/config/src/mkvlrn.omp.jsonc
    test -f $omp_config; or set omp_config https://raw.githubusercontent.com/mkvlrn/ts-tools/main/packages/config/src/mkvlrn.omp.jsonc
    oh-my-posh init fish --config $omp_config | source

    # aliases
    # zed
    # alias zed zeditor
    # eza to ls
    alias ls 'eza --git --git-repos --group-directories-first'
    # repo eza to k
    alias k 'eza -la --git --git-repos --group-directories-first'
    # gpa_vm ssh
    alias gpa-vm 'ssh -i ~/.ssh/gpa_vm -p 22220 mkvlrn@127.0.0.1'
end
