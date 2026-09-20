package steps

import (
	"path/filepath"

	"github.com/mkvlrn/arch-setup/internal/config"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// Xdg returns the step that replaces the default XDG directories.
func Xdg(cfg *config.Config) setup.Step {
	mkdir, rmrf, homeDir := cfg.Xdg.MkDir, cfg.Xdg.RmRf, cfg.Machine.HomeDir
	createPaths := make([]string, 0, len(mkdir))

	for _, directory := range mkdir {
		createPaths = append(createPaths, filepath.Join(homeDir, directory))
	}

	removePaths := make([]string, 0, len(rmrf))
	for _, directory := range rmrf {
		removePaths = append(removePaths, filepath.Join(homeDir, directory))
	}

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
				Args: append([]string{"-p"}, createPaths...),
			},
			{
				Name: "remove old xdg set",
				Path: "rm",
				Args: append([]string{"-rf"}, removePaths...),
			},
		},
	}
}
