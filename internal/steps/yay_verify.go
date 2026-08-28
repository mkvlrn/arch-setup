package steps

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// YayVerify returns the check for the Yay installation and mirror configuration.
func YayVerify(mirrorListPath string, mirrorListCheck string) setup.Check {
	return setup.Check{
		Name: "Verify Yay and mirrors",
		Run: func(ctx context.Context) error {
			return errors.Join(
				verifyYayInstalled(ctx),
				verifyNoDebugPackages(ctx),
				verifyMirrorlist(mirrorListPath, mirrorListCheck),
			)
		},
	}
}

func verifyYayInstalled(ctx context.Context) error {
	_, err := shell.Run(ctx, []shell.Command{
		{
			Name: "get Yay version",
			Path: "yay",
			Args: []string{"--version"},
		},
	})
	if err != nil {
		return fmt.Errorf("yay is not available: %w", err)
	}

	return nil
}

func verifyNoDebugPackages(ctx context.Context) error {
	results, err := shell.Run(ctx, []shell.Command{
		{
			Name: "list installed packages",
			Path: "yay",
			Args: []string{"-Qq"},
		},
	})
	if err != nil {
		return fmt.Errorf("list installed packages: %w", err)
	}

	var debugPackages []string

	for packageName := range strings.SplitSeq(strings.TrimSpace(results[0].Stdout), "\n") {
		if strings.HasSuffix(packageName, "-debug") {
			debugPackages = append(debugPackages, packageName)
		}
	}

	if len(debugPackages) > 0 {
		return fmt.Errorf("debug packages are installed: %s", strings.Join(debugPackages, ", "))
	}

	return nil
}

func verifyMirrorlist(mirrorListPath string, mirrorListCheck string) error {
	// #nosec G304 -- mirrorListPath is a config-controlled string
	content, err := os.ReadFile(mirrorListPath)
	if err != nil {
		return fmt.Errorf("read mirrorlist: %w", err)
	}

	if !strings.Contains(string(content), mirrorListCheck) {
		return errors.New("mirrorlist was not generated with the expected Reflector command")
	}

	return nil
}
