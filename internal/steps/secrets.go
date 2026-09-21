package steps

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"filippo.io/age"
	"golang.org/x/term"

	"github.com/mkvlrn/arch-setup/internal/config"
	"github.com/mkvlrn/arch-setup/internal/setup"
)

// Secrets decrypts the embedded secrets archive into the user's home directory.
func Secrets(cfg *config.Config, secretsData []byte) setup.Step {
	if cfg.Env.CI {
		return setup.Step{Name: "Restore machine secrets (skipped in CI)"}
	}

	return setup.Step{
		Name: "Restore machine secrets",
		Run: func(context.Context) error {
			passphrase, err := readPassphrase()
			if err != nil {
				return fmt.Errorf("read secrets passphrase: %w", err)
			}

			return restoreSecrets(secretsData, cfg.Machine.HomeDir, passphrase)
		},
	}
}

func readPassphrase() ([]byte, error) {
	terminal := os.Stdin
	closeTerminal := func() {}

	if !term.IsTerminal(int(terminal.Fd())) {
		var err error

		terminal, err = os.Open("/dev/tty")
		if err != nil {
			return nil, fmt.Errorf("open controlling terminal: %w", err)
		}

		closeTerminal = func() { _ = terminal.Close() }
	}

	defer closeTerminal()

	_, _ = fmt.Fprint(terminal, "Enter secrets passphrase: ")
	passphrase, err := term.ReadPassword(int(terminal.Fd()))
	_, _ = fmt.Fprintln(terminal)

	return passphrase, err
}

func restoreSecrets(encrypted []byte, homeDir string, passphrase []byte) error {
	identity, err := age.NewScryptIdentity(string(passphrase))
	if err != nil {
		return fmt.Errorf("create age identity: %w", err)
	}

	reader, err := age.Decrypt(bytes.NewReader(encrypted), identity)
	if err != nil {
		return fmt.Errorf("decrypt secrets: %w", err)
	}

	return extractSecrets(tar.NewReader(reader), homeDir)
}

func extractSecrets(archive *tar.Reader, homeDir string) error {
	for {
		header, err := archive.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("read secrets archive: %w", err)
		}

		if err := extractSecretEntry(archive, homeDir, header); err != nil {
			return err
		}
	}
}

func extractSecretEntry(archive *tar.Reader, homeDir string, header *tar.Header) error {
	target, err := secureSecretPath(homeDir, header.Name)
	if err != nil {
		return err
	}

	mode, err := secretFileMode(header.Mode)
	if err != nil {
		return fmt.Errorf("invalid mode for secret archive entry %q: %w", header.Name, err)
	}

	switch header.Typeflag {
	case tar.TypeDir:
		if err := os.MkdirAll(target, mode); err != nil {
			return fmt.Errorf("create secret directory %q: %w", target, err)
		}
	case tar.TypeReg:
		return extractSecretFile(archive, target, mode)
	default:
		return fmt.Errorf("unsupported secret archive entry %q", header.Name)
	}

	return nil
}

func extractSecretFile(archive *tar.Reader, target string, mode os.FileMode) error {
	const secretParentMode = 0o700

	if err := os.MkdirAll(filepath.Dir(target), secretParentMode); err != nil {
		return fmt.Errorf("create secret parent %q: %w", target, err)
	}

	// #nosec G304 -- target is validated by secureSecretPath.
	file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("create secret file %q: %w", target, err)
	}

	_, copyErr := io.Copy(file, archive)
	closeErr := file.Close()

	if copyErr != nil {
		return fmt.Errorf("write secret file %q: %w", target, copyErr)
	}

	if closeErr != nil {
		return fmt.Errorf("close secret file %q: %w", target, closeErr)
	}

	return nil
}

func secretFileMode(mode int64) (os.FileMode, error) {
	const maxFileMode = 0o7777

	if mode < 0 || mode > maxFileMode {
		return 0, fmt.Errorf("mode %o is outside the supported range", mode)
	}

	return os.FileMode(mode), nil
}

func secureSecretPath(homeDir, name string) (string, error) {
	if filepath.IsAbs(name) {
		return "", fmt.Errorf("secret archive contains absolute path %q", name)
	}

	target := filepath.Join(homeDir, filepath.FromSlash(name))

	relative, err := filepath.Rel(homeDir, target)
	if err != nil || relative == ".." || len(relative) > 3 && relative[:4] == ".."+string(filepath.Separator) {
		return "", fmt.Errorf("secret archive path escapes home directory: %q", name)
	}

	return target, nil
}
