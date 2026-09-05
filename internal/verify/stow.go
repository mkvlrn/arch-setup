package verify

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/mkvlrn/arch-setup/internal/setup"
)

type stowPackage string

const (
	// StowSystem selects the system Stow package.
	StowSystem stowPackage = "system"
	// StowUser selects the user Stow package.
	StowUser stowPackage = "user"
)

// Stow returns the check for a stowed package.
func Stow(pkg stowPackage, repoDir string, homeDir string) setup.Check {
	return setup.Check{
		Name: fmt.Sprintf("Verify stowed %s files", pkg),
		Run: func(_ context.Context) error {
			sourceRoot := filepath.Join(repoDir, "stow", string(pkg))

			return verifyStowTree(sourceRoot, stowTarget(pkg, homeDir))
		},
	}
}

func verifyStowTree(sourceRoot string, targetRoot string) error {
	var failures []error

	ignoreStow := []string{".gitkeep", ".stow-local-ignore"}

	err := filepath.WalkDir(
		sourceRoot,
		func(source string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			if entry.IsDir() {
				return nil
			}

			if slices.Contains(ignoreStow, entry.Name()) {
				return nil
			}

			if err := verifyStowEntry(sourceRoot, targetRoot, source); err != nil {
				failures = append(failures, err)
			}

			return nil
		},
	)
	if err != nil {
		failures = append(
			failures,
			fmt.Errorf("walk stow package %q: %w", sourceRoot, err),
		)
	}

	return errors.Join(failures...)
}

func verifyStowEntry(sourceRoot string, targetRoot string, source string) error {
	relative, err := filepath.Rel(sourceRoot, source)
	if err != nil {
		return fmt.Errorf("get relative path for %q: %w", source, err)
	}

	destination := filepath.Join(targetRoot, relative)

	return verifyStowLink(source, destination)
}

func verifyStowLink(source string, destination string) error {
	info, err := os.Lstat(destination)
	if err != nil {
		return fmt.Errorf("inspect %q: %w", destination, err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("%q is not a symlink", destination)
	}

	sourceReal, err := resolvePath(source)
	if err != nil {
		return fmt.Errorf("resolve source %q: %w", source, err)
	}

	destinationReal, err := resolvePath(destination)
	if err != nil {
		return fmt.Errorf("resolve destination %q: %w", destination, err)
	}

	if sourceReal != destinationReal {
		return fmt.Errorf(
			"%q points to %q instead of %q",
			destination,
			destinationReal,
			sourceReal,
		)
	}

	return nil
}

func resolvePath(path string) (string, error) {
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}

	return filepath.Abs(resolved)
}

func stowTarget(dest stowPackage, homeDir string) string {
	switch dest {
	case StowSystem:
		return "/"
	case StowUser:
		return homeDir
	default:
		panic(fmt.Sprintf("unknown stow destination %q", dest))
	}
}
