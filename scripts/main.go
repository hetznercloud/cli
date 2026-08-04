package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	"github.com/hetznercloud/cli/internal/cli"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/cli/internal/version"
)

//go:generate go run $GOFILE docs

func main() {
	if err := run(os.Args[1:]); err != nil {
		if _, writeErr := fmt.Fprintf(os.Stderr, "hcloud documentation generator: %v\n", err); writeErr != nil {
			os.Exit(1)
		}
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) < 1 {
		return errors.New("usage: docs|manpages")
	}

	switch arg := args[0]; arg {
	case "docs":
		return generateDocs()
	case "manpages":
		return generateManPages()
	default:
		return fmt.Errorf("unknown argument: %s", strconv.Quote(arg))
	}
}

func generateDocs() error {
	dir := "../docs/reference/manual"
	if err := ensureEmptyDir(dir); err != nil {
		return err
	}

	return withRootCommand(true, func(cmd *cobra.Command) error {
		return doc.GenMarkdownTree(cmd, dir)
	})
}

func generateManPages() error {
	dir := "./manpages"
	if err := ensureEmptyDir(dir); err != nil {
		return err
	}

	return withRootCommand(true, func(cmd *cobra.Command) error {
		return doc.GenManTree(cmd, &doc.GenManHeader{
			Source: version.Version,
			Manual: "CLI for Hetzner API",
		}, dir)
	})
}

func ensureEmptyDir(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("could not remove directory: %w", err)
	}

	if err := os.MkdirAll(dir, 0755); err != nil { //nolint:gosec
		return fmt.Errorf("error creating directory: %w", err)
	}
	return nil
}

func withRootCommand(withMdTables bool, run func(*cobra.Command) error) (err error) {
	s, err := state.New(config.New(), state.Options{Stderr: os.Stderr})
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, s.Close())
	}()

	cmd, err := cli.NewRootCommand(s, withMdTables)
	if err != nil {
		return err
	}
	return run(cmd)
}
