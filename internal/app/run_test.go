package app_test

import (
	"context"
	"strings"
	"testing"

	"github.com/mkvlrn/arch-setup/internal/app"
	"github.com/mkvlrn/arch-setup/internal/revision"
)

func TestRunRejectsUnstampedInstallerBeforeBootstrap(t *testing.T) {
	previous := revision.Commit
	revision.Commit = "development"

	t.Cleanup(func() { revision.Commit = previous })
	// If bootstrap is reached, HOME would cause a different error.
	t.Setenv("HOME", "")

	for _, verifyOnly := range []bool{false, true} {
		err := app.Run(context.Background(), nil, nil, verifyOnly, false)
		if err == nil || !strings.Contains(err.Error(), "use make build") {
			t.Errorf("expected revision error (verify=%t), got %v", verifyOnly, err)
		}
	}
}
