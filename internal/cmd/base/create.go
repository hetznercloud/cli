package base

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

// CreateCmd allows defining commands for resource creation
type CreateCmd struct {
	BaseCobraCommand func(hcapi2.Client) *cobra.Command
	// Run is the function that will be called when the command is executed.
	// It should return the created resource, the schema of the resource and an error.
	Run           func(context.Context, hcapi2.Client, state.ActionWaiter, *cobra.Command, []string) (any, any, error)
	PrintResource func(context.Context, hcapi2.Client, *cobra.Command, any)
}

// CobraCommand creates a command that can be registered with cobra.
func (cc *CreateCmd) CobraCommand(
	ctx context.Context, client hcapi2.Client, tokenEnsurer state.TokenEnsurer, actionWaiter state.ActionWaiter,
) *cobra.Command {
	cmd := cc.BaseCobraCommand(client)

	output.AddFlag(cmd, output.OptionJSON(), output.OptionYAML())

	if cmd.Args == nil {
		cmd.Args = cobra.NoArgs
	}

	cmd.TraverseChildren = true
	cmd.DisableFlagsInUseLine = true

	if cmd.PreRunE != nil {
		cmd.PreRunE = util.ChainRunE(cmd.PreRunE, tokenEnsurer.EnsureToken)
	} else {
		cmd.PreRunE = tokenEnsurer.EnsureToken
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		outputFlags := output.FlagsForCommand(cmd)

		isSchema := outputFlags.IsSet("json") || outputFlags.IsSet("yaml")
		if isSchema {
			cmd.SetOut(os.Stderr)
		} else {
			cmd.SetOut(os.Stdout)
		}

		resource, schema, err := cc.Run(ctx, client, actionWaiter, cmd, args)
		if err != nil {
			return err
		}

		if isSchema {
			if outputFlags.IsSet("json") {
				return util.DescribeJSON(schema)
			} else {
				return util.DescribeYAML(schema)
			}
		} else if cc.PrintResource != nil && resource != nil {
			cc.PrintResource(ctx, client, cmd, resource)
		}
		return nil
	}

	return cmd
}
