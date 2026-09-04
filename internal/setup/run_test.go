package setup_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

func TestRunWithNoSteps(t *testing.T) {
	var output bytes.Buffer

	err := setup.Run(t.Context(), &output, nil)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	expected := "Done.\n"

	if output.String() != expected {
		t.Fatalf(
			"expected output %q, got %q",
			expected,
			output.String(),
		)
	}
}

func TestRunExecutesStepsInOrder(t *testing.T) {
	var output bytes.Buffer

	dir := t.TempDir()
	result := filepath.Join(dir, "result")
	steps := []setup.Step{
		{
			Name: "first",
			Commands: []shell.Command{
				{
					Name: "write first value",
					Path: "sh",
					Args: []string{"-c", "printf 'first\\n' >> \"$1\"", "sh", result},
				},
			},
		},
		{
			Name: "second",
			Commands: []shell.Command{
				{
					Name: "write second value",
					Path: "sh",
					Args: []string{"-c", "printf 'second\\n' >> \"$1\"", "sh", result},
				},
			},
		},
	}

	if err := setup.Run(t.Context(), &output, steps); err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	// #nosec G304 -- result is a test-controlled path inside t.TempDir
	content, err := os.ReadFile(result)
	if err != nil {
		t.Fatalf("read result: %v", err)
	}

	if string(content) != "first\nsecond\n" {
		t.Fatalf("unexpected execution order: %q", content)
	}

	expectedOutput := strings.Join([]string{
		"[1/2] first",
		"[2/2] second",
		"Done.",
		"",
	}, "\n")

	if output.String() != expectedOutput {
		t.Fatalf(
			"expected output %q, got %q",
			expectedOutput,
			output.String(),
		)
	}
}

func TestRunStopsAfterFailure(t *testing.T) {
	var output bytes.Buffer

	file := filepath.Join(t.TempDir(), "unexpected")
	steps := []setup.Step{
		{
			Name: "failing step",
			Commands: []shell.Command{
				{
					Name: "fail",
					Path: "false",
				},
			},
		},
		{
			Name: "later step",
			Commands: []shell.Command{{
				Name: "must not run",
				Path: "touch",
				Args: []string{file},
			}},
		},
	}

	err := setup.Run(t.Context(), &output, steps)
	if err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(err.Error(), `run step "failing step"`) {
		t.Fatalf("step name missing from error: %v", err)
	}

	if _, err := os.Stat(file); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("later step ran unexpectedly")
	}

	expectedOutput := "[1/2] failing step\n"

	if output.String() != expectedOutput {
		t.Fatalf(
			"expected output %q, got %q",
			expectedOutput,
			output.String(),
		)
	}
}
