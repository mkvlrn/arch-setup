package steps

import (
	"bytes"
	"encoding/base64"
	"path/filepath"

	"github.com/mkvlrn/arch-setup/internal/config"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// Secrets decrypts the embedded secrets archive into the user's home directory.
func Secrets(cfg *config.Config, secretsData []byte, passphrase []byte) setup.Step {
	if cfg.Env.CI {
		return setup.Step{Name: "Restore machine secrets (skipped in CI)"}
	}

	homeDir := cfg.Machine.HomeDir
	encodedSecrets := base64.StdEncoding.EncodeToString(secretsData)

	return setup.Step{
		Name: "Restore machine secrets",
		Commands: []shell.Command{{
			Name: "decrypt and restore machine secrets",
			Path: "sh",
			Args: []string{
				"-c",
				`set -eu
archive=$(mktemp)
trap 'rm -f "$archive"' EXIT
printf '%s' "$1" | base64 -d >"$archive"
age -d "$archive" | tar -C "$2" -xf -`,
				"restore-secrets",
				encodedSecrets,
				filepath.Clean(homeDir),
			},
			Stdin: bytes.NewReader(append(append([]byte{}, passphrase...), '\n')),
		}},
	}
}
