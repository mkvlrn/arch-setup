package execute

import (
	"path/filepath"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// Yay installs yay-bin and updates mirrors with reflector.
func Yay(tempDir string, mirrorListPath string) setup.Step {
	yaySrcDir := filepath.Join(tempDir, "yay-bin")

	return setup.Step{
		Name: "Install yay and update mirrors",
		Commands: []shell.Command{
			{
				Name: "clone yay-bin",
				Path: "git",
				Args: []string{"clone", "https://aur.archlinux.org/yay-bin", yaySrcDir},
			},
			{
				Name: "build yay",
				Path: "makepkg",
				Args: []string{"-si", "--noconfirm"},
				Dir:  yaySrcDir,
			},
			{
				Name: "track git packages",
				Path: "yay",
				Args: []string{"-Y", "--gendb"},
			},
			{
				Name: "enable dev packages updates",
				Path: "yay",
				Args: []string{"-Y", "--devel", "--save"},
			},
			{
				Name: "get best mirrors list",
				Path: "reflector",
				Args: []string{
					"--latest", "20",
					"--protocol", "https",
					"--sort", "rate",
					"--save", mirrorListPath,
				},
				Sudo: true,
			},
			{
				Name: "update package data",
				Path: "yay",
				Args: []string{"-Syu", "--noconfirm"},
			},
			{
				Name: "remove debug packages",
				Path: "sh",
				Args: []string{"-c", "yay -Qq | grep -- '-debug$' | xargs -r yay -Rnsu"},
			},
		},
	}
}
