package checks

import (
	"context"
	"path/filepath"

	"github.com/mkvlrn/arch-setup/internal/config"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// RcloneProton verifies that the Proton Drive remote is configured.
func RcloneProton(cfg *config.Config) setup.Check {
	misePath := filepath.Join(cfg.Machine.HomeDir, ".local", "bin", "mise")

	return setup.Check{
		Name: "Verify rclone Proton Drive configuration",
		Run: func(ctx context.Context) error {
			_, err := shell.Run(ctx, []shell.Command{{
				Name: "get rclone Proton Drive type",
				Path: "sh",
				Args: []string{
					"-c",
					`set -eu
"$1" exec -- rclone config show proton | grep -q '^type[[:space:]]*=[[:space:]]*protondrive$'`,
					"check-rclone-proton",
					misePath,
				},
			}})

			return err
		},
	}
}
