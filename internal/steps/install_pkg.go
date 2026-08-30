package steps

import (
	"fmt"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

type packageManager string

const (
	// UsePacman picks pacman, for early installs.
	UsePacman packageManager = "pacman"
	// UseYay picks yay, the superior wrapper.
	UseYay packageManager = "yay"
)

// InstallPkg returns a package installation step.
func InstallPkg(pm packageManager, packages []string) setup.Step {
	stepName := fmt.Sprintf("Installing %d packages with %s", len(packages), string(pm))
	operation := "-S"
	sudo := false

	if pm == UsePacman {
		operation = "-Syu"
		sudo = true
	}

	args := append([]string{operation, "--noconfirm", "--needed"}, packages...)

	return setup.Step{
		Name: stepName,
		Commands: []shell.Command{{
			Name: "install packages",
			Path: string(pm),
			Args: args,
			Sudo: sudo,
		}},
	}
}
