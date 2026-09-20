// Package main is the application entrypoint.
package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/mkvlrn/arch-setup/internal/execution"
)

//go:embed config.json
var configData []byte

func main() {
	verifyOnly := flag.Bool("verify", false, "verify the installed system without modifying it")

	flag.Parse()

	if err := execution.Run(context.Background(), configData, *verifyOnly); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
