package setup

import "context"

// Check describes one machine-state verification.
type Check struct {
	// Name identifies the verification in progress messages and errors.
	Name string

	// Run performs the verification without modifying the machine.
	Run func(context.Context) error
}
