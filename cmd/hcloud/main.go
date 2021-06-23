package main

import (
	"log"
	"os"

	"github.com/hetznercloud/cli/internal/cli"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("hcloud: ")
	log.SetOutput(os.Stderr)
}

func main() {
	cliState := state.New()

	if cliState.ConfigPath != "" {
		_, err := os.Stat(cliState.ConfigPath)
		switch {
		case err == nil:
			if err := cliState.ReadConfig(); err != nil {
				log.Fatalf("unable to read config file %q: %s\n", cliState.ConfigPath, err)
			}
		case os.IsNotExist(err):
			break
		default:
			log.Fatalf("unable to read config file %q: %s\n", cliState.ConfigPath, err)
		}
	}

	cliState.ReadEnv()
	apiClient := hcapi2.NewClient(cliState.Client())
	rootCommand := cli.NewRootCommand(cliState, apiClient)
	if err := rootCommand.Execute(); err != nil {
		log.Fatalln(err)
	}
}
