package main

import (
	"log"
	"os"

	"github.com/hetznercloud/cli/internal/cli"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("hcloud: ")
	log.SetOutput(os.Stderr)
}

func main() {
	c := cli.NewCLI()

	if c.State.ConfigPath != "" {
		_, err := os.Stat(c.State.ConfigPath)
		switch {
		case err == nil:
			if err := c.State.ReadConfig(); err != nil {
				log.Fatalf("unable to read config file %q: %s\n", c.State.ConfigPath, err)
			}
		case os.IsNotExist(err):
			break
		default:
			log.Fatalf("unable to read config file %q: %s\n", c.State.ConfigPath, err)
		}
	}

	c.State.ReadEnv()

	if err := c.RootCommand.Execute(); err != nil {
		log.Fatalln(err)
	}
}
