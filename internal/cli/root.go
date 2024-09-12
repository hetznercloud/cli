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

	for _, opt := range config.Options {
		f := opt.GetFlagCompletionFunc()
		if !opt.HasFlags(config.OptionFlagPFlag) || f == nil {
			continue
		}
		// opt.FlagName() is prefixed with --
		flagName := opt.FlagName()[2:]
		_ = cmd.RegisterFlagCompletionFunc(flagName, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return f(s.Client(), s.Config(), cmd, args, toComplete)
		})
	}

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
