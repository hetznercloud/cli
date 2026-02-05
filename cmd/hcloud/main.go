package main

import (
	"context"
	"log"
	"os"

	"github.com/hetznercloud/cli/internal/cli"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("hcloud: ")
	log.SetOutput(os.Stderr)
}

func main() {
	ctx := context.Background()

	cfg := config.New()
	if err := cfg.Read(nil); err != nil {
		log.Fatalf("unable to read config file \"%s\": %s\n", cfg.Path(), err)
	}

	s, err := state.New(cfg, ctx)
	if err != nil {
		log.Fatalln(err)
	}

	rootCommand := cli.NewRootCommand(s)

	if err := rootCommand.Execute(); err != nil {
		log.Fatalln(util.FormatHcloudError(err))
	}
}
