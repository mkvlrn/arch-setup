package execute

import (
	"path/filepath"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// User configures settings needed for normal usage after install.
//
//nolint:funlen // The function is long because it declaratively lists setup commands.
func User(username string, homeDir string) setup.Step {
	baseCompletion := filepath.Join(homeDir, ".config", "fish", "completions")
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
				Name: "start paccache service",
				Path: "systemctl",
				Args: []string{"enable", "--now", "paccache.timer"},
				Sudo: true,
			},
			{
				Name: "create completions directory",
				Path: "mkdir",
				Args: []string{"-p", baseCompletion},
			},
			{
				Name: "generate mise completions",
				Path: "sh",
				Args: []string{"-c", `"$1" completion fish > "$2"`, "sh", misePath, filepath.Join(baseCompletion, "mise.fish")},
			},
			{
				Name: "generate gh completions",
				Path: "sh",
				Args: []string{"-c", `"$1" completion -s fish > "$2"`, "sh", ghPath, filepath.Join(baseCompletion, "gh.fish")},
			},
			{
				Name: "generate glab completions",
				Path: "sh",
				Args: []string{"-c", `"$1" completion -s fish > "$2"`, "sh", glabPath, filepath.Join(baseCompletion, "glab.fish")},
			},
		},
	}
}
