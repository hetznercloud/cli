package base

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

// ExperimentalWrapper create a command wrapper that appends a notice to the command
// descriptions and logs a warning when it is used.
//
// Usage:
//
//	var (
//		ExperimentalProduct = ExperimentalWrapper("Product name", "beta", "https://docs.hetzner.cloud/changelog#new-product")
//	)
//
//	func (c) CobraCommand(s state.State) *cobra.Command {
//		cmd := &cobra.Command{
//			Use:     "command",
//			Short:   "My experimental command",
//			Long:    "This is an experimental command.",
//			PreRunE: s.EnsureToken,
//		}
//
//		cmd.Run = func(cmd *cobra.Command, _ []string) {}
//
//		return ExperimentalProduct(s, cmd)
//	}
func ExperimentalWrapper(product, maturity, url string) func(state.State, *cobra.Command) *cobra.Command {
	return func(s state.State, cmd *cobra.Command) *cobra.Command {
		cmd.Long = strings.TrimLeft(cmd.Long, "\n")

		if cmd.Long == "" {
			cmd.Long = cmd.Short
		}

		cmd.Short = "[experimental] " + cmd.Short
		cmd.Long += fmt.Sprintf(`

Experimental: %s is %s, breaking changes may occur within minor releases.
See %s for more details.
`, product, maturity, url)

		cmd.PreRunE = util.ChainRunE(cmd.PreRunE, func(cmd *cobra.Command, _ []string) error {
			hideWarning, err := config.OptionNoExperimentalWarning.Get(s.Config())
			if err != nil {
				return err
			}
			if !hideWarning {
				cmd.PrintErrf("Warning: %s is %s and may change in the future. Use --no-experimental-warnings to suppress this warning.\n", product, maturity)
			}
			return nil
		})

		return cmd
	}
}
