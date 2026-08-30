package steps

import (
	"path/filepath"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// User configures settings needed for normal usage after install.
func User(username string, homeDir string) setup.Step {
	misePath := filepath.Join(homeDir, ".local", "bin", "mise")
	ghPath := filepath.Join(homeDir, ".local", "share", "mise", "shims", "gh")
	glabPath := filepath.Join(homeDir, ".local", "share", "mise", "shims", "glab")

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
				Name: "set anonymous ftp user root dir",
				Path: "usermod",
				Args: []string{"-d", filepath.Join(homeDir, "torrents"), "ftp"},
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
			{
				Name: "generate mise completions",
				Path: misePath,
				Args: []string{"completion", "fish", ">~/.config/fish/completions/mise.fish"},
			},
			{
				Name: "generate gh completions",
				Path: ghPath,
				Args: []string{"completion", "-s", "fish", ">~/.config/fish/completions/gh.fish"},
			},
			{
				Name: "generate glab completions",
				Path: glabPath,
				Args: []string{"completion", "-s", "fish", ">~/.config/fish/completions/glab.fish"},
			},
		},
	}
}
