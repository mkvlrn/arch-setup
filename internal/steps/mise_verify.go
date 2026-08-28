package steps

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// MiseVerify returns a check for mise and all tools in its global manifest.
func MiseVerify(homeDir string) setup.Check {
	misePath := filepath.Join(homeDir, ".local", "bin", "mise")

	return setup.Check{
		Name: "Verify mise and its managed tools",
		Run: func(ctx context.Context) error {
			results, err := shell.Run(ctx, []shell.Command{
				{
					Name: "get mise version",
					Path: misePath,
					Args: []string{"--version"},
				},
				{
					Name: "list missing mise tools",
					Path: misePath,
					Args: []string{"ls", "--global", "--missing", "--no-header"},
				},
			})
			if err != nil {
				return err
			}

			missing := strings.TrimSpace(results[1].Stdout)
			if missing != "" {
				return fmt.Errorf("mise tools are missing:\n%s", missing)
			}

			return nil
		},
	}
}
