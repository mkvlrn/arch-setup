package execution

import (
	"context"

	"github.com/mkvlrn/arch-setup/internal/checks"
	"github.com/mkvlrn/arch-setup/internal/config"
	"github.com/mkvlrn/arch-setup/internal/revision"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/steps"
)

func runPlan(ctx context.Context, config config.Config, secretsData []byte) []setup.Step {
	repoStep := steps.CloneRepo(
		config.Repo.HTTP,
		config.Repo.SSH,
		config.Machine.RepoDir,
		revision.Commit,
	)

	if config.Env.CI {
		repoStep = steps.ExistingRepo(
			config.Repo.SSH,
			config.Machine.RepoDir,
			revision.Commit,
		)
	}

	plan := []setup.Step{
		steps.InstallPkg(steps.UsePacman, config.Pacman.Install),
		steps.RemovePkg(config.Pacman.Uninstall),
		repoStep,
		steps.Stow(steps.StowMakepkg, config.Machine.RepoDir, config.Machine.HomeDir),
		steps.Stow(steps.StowSystem, config.Machine.RepoDir, config.Machine.HomeDir),
		steps.Yay(config.Machine.TempDir, config.Yay.MirrorListPath),
		steps.InstallPkg(steps.UseYay, config.Yay.Packages),
	}

	if !config.Env.CI {
		plan = append(plan, steps.Secrets(secretsData, config.Machine.HomeDir))
	}

	plan = append(
		plan,
		steps.Xdg(config.Xdg.MkDir, config.Xdg.RmRf, config.Machine.HomeDir),
		steps.Mise(config.Machine.HomeDir, config.Mise.Tools, config.Mise.Settings),
		steps.Fonts(config.GetNF),
		steps.Stow(steps.StowUser, config.Machine.RepoDir, config.Machine.HomeDir),
	)

	if !config.Env.CI {
		plan = append(plan, steps.RcloneProton(ctx, config.Machine.HomeDir))
	}

	plan = append(plan, steps.User(config.Machine.Username, config.Machine.HomeDir))

	return plan
}

func verifyPlan(config config.Config) []setup.Check {
	packages := append(append([]string{}, config.Pacman.Install...), config.Yay.Packages...)

	planChecks := []setup.Check{
		checks.Repo(config.Repo.SSH, config.Machine.RepoDir, revision.Commit),
		checks.Stow(steps.StowMakepkg, config.Machine.RepoDir, config.Machine.HomeDir),
		checks.Stow(steps.StowSystem, config.Machine.RepoDir, config.Machine.HomeDir),
		checks.Yay(config.Yay.MirrorListPath, config.Yay.MirrorListCheck),
		checks.InstalledPkg(packages),
		checks.RemovedPkg(config.Pacman.Uninstall),
		checks.Xdg(config.Xdg.MkDir, config.Xdg.RmRf, config.Machine.HomeDir),
		checks.Stow(steps.StowUser, config.Machine.RepoDir, config.Machine.HomeDir),
		checks.Mise(config.Machine.HomeDir),
	}

	return append(planChecks, checks.User(
		config.Machine.Username,
		config.Machine.HomeDir,
		config.Env.CI,
	))
}
