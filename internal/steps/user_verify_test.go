package steps //nolint:testpackage // Tests exercise unexported filesystem verification helpers.

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestVerifyHomeTraversalFilesystem(t *testing.T) {
	for _, mode := range []os.FileMode{0o700, 0o710, 0o750, 0o744, 0o701, 0o711, 0o755} {
		t.Run(fmt.Sprintf("%04o", mode), func(t *testing.T) {
			home := t.TempDir()
			if err := os.Chmod(home, mode); err != nil {
				t.Fatal(err)
			}

			info, err := os.Stat(home)
			if err != nil {
				t.Fatal(err)
			}

			if info.Mode().Perm() != mode {
				t.Fatalf("mode = %04o, want %04o", info.Mode().Perm(), mode)
			}

			err = verifyHomeTraversal(home)
			if mode&0o001 == 0 {
				filesystemVerifyError(t, err, home, "not traversable by other users")

				return
			}

			if err != nil {
				t.Fatal(err)
			}
		})
	}

	t.Run("missing", func(t *testing.T) {
		home := filepath.Join(t.TempDir(), "absent")
		err := verifyHomeTraversal(home)
		filesystemVerifyError(t, err, "inspect home directory", home)

		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("expected wrapped ErrNotExist, got %v", err)
		}
	})
}

func filesystemVerifyCompletionFixture(t *testing.T, home, kind, selected string) {
	t.Helper()

	base := filepath.Join(home, ".config", "fish", "completions")

	for _, name := range []string{"mise", "gh", "glab"} {
		path := filepath.Join(base, name+".fish")
		if name != selected || kind == "regular" {
			filesystemVerifyFile(t, path)

			continue
		}

		switch kind {
		case "symlink":
			target := filepath.Join(home, "generated completion")
			filesystemVerifyFile(t, target)
			filesystemVerifyLink(t, target, path)
		case "dangling":
			filesystemVerifyLink(t, "absent", path)
		}
	}
}

func TestVerifyCompletionsFilesystem(t *testing.T) {
	for _, kind := range []string{"regular", "symlink", "missing", "dangling"} {
		for _, selected := range []string{"mise", "gh", "glab"} {
			t.Run(kind+"/"+selected, func(t *testing.T) {
				home := t.TempDir()
				filesystemVerifyCompletionFixture(t, home, kind, selected)

				err := verifyCompletions(home)

				if kind == "missing" || kind == "dangling" {
					want := "completion file for " + selected + " not generated"
					if err == nil || err.Error() != want {
						t.Fatalf("error = %v, want %q", err, want)
					}

					return
				}

				if err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}
