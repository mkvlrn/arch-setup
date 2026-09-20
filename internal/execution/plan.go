package execution

import (
	"context"

	"github.com/mkvlrn/arch-setup/internal/checks"
	"github.com/mkvlrn/arch-setup/internal/config"
	"github.com/mkvlrn/arch-setup/internal/revision"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/steps"
)

func runPlan(ctx context.Context, cfg *config.Config, secretsData []byte) []setup.Step {
	plan := []setup.Step{
		steps.InstallPkg(cfg, steps.UsePacman),
		steps.RemovePkg(cfg),
		steps.Repo(cfg, revision.Commit),
		steps.Stow(cfg, steps.StowMakepkg),
		steps.Stow(cfg, steps.StowSystem),
		steps.Yay(cfg),
		steps.InstallPkg(cfg, steps.UseYay),
		steps.Secrets(cfg, secretsData),
		steps.Xdg(cfg),
		steps.Mise(cfg),
		steps.Fonts(cfg),
		steps.Stow(cfg, steps.StowUser),
		steps.RcloneProton(ctx, cfg),
		steps.User(cfg),
	}

	return plan
}

func verifyPlan(cfg *config.Config) []setup.Check {
	return []setup.Check{
		checks.Repo(cfg, revision.Commit),
		checks.Stow(cfg, steps.StowMakepkg),
		checks.Stow(cfg, steps.StowSystem),
		checks.Yay(cfg),
		checks.InstalledPkg(cfg),
		checks.RemovedPkg(cfg),
		checks.Xdg(cfg),
		checks.Stow(cfg, steps.StowUser),
		checks.Mise(cfg),
		checks.Fonts(cfg),
		checks.User(cfg),
	}
}
