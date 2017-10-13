package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/hetznercloud/cli"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("hcloud: ")
	log.SetOutput(os.Stderr)
}

func main() {
	cli := cli.NewCLI()

	if p := configPath(); p != "" {
		_, err := os.Stat(p)
		switch {
		case err == nil:
			if err := cli.ReadConfig(p); err != nil {
				log.Fatalf("unable to read config file %q: %s\n", p, err)
			}
		case os.IsNotExist(err):
			break
		default:
			log.Fatalf("unable to read config file %q: %s\n", p, err)
		}
	}

	if err := cli.RootCommand.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func configPath() string {
	home := os.Getenv("HOME")
	if home == "" {
		return ""
	}
	return filepath.Join(home, ".config", "hcloud", "config.toml")
}
