package main

import (
	"log"
	"os"

	"github.com/hetznercloud/cli/internal/cli"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("hcloud: ")
	log.SetOutput(os.Stderr)
}

func main() {
	cfg := config.New()
	if err := cfg.Read(nil); err != nil {
		log.Fatalf("unable to read config file \"%s\": %s\n", cfg.Path(), err)
	}

	s, err := state.New(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	rootCommand := cli.NewRootCommand(s)

	if err := rootCommand.Execute(); err != nil {
		log.Fatalln(err)
	}
}
