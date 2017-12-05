package main

import (
	"log"
	"os"

	"github.com/hetznercloud/cli"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("hcloud: ")
	log.SetOutput(os.Stderr)
}

func main() {
	c := cli.NewCLI()

	if cli.DefaultConfigPath != "" {
		_, err := os.Stat(cli.DefaultConfigPath)
		switch {
		case err == nil:
			if err := c.ReadConfig(cli.DefaultConfigPath); err != nil {
				log.Fatalf("unable to read config file %q: %s\n", cli.DefaultConfigPath, err)
			}
		case os.IsNotExist(err):
			break
		default:
			log.Fatalf("unable to read config file %q: %s\n", cli.DefaultConfigPath, err)
		}
	}

	c.ReadEnv()

	if err := c.RootCommand.Execute(); err != nil {
		log.Fatalln(err)
	}
}
