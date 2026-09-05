package verify

import (
	"context"
	"errors"
	"fmt"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// RemovedPackages returns the check for packages that should not be installed.
func RemovedPackages(packages []string) setup.Check {
	return setup.Check{
		Name: "Verify removed packages",
		Run: func(ctx context.Context) error {
			var failures []error

			for _, pkg := range packages {
				if _, err := shell.Run(ctx, []shell.Command{
					{
						Name: "query package " + pkg,
						Path: "pacman",
						Args: []string{"-Q", pkg},
					},
				}); err == nil {
					failures = append(
						failures,
						fmt.Errorf("package %q is still installed", pkg),
					)
				}
			}

			return errors.Join(failures...)
		},
	}
}
