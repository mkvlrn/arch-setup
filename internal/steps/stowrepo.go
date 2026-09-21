package steps

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mkvlrn/arch-setup/internal/config"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

const (
	bootstrapKeyMode  os.FileMode = 0o600
	askPassScriptMode os.FileMode = 0o700
)

// StowRepo clones the private repository using the embedded, passphrase-protected key.
func StowRepo(cfg *config.Config, key, passphrase []byte) setup.Step {
	if cfg.Env.CI {
		return setup.Step{
			Name: "Use existing stow repository",
			Run: func(context.Context) error {
				if _, err := os.Stat(filepath.Join(cfg.Machine.StowRepoDir, ".git")); err != nil {
					return fmt.Errorf("find stow repository at %q: %w", cfg.Machine.StowRepoDir, err)
				}

				return nil
			},
		}
	}

	return setup.Step{
		Name: fmt.Sprintf("Clone %s to %s", cfg.Repo.StowSSH, cfg.Machine.StowRepoDir),
		Run: func(ctx context.Context) error {
			keyPath := filepath.Join(cfg.Machine.TempDir, "arch-stow-bootstrap")
			if err := os.WriteFile(keyPath, key, bootstrapKeyMode); err != nil {
				return fmt.Errorf("write bootstrap SSH key: %w", err)
			}
			defer os.Remove(keyPath) //nolint:errcheck // best-effort cleanup

			askPassPath := filepath.Join(cfg.Machine.TempDir, "arch-stow-askpass")
			askPassScript := []byte("#!/bin/sh\nprintf '%s\\n' \"$ARCH_STOW_PASSPHRASE\"\n")

			// #nosec G306 -- SSH_ASKPASS requires an executable helper.
			if err := os.WriteFile(askPassPath, askPassScript, askPassScriptMode); err != nil {
				return fmt.Errorf("write SSH askpass helper: %w", err)
			}
			defer os.Remove(askPassPath) //nolint:errcheck // best-effort cleanup

			sshCommand := fmt.Sprintf("ssh -i %s -o IdentitiesOnly=yes -o StrictHostKeyChecking=accept-new", keyPath)

			_, err := shell.Run(ctx, []shell.Command{{
				Name: "clone stow repository",
				Path: "git",
				Args: []string{"clone", cfg.Repo.StowSSH, cfg.Machine.StowRepoDir},
				Env: []string{
					"GIT_SSH_COMMAND=" + sshCommand,
					"SSH_ASKPASS=" + askPassPath,
					"SSH_ASKPASS_REQUIRE=force",
					"DISPLAY=:0",
					"ARCH_STOW_PASSPHRASE=" + string(passphrase),
				},
			}})
			if err != nil {
				return err
			}

			return nil
		},
	}
}
