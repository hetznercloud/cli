package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hetznercloud/cli/internal/cli"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

var defaultConfigPath string

func main() {
	os.Exit(execute())
}

func execute() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		if !util.IsSilent(err) {
			if _, writeErr := fmt.Fprintf(os.Stderr, "hcloud: %s\n", util.FormatHcloudError(err)); writeErr != nil {
				return 1
			}
		}
		return util.ExitCode(err)
	}
	return 0
}

func run(ctx context.Context) (err error) {
	cfg := config.New(
		config.WithArgs(os.Args[1:]),
		config.WithWarningWriter(os.Stderr),
		config.WithDefaultPath(defaultConfigPath),
	)
	if err := cfg.LoadDefault(); err != nil {
		return fmt.Errorf("unable to read config file %q: %w", cfg.Path(), err)
	}

	s, err := state.New(cfg, state.Options{Stderr: os.Stderr})
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, s.Close())
	}()

	rootCommand, err := cli.NewRootCommand(s, false)
	if err != nil {
		return err
	}

	return rootCommand.ExecuteContext(ctx)
}
