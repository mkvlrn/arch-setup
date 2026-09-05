package verify

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// Repository returns the check for the cloned arch-setup repository.
func Repository(repoSSH, repoDir, revision string) setup.Check {
	return setup.Check{
		Name: "Verify cloned repository",
		Run: func(ctx context.Context) error {
			gitDir := filepath.Join(repoDir, ".git")

			info, err := os.Stat(gitDir)
			if err != nil {
				return fmt.Errorf("find repository metadata at %q: %w", gitDir, err)
			}

			if !info.IsDir() {
				return fmt.Errorf("repository metadata path %q is not a directory", gitDir)
			}

			results, err := shell.Run(ctx, []shell.Command{
				{
					Name: "get repository remote",
					Path: "git",
					Args: []string{"remote", "get-url", "origin"},
					Dir:  repoDir,
				},
				{
					Name: "get repository HEAD",
					Path: "git",
					Args: []string{"rev-parse", "HEAD"},
					Dir:  repoDir,
				},
				{
					Name: "get repository status",
					Path: "git",
					Args: []string{"status", "--porcelain"},
					Dir:  repoDir,
				},
			})
			if err != nil {
				return err
			}

			return verifyRepoResults(results, repoSSH, revision)
		},
	}
}

func verifyRepoResults(results []shell.Result, repoSSH, revision string) error {
	var failures []error

	remote := strings.TrimSpace(results[0].Stdout)
	if remote != repoSSH {
		failures = append(failures, fmt.Errorf("expected origin %q, got %q", repoSSH, remote))
	}

	head := strings.TrimSpace(results[1].Stdout)
	if head != revision {
		failures = append(failures, fmt.Errorf("expected HEAD %q, got %q", revision, head))
	}

	status := strings.TrimSpace(results[2].Stdout)
	if status != "" {
		failures = append(failures, fmt.Errorf("repository is dirty:\n%s", status))
	}

	return errors.Join(failures...)
}
