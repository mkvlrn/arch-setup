package execute

import (
	"fmt"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// RemovePkg returns a package-removal step.
func RemovePkg(packages []string) setup.Step {
	args := append([]string{"-Rns", "--noconfirm"}, packages...)

	return setup.Step{
		Name: fmt.Sprintf("Removing %d unused packages", len(packages)),
		Commands: []shell.Command{{
			Name: "remove unused packages",
			Path: "pacman",
			Args: args,
			Sudo: true,
		}},
	}
}
