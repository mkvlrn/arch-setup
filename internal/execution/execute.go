package execution

import (
	"context"
	"os"

	"github.com/mkvlrn/arch-setup/internal/config"
	"github.com/mkvlrn/arch-setup/internal/revision"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/sudo"
)

// Run bootstraps the embedded configuration and performs setup or verification.
func Run(ctx context.Context, configData []byte, secretsData []byte, verifyOnly bool) error {
	if err := revision.Validate(revision.Commit); err != nil {
		return err
	}

	config, err := config.Load(configData)
	if err != nil {
		return err
	}

	if verifyOnly {
		return setup.Verify(ctx, os.Stdout, verifyPlan(&config))
	}

	stopSudo, err := sudo.KeepAlive(ctx)
	if err != nil {
		return err
	}
	defer stopSudo()

	return setup.Run(ctx, os.Stdout, runPlan(ctx, &config, secretsData))
}
