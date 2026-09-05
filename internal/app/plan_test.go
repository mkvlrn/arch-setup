package app //nolint:testpackage // Test plan construction without exporting it.

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/mkvlrn/arch-setup/internal/revision"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/start"
	"github.com/mkvlrn/arch-setup/internal/steps"
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
		want := []setup.Step{steps.InstallPkg(steps.UsePacman, config.BasePackages)}

		if ci {
			want = append(want, steps.ExistingRepo(config.RepoSSH, config.RepoDir, revision.Commit))
		} else {
			want = append(want,
				steps.RemovePkg(config.RemovePackages),
				steps.CloneRepo(config.RepoHTTP, config.RepoSSH, config.RepoDir, revision.Commit),
			)
		}

		want = append(want,
			steps.Stow(steps.StowSystem, config.RepoDir, config.HomeDir),
			steps.Yay(config.TempDir, config.MirrorListPath),
			steps.InstallPkg(steps.UseYay, config.MainPackages),
			steps.Xdg(config.XdgMkDir, config.XdgRmRf, config.HomeDir),
			steps.Stow(steps.StowUser, config.RepoDir, config.HomeDir),
			steps.Mise(config.HomeDir, config.MiseTools),
			steps.User(config.Username, config.HomeDir),
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
			steps.CloneRepoVerify(config.RepoSSH, config.RepoDir, revision.Commit),
			steps.StowVerify(steps.StowSystem, config.RepoDir, config.HomeDir),
			steps.YayVerify(config.MirrorListPath, config.MirrorListCheck),
			steps.InstallPkgVerify(append(slices.Clone(config.BasePackages), config.MainPackages...)),
		}
		if !ci {
			want = append(want, steps.RemovePkgVerify(config.RemovePackages))
		}

		want = append(want,
			steps.XdgVerify(config.XdgMkDir, config.XdgRmRf, config.HomeDir),
			steps.StowVerify(steps.StowUser, config.RepoDir, config.HomeDir),
			steps.MiseVerify(config.HomeDir),
			steps.UserVerify(config.Username, config.HomeDir),
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
