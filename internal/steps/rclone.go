package steps

import (
	"bufio"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

type protonSecrets struct {
	username string
	password string
	totp     string
}

// RcloneProton sets up the rclone configuration for the proton drive.
func RcloneProton(ctx context.Context, homeDir string) setup.Step {
	secrets, err := loadSecrets(filepath.Join(homeDir, ".config", "rclone", "proton-secrets"))
	if err != nil {
		panic(err)
	}

	otpSecret, err := getOtpSecret(ctx, secrets.totp)
	if err != nil {
		panic(err)
	}

	args := []string{"config", "create", "proton", "protondrive"}
	args = append(args, "username="+secrets.username, "password="+secrets.password, "otp_secret="+otpSecret[0].Stdout)

	return setup.Step{
		Name: "rclone-proton",
		Commands: []shell.Command{
			{
				Name: "rclone",
				Args: args,
			},
		},
	}
}

func loadSecrets(path string) (protonSecrets, error) {
	// #nosec G304 -- path is controlled
	file, err := os.Open(path)
	if err != nil {
		return protonSecrets{}, err
	}
	defer func() {
		_ = file.Close()
	}()

	envMap := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, val, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		envMap[strings.TrimSpace(key)] = strings.Trim(strings.TrimSpace(val), `'"`)
	}

	if err := scanner.Err(); err != nil {
		return protonSecrets{}, err
	}

	username, ok := envMap["username"]
	if !ok {
		return protonSecrets{}, errors.New("reading proton username from secrets")
	}

	password, ok := envMap["password"]
	if !ok {
		return protonSecrets{}, errors.New("reading proton password from secrets")
	}

	totp, ok := envMap["totp"]
	if !ok {
		return protonSecrets{}, errors.New("reading proton totp from secrets")
	}

	return protonSecrets{
		username: username,
		password: password,
		totp:     totp,
	}, nil
}

func getOtpSecret(ctx context.Context, totp string) ([]shell.Result, error) {
	return shell.Run(ctx, []shell.Command{
		{
			Name: "rclone",
			Args: []string{"obscure", totp},
		},
	})
}
