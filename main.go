// Package main is the application entrypoint.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/mkvlrn/arch-setup/internal/setup"
	"github.com/mkvlrn/arch-setup/internal/start"
	"github.com/mkvlrn/arch-setup/internal/steps"
	"github.com/mkvlrn/arch-setup/internal/sudo"
)

func main() {
	repoReady := os.Getenv("GITHUB_ACTIONS") == "true"
	verifyOnly := flag.Bool("verify", false, "verify the installed system without modifying it")

	flag.Parse()

	if err := run(context.Background(), *verifyOnly, repoReady); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}

func run(ctx context.Context, verifyOnly bool, repoReady bool) error {
	config, err := start.Bootstrap()
	if err != nil {
		return err
	}

	if verifyOnly {
		return runVerification(ctx, config)
	}

	stopSudo, err := sudo.KeepAlive(ctx)
	if err != nil {
		return err
	}
	defer stopSudo()

	return runSetup(ctx, config, repoReady)
}

func runSetup(ctx context.Context, config start.Config, repoReady bool) error {
	plan := []setup.Step{
		steps.InstallPkg(steps.UsePacman, config.BasePackages),
	}

	if repoReady {
		plan = append(
			plan,
			steps.ExistingRepo(config.RepoSSH, config.RepoDir),
		)
	} else {
		plan = append(
			plan,
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
		steps.Mise(config.HomeDir),
		steps.User(config.Username, config.HomeDir),
	)

	return setup.Run(ctx, os.Stdout, plan)
}

func runVerification(ctx context.Context, config start.Config) error {
	packages := append(
		append([]string{}, config.BasePackages...),
		config.MainPackages...,
	)

	checks := []setup.Check{
		steps.CloneRepoVerify(config.RepoSSH, config.RepoDir),
		steps.StowVerify(steps.StowSystem, config.RepoDir, config.HomeDir),
		steps.YayVerify(config.MirrorListPath, config.MirrorListCheck),
		steps.InstallPkgVerify(packages),
		steps.XdgVerify(config.XdgMkDir, config.XdgRmRf, config.HomeDir),
		steps.StowVerify(steps.StowUser, config.RepoDir, config.HomeDir),
		steps.MiseVerify(config.HomeDir),
		steps.UserVerify(config.Username, config.HomeDir),
	}

	return setup.Verify(ctx, os.Stdout, checks)
}
