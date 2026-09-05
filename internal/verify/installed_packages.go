package verify

import (
	"context"
	"errors"
	"fmt"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// InstalledPackages returns the check for packages installed as base and main packages.
func InstalledPackages(packages []string) setup.Check {
	return setup.Check{
		Name: "Verify installed packages",
		Run: func(ctx context.Context) error {
			var failures []error

			for _, pkg := range packages {
				_, err := shell.Run(ctx, []shell.Command{
					{
						Name: "query package " + pkg,
						Path: "pacman",
						Args: []string{"-Q", pkg},
					},
				})
				if err != nil {
					failures = append(failures, fmt.Errorf("package %q is not installed: %w", pkg, err))
				}
			}

			return errors.Join(failures...)
		},
	}
}
