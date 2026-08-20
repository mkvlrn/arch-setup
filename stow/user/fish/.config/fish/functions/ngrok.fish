function ngrok
    set -l port 3000
    if test -n "$argv[1]"
        set port $argv[1]
    end
    docker run -it --rm -e NGROK_AUTHTOKEN="$NGROK_AUTHTOKEN" --net=host ngrok/ngrok:latest http "$port"
end
