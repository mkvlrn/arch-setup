package steps

import (
	"path/filepath"

	"github.com/mkvlrn/arch-setup/internal/config"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// Yay installs yay-bin and updates mirrors with reflector.
//
//nolint:funlen // The function is long because it declaratively lists setup commands.
func Yay(cfg *config.Config) setup.Step {
	yaySrcDir := filepath.Join(cfg.Machine.TempDir, "yay-bin")

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
				Args: []string{"--noconfirm"},
				Dir:  yaySrcDir,
			},
			{
				Name: "install yay",
				Path: "sh",
				Args: []string{
					"-c",
					`set -eu
pacman -U --noconfirm "$1"/*.pkg.tar.zst`,
					"install-yay",
					yaySrcDir,
				},
				Sudo: true,
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
					"--save", cfg.Yay.MirrorListPath,
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
