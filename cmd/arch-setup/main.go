// Package main is the application entrypoint.
package main

import (
	"bufio"
	"context"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mkvlrn/arch-setup/internal/execution"
)

//go:embed config.json
var configData []byte

//go:embed arch-stow-bootstrap
var archStowBootstrap []byte

func main() {
	verifyOnly := flag.Bool("verify", false, "verify the installed system without modifying it")

	flag.Parse()

	var passphrase []byte

	if !*verifyOnly && os.Getenv("CI") != "true" {
		var err error

		passphrase, err = readPassphrase()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)

			os.Exit(1)
		}
	}

	if err := execution.Run(context.Background(), configData, archStowBootstrap, passphrase, *verifyOnly); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}

func readPassphrase() ([]byte, error) {
	if err := exec.Command("stty", "-echo").Run(); err != nil {
		return nil, fmt.Errorf("disable terminal echo: %w", err)
	}
	defer func() { _ = exec.Command("stty", "echo").Run() }()

	_, _ = fmt.Fprint(os.Stderr, "Enter arch-stow SSH key passphrase: ")
	passphrase, err := bufio.NewReader(os.Stdin).ReadString('\n')
	_, _ = fmt.Fprintln(os.Stderr)

	return []byte(strings.TrimSpace(passphrase)), err
}
