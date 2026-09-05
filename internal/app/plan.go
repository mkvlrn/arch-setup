package app

import (
	"github.com/mkvlrn/arch-setup/internal/revision"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/start"
	"github.com/mkvlrn/arch-setup/internal/steps"
)

func setupPlan(config start.Config, ci bool) []setup.Step {
	plan := []setup.Step{
		steps.InstallPkg(steps.UsePacman, config.BasePackages),
	}

	if ci {
		plan = append(
			plan,
			steps.ExistingRepo(config.RepoSSH, config.RepoDir, revision.Commit),
		)
	} else {
		plan = append(
			plan,
			steps.RemovePkg(config.RemovePackages),
			steps.CloneRepo(config.RepoHTTP, config.RepoSSH, config.RepoDir, revision.Commit),
		)
	}

	plan = append(
		plan,
		steps.Stow(steps.StowSystem, config.RepoDir, config.HomeDir),
		steps.Yay(config.TempDir, config.MirrorListPath),
		steps.InstallPkg(steps.UseYay, config.MainPackages),
		steps.Xdg(config.XdgMkDir, config.XdgRmRf, config.HomeDir),
		steps.Stow(steps.StowUser, config.RepoDir, config.HomeDir),
		steps.Mise(config.HomeDir, config.MiseTools),
		steps.User(config.Username, config.HomeDir),
	)

	return plan
}

func verificationPlan(config start.Config, ci bool) []setup.Check {
	packages := append(
		append([]string{}, config.BasePackages...),
		config.MainPackages...,
	)

	checks := []setup.Check{
		steps.CloneRepoVerify(config.RepoSSH, config.RepoDir, revision.Commit),
		steps.StowVerify(steps.StowSystem, config.RepoDir, config.HomeDir),
		steps.YayVerify(config.MirrorListPath, config.MirrorListCheck),
		steps.InstallPkgVerify(packages),
	}

	if !ci {
		checks = append(checks, steps.RemovePkgVerify(config.RemovePackages))
	}

	checks = append(
		checks,
		steps.XdgVerify(config.XdgMkDir, config.XdgRmRf, config.HomeDir),
		steps.StowVerify(steps.StowUser, config.RepoDir, config.HomeDir),
		steps.MiseVerify(config.HomeDir),
		steps.UserVerify(config.Username, config.HomeDir),
	)

	return checks
}
