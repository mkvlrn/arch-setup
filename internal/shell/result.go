package shell

// Result contains output captured from a successful command.
type Result struct {
	// Stdout contains the command's standard output.
	Stdout string

	// Stderr contains the command's standard error.
	Stderr string
}
