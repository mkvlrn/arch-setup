package steps

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// FlatpakVerify returns checks for Flatpak remotes and installed applets.
func FlatpakVerify(applets []string) setup.Check {
	return setup.Check{
		Name: "Verify Flatpak configuration",
		Run: func(ctx context.Context) error {
			return errors.Join(
				verifyFlatpakRemotes(ctx),
				verifyFlatpakApplets(ctx, applets),
			)
		},
	}
}

func verifyFlatpakRemotes(ctx context.Context) error {
	results, err := shell.Run(ctx, []shell.Command{
		{
			Name: "list Flatpak remotes",
			Path: "flatpak",
			Args: []string{"remotes", "--user", "--columns=name"},
		},
	})
	if err != nil {
		return err
	}

	remotes := strings.Fields(results[0].Stdout)

	var failures []error

	for _, expected := range []string{"flathub", "cosmic"} {
		found := slices.Contains(remotes, expected)

		if !found {
			failures = append(
				failures,
				fmt.Errorf("Flatpak remote %q is missing", expected),
			)
		}
	}

	return errors.Join(failures...)
}

func verifyFlatpakApplets(ctx context.Context, applets []string) error {
	var failures []error

	for _, applet := range applets {
		_, err := shell.Run(ctx, []shell.Command{
			{
				Name: "inspect Flatpak applet " + applet,
				Path: "flatpak",
				Args: []string{"info", "--user", applet},
			},
		})
		if err != nil {
			failures = append(
				failures,
				fmt.Errorf("Flatpak applet %q is not installed: %w", applet, err),
			)
		}
	}

	return errors.Join(failures...)
}
