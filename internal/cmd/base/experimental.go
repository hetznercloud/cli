package base

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

func Experimental(s state.State, cmd *cobra.Command, slug string) *cobra.Command {

	if cmd.Long == "" {
		cmd.Long = cmd.Short
	}
	cmd.Long += fmt.Sprintf("\n\nExperimental: %s is experimental, breaking changes may occur within minor releases. See %s for more details.", product, url)
	cmd.Short = "[experimental] " + cmd.Short

	cmd.PreRunE = util.ChainRunE(cmd.PreRunE, func(cmd *cobra.Command, _ []string) error {
		hideWarning, err := config.OptionExperimental.Get(s.Config())
		if err != nil {
			return err
		}
		if !hideWarning {
			cmd.PrintErrln("Warning: This command is experimental and may change in the future. Use --experimental to suppress this warning.")
		}
		return nil
	})
	return cmd
}
