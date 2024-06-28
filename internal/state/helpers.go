package state

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/state/config"
)

func Wrap(s State, f func(State, *cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return f(s, cmd, args)
	}
}

func (c *state) EnsureToken(_ *cobra.Command, _ []string) error {
	token, err := config.OptionToken.Get(c.config)
	if err != nil {
		return err
	}
	if token == "" {
		return errors.New("no active context or token (see `hcloud context --help`)")
	}
	return nil
}
