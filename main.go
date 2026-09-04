// Package main is the application entrypoint.
package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/start"
	"github.com/mkvlrn/arch-setup/internal/steps"
	"github.com/mkvlrn/arch-setup/internal/sudo"
)

//go:embed config.json
var configData []byte

//go:embed stow/user/.config/mise/config.toml
var miseConfigData []byte

func main() {
	ci := os.Getenv("GITHUB_ACTIONS") == "true"
	verifyOnly := flag.Bool("verify", false, "verify the installed system without modifying it")

	flag.Parse()

	if err := run(context.Background(), *verifyOnly, ci); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}

func run(ctx context.Context, verifyOnly bool, ci bool) error {
	config, err := start.Bootstrap(configData, miseConfigData)
	if err != nil {
		return err
	}

	if verifyOnly {
		return runVerification(ctx, config, ci)
	}

	stopSudo, err := sudo.KeepAlive(ctx)
	if err != nil {
		return err
	}
	defer stopSudo()

	return runSetup(ctx, config, ci)
}

func runSetup(ctx context.Context, config start.Config, ci bool) error {
	plan := []setup.Step{
		steps.InstallPkg(steps.UsePacman, config.BasePackages),
	}

	if ci {
		plan = append(
			plan,
			steps.ExistingRepo(config.RepoSSH, config.RepoDir),
		)
	} else {
		plan = append(
			plan,
			steps.RemovePkg(config.RemovePackages),
			steps.CloneRepo(config.RepoHTTP, config.RepoSSH, config.RepoDir),
		)
	}

	plan = append(
		plan,
		steps.Stow(steps.StowSystem, config.RepoDir, config.HomeDir),
		steps.Yay(config.TempDir, config.MirrorListPath),
		steps.InstallPkg(steps.UseYay, config.MainPackages),
		steps.Xdg(config.XdgMkDir, config.XdgRmRf, config.HomeDir),
		steps.Stow(steps.StowUser, config.RepoDir, config.HomeDir),
		steps.Mise(config.HomeDir, config.MiseTools),
		steps.User(config.Username, config.HomeDir),
	)

	return setup.Run(ctx, os.Stdout, plan)
}

func runVerification(ctx context.Context, config start.Config, ci bool) error {
	packages := append(
		append([]string{}, config.BasePackages...),
		config.MainPackages...,
	)

	checks := []setup.Check{
		steps.CloneRepoVerify(config.RepoSSH, config.RepoDir),
		steps.StowVerify(steps.StowSystem, config.RepoDir, config.HomeDir),
		steps.YayVerify(config.MirrorListPath, config.MirrorListCheck),
		steps.InstallPkgVerify(packages),
	}

	if !ci {
		checks = append(checks, steps.RemovePkgVerify(config.RemovePackages))
	}

	checks = append(
		checks,
		steps.XdgVerify(config.XdgMkDir, config.XdgRmRf, config.HomeDir),
		steps.StowVerify(steps.StowUser, config.RepoDir, config.HomeDir),
		steps.MiseVerify(config.HomeDir),
		steps.UserVerify(config.Username, config.HomeDir),
	)

	return setup.Verify(ctx, os.Stdout, checks)
}
