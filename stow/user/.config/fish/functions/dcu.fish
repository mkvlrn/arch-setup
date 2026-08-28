function dcu
    argparse 'p=' -- $argv
    or return

    set -l profile
    if set -q _flag_p
        set profile "--profile=$_flag_p"
    end

    docker compose $profile up -d $argv
end
