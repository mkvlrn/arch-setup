package steps

import (
	"fmt"
	"path/filepath"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// Mise installs the mise binary and uses it to globally install the tools in its manifest.
func Mise(homeDir string, tools []string, settings [][]string) setup.Step {
	stepName := fmt.Sprintf("Installing mise and %d tools", len(tools))
	misePath := filepath.Join(homeDir, ".local", "bin", "mise")
	settingsCmds := configSettings(settings, misePath)

	commands := []shell.Command{
		{
			Name: "install mise",
			Path: "sh",
			Args: []string{"-c", "curl https://mise.run | sh"},
		},
	}

	commands = append(commands, settingsCmds...)

	commands = append(commands, shell.Command{
		Name: "install tools managed by mise",
		Path: misePath,
		Args: []string{"install"},
		Env:  []string{"GOPATH=" + filepath.Join(homeDir, ".go")},
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
