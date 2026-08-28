package shell_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/mkvlrn/arch-setup/internal/shell"
)

func TestRunSucceeds(t *testing.T) {
	results, err := shell.Run(t.Context(), []shell.Command{
		{
			Name: "print value",
			Path: "printf",
			Args: []string{"some stdout"},
		},
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected one result, got %d", len(results))
	}

	if results[0].Stdout != "some stdout" {
		t.Fatalf("expected stdout %q, got %q", "some stdout", results[0].Stdout)
	}
}

func TestRunReturnsCommandError(t *testing.T) {
	_, err := shell.Run(t.Context(), []shell.Command{{
		Name: "failing command",
		Path: "false",
	}})
	if err == nil {
		t.Fatal("expected an error")
	}

	var commandErr *shell.CommandError
	if !errors.As(err, &commandErr) {
		t.Fatalf("expected CommandError, got %T", err)
	}

	if commandErr.Name != "failing command" {
		t.Errorf("expected command name %q, got %q", "failing command", commandErr.Name)
	}
}

func TestRunCapturesOutput(t *testing.T) {
	_, err := shell.Run(t.Context(), []shell.Command{
		{
			Name: "failing shell command",
			Path: "sh",
			Args: []string{"-c", "printf 'some stdout'; printf 'some stderr' >&2; exit 1"},
		},
	})
	if err == nil {
		t.Fatal("expected an error")
	}

	var commandErr *shell.CommandError
	if !errors.As(err, &commandErr) {
		t.Fatalf("expected CommandError, got %T", err)
	}

	if commandErr.Stdout != "some stdout" {
		t.Errorf("expected stdout %q, got %q", "some stdout", commandErr.Stdout)
	}

	if commandErr.Stderr != "some stderr" {
		t.Errorf("expected stderr %q, got %q", "some stderr", commandErr.Stderr)
	}
}

func TestRunUsesWorkingDirectory(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "created")

	_, err := shell.Run(t.Context(), []shell.Command{
		{
			Name: "create file",
			Path: "touch",
			Args: []string{"created"},
			Dir:  dir,
		},
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	if _, err := os.Stat(file); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestRunStopsAfterFailure(t *testing.T) {
	file := filepath.Join(t.TempDir(), "unexpected")

	_, err := shell.Run(t.Context(), []shell.Command{
		{
			Name: "fail",
			Path: "false",
		},
		{
			Name: "must not run",
			Path: "touch",
			Args: []string{file},
		},
	})
	if err == nil {
		t.Fatal("expected an error")
	}

	if _, err := os.Stat(file); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("second command ran unexpectedly")
	}
}
