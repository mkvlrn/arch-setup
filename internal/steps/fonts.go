package steps

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mkvlrn/arch-setup/internal/config"
	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/shell"
)

// Fonts installs nerd fonts with getnf.
func Fonts(cfg *config.Config) setup.Step {
	fonts := cfg.GetNF
	misePath := filepath.Join(cfg.Machine.HomeDir, ".local", "bin", "mise")
	stepName := fmt.Sprintf("Installing %d fonts with getnf", len(fonts))
	args := []string{"-i", strings.Join(fonts, ",")}

	return setup.Step{
		Name: stepName,
		Commands: []shell.Command{
			{
				Name: "install fonts managed by getnf",
				Path: misePath,
				Args: append([]string{"exec", "--", "getnf"}, args...),
				Env:  []string{"TERM=xterm-256color"},
			},
			{
				Name: "clear fonts cache",
				Path: "fc-cache",
				Args: []string{"-f"},
			},
		},
	}
}
