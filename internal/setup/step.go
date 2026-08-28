package setup

import "github.com/mkvlrn/arch-setup/internal/shell"

// Step groups the commands belonging to one stage of system setup.
type Step struct {
	// Name identifies the step in progress messages and errors.
	Name string

	// Commands contains the external commands executed by the step.
	Commands []shell.Command
}
