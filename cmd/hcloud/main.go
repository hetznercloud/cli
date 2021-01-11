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
	state := state.New()

	if state.ConfigPath != "" {
		_, err := os.Stat(state.ConfigPath)
		switch {
		case err == nil:
			if err := state.ReadConfig(); err != nil {
				log.Fatalf("unable to read config file %q: %s\n", state.ConfigPath, err)
			}
		case os.IsNotExist(err):
			break
		default:
			log.Fatalf("unable to read config file %q: %s\n", state.ConfigPath, err)
		}
	}

	state.ReadEnv()

	rootCommand := cli.NewRootCommand(state)
	if err := rootCommand.Execute(); err != nil {
		log.Fatalln(err)
	}
}
