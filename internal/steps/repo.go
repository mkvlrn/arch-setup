package steps

import (
	"fmt"

	"github.com/mkvlrn/arch-setup/internal/config"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// Repo returns the repository setup step for the current environment.
func Repo(cfg *config.Config, revision string) setup.Step {
	if cfg.Env.CI {
		return ExistingRepo(cfg, revision)
	}

	return CloneRepo(cfg, revision)
}

// CloneRepo aligns the fresh checkout with the installer while keeping main
// attached to origin/main for ordinary pulls after installation.
func CloneRepo(cfg *config.Config, revision string) setup.Step {
	return setup.Step{
		Name: fmt.Sprintf("Clone %s to %s", cfg.Repo.HTTP, cfg.Machine.RepoDir),
		Commands: []shell.Command{
			{
				Name: "clone repo",
				Path: "git",
				Args: []string{"clone", cfg.Repo.HTTP, cfg.Machine.RepoDir},
			},
			{
				Name: "checkout pinned revision",
				Path: "git",
				Args: []string{"checkout", "-B", "main", revision},
				Dir:  cfg.Machine.RepoDir,
			},
			{
				Name: "track origin main",
				Path: "git",
				Args: []string{"branch", "--set-upstream-to=origin/main", "main"},
				Dir:  cfg.Machine.RepoDir,
			},
			setRepoUpstream(cfg.Repo.SSH, cfg.Machine.RepoDir),
		},
	}
}

// ExistingRepo returns the step that configures a repository already in place.
func ExistingRepo(cfg *config.Config, revision string) setup.Step {
	return setup.Step{
		Name: fmt.Sprintf("Use existing repository at %s", cfg.Machine.RepoDir),
		Commands: []shell.Command{
			{
				Name: "assert repository revision",
				Path: "sh",
				Args: []string{"-c", `head=$(git rev-parse HEAD) || exit
if [ "$head" != "$1" ]; then
    printf 'expected HEAD %s, got %s\n' "$1" "$head" >&2
    exit 1
fi`, "assert-repository-revision", revision},
				Dir: cfg.Machine.RepoDir,
			},
			setRepoUpstream(cfg.Repo.SSH, cfg.Machine.RepoDir),
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
