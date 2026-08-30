package steps

import (
	"path/filepath"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// User configures settings needed for normal usage after install.
func User(username string, homeDir string) setup.Step {
	ftpRootDir := filepath.Join(homeDir, "torrents")
	fishComplDir := filepath.Join(homeDir, ".config", "fish", "completions")

	return setup.Step{
		Name: "Config misc user settings",
		Commands: []shell.Command{
			{
				Name: "set user shell",
				Path: "chsh",
				Args: []string{"-s", "/usr/bin/fish", username},
				Sudo: true,
			},
			{
				Name: "add user to docker group",
				Path: "usermod",
				Args: []string{"-aG", "docker", username},
				Sudo: true,
			},
			{
				Name: "create fish compl dir",
				Path: "mkdir",
				Args: []string{"-p", fishComplDir},
			},
			{
				Name: "set anonymous ftp user root dir",
				Path: "usermod",
				Args: []string{"-d", ftpRootDir, "ftp"},
				Sudo: true,
			},
			{
				Name: "allow ftp user to traverse to download dir",
				Path: "chmod",
				Args: []string{"o+x", homeDir},
			},
			{
				Name: "start docker service",
				Path: "systemctl",
				Args: []string{"enable", "--now", "docker.socket"},
				Sudo: true,
			},
			{
				Name: "start pure-ftpd service",
				Path: "systemctl",
				Args: []string{"enable", "--now", "pure-ftpd.service"},
				Sudo: true,
			},
		},
	}
}
