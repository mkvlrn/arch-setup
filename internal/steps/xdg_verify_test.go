package steps //nolint:testpackage // Tests share fixtures with unexported filesystem helper tests.

import (
	"os"
	"path/filepath"
	"testing"
)

func filesystemVerifyXdgFixture(t *testing.T, kind, path string) {
	t.Helper()

	switch kind {
	case "directory":
		if err := os.MkdirAll(path, 0o700); err != nil {
			t.Fatal(err)
		}
	case "file":
		filesystemVerifyFile(t, path)
	case "directory link":
		filesystemVerifyLink(t, t.TempDir(), path)
	case "file link":
		target := filepath.Join(t.TempDir(), "file")
		filesystemVerifyFile(t, target)
		filesystemVerifyLink(t, target, path)
	case "dangling":
		filesystemVerifyLink(t, "absent", path)
	case "cycle":
		filesystemVerifyLink(t, path, path)
	}
}

func TestXdgVerifyFilesystem(t *testing.T) {
	for _, tc := range []struct {
		name    string
		kind    string
		removed bool
		want    string
	}{
		{name: "required directory", kind: "directory"},
		{name: "required symlink directory", kind: "directory link"},
		{name: "required missing", want: "inspect directory"},
		{name: "required regular file", kind: "file", want: "not a directory"},
		{name: "required symlink file", kind: "file link", want: "not a directory"},
		{name: "required dangling", kind: "dangling", want: "inspect directory"},
		{name: "removed absent", removed: true},
		{name: "removed dangling is absent", kind: "dangling", removed: true},
		{name: "removed directory remains", kind: "directory", removed: true, want: "still exists"},
		{name: "removed file remains", kind: "file", removed: true, want: "still exists"},
		{name: "removed link remains", kind: "directory link", removed: true, want: "still exists"},
		{name: "removed cyclic link", kind: "cycle", removed: true, want: "inspect path"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			path := filepath.Join(home, "nested", "directory with spaces")

			filesystemVerifyXdgFixture(t, tc.kind, path)

			required, removed := []string{"nested/directory with spaces"}, []string(nil)
			if tc.removed {
				required, removed = removed, required
			}

			check := XdgVerify(required, removed, home)
			if check.Name != "Verify XDG user directories" {
				t.Errorf("unexpected check name %q", check.Name)
			}

			err := check.Run(t.Context())
			if tc.want == "" {
				if err != nil {
					t.Fatal(err)
				}

				return
			}

			filesystemVerifyError(t, err, path, tc.want)
		})
	}
}

func TestXdgVerifyEmptyFilesystem(t *testing.T) {
	if err := XdgVerify(nil, nil, filepath.Join(t.TempDir(), "absent")).Run(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func TestXdgVerifyAggregatesFilesystem(t *testing.T) {
	home := t.TempDir()
	filesystemVerifyFile(t, filepath.Join(home, "wrong type"))
	filesystemVerifyFile(t, filepath.Join(home, "old one"))
	filesystemVerifyFile(t, filepath.Join(home, "old two"))

	err := XdgVerify([]string{"missing", "wrong type"}, []string{"old one", "old two"}, home).Run(t.Context())
	filesystemVerifyError(t, err,
		filepath.Join(home, "missing"), filepath.Join(home, "wrong type"),
		filepath.Join(home, "old one"), filepath.Join(home, "old two"),
		"not a directory", "still exists")
}
