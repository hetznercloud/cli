package base

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// CreateCmd allows defining commands for resource creation
type CreateCmd struct {
	BaseCobraCommand func(hcapi2.Client) *cobra.Command
	Run              func(context.Context, hcapi2.Client, state.ActionWaiter, *cobra.Command, []string) (*hcloud.Response, any, error)
	PrintResource    func(context.Context, hcapi2.Client, *cobra.Command, any)
}

// CobraCommand creates a command that can be registered with cobra.
func (cc *CreateCmd) CobraCommand(
	ctx context.Context, client hcapi2.Client, tokenEnsurer state.TokenEnsurer, actionWaiter state.ActionWaiter,
) *cobra.Command {
	cmd := cc.BaseCobraCommand(client)

	output.AddFlag(cmd, output.OptionJSON())

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

		isJson := outputFlags.IsSet("json")
		if isJson {
			cmd.SetOut(os.Stderr)
		} else {
			cmd.SetOut(os.Stdout)
		}

		response, resource, err := cc.Run(ctx, client, actionWaiter, cmd, args)
		if err != nil {
			return err
		}

		if isJson {
			bytes, _ := io.ReadAll(response.Body)

			var data map[string]any
			if err := json.Unmarshal(bytes, &data); err != nil {
				return err
			}

			delete(data, "action")
			delete(data, "actions")
			delete(data, "next_actions")

			return util.DescribeJSON(data)
		} else if resource != nil {
			cc.PrintResource(ctx, client, cmd, resource)
		}
		return nil
	}

	return cmd
}
