package context

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/state"
)

func newCreateCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create [FLAGS] [NAME]",
		Short:                 "Create a new context",
		Args:                  cobra.MaximumNArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.Wrap(runCreate),
	}
	return cmd
}

func runCreate(cli *state.State, cmd *cobra.Command, args []string) error {
	if !state.StdoutIsTerminal() {
		return errors.New("context create is an interactive command")
	}

	var name string
	if len(args) > 0 {
		name = strings.TrimSpace(args[0])
		if err := validateContextName(cli, name); err != nil {
			return err
		}
	}

	if name == "" {
		err := huh.NewForm(huh.NewGroup(
			huh.NewInput().
				Title("Context name").
				Validate(func(s string) error {
					return validateContextName(cli, s)
				}).
				Value(&name),
		)).Run()
		if err != nil {
			return err
		}
	}

	context := &state.ConfigContext{Name: name}

	if envToken := os.Getenv("HCLOUD_TOKEN"); envToken != "" {
		var (
			group       *huh.Group
			useEnvToken bool
		)
		if len(envToken) != 64 {
			group = huh.NewGroup(
				huh.NewNote().
					Title("Warning").
					Description("The HCLOUD_TOKEN environment variable is set, but token is invalid (must be exactly 64 characters long)"),
			)
		} else {
			useEnvToken = true
			group = huh.NewGroup(
				huh.NewConfirm().
					Title("Use HCLOUD_TOKEN").
					Description("The HCLOUD_TOKEN environment variable is set. Do you want to use the token from HCLOUD_TOKEN for the new context?").
					Value(&useEnvToken),
			)
		}
		err := huh.NewForm(group).Run()
		if err != nil {
			return err
		}
		if useEnvToken {
			context.Token = envToken
		}
	}

	if context.Token == "" {
		err := huh.NewForm(huh.NewGroup(
			huh.NewInput().
				Title("API Token").
				Description("Get your Hetzner Cloud API Token here:\nhttps://console.hetzner.cloud/").
				Password(true).
				Validate(func(s string) error {
					if len(s) != 64 {
						return fmt.Errorf("invalid token")
					}
					return nil
				}).
				Value(&context.Token),
		)).Run()
		if err != nil {
			return err
		}
	}

	cli.Config.Contexts = append(cli.Config.Contexts, context)

	activateContext := true
	if len(cli.Config.Contexts) > 1 {
		err := huh.NewForm(huh.NewGroup(
			huh.NewConfirm().
				Title("Activate context").
				Description("Do you want to activate the new context?\nIf not, you can activate it later using `hcloud context use`.").
				Value(&activateContext),
		)).Run()
		if err != nil {
			return err
		}
	}

	msg := "Context \"%s\" created"
	if activateContext {
		cli.Config.ActiveContext = context
		msg += " and activated"
	}

	if err := cli.WriteConfig(); err != nil {
		return err
	}

	return huh.Run(huh.NewNote().
		Title("Success").
		Description(fmt.Sprintf(msg, name)),
	)
}

func validateContextName(cli *state.State, name string) error {
	if name == "" {
		return errors.New("name must not be empty")
	}
	if cli.Config.ContextByName(name) != nil {
		return errors.New("name already used")
	}
	return nil
}
