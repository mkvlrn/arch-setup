package steps

import (
	"fmt"
	"path/filepath"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

type stowDestination string

const (
	// StowSystem selects the system Stow package.
	StowSystem stowDestination = "system"
	// StowUser selects the user Stow package.
	StowUser stowDestination = "user"
)

// Stow stows system and user stow packages to the correct paths.
func Stow(dest stowDestination, repoDir string, homeDir string) setup.Step {
	var commands []shell.Command

	stowDir := filepath.Join(repoDir, "stow")
	targetRoot := stowTarget(dest, homeDir)

	switch dest {
	case StowSystem:
		commands = []shell.Command{
			{
				Name: "remove existing confs",
				Path: "rm",
				Args: []string{"-f", "/etc/pacman.conf", "/etc/makepkg.conf"},
				Sudo: true,
			},
			{
				Name: "stow system files",
				Path: "stow",
				Args: []string{"-R", "--no-folding", "-d", stowDir, "-t", targetRoot, string(dest)},
				Sudo: true,
			},
		}
	case StowUser:
		commands = []shell.Command{
			{
				Name: "stow user files",
				Path: "stow",
				Args: []string{"-R", "--no-folding", "--adopt", "-d", stowDir, "-t", targetRoot, string(dest)},
			},
			{
				Name: "restore adopted files",
				Path: "git",
				Args: []string{"restore", "."},
				Dir:  repoDir,
			},
			{
				Name: "clean git state",
				Path: "git",
				Args: []string{"clean", "-fd"},
				Dir:  repoDir,
			},
		}
	}

	return setup.Step{
		Name:     fmt.Sprintf("Stow %s files", dest),
		Commands: commands,
	}
}

func stowTarget(dest stowDestination, homeDir string) string {
	switch dest {
	case StowSystem:
		return "/"
	case StowUser:
		return homeDir
	default:
		panic(fmt.Sprintf("unknown stow destination %q", dest))
	}
}
