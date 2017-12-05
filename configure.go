package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func newConfigureCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "configure",
		Short:            "Configure the CLI",
		Args:             cobra.NoArgs,
		TraverseChildren: true,
		RunE:             cli.wrap(runConfigure),
	}
	return cmd
}

func runConfigure(cli *CLI, cmd *cobra.Command, args []string) error {
	if !cli.Terminal() {
		return errors.New("configure is an interactive command")
	}
	if DefaultConfigPath == "" {
		return errors.New("could not determine config file path")
	}
	if cli.Config == nil {
		cli.Config = &Config{}
	}

	r := bufio.NewReader(os.Stdin)

	fmt.Printf("Token: ")
	token, err := r.ReadString('\n')
	if err != nil {
		return err
	}
	token = strings.TrimSpace(token)
	if token != "" {
		cli.Config.Token = token
		if err := cli.WriteConfig(DefaultConfigPath); err != nil {
			return err
		}
	} else {
		fmt.Println("no token set")
	}

	return nil
}
