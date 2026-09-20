package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

// Load reads config.json for setup data.
func Load(configData []byte) (Config, error) {
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

	config.Env.CI = os.Getenv("CI") != ""
	config.Machine.HomeDir = homeDir
	config.Machine.Username = currentUser.Username
	config.Machine.RepoDir = filepath.Join(homeDir, "repos", "arch-setup")
	config.Machine.TempDir = os.TempDir()

	return config, nil
}
