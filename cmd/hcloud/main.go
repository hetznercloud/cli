package main

import (
	"log"
	"os"

	"github.com/hetznercloud/cli/internal/cli"
	"github.com/hetznercloud/cli/internal/state"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("hcloud: ")
	log.SetOutput(os.Stderr)
}

func main() {
	configPath := os.Getenv("HCLOUD_CONFIG")
	if configPath == "" {
		configPath = state.DefaultConfigPath()
	}

	cfg, err := state.ReadConfig(configPath)
	if err != nil {
		log.Fatalf("unable to read config file %q: %s\n", configPath, err)
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
