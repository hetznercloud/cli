package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/state"
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
	cmd.PersistentFlags().String("context", "", "Context to use")
	cmd.PersistentFlags().String("config", "", "Config file path")

	cmd.PersistentFlags().String("token", "", "Hetzner Cloud API token")
	cmd.PersistentFlags().String("endpoint", "", "Hetzner Cloud API endpoint")
	cmd.PersistentFlags().Bool("debug", false, "Enable debug output")
	cmd.PersistentFlags().String("debug-file", "", "Write debug output to file")
	cmd.PersistentFlags().Bool("quiet", false, "Only show command output")
	cmd.PersistentFlags().Duration("poll-interval", 0, "Interval for polling resources")

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// flags have now been parsed into viper; active context might have changed
		/*if err := s.Config().LoadActiveContext(); err != nil {
			return err
		}
		s.Client().FromConfig(s.Config())*/

		var err error

		fmt.Println("ssh keys:", s.Config().GetStringSlice("ssh-keys"))

		out := os.Stdout
		if quiet := s.Config().GetBool("quiet"); quiet {
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
