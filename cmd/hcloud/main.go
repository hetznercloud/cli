package main

import (
	"log"
	"os"

	"github.com/hetznercloud/cli/cli"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("hcloud: ")
	log.SetOutput(os.Stderr)
}

func main() {
	c := cli.NewCLI()

	if c.ConfigPath != "" {
		_, err := os.Stat(c.ConfigPath)
		switch {
		case err == nil:
			if err := c.ReadConfig(); err != nil {
				log.Fatalf("unable to read config file %q: %s\n", c.ConfigPath, err)
			}
		case os.IsNotExist(err):
			break
		default:
			log.Fatalf("unable to read config file %q: %s\n", c.ConfigPath, err)
		}
	}

	c.ReadContextFile()
	c.ReadEnv()

	if err := c.RootCommand.Execute(); err != nil {
		log.Fatalln(err)
	}
}
