package setup

import (
	"context"
	"fmt"
	"io"

	"github.com/mkvlrn/arch-setup/internal/shell"
)

// Run executes setup steps sequentially and stops at the first failure.
func Run(ctx context.Context, output io.Writer, steps []Step) error {
	for index, step := range steps {
		_, _ = fmt.Fprintf(output, "[%d/%d] %s\n", index+1, len(steps), step.Name)

		var err error
		if step.Run != nil {
			err = step.Run(ctx)
		} else {
			_, err = shell.Run(ctx, step.Commands)
		}

		if err != nil {
			return fmt.Errorf("run step %q: %w", step.Name, err)
		}
	}

	_, _ = fmt.Fprintln(output, "Done.")

	return nil
}
