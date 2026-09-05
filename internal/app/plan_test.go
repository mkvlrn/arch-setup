package app //nolint:testpackage // Test plan construction without exporting it.

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/mkvlrn/arch-setup/internal/execute"
	"github.com/mkvlrn/arch-setup/internal/revision"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/start"
	"github.com/mkvlrn/arch-setup/internal/verify"
)

func planConfig() start.Config {
	return start.Config{
		BasePackages: []string{"base"}, MainPackages: []string{"main"},
		RemovePackages: []string{"removed"}, MiseTools: []string{"go"},
		RepoHTTP: "https://example.invalid/setup", RepoSSH: "git@example.invalid:setup",
		RepoDir: "/home/test user/repos/setup", HomeDir: "/home/test user",
		TempDir: "/tmp/build dir", Username: "test-user",
		MirrorListPath: "/etc/mirrorlist", MirrorListCheck: "reflector marker",
		XdgMkDir: []string{"documents"}, XdgRmRf: []string{"Documents"},
	}
}

func TestSetupPlan(t *testing.T) {
	previous := revision.Commit
	revision.Commit = strings.Repeat("a", 40)

	t.Cleanup(func() { revision.Commit = previous })

	for _, ci := range []bool{false, true} {
		config := planConfig()
		want := []setup.Step{execute.InstallPkg(execute.UsePacman, config.BasePackages)}

		if ci {
			want = append(want, execute.ExistingRepo(config.RepoSSH, config.RepoDir, revision.Commit))
		} else {
			want = append(want,
				execute.RemovePkg(config.RemovePackages),
				execute.CloneRepo(config.RepoHTTP, config.RepoSSH, config.RepoDir, revision.Commit),
			)
		}

		want = append(want,
			execute.Stow(execute.StowSystem, config.RepoDir, config.HomeDir),
			execute.Yay(config.TempDir, config.MirrorListPath),
			execute.InstallPkg(execute.UseYay, config.MainPackages),
			execute.Xdg(config.XdgMkDir, config.XdgRmRf, config.HomeDir),
			execute.Stow(execute.StowUser, config.RepoDir, config.HomeDir),
			execute.Mise(config.HomeDir, config.MiseTools),
			execute.User(config.Username, config.HomeDir),
		)

		got := setupPlan(config, ci)
		if len(got) != len(want) {
			t.Fatalf("CI=%t: expected %d steps, got %d", ci, len(want), len(got))
		}

		for i := range want {
			if !reflect.DeepEqual(got[i].Commands, want[i].Commands) {
				t.Errorf("CI=%t: step %d differs from %s", ci, i, want[i].Name)
			}
		}
	}
}

func TestVerificationPlan(t *testing.T) {
	for _, ci := range []bool{false, true} {
		config := planConfig()

		want := []setup.Check{
			verify.Repository(config.RepoSSH, config.RepoDir, revision.Commit),
			verify.Stow(verify.StowSystem, config.RepoDir, config.HomeDir),
			verify.Yay(config.MirrorListPath, config.MirrorListCheck),
			verify.InstalledPackages(append(slices.Clone(config.BasePackages), config.MainPackages...)),
		}
		if !ci {
			want = append(want, verify.RemovedPackages(config.RemovePackages))
		}

		want = append(want,
			verify.Xdg(config.XdgMkDir, config.XdgRmRf, config.HomeDir),
			verify.Stow(verify.StowUser, config.RepoDir, config.HomeDir),
			verify.Mise(config.HomeDir),
			verify.User(config.Username, config.HomeDir),
		)

		got := verificationPlan(config, ci)
		if len(got) != len(want) {
			t.Fatalf("CI=%t: expected %d checks, got %d", ci, len(want), len(got))
		}

		// Check closures are opaque; compare their constructor-provided identities,
		// not literal progress text, and never execute external-command checks.
		for i := range want {
			if got[i].Name != want[i].Name || got[i].Run == nil {
				t.Errorf("CI=%t: check %d differs from %s", ci, i, want[i].Name)
			}
		}
	}
}

func TestPlansPreserveConfig(t *testing.T) {
	config := planConfig()
	// Spare capacity catches append operations that overwrite the caller's backing array.
	backing := []string{"base", "sentinel", "untouched"}
	config.BasePackages = backing[:1]
	before := planConfig()

	for _, ci := range []bool{false, true} {
		setupPlan(config, ci)
		verificationPlan(config, ci)
	}

	if !reflect.DeepEqual(config, before) || !slices.Equal(backing, []string{"base", "sentinel", "untouched"}) {
		t.Fatal("plan construction mutated configuration")
	}
}
