package shell_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/mkvlrn/arch-setup/internal/shell"
)

func TestCommandError(t *testing.T) {
	err := &shell.CommandError{
		Name:   "test command",
		Err:    errors.New("exit status 1"),
		Stdout: "some stdout\n",
		Stderr: "some stderr\n",
	}
	formatted := err.Error()
	expected := []string{
		"could not run test command: exit status 1",
		"stdout:\nsome stdout",
		"stderr:\nsome stderr",
	}

	for _, value := range expected {
		if !strings.Contains(formatted, value) {
			t.Errorf(
				"expected error to contain %q:\n%s",
				value,
				formatted,
			)
		}
	}
}

func TestCommandErrorIgnoresEmptyOutput(t *testing.T) {
	err := &shell.CommandError{
		Name:   "test command",
		Err:    errors.New("exit status 1"),
		Stdout: "\n",
		Stderr: "",
	}
	formatted := err.Error()

	if strings.Contains(formatted, "stdout:") {
		t.Errorf("unexpected stdout section:\n%s", formatted)
	}

	if strings.Contains(formatted, "stderr:") {
		t.Errorf("unexpected stderr section:\n%s", formatted)
	}
}

func TestCommandErrorUnwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &shell.CommandError{
		Name: "test command",
		Err:  cause,
	}

	if !errors.Is(err, cause) {
		t.Fatal("expected CommandError to wrap its cause")
	}
}
