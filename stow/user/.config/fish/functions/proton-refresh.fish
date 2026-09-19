function proton-refresh
    rclone rc vfs/refresh recursive=true 2>&1 | string match -r '"": "OK"' | string replace -r '.*"": "([^"]+)".*' '$1'
end
