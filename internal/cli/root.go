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

	cmd.PersistentFlags().AddFlagSet(config.FlagSet)

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		var err error
		out := os.Stdout
		if quiet := config.OptionQuiet.Value(); quiet {
			//if quiet := viper.GetBool("quiet"); quiet {
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
