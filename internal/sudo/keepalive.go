package sudo

import (
	"context"
	"fmt"
	"time"

	"github.com/mkvlrn/arch-setup/internal/shell"
)

const refreshInterval = 10 * time.Second

// KeepAlive validates sudo credentials and periodically refreshes them.
func KeepAlive(ctx context.Context) (context.Context, context.CancelFunc, error) {
	if _, err := shell.Run(ctx, []shell.Command{
		{
			Name: "validate sudo credentials",
			Path: "sudo",
			Args: []string{"-v"},
		},
	}); err != nil {
		return nil, nil, fmt.Errorf(
			"validate sudo credentials: %w",
			err,
		)
	}

	keepaliveCtx, cancel := context.WithCancel(ctx)

	go refresh(keepaliveCtx, cancel)

	return keepaliveCtx, cancel, nil
}

func refresh(ctx context.Context, cancel context.CancelFunc) {
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			if _, err := shell.Run(ctx, []shell.Command{
				{
					Name: "refresh sudo credentials",
					Path: "sudo",
					Args: []string{"-n", "-v"},
				},
			}); err != nil {
				cancel()

				return
			}
		}
	}
}
