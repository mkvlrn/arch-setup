package checks

import (
	"context"
	"path/filepath"

	"github.com/mkvlrn/arch-setup/internal/config"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// RcloneProton verifies that the Proton Drive remote authenticates successfully.
func RcloneProton(cfg *config.Config) setup.Check {
	misePath := filepath.Join(cfg.Machine.HomeDir, ".local", "bin", "mise")

	return setup.Check{
		Name: "Verify rclone Proton Drive authentication",
		Run: func(ctx context.Context) error {
			_, err := shell.Run(ctx, []shell.Command{{
				Name: "list rclone Proton Drive root",
				Path: misePath,
				Args: []string{"exec", "--", "rclone", "lsd", "proton:"},
			}})

			return err
		},
	}
}
