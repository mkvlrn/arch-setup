// Package app builds and executes the installer and verification plans.
package app

import (
	"context"
	"os"

	"github.com/mkvlrn/arch-setup/internal/revision"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/start"
	"github.com/mkvlrn/arch-setup/internal/sudo"
)

// Run bootstraps the embedded configuration and performs setup or verification.
func Run(ctx context.Context, configData, miseConfigData []byte, verifyOnly, ci bool) error {
	if err := revision.Validate(revision.Commit); err != nil {
		return err
	}

	config, err := start.Bootstrap(configData, miseConfigData)
	if err != nil {
		return err
	}

	if verifyOnly {
		return setup.Verify(ctx, os.Stdout, verificationPlan(config, ci))
	}

	stopSudo, err := sudo.KeepAlive(ctx)
	if err != nil {
		return err
	}
	defer stopSudo()

	return setup.Run(ctx, os.Stdout, setupPlan(config, ci))
}
