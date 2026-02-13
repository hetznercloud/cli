package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra/doc"

	"github.com/hetznercloud/cli/internal/cli"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/cli/internal/version"
)

const directory = "./manpages"

func main() {
	//nolint:gosec
	if err := os.MkdirAll(directory, 0755); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, state.ContextKeyMarkdownTables{}, true)
	s, _ := state.New(ctx, config.New())
	rootCommand := cli.NewRootCommand(s)

	err := doc.GenManTree(rootCommand, &doc.GenManHeader{
		Source: version.Version,
		Manual: "CLI for Hetzner API",
	}, directory)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Man pages generated successfully")
}
