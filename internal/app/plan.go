package app

import (
	"github.com/mkvlrn/arch-setup/internal/execute"
	"github.com/mkvlrn/arch-setup/internal/revision"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/start"
	"github.com/mkvlrn/arch-setup/internal/verify"
)

func setupPlan(config start.Config, ci bool) []setup.Step {
	plan := []setup.Step{
		execute.InstallPkg(execute.UsePacman, config.BasePackages),
	}

	if ci {
		plan = append(
			plan,
			execute.ExistingRepo(config.RepoSSH, config.RepoDir, revision.Commit),
		)
	} else {
		plan = append(
			plan,
			execute.RemovePkg(config.RemovePackages),
			execute.CloneRepo(config.RepoHTTP, config.RepoSSH, config.RepoDir, revision.Commit),
		)
	}

	plan = append(
		plan,
		execute.Stow(execute.StowSystem, config.RepoDir, config.HomeDir),
		execute.Yay(config.TempDir, config.MirrorListPath),
		execute.InstallPkg(execute.UseYay, config.MainPackages),
		execute.Xdg(config.XdgMkDir, config.XdgRmRf, config.HomeDir),
		execute.Stow(execute.StowUser, config.RepoDir, config.HomeDir),
		execute.Mise(config.HomeDir, config.MiseTools),
		execute.User(config.Username, config.HomeDir),
	)

	return plan
}

func verificationPlan(config start.Config, ci bool) []setup.Check {
	packages := append(
		append([]string{}, config.BasePackages...),
		config.MainPackages...,
	)

	checks := []setup.Check{
		verify.Repository(config.RepoSSH, config.RepoDir, revision.Commit),
		verify.Stow(verify.StowSystem, config.RepoDir, config.HomeDir),
		verify.Yay(config.MirrorListPath, config.MirrorListCheck),
		verify.InstalledPackages(packages),
	}

	if !ci {
		checks = append(checks, verify.RemovedPackages(config.RemovePackages))
	}

	checks = append(
		checks,
		verify.Xdg(config.XdgMkDir, config.XdgRmRf, config.HomeDir),
		verify.Stow(verify.StowUser, config.RepoDir, config.HomeDir),
		verify.Mise(config.HomeDir),
		verify.User(config.Username, config.HomeDir),
	)

	return checks
}
