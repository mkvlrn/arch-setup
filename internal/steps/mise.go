package steps

import (
	"fmt"
	"path/filepath"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// Mise installs the mise binary and uses it to globally install the tools in its manifest.
func Mise(homeDir string, tools []string) setup.Step {
	stepName := fmt.Sprintf("Installing mise and %d tools", len(tools))
	misePath := filepath.Join(homeDir, ".local", "bin", "mise")

	return setup.Step{
		Name: stepName,
		Commands: []shell.Command{
			{
				Name: "install mise",
				Path: "sh",
				Args: []string{"-c", "curl https://mise.run | sh"},
			},
			{
				Name: "install tools managed by mise",
				Path: misePath,
				Args: []string{"install"},
			},
		},
	}
}
