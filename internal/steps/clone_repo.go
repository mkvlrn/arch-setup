package steps

import (
	"fmt"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// CloneRepo aligns the fresh checkout with the installer while keeping main
// attached to origin/main for ordinary pulls after installation.
func CloneRepo(repoHTTP, repoSSH, repoDir, revision string) setup.Step {
	return setup.Step{
		Name: fmt.Sprintf("Clone %s to %s", repoHTTP, repoDir),
		Commands: []shell.Command{
			{
				Name: "clone repo",
				Path: "git",
				Args: []string{"clone", repoHTTP, repoDir},
			},
			{
				Name: "checkout pinned revision",
				Path: "git",
				Args: []string{"checkout", "-B", "main", revision},
				Dir:  repoDir,
			},
			{
				Name: "track origin main",
				Path: "git",
				Args: []string{"branch", "--set-upstream-to=origin/main", "main"},
				Dir:  repoDir,
			},
			setRepoUpstream(repoSSH, repoDir),
		},
	}
}

// ExistingRepo returns the step that configures a repository already in place.
func ExistingRepo(repoSSH, repoDir, revision string) setup.Step {
	return setup.Step{
		Name: fmt.Sprintf("Use existing repository at %s", repoDir),
		Commands: []shell.Command{
			{
				Name: "assert repository revision",
				Path: "sh",
				Args: []string{"-c", `head=$(git rev-parse HEAD) || exit
if [ "$head" != "$1" ]; then
    printf 'expected HEAD %s, got %s\n' "$1" "$head" >&2
    exit 1
fi`, "assert-repository-revision", revision},
				Dir: repoDir,
			},
			setRepoUpstream(repoSSH, repoDir),
		},
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
