package setup

import (
	"context"
	"fmt"
	"io"

	"github.com/mkvlrn/arch-setup/internal/shell"
)

// Run executes setup steps sequentially and stops at the first failure.
func Run(ctx context.Context, output io.Writer, steps []Step) error {
	_, _ = fmt.Fprintln(output, "Starting.")

	defer func() {
		_, _ = fmt.Fprintln(output, "Done.")
	}()

	for index, step := range steps {
		_, _ = fmt.Fprintf(output, "[%d/%d] %s\n", index+1, len(steps), step.Name)

		if _, err := shell.Run(ctx, step.Commands); err != nil {
			return fmt.Errorf("run step %q: %w", step.Name, err)
		}
	}

	return nil
}
