package steps_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mkvlrn/arch-setup/internal/shell"
	"github.com/mkvlrn/arch-setup/internal/steps"
)

const testRemote = "git@example.invalid:setup.git"

func gitOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()

	results, err := shell.Run(t.Context(), []shell.Command{{
		Name: "test git command",
		Path: "git",
		Args: args,
		Dir:  dir,
		Env: []string{
			"GIT_CONFIG_NOSYSTEM=1",
			"GIT_CONFIG_GLOBAL=/dev/null",
			"GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=test@example.invalid",
			"GIT_COMMITTER_NAME=Test", "GIT_COMMITTER_EMAIL=test@example.invalid",
		},
	}})
	if err != nil {
		t.Fatalf("git %v: %v", args, err)
	}

	return strings.TrimSpace(results[0].Stdout)
}

func revisionRepo(t *testing.T) (dir, older, newer string) {
	t.Helper()

	dir = t.TempDir()
	gitOutput(t, dir, "init", "--initial-branch=main")

	if err := os.WriteFile(filepath.Join(dir, "tracked"), []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}

	gitOutput(t, dir, "add", "tracked")
	gitOutput(t, dir, "commit", "-m", "older")
	older = gitOutput(t, dir, "rev-parse", "HEAD")
	gitOutput(t, dir, "commit", "--allow-empty", "-m", "newer")
	newer = gitOutput(t, dir, "rev-parse", "HEAD")

	return dir, older, newer
}

func TestCloneRepoPinsOlderRevision(t *testing.T) {
	source, older, _ := revisionRepo(t)
	destination := filepath.Join(t.TempDir(), "clone with spaces")
	step := steps.CloneRepo(source, testRemote, destination, older)

	if _, err := shell.Run(t.Context(), step.Commands); err != nil {
		t.Fatal(err)
	}

	if head := gitOutput(t, destination, "rev-parse", "HEAD"); head != older {
		t.Fatalf("expected HEAD %s, got %s", older, head)
	}

	if branch := gitOutput(t, destination, "rev-parse", "--abbrev-ref", "HEAD"); branch != "main" {
		t.Fatalf("expected main branch, got %s", branch)
	}

	if err := steps.CloneRepoVerify(testRemote, destination, older).Run(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func TestCloneRepoTracksMainForPull(t *testing.T) {
	source, older, newer := revisionRepo(t)
	destination := filepath.Join(t.TempDir(), "clone")
	// Keep the test origin local so pull exercises tracking without network access.
	step := steps.CloneRepo(source, source, destination, older)

	if _, err := shell.Run(t.Context(), step.Commands); err != nil {
		t.Fatal(err)
	}

	if upstream := gitOutput(t, destination, "rev-parse", "--abbrev-ref", "@{upstream}"); upstream != "origin/main" {
		t.Fatalf("expected origin/main upstream, got %s", upstream)
	}

	gitOutput(t, destination, "pull", "--ff-only")

	if head := gitOutput(t, destination, "rev-parse", "HEAD"); head != newer {
		t.Fatalf("expected pull to advance to %s, got %s", newer, head)
	}
}

func TestCloneRepoUnavailableRevision(t *testing.T) {
	source, _, _ := revisionRepo(t)
	destination := filepath.Join(t.TempDir(), "clone")
	step := steps.CloneRepo(source, testRemote, destination, strings.Repeat("0", 40))

	if _, err := shell.Run(t.Context(), step.Commands); err == nil {
		t.Fatal("expected unavailable revision to fail")
	}

	if remote := gitOutput(t, destination, "remote", "get-url", "origin"); remote != source {
		t.Fatalf("remote changed after checkout failure: %s", remote)
	}
}

func TestExistingRepoRevision(t *testing.T) {
	for _, match := range []bool{true, false} {
		name := "mismatch"
		if match {
			name = "match"
		}

		t.Run(name, func(t *testing.T) {
			dir, older, newer := revisionRepo(t)
			gitOutput(t, dir, "remote", "add", "origin", "original")

			revision := older
			if match {
				revision = newer
			}

			_, err := shell.Run(t.Context(), steps.ExistingRepo(testRemote, dir, revision).Commands)
			if match && err != nil {
				t.Fatal(err)
			}

			remote := "original"
			if match {
				remote = testRemote
			} else {
				assertRevisionMismatch(t, err, older, newer)
			}

			if got := gitOutput(t, dir, "remote", "get-url", "origin"); got != remote {
				t.Errorf("expected origin %s, got %s", remote, got)
			}

			if got := gitOutput(t, dir, "rev-parse", "HEAD"); got != newer {
				t.Errorf("existing HEAD changed: %s", got)
			}
		})
	}
}

func assertRevisionMismatch(t *testing.T, err error, expected, actual string) {
	t.Helper()

	var commandErr *shell.CommandError
	if !errors.As(err, &commandErr) {
		t.Fatalf("expected command error, got %v", err)
	}

	if !strings.Contains(commandErr.Stderr, expected) || !strings.Contains(commandErr.Stderr, actual) {
		t.Fatalf("missing revision diagnostic: %s", commandErr.Stderr)
	}
}

func makeRepoDirty(t *testing.T, dir, kind string) {
	t.Helper()

	file := "untracked"
	if kind == "tracked dirty" || kind == "staged dirty" {
		file = "tracked"
	}

	if err := os.WriteFile(filepath.Join(dir, file), []byte("dirty"), 0o600); err != nil {
		t.Fatal(err)
	}

	if kind == "staged dirty" {
		gitOutput(t, dir, "add", file)
	}
}

func TestCloneRepoVerify(t *testing.T) {
	tests := []struct {
		name          string
		wrongRevision bool
		remote        string
		dirty         string
		want          []string
	}{
		{name: "match", remote: testRemote},
		{name: "wrong revision", wrongRevision: true, remote: testRemote, want: []string{"expected HEAD"}},
		{name: "wrong remote", remote: "different", want: []string{"expected origin"}},
		{name: "dirty", remote: testRemote, dirty: "dirty", want: []string{"repository is dirty"}},
		{name: "tracked dirty", remote: testRemote, dirty: "tracked dirty", want: []string{"repository is dirty"}},
		{name: "staged dirty", remote: testRemote, dirty: "staged dirty", want: []string{"repository is dirty"}},
		{
			name: "all wrong", wrongRevision: true, remote: "different", dirty: "dirty",
			want: []string{"expected HEAD", "expected origin", "repository is dirty"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir, older, newer := revisionRepo(t)
			gitOutput(t, dir, "remote", "add", "origin", testRemote)

			revision := newer
			if test.wrongRevision {
				revision = older
			}

			if test.dirty != "" {
				makeRepoDirty(t, dir, test.dirty)
			}

			err := steps.CloneRepoVerify(test.remote, dir, revision).Run(t.Context())
			if len(test.want) == 0 && err != nil {
				t.Fatal(err)
			}

			for _, message := range test.want {
				if err == nil || !strings.Contains(err.Error(), message) {
					t.Errorf("expected %q, got %v", message, err)
				}
			}
		})
	}
}
