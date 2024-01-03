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
	s, err := state.New()
	if err != nil {
		log.Fatalln(err)
	}

	rootCommand := cli.NewRootCommand(s)
	if err := rootCommand.Execute(); err != nil {
		log.Fatalln(err)
	}
}
