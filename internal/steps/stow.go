package steps

import (
	"fmt"
	"path/filepath"

	"github.com/mkvlrn/arch-setup/internal/config"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// StowPackage represents the Stow package to be stowed.
type StowPackage string

const (
	// StowSystem selects the system Stow package.
	StowSystem StowPackage = "system"
	// StowUser selects the user Stow package.
	StowUser StowPackage = "user"
	// StowMakepkg selects the makepkg Stow package.
	StowMakepkg StowPackage = "makepkg"
)

// Stow symlink packages to the correct paths.
func Stow(cfg *config.Config, pkg StowPackage) setup.Step {
	stowDir := filepath.Join(cfg.Machine.RepoDir, "stow")
	targetRoot := stowTarget(pkg, cfg.Machine.HomeDir)
	args := []string{"--no-folding"}

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
			Name: "set pacman to display color",
			Path: "sed",
			Args: []string{"-i", "s/^#Color/Color/", "/etc/pacman.conf"},
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
		commands = append(commands, restoreRepo(cfg.Machine.RepoDir)...)
	}

	return setup.Step{
		Name:     fmt.Sprintf("Stow %s files", pkg),
		Commands: commands,
	}
}

func stowTarget(dest StowPackage, homeDir string) string {
	switch dest {
	case StowSystem:
		return "/"
	case StowUser, StowMakepkg:
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
