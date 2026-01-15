package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra/doc"

	"github.com/hetznercloud/cli/internal/cli"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

func main() {
	//nolint:gosec
	if err := os.MkdirAll("./manpages", 0755); err != nil {
		log.Fatal(err)
	}

	s, _ := state.New(config.New())
	rootCommand := cli.NewRootCommand(s)

	err := doc.GenManTree(rootCommand, nil, "./manpages")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Man pages generated successfully")
}
