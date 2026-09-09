function gpa
    switch "$argv[1]"
        case on
            sudo systemctl start wapptunnel.service
            sudo systemctl start gpawatchdog.service

            if not pgrep -x GPAClient >/dev/null
                /opt/GuardicorePlatformAgent/gui/GPAClient >/dev/null 2>&1 &
                disown
            end

        case off
            pkill -TERM -x GPAClient 2>/dev/null

            sudo systemctl stop gpawatchdog.service
            sudo systemctl stop gpaservice.service
            sudo systemctl stop wapptunnel.service

        case status
            systemctl --no-pager --full status \
                wapptunnel.service \
                gpawatchdog.service \
                gpaservice.service

            if pgrep -x GPAClient >/dev/null
                echo "GPAClient: running"
            else
                echo "GPAClient: stopped"
            end

        case '*'
            echo "usage: gpa {on|off|status}"
            return 1
    end
end
