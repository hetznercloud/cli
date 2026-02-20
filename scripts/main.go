package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	"github.com/hetznercloud/cli/internal/cli"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/cli/internal/version"
)

//go:generate go run $GOFILE docs

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: docs|manpages")
	}

	var err error
	switch arg := os.Args[1]; arg {
	case "docs":
		err = generateDocs()
	case "manpages":
		err = generateManPages()
	default:
		log.Fatalln("Unknown argument:", arg)
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func generateDocs() error {
	dir := "../docs/reference/manual"
	if err := ensureEmptyDir(dir); err != nil {
		return err
	}

	cmd, err := newRootCommand(true)
	if err != nil {
		return err
	}

	// Generate the docs
	return doc.GenMarkdownTree(cmd, dir)
}

func generateManPages() error {
	dir := "./manpages"
	if err := ensureEmptyDir(dir); err != nil {
		return err
	}

	cmd, err := newRootCommand(true)
	if err != nil {
		return err
	}

	return doc.GenManTree(cmd, &doc.GenManHeader{
		Source: version.Version,
		Manual: "CLI for Hetzner API",
	}, dir)
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

func newRootCommand(withMdTables bool) (*cobra.Command, error) {
	ctx := context.Background()
	if withMdTables {
		ctx = context.WithValue(ctx, state.ContextKeyMarkdownTables{}, true)
	}
	s, err := state.New(ctx, config.New())
	if err != nil {
		return nil, err
	}
	return cli.NewRootCommand(s), nil
}
