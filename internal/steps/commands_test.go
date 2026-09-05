package steps_test

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/mkvlrn/arch-setup/internal/shell"
	"github.com/mkvlrn/arch-setup/internal/steps"
)

const (
	commandHome = "/home/test user"
	commandRepo = "/home/test user/setup repo"
)

func requireCommandCount(t *testing.T, commands []shell.Command, want int) {
	t.Helper()

	if len(commands) != want {
		t.Fatalf("command count = %d, want %d", len(commands), want)
	}
}

func checkCommand(t *testing.T, got shell.Command, path string, sudo bool, dir string, args ...string) {
	t.Helper()

	if got.Path != path {
		t.Errorf("executable = %q, want %q", got.Path, path)
	}

	if got.Sudo != sudo {
		t.Errorf("%s sudo = %t, want %t", path, got.Sudo, sudo)
	}

	if got.Dir != dir {
		t.Errorf("%s cwd = %q, want %q", path, got.Dir, dir)
	}

	if !slices.Equal(got.Args, args) {
		t.Errorf("%s arguments = %q, want %q", path, got.Args, args)
	}
}

func TestCommandsInstallPackages(t *testing.T) {
	packages := []string{"git", "base-devel"}

	t.Run("pacman upgrades with sudo", func(t *testing.T) {
		commands := steps.InstallPkg(steps.UsePacman, packages).Commands
		requireCommandCount(t, commands, 1)
		checkCommand(t, commands[0], "pacman", true, "", "-Syu", "--noconfirm", "--needed", "git", "base-devel")
	})
	t.Run("yay installs without sudo", func(t *testing.T) {
		commands := steps.InstallPkg(steps.UseYay, packages).Commands
		requireCommandCount(t, commands, 1)
		checkCommand(t, commands[0], "yay", false, "", "-S", "--noconfirm", "--needed", "git", "base-devel")
	})
}

func TestCommandsRemovePackages(t *testing.T) {
	commands := steps.RemovePkg([]string{"nano", "vim"}).Commands
	requireCommandCount(t, commands, 1)
	checkCommand(t, commands[0], "pacman", true, "", "-Rns", "--noconfirm", "nano", "vim")
}

func TestCommandsStowSystem(t *testing.T) {
	commands := steps.Stow(steps.StowSystem, commandRepo, commandHome).Commands
	requireCommandCount(t, commands, 2)
	checkCommand(t, commands[0], "rm", true, "", "-f", "/etc/pacman.conf", "/etc/makepkg.conf")
	checkCommand(t, commands[1], "stow", true, "",
		"-R", "--no-folding", "-d", filepath.Join(commandRepo, "stow"), "-t", "/", "system")
}

func TestCommandsStowUser(t *testing.T) {
	commands := steps.Stow(steps.StowUser, commandRepo, commandHome).Commands
	requireCommandCount(t, commands, 3)
	checkCommand(t, commands[0], "stow", false, "",
		"-R", "--no-folding", "--adopt", "-d", filepath.Join(commandRepo, "stow"), "-t", commandHome, "user")
	checkCommand(t, commands[1], "git", false, commandRepo, "restore", ".")
	checkCommand(t, commands[2], "git", false, commandRepo, "clean", "-fd")
}

func TestCommandsXdg(t *testing.T) {
	commands := steps.Xdg(
		[]string{"new downloads", "new documents"},
		[]string{"Old Downloads", "Old Documents"},
		commandHome,
	).Commands
	requireCommandCount(t, commands, 3)
	checkCommand(t, commands[0], "xdg-user-dirs-update", false, "")
	checkCommand(t, commands[1], "mkdir", false, commandHome, "-p", "new downloads", "new documents")
	checkCommand(t, commands[2], "rm", false, commandHome, "-rf", "Old Downloads", "Old Documents")
}

func TestCommandsMise(t *testing.T) {
	commands := steps.Mise(commandHome, []string{"go", "node"}).Commands
	requireCommandCount(t, commands, 2)
	checkCommand(t, commands[0], "sh", false, "", "-c", "curl https://mise.run | sh")
	checkCommand(t, commands[1], filepath.Join(commandHome, ".local", "bin", "mise"), false, "", "install")

	wantEnv := []string{"GOPATH=" + filepath.Join(commandHome, ".go")}
	if !slices.Equal(commands[1].Env, wantEnv) {
		t.Errorf("mise environment = %q, want %q", commands[1].Env, wantEnv)
	}
}

func TestCommandsUserSettings(t *testing.T) {
	commands := steps.User("test-user", commandHome).Commands
	requireCommandCount(t, commands, 11)
	checkCommand(t, commands[0], "chsh", true, "", "-s", "/usr/bin/fish", "test-user")
	checkCommand(t, commands[1], "usermod", true, "", "-aG", "docker", "test-user")
	checkCommand(t, commands[2], "usermod", true, "", "-d", filepath.Join(commandHome, "torrents"), "ftp")
	checkCommand(t, commands[3], "chmod", false, "", "o+x", commandHome)

	for i, service := range []string{"docker.socket", "pure-ftpd.service", "paccache.timer"} {
		checkCommand(t, commands[4+i], "systemctl", true, "", "enable", "--now", service)
	}
}

func TestCommandsUserCompletions(t *testing.T) {
	commands := steps.User("test-user", commandHome).Commands
	requireCommandCount(t, commands, 11)

	completionDir := filepath.Join(commandHome, ".config", "fish", "completions")
	checkCommand(t, commands[7], "mkdir", false, "", "-p", completionDir)

	for i, tool := range []string{"mise", "gh", "glab"} {
		t.Run(tool, func(t *testing.T) {
			binary := filepath.Join(commandHome, ".local", "share", "mise", "shims", tool)
			script := `"$1" completion -s fish > "$2"`

			if tool == "mise" {
				binary = filepath.Join(commandHome, ".local", "bin", tool)
				script = `"$1" completion fish > "$2"`
			}

			checkCommand(t, commands[8+i], "sh", false, "",
				"-c", script, "sh", binary, filepath.Join(completionDir, tool+".fish"))
		})
	}
}

func TestCommandsYay(t *testing.T) {
	const (
		tempDir    = "/tmp/build directory"
		mirrorList = "/etc/test mirrors/mirrorlist"
	)

	commands := steps.Yay(tempDir, mirrorList).Commands
	requireCommandCount(t, commands, 7)

	source := filepath.Join(tempDir, "yay-bin")
	checkCommand(t, commands[0], "git", false, "", "clone", "https://aur.archlinux.org/yay-bin", source)
	checkCommand(t, commands[1], "makepkg", false, source, "-si", "--noconfirm")
	checkCommand(t, commands[2], "yay", false, "", "-Y", "--gendb")
	checkCommand(t, commands[3], "yay", false, "", "-Y", "--devel", "--save")
	checkCommand(t, commands[4], "reflector", true, "",
		"--latest", "20", "--protocol", "https", "--sort", "rate", "--save", mirrorList)
	checkCommand(t, commands[5], "yay", false, "", "-Syu", "--noconfirm")
	checkCommand(t, commands[6], "sh", false, "", "-c", "yay -Qq | grep -- '-debug$' | xargs -r yay -Rnsu")
}

func TestCommandsCloneRepo(t *testing.T) {
	const (
		remoteHTTP = "https://example.invalid/setup.git"
		remoteSSH  = "git@example.invalid:setup.git"
		revision   = "0123456789abcdef0123456789abcdef01234567"
	)

	commands := steps.CloneRepo(remoteHTTP, remoteSSH, commandRepo, revision).Commands
	requireCommandCount(t, commands, 4)
	checkCommand(t, commands[0], "git", false, "", "clone", remoteHTTP, commandRepo)
	checkCommand(t, commands[1], "git", false, commandRepo, "checkout", "-B", "main", revision)
	checkCommand(t, commands[2], "git", false, commandRepo, "branch", "--set-upstream-to=origin/main", "main")
	checkCommand(t, commands[3], "git", false, commandRepo, "remote", "set-url", "origin", remoteSSH)
}

func TestCommandsExistingRepoAssertsRevisionBeforeRemote(t *testing.T) {
	const (
		remote   = "git@example.invalid:setup.git"
		revision = "revision with spaces"
	)

	commands := steps.ExistingRepo(remote, commandRepo, revision).Commands
	requireCommandCount(t, commands, 2)

	assertion := commands[0]
	if len(assertion.Args) != 4 {
		t.Fatalf("revision assertion arguments = %q, want four shell arguments", assertion.Args)
	}

	script := assertion.Args[1]
	checkCommand(t, assertion, "sh", false, commandRepo, "-c", script, "assert-repository-revision", revision)

	for _, fragment := range []string{"git rev-parse HEAD", "|| exit", `if [ "$head" != "$1" ]; then`, "exit 1"} {
		if !strings.Contains(script, fragment) {
			t.Errorf("revision assertion lacks %q", fragment)
		}
	}

	if strings.Contains(script, revision) {
		t.Error("revision must be passed as a positional argument, not interpolated into shell code")
	}

	checkCommand(t, commands[1], "git", false, commandRepo, "remote", "set-url", "origin", remote)
}
