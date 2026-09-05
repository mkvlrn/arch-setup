package execute

import (
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// Xdg returns the step that replaces the default XDG directories.
func Xdg(mkdir []string, rmrf []string, homeDir string) setup.Step {
	return setup.Step{
		Name: "Rework XDG user dirs",
		Commands: []shell.Command{
			{
				Name: "update xdg dirs",
				Path: "xdg-user-dirs-update",
			},
			{
				Name: "create new xdg set",
				Path: "mkdir",
				Args: append([]string{"-p"}, mkdir...),
				Dir:  homeDir,
			},
			{
				Name: "remove old xdg set",
				Path: "rm",
				Args: append([]string{"-rf"}, rmrf...),
				Dir:  homeDir,
			},
		},
	}
}
