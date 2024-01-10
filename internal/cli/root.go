package cli

import (
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
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
	cmd.PersistentFlags().Duration("poll-interval", 500*time.Millisecond, "Interval at which to poll information, for example action progress")
	cmd.PersistentFlags().Bool("quiet", false, "Only print error messages")

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		pollInterval, err := cmd.Flags().GetDuration("poll-interval")
		if err != nil {
			return err
		}
		s.Client().WithOpts(hcloud.WithPollBackoffFunc(hcloud.ConstantBackoff(pollInterval)))

		out := os.Stdout
		if quiet, _ := cmd.Flags().GetBool("quiet"); quiet {
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
