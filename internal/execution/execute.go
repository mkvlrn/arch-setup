package execution

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/term"

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

	var passphrase []byte

	if !config.Env.CI {
		_, _ = fmt.Fprint(os.Stderr, "Enter secrets passphrase: ")
		passphrase, err = term.ReadPassword(int(os.Stdin.Fd()))
		_, _ = fmt.Fprintln(os.Stderr)

		if err != nil {
			return fmt.Errorf("read secrets passphrase: %w", err)
		}
	}

	stopSudo, err := sudo.KeepAlive(ctx)
	if err != nil {
		return err
	}
	defer stopSudo()

	return setup.Run(ctx, os.Stdout, runPlan(&config, secretsData, passphrase))
}
