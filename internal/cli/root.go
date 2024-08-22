package cli

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

func NewRootCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "hcloud",
		Short:                 "Hetzner Cloud CLI",
		Long:                  "A command-line interface for Hetzner Cloud",
		TraverseChildren:      true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		DisableFlagsInUseLine: true,
	}

	cmd.PersistentFlags().AddFlagSet(s.Config().FlagSet())

	cmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		var err error
		out := os.Stdout
		quiet, err := config.OptionQuiet.Get(s.Config())
		if err != nil {
			return err
		}
		if quiet {
			out, err = os.Open(os.DevNull)
			if err != nil {
				return err
			}
		}
		cmd.SetOut(out)
		return nil
	}
	return cmd
}
