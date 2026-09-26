package checks

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mkvlrn/arch-setup/internal/config"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// User returns the check for miscellaneous user settings.
func User(cfg *config.Config) setup.Check {
	username, homeDir := cfg.Machine.Username, cfg.Machine.HomeDir

	return setup.Check{
		Name: "Verify miscellaneous user settings",
		Run: func(ctx context.Context) error {
			failures := []error{
				verifyUserShell(ctx, username),
				verifyDockerGroup(ctx, username),
				verifyHomeTraversal(homeDir),
				verifySystemdUnit(ctx, "docker.socket"),
				verifySystemdUnit(ctx, "paccache.timer"),
				verifyUserSystemdUnit(ctx, "ssh-agent.service"),
				verifyUserSystemdUnit(ctx, "proton-drive.service"),
			}

			failures = append(failures, verifyCompletions(homeDir))

			return errors.Join(failures...)
		},
	}
}

func verifyUserShell(ctx context.Context, username string) error {
	entry, err := passwdEntry(ctx, username)
	if err != nil {
		return err
	}

	const shellField = 6

	if entry[shellField] != "/usr/bin/zsh" {
		return fmt.Errorf(
			"user %q has shell %q instead of %q",
			username,
			entry[shellField],
			"/usr/bin/zsh",
		)
	}

	return nil
}

func verifyDockerGroup(ctx context.Context, username string) error {
	results, err := shell.Run(ctx, []shell.Command{
		{
			Name: "get user groups",
			Path: "id",
			Args: []string{"-nG", username},
		},
	})
	if err != nil {
		return err
	}

	if slices.Contains(strings.Fields(results[0].Stdout), "docker") {
		return nil
	}

	return fmt.Errorf("user %q is not in the docker group", username)
}

func passwdEntry(ctx context.Context, username string) ([]string, error) {
	results, err := shell.Run(ctx, []shell.Command{
		{
			Name: "get passwd entry for " + username,
			Path: "getent",
			Args: []string{"passwd", username},
		},
	})
	if err != nil {
		return nil, err
	}

	const expectedLength = 7

	entry := strings.Split(strings.TrimSpace(results[0].Stdout), ":")
	if len(entry) != expectedLength {
		return nil, fmt.Errorf("unexpected passwd entry for %q", username)
	}

	return entry, nil
}

func verifyHomeTraversal(homeDir string) error {
	info, err := os.Stat(homeDir)
	if err != nil {
		return fmt.Errorf("inspect home directory %q: %w", homeDir, err)
	}

	if info.Mode().Perm()&0o001 == 0 {
		return fmt.Errorf("%q is not traversable by other users", homeDir)
	}

	return nil
}

func verifySystemdUnit(ctx context.Context, unit string) error {
	_, err := shell.Run(ctx, []shell.Command{
		{
			Name: "check that " + unit + " is enabled",
			Path: "systemctl",
			Args: []string{"is-enabled", "--quiet", unit},
		},
		{
			Name: "check that " + unit + " is active",
			Path: "systemctl",
			Args: []string{"is-active", "--quiet", unit},
		},
	})
	if err != nil {
		return fmt.Errorf("verify systemd unit %q: %w", unit, err)
	}

	return nil
}

func verifyUserSystemdUnit(ctx context.Context, unit string) error {
	_, err := shell.Run(ctx, []shell.Command{
		{
			Name: "check that user " + unit + " is enabled",
			Path: "systemctl",
			Args: []string{"--user", "is-enabled", "--quiet", unit},
		},
		{
			Name: "check that user " + unit + " is active",
			Path: "systemctl",
			Args: []string{"--user", "is-active", "--quiet", unit},
		},
	})
	if err == nil {
		return nil
	}

	status, statusErr := shell.Run(ctx, []shell.Command{
		{
			Name: "get user " + unit + " status",
			Path: "sh",
			Args: []string{"-c", "systemctl --user status --no-pager --full \"$1\" || true", "get-user-unit-status", unit},
		},
		{
			Name: "get user " + unit + " journal",
			Path: "sh",
			Args: []string{"-c", "journalctl --user -u \"$1\" -n 50 --no-pager || true", "get-user-unit-journal", unit},
		},
	})
	if statusErr == nil {
		return fmt.Errorf("verify user systemd unit %q: %w\n%s\n%s", unit, err, status[0].Stdout, status[1].Stdout)
	}

	return fmt.Errorf("verify user systemd unit %q: %w", unit, err)
}

func verifyCompletions(homeDir string) error {
	base := filepath.Join(homeDir, ".config", "zsh", "completions")

	for _, file := range []string{"mise", "gh"} {
		if _, err := os.Stat(filepath.Join(base, "_"+file)); err != nil {
			return fmt.Errorf("completion file for %s not generated", file)
		}
	}

	return nil
}
