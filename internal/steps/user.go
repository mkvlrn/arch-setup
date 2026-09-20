package steps

import (
	"os"
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
	devKey := filepath.Join(homeDir, ".ssh", "dev")
	cbKey := filepath.Join(homeDir, ".ssh", "cb")

	userSteps := setup.Step{
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
				Name: "allow user dir to be executable",
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
				Name: "start paccache service",
				Path: "systemctl",
				Args: []string{"enable", "--now", "paccache.timer"},
				Sudo: true,
			},
			{
				Name: "reload services daemon",
				Path: "systemctl",
				Args: []string{"--user", "daemon-reload"},
			},
			{
				Name: "start ssh-agent service",
				Path: "systemctl",
				Args: []string{"--user", "enable", "--now", "ssh-agent.service"},
			},
			{
				Name: "start proton drive rclone service",
				Path: "systemctl",
				Args: []string{"--user", "enable", "--now", "proton-drive.service"},
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
				Args: []string{
					"-c",
					`"$1" exec -- gh completion -s fish > "$2"`,
					"sh",
					misePath,
					filepath.Join(baseCompletion, "gh.fish"),
				},
			},
		},
	}

	if os.Getenv("SSH_AUTH_SOCK") != "" {
		userSteps.Commands = append(userSteps.Commands, []shell.Command{
			{
				Name: "add dev key to agent for the first time",
				Path: "ssh-add",
				Args: []string{"-q", devKey, "</dev/null"},
			},
			{
				Name: "add cb key to agent for the first time",
				Path: "ssh-add",
				Args: []string{"-q", cbKey, "</dev/null"},
			},
		}...)
	}

	return userSteps
}
