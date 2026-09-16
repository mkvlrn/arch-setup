alias ztna-up="ssh -fN -p 2222 -D 1080 mkvlrn@localhost"
alias ztna-down="pkill -f 'ssh -fN -p 2222 -D 1080'"
set -gx ALL_PROXY "socks5h://127.0.0.1:1080"
