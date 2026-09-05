package bootstrap

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Bootstrap reads config.json for setup data.
func Bootstrap(configData []byte, miseConfigData []byte) (Config, error) {
	var config Config

	var mise miseConfig

	if err := json.Unmarshal(configData, &config); err != nil {
		return Config{}, fmt.Errorf("decode embed config file: %w", err)
	}

	if err := toml.Unmarshal(miseConfigData, &mise); err != nil {
		return Config{}, fmt.Errorf("decode embed mise config file: %w", err)
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

	config.MiseTools = make([]string, 0, len(mise.Tools))
	for tool := range mise.Tools {
		config.MiseTools = append(config.MiseTools, tool)
	}

	return config, nil
}
