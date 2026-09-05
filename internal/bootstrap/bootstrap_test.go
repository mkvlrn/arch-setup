package bootstrap_test

import (
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/mkvlrn/arch-setup/internal/bootstrap"
)

func TestBootstrapLoadsConfig(t *testing.T) {
	homeDir := t.TempDir()
	tempDir := t.TempDir()

	t.Setenv("HOME", homeDir)
	t.Setenv("TMPDIR", tempDir)

	currentUser, err := user.Current()
	if err != nil {
		t.Fatalf("get current user: %v", err)
	}

	configData := []byte(`{
		"basePackages": ["git", "stow"],
		"removePackages": ["unwanted"],
		"mainPackages": ["fish"],
		"repoHttp": "https://example.com/config.git",
		"repoSsh": "git@example.com:config.git",
		"mirrorListPath": "/etc/pacman.d/mirrorlist",
		"mirrorListCheck": "example-mirror",
		"xdgMkDir": [".config", ".local/share"],
		"xdgRmRf": ["Desktop"],
		"Username": "ignored-user",
		"HomeDir": "/ignored-home",
		"RepoDir": "/ignored-repo",
		"TempDir": "/ignored-temp",
		"MiseTools": ["ignored-tool"]
	}`)
	miseData := []byte(`
[tools]
go = "latest"
node = "lts"
`)

	got, err := bootstrap.Bootstrap(configData, miseData)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	slices.Sort(got.MiseTools)

	want := bootstrap.Config{
		BasePackages:    []string{"git", "stow"},
		RemovePackages:  []string{"unwanted"},
		MainPackages:    []string{"fish"},
		MiseTools:       []string{"go", "node"},
		RepoHTTP:        "https://example.com/config.git",
		RepoSSH:         "git@example.com:config.git",
		MirrorListPath:  "/etc/pacman.d/mirrorlist",
		MirrorListCheck: "example-mirror",
		XdgMkDir:        []string{".config", ".local/share"},
		XdgRmRf:         []string{"Desktop"},
		Username:        currentUser.Username,
		HomeDir:         homeDir,
		RepoDir:         filepath.Join(homeDir, "repos", "arch-setup"),
		TempDir:         tempDir,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected config:\ngot:  %#v\nwant: %#v", got, want)
	}
}

func TestBootstrapWithNoMiseTools(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	tests := []struct {
		name string
		data string
	}{
		{
			name: "missing tools table",
			data: "",
		},
		{
			name: "empty tools table",
			data: "[tools]\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := bootstrap.Bootstrap([]byte(`{}`), []byte(tt.data))
			if err != nil {
				t.Fatalf("expected success, got %v", err)
			}

			if len(config.MiseTools) != 0 {
				t.Fatalf("expected no mise tools, got %v", config.MiseTools)
			}
		})
	}
}

func TestBootstrapRejectsInvalidConfig(t *testing.T) {
	tests := []struct {
		name       string
		configData string
		miseData   string
		wantErr    string
	}{
		{
			name:       "malformed JSON",
			configData: `{`,
			wantErr:    "decode embed config file",
		},
		{
			name:       "wrong JSON field type",
			configData: `{"basePackages": 123}`,
			wantErr:    "decode embed config file",
		},
		{
			name:       "malformed TOML",
			configData: `{}`,
			miseData:   "[tools",
			wantErr:    "decode embed mise config file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := bootstrap.Bootstrap(
				[]byte(tt.configData),
				[]byte(tt.miseData),
			)
			if err == nil {
				t.Fatal("expected an error")
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(config, bootstrap.Config{}) {
				t.Errorf("expected zero config on failure, got %#v", config)
			}
		})
	}
}

func TestBootstrapWithoutHomeDirectory(t *testing.T) {
	t.Setenv("HOME", "")

	_, err := bootstrap.Bootstrap([]byte(`{}`), nil)
	if err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(err.Error(), "get user home dir") {
		t.Fatalf("expected home directory error, got %v", err)
	}
}

// Keep the test's environment assumptions explicit: Bootstrap uses the
// platform temporary directory, which TMPDIR controls on Linux.
func TestBootstrapUsesTemporaryDirectory(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("TMPDIR", t.TempDir())

	config, err := bootstrap.Bootstrap([]byte(`{}`), nil)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	if config.TempDir != os.TempDir() {
		t.Fatalf("expected temp directory %q, got %q", os.TempDir(), config.TempDir)
	}
}
