function doc
    set -l dockerFile "$HOME/.config/docker-compose.yml"

    if not test -f "$dockerFile"
        echo "Docker compose file not found: $dockerFile"
        return 1
    end

    switch $argv[1]
        case up
            docker compose -f "$dockerFile" up -d $argv[2..]
        case stop
            docker compose -f "$dockerFile" stop $argv[2..]
            docker compose -f "$dockerFile" rm -f $argv[2..]
        case down
            docker compose -f "$dockerFile" down
        case '*'
            echo "Usage: doc {up|stop|down} [options]"
            return 1
    end
end
