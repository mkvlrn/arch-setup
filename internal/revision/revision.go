// Package revision identifies and validates the installer's source commit.
package revision

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// Commit is set by make build to the source checkout's commit rather than a branch name.
var Commit = "development"

// Validate rejects unstamped builds and anything other than a full Git object ID.
func Validate(value string) error {
	// Git supports SHA-1 and SHA-256 repositories; only full object IDs are accepted.
	const (
		sha1Length   = 40
		sha256Length = 64
	)

	if len(value) == sha1Length || len(value) == sha256Length {
		if _, err := hex.DecodeString(value); err == nil && value == strings.ToLower(value) {
			return nil
		}
	}

	return fmt.Errorf("invalid installer revision %q: use make build from a clean checkout", value)
}
