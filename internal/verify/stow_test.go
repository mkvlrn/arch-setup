package verify //nolint:testpackage // Tests exercise unexported filesystem verification helpers.

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func filesystemVerifyFile(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path, []byte("fixture"), 0o600); err != nil {
		t.Fatal(err)
	}
}

func filesystemVerifyLink(t *testing.T, target, path string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}

	if err := os.Symlink(target, path); err != nil {
		t.Fatal(err)
	}
}

func filesystemVerifyError(t *testing.T, err error, fragments ...string) {
	t.Helper()

	if err == nil {
		t.Fatal("expected an error")
	}

	for _, fragment := range fragments {
		if !strings.Contains(err.Error(), fragment) {
			t.Errorf("error %q does not contain %q", err, fragment)
		}
	}
}

func TestStowFilesystem(t *testing.T) {
	for _, relative := range []bool{false, true} {
		name := "absolute"
		if relative {
			name = "relative"
		}

		t.Run(name, func(t *testing.T) {
			repo, home := t.TempDir(), t.TempDir()
			root := filepath.Join(repo, "stow", "user")
			source := filepath.Join(root, "nested directory", "config file")
			destination := filepath.Join(home, "nested directory", "config file")

			filesystemVerifyFile(t, source)

			for _, ignored := range []string{
				".gitkeep", ".stow-local-ignore",
				"nested directory/.gitkeep", "nested directory/.stow-local-ignore",
			} {
				filesystemVerifyFile(t, filepath.Join(root, ignored))
			}

			target := source

			if relative {
				var err error

				target, err = filepath.Rel(filepath.Dir(destination), source)
				if err != nil {
					t.Fatal(err)
				}
			}

			filesystemVerifyLink(t, target, destination)

			check := Stow(StowUser, repo, home)
			if check.Name != "Verify stowed user files" {
				t.Errorf("unexpected check name %q", check.Name)
			}

			if err := check.Run(t.Context()); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestVerifyStowTreeFilesystem(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		if err := verifyStowTree(t.TempDir(), t.TempDir()); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("missing root", func(t *testing.T) {
		root := filepath.Join(t.TempDir(), "missing")
		err := verifyStowTree(root, t.TempDir())
		filesystemVerifyError(t, err, "walk stow package", root)

		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("expected wrapped ErrNotExist, got %v", err)
		}
	})
	t.Run("aggregation", func(t *testing.T) {
		source, destination := t.TempDir(), t.TempDir()
		filesystemVerifyFile(t, filepath.Join(source, "missing"))
		filesystemVerifyFile(t, filepath.Join(source, "regular"))
		filesystemVerifyFile(t, filepath.Join(destination, "regular"))
		filesystemVerifyError(t, verifyStowTree(source, destination),
			filepath.Join(destination, "missing"), filepath.Join(destination, "regular"), "not a symlink")
	})
}

func filesystemVerifyStowSource(t *testing.T, kind, source, other string) {
	t.Helper()

	switch kind {
	case "dangling source":
		filesystemVerifyLink(t, "absent", source)
	case "cyclic source":
		filesystemVerifyLink(t, "source", source)
	case "resolved source":
		filesystemVerifyLink(t, other, source)
	default:
		filesystemVerifyFile(t, source)
	}
}

func TestVerifyStowLinkFilesystem(t *testing.T) {
	for _, tc := range []struct {
		name string
		want string
	}{
		{name: "missing", want: "inspect"},
		{name: "regular", want: "not a symlink"},
		{name: "wrong", want: "instead of"},
		{name: "dangling destination", want: "resolve destination"},
		{name: "cyclic destination", want: "resolve destination"},
		{name: "dangling source", want: "resolve source"},
		{name: "cyclic source", want: "resolve source"},
		{name: "resolved source"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			source, destination := filepath.Join(root, "source"), filepath.Join(root, "destination")
			other := filepath.Join(root, "other")
			filesystemVerifyFile(t, other)

			filesystemVerifyStowSource(t, tc.name, source, other)

			switch tc.name {
			case "missing":
			case "regular":
				filesystemVerifyFile(t, destination)
			case "wrong", "resolved source":
				filesystemVerifyLink(t, other, destination)
			case "dangling destination":
				filesystemVerifyLink(t, "absent", destination)
			case "cyclic destination":
				filesystemVerifyLink(t, "cycle", destination)
				filesystemVerifyLink(t, "destination", filepath.Join(root, "cycle"))
			default:
				filesystemVerifyLink(t, source, destination)
			}

			err := verifyStowLink(source, destination)
			if tc.want == "" {
				if err != nil {
					t.Fatal(err)
				}

				return
			}

			filesystemVerifyError(t, err, tc.want)
		})
	}
}
