package context

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"syscall"

	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func newCreateCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create [FLAGS] NAME",
		Short:                 "Create a new context",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.Wrap(runContextCreate),
	}
	return cmd
}

func runContextCreate(cli *state.State, cmd *cobra.Command, args []string) error {
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

	context := &state.ConfigContext{Name: name}

	for {
		fmt.Printf("Token: ")
		btoken, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Print("\n")
		if err != nil {
			return err
		}
		token := string(bytes.TrimSpace(btoken))
		if token == "" {
			continue
		}
		if len(token) != 64 {
			fmt.Print("Entered token is invalid (must be exactly 64 characters long)\n")
			continue
		}
		context.Token = token
		break
	}

	cli.Config.Contexts = append(cli.Config.Contexts, context)
	cli.Config.ActiveContext = context

	if err := cli.WriteConfig(); err != nil {
		return err
	}

	fmt.Printf("Context %s created and activated\n", name)

	return nil
}
