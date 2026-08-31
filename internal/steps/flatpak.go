package steps

import (
	"fmt"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

const (
	flathubURL = "https://dl.flathub.org/repo/flathub.flatpakrepo"
	cosmicURL  = "https://apt.pop-os.org/cosmic/cosmic.flatpakrepo"
)

// Flatpak configures Flatpak remotes and installs COSMIC applets.
func Flatpak(applets []string) setup.Step {
	installArgs := append(
		[]string{
			"install",
			"--user",
			"--noninteractive",
			"--assumeyes",
			"cosmic",
		},
		applets...,
	)

	return setup.Step{
		Name: fmt.Sprintf("Configure Flatpak and install %d applets", len(applets)),
		Commands: []shell.Command{
			{
				Name: "add Flathub remote",
				Path: "flatpak",
				Args: []string{"remote-add", "--user", "--if-not-exists", "flathub", flathubURL},
			},
			{
				Name: "add COSMIC remote",
				Path: "flatpak",
				Args: []string{"remote-add", "--user", "--if-not-exists", "cosmic", cosmicURL},
			},
			{
				Name: "install COSMIC applets",
				Path: "flatpak",
				Args: installArgs,
			},
		},
	}
}
