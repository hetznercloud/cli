package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func newContextCreateCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create [FLAGS] NAME",
		Short:                 "Create a new context",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runContextCreate),
	}
	return cmd
}

func runContextCreate(cli *CLI, cmd *cobra.Command, args []string) error {
	if !cli.Terminal() {
		return errors.New("context create is an interactive command")
	}

	name := strings.TrimSpace(args[0])
	if name == "" {
		return errors.New("invalid name")
	}
	if cli.Config.ContextByName(name) != nil {
		return errors.New("name already used")
	}

	r := bufio.NewReader(os.Stdin)
	context := &ConfigContext{Name: name}

	for {
		fmt.Printf("Token: ")
		token, err := r.ReadString('\n')
		if err != nil {
			return err
		}
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}
		context.Token = token
		break
	}

	cli.Config.Contexts = append(cli.Config.Contexts, context)
	if len(cli.Config.Contexts) == 1 {
		cli.Config.ActiveContext = context
	}
	if err := cli.WriteConfig(); err != nil {
		return err
	}

	return nil
}
