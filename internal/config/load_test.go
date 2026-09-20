package config_test

import (
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mkvlrn/arch-setup/internal/config"
)

func TestLoadConfig(t *testing.T) {
	homeDir := t.TempDir()
	tempDir := t.TempDir()

	t.Setenv("HOME", homeDir)
	t.Setenv("TMPDIR", tempDir)

	currentUser, err := user.Current()
	if err != nil {
		t.Fatalf("get current user: %v", err)
	}

	configData, err := os.ReadFile("testdata/sample-config.json")
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}

	got, err := config.Load(configData)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	want := config.Config{
		Repo: config.Repo{
			HTTP: "https://github.com/mkvlrn/arch-setup",
			SSH:  "git@github.com:mkvlrn/arch-setup",
		},
		Machine: config.Machine{
			Username: currentUser.Username,
			HomeDir:  homeDir,
			RepoDir:  filepath.Join(homeDir, "repos", "arch-setup"),
			TempDir:  tempDir,
		},
		Mise: config.Mise{
			Tools:    []string{"go", "bun"},
			Settings: [][]string{{"auto_update", "true"}, {"auto_update_check_duration", "3d"}},
		},
		Pacman: config.Pacman{
			Install:   []string{"git", "stow"},
			Uninstall: []string{"vim"},
		},
		Yay: config.Yay{
			MirrorListPath:  "path",
			MirrorListCheck: "check",
			Packages:        []string{"zed", "zen-browser-bin"},
		},
		Xdg: config.Xdg{
			MkDir: []string{"documents", "downloads"},
			RmRf:  []string{"Documents", "Downloads"},
		},
		GetNF: []string{"Hack", "IosevkaTerm"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestLoadRejectsInvalidConfig(t *testing.T) {
	tests := []struct {
		name       string
		configData string
		wantErr    string
	}{
		{
			name:       "malformed JSON",
			configData: `{`,
			wantErr:    "decode embed config file",
		},
		{
			name:       "wrong JSON field type",
			configData: `{"mise": 123}`,
			wantErr:    "decode embed config file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.Load([]byte(tt.configData))
			if err == nil {
				t.Fatal("expected an error")
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(cfg, config.Config{}) {
				t.Errorf("expected zero config on failure, got %#v", cfg)
			}
		})
	}
}

func TestLoadWithoutHomeDirectory(t *testing.T) {
	t.Setenv("HOME", "")

	_, err := config.Load([]byte(`{}`))
	if err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(err.Error(), "get user home dir") {
		t.Fatalf("expected home directory error, got %v", err)
	}
}

func TestBootstrapUsesTemporaryDirectory(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("TMPDIR", t.TempDir())

	config, err := config.Load([]byte(`{}`))
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	if config.Machine.TempDir != os.TempDir() {
		t.Fatalf("expected temp directory %q, got %q", os.TempDir(), config.Machine.TempDir)
	}
}
