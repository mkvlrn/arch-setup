package checks

import (
	"context"
	"fmt"
	"strings"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// Fonts returns a check for the Nerd Fonts managed by getnf.
func Fonts(fonts []string) setup.Check {
	return setup.Check{
		Name: "Verify installed Nerd Fonts",
		Run: func(ctx context.Context) error {
			results, err := shell.Run(ctx, []shell.Command{
				{
					Name: "list installed Nerd Fonts",
					Path: "getnf",
					Args: []string{"-l"},
				},
			})
			if err != nil {
				return fmt.Errorf("list installed Nerd Fonts: %w", err)
			}

			installed := make(map[string]struct{})
			for _, field := range strings.Fields(results[0].Stdout) {
				installed[field] = struct{}{}
			}

			var missing []string

			for _, font := range fonts {
				if _, ok := installed[font]; !ok {
					missing = append(missing, font)
				}
			}

			if len(missing) > 0 {
				return fmt.Errorf("fonts are not installed: %s", strings.Join(missing, ", "))
			}

			return nil
		},
	}
}
