package steps

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// UserVerify returns the check for miscellaneous user settings.
func UserVerify(username string, homeDir string) setup.Check {
	return setup.Check{
		Name: "Verify miscellaneous user settings",
		Run: func(ctx context.Context) error {
			return errors.Join(
				verifyUserShell(ctx, username),
				verifyDockerGroup(ctx, username),
				verifyFtpHome(ctx, homeDir),
				verifyHomeTraversal(homeDir),
				verifySystemdUnit(ctx, "docker.socket"),
				verifySystemdUnit(ctx, "pure-ftpd.service"),
				verifyCompletions(homeDir),
			)
		},
	}
}

func verifyUserShell(ctx context.Context, username string) error {
	entry, err := passwdEntry(ctx, username)
	if err != nil {
		return err
	}

	const shellField = 6

	if entry[shellField] != "/usr/bin/fish" {
		return fmt.Errorf(
			"user %q has shell %q instead of %q",
			username,
			entry[shellField],
			"/usr/bin/fish",
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

func verifyFtpHome(ctx context.Context, homeDir string) error {
	entry, err := passwdEntry(ctx, "ftp")
	if err != nil {
		return err
	}

	const homeField = 5

	expected := filepath.Join(homeDir, "torrents")

	if entry[homeField] != expected {
		return fmt.Errorf(
			"ftp home is %q instead of %q",
			entry[homeField],
			expected,
		)
	}

	return nil
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

func verifyCompletions(homeDir string) error {
	base := filepath.Join(homeDir, ".config", "fish", "completions")
	files := []string{"mise", "gh", "glab"}

	for _, file := range files {
		_, err := os.Stat(filepath.Join(base, file+".fish"))
		if err != nil {
			return fmt.Errorf("completion file for %s not generated", file)
		}
	}

	return nil
}
