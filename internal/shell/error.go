package shell

import (
	"fmt"
	"strings"
)

// CommandError describes a failed external command.
type CommandError struct {
	// Name identifies the command that failed.
	Name string

	// Err is the underlying execution error.
	Err error

	// Stdout contains output captured from standard output.
	Stdout string

	// Stderr contains output captured from standard error.
	Stderr string
}

// Error returns the command failure and its captured output.
func (e *CommandError) Error() string {
	output := []string{fmt.Sprintf("could not run %s: %v", e.Name, e.Err)}

	if stdout := strings.TrimSpace(e.Stdout); stdout != "" {
		output = append(output, "stdout:\n"+stdout)
	}

	if stderr := strings.TrimSpace(e.Stderr); stderr != "" {
		output = append(output, "stderr:\n"+stderr)
	}

	return strings.Join(output, "\n\n")
}

// Unwrap returns the underlying execution error.
func (e *CommandError) Unwrap() error {
	return e.Err
}
