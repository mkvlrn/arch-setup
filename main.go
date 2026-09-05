// Package main is the application entrypoint.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
)

func main() {
	ci := os.Getenv("GITHUB_ACTIONS") == "true"
	verifyOnly := flag.Bool("verify", false, "verify the installed system without modifying it")

	flag.Parse()

	if err := run(context.Background(), *verifyOnly, ci); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
