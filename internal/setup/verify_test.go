package setup_test

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/mkvlrn/arch-setup/internal/setup"
)

func TestVerifyWithNoChecks(t *testing.T) {
	var output bytes.Buffer

	if err := setup.Verify(t.Context(), &output, nil); err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	if got := output.String(); got != "Verification passed.\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestVerifyExecutesChecksInOrder(t *testing.T) {
	var output bytes.Buffer

	var calls []string

	ctx := t.Context()
	checks := []setup.Check{
		{
			Name: "first",
			Run: func(got context.Context) error {
				if got != ctx {
					t.Error("first check did not receive the caller's context")
				}

				calls = append(calls, "first")

				return nil
			},
		},
		{
			Name: "second",
			Run: func(got context.Context) error {
				if got != ctx {
					t.Error("second check did not receive the caller's context")
				}

				calls = append(calls, "second")

				return nil
			},
		},
	}

	if err := setup.Verify(ctx, &output, checks); err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	if !reflect.DeepEqual(calls, []string{"first", "second"}) {
		t.Fatalf("unexpected execution order: %v", calls)
	}

	expected := "[1/2] first\n[2/2] second\nVerification passed.\n"

	if got := output.String(); got != expected {
		t.Fatalf("expected output %q, got %q", expected, got)
	}
}

func TestVerifyContinuesAfterFailuresAndCombinesErrors(t *testing.T) {
	var output bytes.Buffer

	var calls []string

	firstErr := errors.New("first failure")
	lastErr := errors.New("last failure")

	checks := []setup.Check{
		{
			Name: "first",
			Run: func(context.Context) error {
				calls = append(calls, "first")

				return firstErr
			},
		},
		{
			Name: "middle",
			Run: func(context.Context) error {
				calls = append(calls, "middle")

				return nil
			},
		},
		{
			Name: "last",
			Run: func(context.Context) error {
				calls = append(calls, "last")

				return lastErr
			},
		},
	}

	err := setup.Verify(t.Context(), &output, checks)
	if err == nil {
		t.Fatal("expected an error")
	}

	if !reflect.DeepEqual(calls, []string{"first", "middle", "last"}) {
		t.Fatalf("not all checks ran in order: %v", calls)
	}

	for _, cause := range []error{firstErr, lastErr} {
		if !errors.Is(err, cause) {
			t.Errorf("expected combined error to wrap %v, got %v", cause, err)
		}
	}

	for _, name := range []string{`verify "first"`, `verify "last"`} {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("expected error to contain %q, got %v", name, err)
		}
	}

	expected := "[1/3] first\n[2/3] middle\n[3/3] last\n"

	if got := output.String(); got != expected {
		t.Fatalf("expected output %q, got %q", expected, got)
	}
}
