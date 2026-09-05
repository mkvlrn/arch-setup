package verify

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mkvlrn/arch-setup/internal/setup"
)

// Xdg returns the check for the configured XDG directories.
func Xdg(mkdir []string, rmrf []string, homeDir string) setup.Check {
	return setup.Check{
		Name: "Verify XDG user directories",
		Run: func(_ context.Context) error {
			return errors.Join(
				verifyDirectoriesExist(homeDir, mkdir),
				verifyDirectoriesRemoved(homeDir, rmrf),
			)
		},
	}
}

func verifyDirectoriesExist(homeDir string, directories []string) error {
	var failures []error

	for _, directory := range directories {
		path := filepath.Join(homeDir, directory)

		info, err := os.Stat(path)
		if err != nil {
			failures = append(failures, fmt.Errorf("inspect directory %q: %w", path, err))

			continue
		}

		if !info.IsDir() {
			failures = append(failures, fmt.Errorf("%q is not a directory", path))
		}
	}

	return errors.Join(failures...)
}

func verifyDirectoriesRemoved(homeDir string, directories []string) error {
	var failures []error

	for _, directory := range directories {
		path := filepath.Join(homeDir, directory)
		_, err := os.Stat(path)

		switch {
		case err == nil:
			failures = append(failures, fmt.Errorf("%q still exists", path))

		case !errors.Is(err, os.ErrNotExist):
			failures = append(failures, fmt.Errorf("inspect path %q: %w", path, err))
		}
	}

	return errors.Join(failures...)
}
