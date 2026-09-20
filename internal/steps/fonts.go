package steps

import (
	"fmt"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// Fonts installs nerd fonts with getnf.
func Fonts(fonts []string) setup.Step {
	stepName := fmt.Sprintf("Installing %d fonts managed by getnf", len(fonts))
	args := append([]string{"-i"}, fonts...)

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
