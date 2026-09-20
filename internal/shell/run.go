package shell

import (
	"bytes"
	"context"
	"os"
	"os/exec"
)

// Run executes commands sequentially and stops at the first failure.
func Run(ctx context.Context, commands []Command) ([]Result, error) {
	results := make([]Result, 0, len(commands))

	for _, command := range commands {
		var stdout bytes.Buffer

		var stderr bytes.Buffer

		path := command.Path
		args := command.Args

		if command.Sudo {
			path = "sudo"

			args = append([]string{command.Path}, command.Args...)
		}

		// #nosec G204 -- commands are defined internally, not derived from user input
		cmd := exec.CommandContext(ctx, path, args...)

		cmd.Dir = command.Dir
		cmd.Env = append(os.Environ(), command.Env...)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return results, &CommandError{
				Name:   command.Name,
				Err:    err,
				Stdout: stdout.String(),
				Stderr: stderr.String(),
			}
		}

		results = append(results, Result{Stdout: stdout.String(), Stderr: stderr.String()})
	}

	return results, nil
}
