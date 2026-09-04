package steps

import (
	"fmt"
	"path/filepath"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

type stowPackage string

const (
	// StowSystem selects the system Stow package.
	StowSystem stowPackage = "system"
	// StowUser selects the user Stow package.
	StowUser stowPackage = "user"
)

// Stow symlink packages to the correct paths.
func Stow(pkg stowPackage, repoDir string, homeDir string) setup.Step {
	stowDir := filepath.Join(repoDir, "stow")
	targetRoot := stowTarget(pkg, homeDir)
	args := []string{"-R", "--no-folding"}

	if pkg != StowSystem {
		args = append(args, "--adopt")
	}

	args = append(
		args,
		"-d", stowDir,
		"-t", targetRoot,
		string(pkg),
	)

	var commands []shell.Command

	if pkg == StowSystem {
		commands = append(commands, shell.Command{
			Name: "remove existing confs",
			Path: "rm",
			Args: []string{"-f", "/etc/pacman.conf", "/etc/makepkg.conf"},
			Sudo: true,
		})
	}

	commands = append(commands, shell.Command{
		Name: fmt.Sprintf("stow %s files", pkg),
		Path: "stow",
		Args: args,
		Sudo: pkg == StowSystem,
	})

	if pkg != StowSystem {
		commands = append(commands, restoreRepo(repoDir)...)
	}

	return setup.Step{
		Name:     fmt.Sprintf("Stow %s files", pkg),
		Commands: commands,
	}
}

func stowTarget(dest stowPackage, homeDir string) string {
	switch dest {
	case StowSystem:
		return "/"
	case StowUser:
		return homeDir
	default:
		panic(fmt.Sprintf("unknown stow destination %q", dest))
	}
}

func restoreRepo(repoDir string) []shell.Command {
	return []shell.Command{
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
