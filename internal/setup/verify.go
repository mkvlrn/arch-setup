package setup

import (
	"context"
	"errors"
	"fmt"
	"io"
)

// Verify executes every verification check and combines their failures.
func Verify(ctx context.Context, output io.Writer, checks []Check) error {
	var failures []error

	for index, check := range checks {
		_, _ = fmt.Fprintf(output, "[%d/%d] %s\n", index+1, len(checks), check.Name)

		if err := check.Run(ctx); err != nil {
			failures = append(failures, fmt.Errorf("verify %q: %w", check.Name, err))
		}
	}

	if len(failures) > 0 {
		return errors.Join(failures...)
	}

	_, _ = fmt.Fprintln(output, "Verification passed.")

	return nil
}
