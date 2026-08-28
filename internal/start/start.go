package start

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

//go:embed config.json
var configData []byte

// Bootstrap reads config.json for setup data.
func Bootstrap() (Config, error) {
	var config Config

	if err := json.Unmarshal(configData, &config); err != nil {
		return Config{}, fmt.Errorf("decode embed config file: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("get user home dir: %w", err)
	}

	currentUser, err := user.Current()
	if err != nil {
		return Config{}, fmt.Errorf("get username: %w", err)
	}

	config.Username = currentUser.Username
	config.HomeDir = homeDir
	config.RepoDir = filepath.Join(homeDir, "repos", "arch-setup")
	config.TempDir = os.TempDir()

	return config, nil
}
