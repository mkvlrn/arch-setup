package steps

import (
	"fmt"
	"strings"

	"github.com/mkvlrn/arch-setup/internal/config"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// Fonts installs nerd fonts with getnf.
func Fonts(cfg *config.Config) setup.Step {
	fonts := cfg.GetNF
	stepName := fmt.Sprintf("Installing %d fonts managed by getnf", len(fonts))
	args := []string{"-i", strings.Join(fonts, ",")}

	return setup.Step{
		Name: stepName,
		Commands: []shell.Command{
			{
				Name: "install fonts managed by getnf",
				Path: "getnf",
				Args: args,
				Env:  []string{"TERM=xterm-256color"},
			},
			{
				Name: "clear fonts cache",
				Path: "fc-cache",
				Args: []string{"-f"},
			},
		},
	}
}
