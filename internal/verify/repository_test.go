package verify_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mkvlrn/arch-setup/internal/shell"
	"github.com/mkvlrn/arch-setup/internal/verify"
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

func TestRepository(t *testing.T) {
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

			err := verify.Repository(test.remote, dir, revision).Run(t.Context())
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
