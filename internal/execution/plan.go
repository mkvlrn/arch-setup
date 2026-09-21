package execution

import (
	"github.com/mkvlrn/arch-setup/internal/checks"
	"github.com/mkvlrn/arch-setup/internal/config"
	"github.com/mkvlrn/arch-setup/internal/revision"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/steps"
)

func runPlan(cfg *config.Config, archStowKey, passphrase []byte) []setup.Step {
	plan := []setup.Step{
		steps.InstallPkg(cfg, steps.UsePacman),
		steps.RemovePkg(cfg),
		steps.Repo(cfg, revision.Commit),
		steps.StowRepo(cfg, archStowKey, passphrase),
		steps.Stow(cfg, steps.StowMakepkg),
		steps.Stow(cfg, steps.StowSystem),
		steps.Stow(cfg, steps.StowSecrets),
		steps.Yay(cfg),
		steps.InstallPkg(cfg, steps.UseYay),
		steps.Mise(cfg),
		steps.Fonts(cfg),
		steps.Xdg(cfg),
		steps.Stow(cfg, steps.StowUser),
		steps.RcloneProton(cfg),
		steps.User(cfg),
	}

	return plan
}

func verifyPlan(cfg *config.Config) []setup.Check {
	return []setup.Check{
		checks.Repo(cfg, revision.Commit),
		checks.Stow(cfg, steps.StowMakepkg),
		checks.Stow(cfg, steps.StowSystem),
		checks.Stow(cfg, steps.StowSecrets),
		checks.Yay(cfg),
		checks.InstalledPkg(cfg),
		checks.RemovedPkg(cfg),
		checks.Xdg(cfg),
		checks.Stow(cfg, steps.StowUser),
		checks.Mise(cfg),
		checks.RcloneProton(cfg),
		checks.Fonts(cfg),
		checks.User(cfg),
	}
}
