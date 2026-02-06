package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra/doc"

	"github.com/hetznercloud/cli/internal/cli"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

//go:generate go run $GOFILE

func run() error {
	// Define the directory where the docs will be generated
	dir := "../docs/reference/manual"

	// Clean the directory to make sure outdated files don't persist
	err := os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("could not remove docs directory: %w", err)
	}

	// Create directory again
	if err := os.MkdirAll(dir, 0755); err != nil { //nolint:gosec
		return fmt.Errorf("error creating docs directory: %w", err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, state.ContextKeyMarkdownTables{}, true)
	cfg := config.New()
	s, err := state.New(ctx, cfg)
	if err != nil {
		return fmt.Errorf("could not create state: %w", err)
	}

	cmd := cli.NewRootCommand(s)

	// Generate the docs
	return doc.GenMarkdownTree(cmd, dir)
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
