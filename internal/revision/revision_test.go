package revision_test

import (
	"strings"
	"testing"

	"github.com/mkvlrn/arch-setup/internal/revision"
)

func TestValidate(t *testing.T) {
	for _, value := range []string{strings.Repeat("a", 40), strings.Repeat("b", 64)} {
		if err := revision.Validate(value); err != nil {
			t.Errorf("valid revision %q rejected: %v", value, err)
		}
	}

	for _, value := range []string{
		"", "development", "main", "abc1234", "--help",
		strings.Repeat("g", 40), strings.Repeat("A", 40), strings.Repeat("a", 41),
	} {
		if err := revision.Validate(value); err == nil {
			t.Errorf("invalid revision %q accepted", value)
		}
	}
}
