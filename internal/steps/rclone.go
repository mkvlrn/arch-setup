package steps

import (
	"path/filepath"

	"github.com/mkvlrn/arch-setup/internal/config"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// RcloneProton sets up the rclone configuration for Proton Drive.
func RcloneProton(cfg *config.Config) setup.Step {
	if cfg.Env.CI {
		return setup.Step{Name: "Configure rclone Proton Drive (skipped in CI)"}
	}

	secretsPath := filepath.Join(cfg.Machine.HomeDir, ".config", "rclone", "proton-secrets")

	return setup.Step{
		Name: "Configure rclone Proton Drive",
		Commands: []shell.Command{{
			Name: "configure rclone Proton Drive",
			Path: "sh",
			Args: []string{
				"-c",
				`set -eu
secret_file="$1"
read_secret() {
    sed -n "s/^$1[[:space:]]*=[[:space:]]*//p" "$secret_file" |
        sed "s/^['\"]//; s/['\"]$//"
}
username=$(read_secret username)
totp=$(read_secret totp)
test -n "$username"
test -n "$totp"
otp_secret=$(rclone obscure "$totp")
exec rclone config create proton protondrive \
    "username=$username" \
    "otp_secret=$otp_secret"`,
				"configure-rclone-proton",
				secretsPath,
			},
		}},
	}
}
