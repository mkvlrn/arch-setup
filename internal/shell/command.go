package shell

// Command describes an external command to execute.
type Command struct {
	// Name identifies the command in errors and diagnostic output.
	Name string

	// Path is the executable name or path.
	Path string

	// Args contains the arguments passed to the executable.
	Args []string

	// Dir sets the command's working directory when non-empty.
	Dir string

	// Sudo runs the command through sudo.
	Sudo bool
}
