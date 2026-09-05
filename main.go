// Package main is the application entrypoint.
package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/mkvlrn/arch-setup/internal/app"
)

//go:embed config.json
var configData []byte

//go:embed stow/user/.config/mise/config.toml
var miseConfigData []byte

func main() {
	ci := os.Getenv("GITHUB_ACTIONS") == "true"
	verifyOnly := flag.Bool("verify", false, "verify the installed system without modifying it")

	flag.Parse()

	if err := app.Run(context.Background(), configData, miseConfigData, *verifyOnly, ci); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
