package steps

import (
	"fmt"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// CloneRepo returns the step that clones the arch-setup repository.
func CloneRepo(repoHTTP string, repoSSH string, repoDir string) setup.Step {
	return setup.Step{
		Name: fmt.Sprintf("Clone %s to %s", repoHTTP, repoDir),
		Commands: []shell.Command{
			{
				Name: "clone repo",
				Path: "git",
				Args: []string{"clone", repoHTTP, repoDir},
			},
			setRepoUpstream(repoSSH, repoDir),
		},
	}
}

// ExistingRepo returns the step that configures a repository already in place.
func ExistingRepo(repoSSH string, repoDir string) setup.Step {
	return setup.Step{
		Name:     fmt.Sprintf("Use existing repository at %s", repoDir),
		Commands: []shell.Command{setRepoUpstream(repoSSH, repoDir)},
	}
}

func setRepoUpstream(repoSSH string, repoDir string) shell.Command {
	return shell.Command{
		Name: "set ssh upstream",
		Path: "git",
		Args: []string{"remote", "set-url", "origin", repoSSH},
		Dir:  repoDir,
	}
}
