package steps

import (
	"fmt"
	"path/filepath"

	"github.com/mkvlrn/arch-setup/internal/config"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// Mise installs the mise binary and uses it to globally install the tools in its manifest.
func Mise(cfg *config.Config) setup.Step {
	homeDir, tools, settings := cfg.Machine.HomeDir, cfg.Mise.Tools, cfg.Mise.Settings
	stepName := fmt.Sprintf("Installing mise and %d tools", len(tools))
	misePath := filepath.Join(homeDir, ".local", "bin", "mise")
	settingsCmds := configSettings(settings, misePath)

	installArgs := append([]string{"install"}, tools...)
	environmentFile := filepath.Join(homeDir, ".config", "environment.d", "10-secrets.conf")

	commands := []shell.Command{
		{
			Name: "install mise",
			Path: "sh",
			Args: []string{
				"-c",
				`export MISE_INSTALL_PATH="$1"
curl https://mise.run | sh`,
				"install-mise",
				misePath,
			},
		},
	}
	commands = append(commands, settingsCmds...)
	commands = append(commands, shell.Command{
		Name: "install tools managed by mise",
		Path: "sh",
		Args: append([]string{
			"-c",
			`set -eu
			if [ -f "$1" ]; then
			    set -a
			    . "$1"
			    set +a
			fi
			shift
			exec "$@"`,
			"install-mise-tools",
			environmentFile,
			misePath,
		}, installArgs...),
		Env: []string{"GOPATH=" + filepath.Join(homeDir, ".go")},
	})

	return setup.Step{
		Name:     stepName,
		Commands: commands,
	}
}

func configSettings(settings [][]string, misePath string) []shell.Command {
	var cmds []shell.Command

	for _, opt := range settings {
		cmds = append(cmds, shell.Command{
			Name: "mise settings: " + opt[0],
			Path: misePath,
			Args: []string{"settings", "set", opt[0], opt[1]},
		})
	}

	return cmds
}
