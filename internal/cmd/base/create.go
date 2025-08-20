package base

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

// CreateCmd allows defining commands for resource creation
type CreateCmd struct {
	BaseCobraCommand func(hcapi2.Client) *cobra.Command
	// Run is the function that will be called when the command is executed.
	// It should return the created resource, the schema of the resource and an error.
	Run           func(state.State, *cobra.Command, []string) (any, any, error)
	PrintResource func(state.State, *cobra.Command, any)
	// ExperimentalF is a function that will be used to mark the command as experimental.
	ExperimentalF func(state.State, *cobra.Command) *cobra.Command
}

// CobraCommand creates a command that can be registered with cobra.
func (cc *CreateCmd) CobraCommand(s state.State) *cobra.Command {
	cmd := cc.BaseCobraCommand(s.Client())

	output.AddFlag(cmd, output.OptionJSON(), output.OptionYAML())

	if cmd.Args == nil {
		cmd.Args = util.Validate
	}

	cmd.TraverseChildren = true
	cmd.DisableFlagsInUseLine = true

	if cmd.PreRunE != nil {
		cmd.PreRunE = util.ChainRunE(cmd.PreRunE, s.EnsureToken)
	} else {
		cmd.PreRunE = s.EnsureToken
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		outputFlags := output.FlagsForCommand(cmd)

		quiet, err := config.OptionQuiet.Get(s.Config())
		if err != nil {
			return err
		}

		schemaOut := cmd.OutOrStdout()
		isSchema := outputFlags.IsSet("json") || outputFlags.IsSet("yaml")
		if isSchema {
			if quiet {
				// If we are in quiet mode, we saved the original output in cmd.errWriter. We can now restore it.
				schemaOut = cmd.ErrOrStderr()
			} else {
				// We don't want anything other than the schema in stdout, so we set the default to stderr
				cmd.SetOut(os.Stderr)
			}
		}

		resource, schema, err := cc.Run(s, cmd, args)
		if err != nil {
			return err
		}

		if isSchema {
			if outputFlags.IsSet("json") {
				return util.DescribeJSON(schemaOut, schema)
			}
			return util.DescribeYAML(schemaOut, schema)
		} else if cc.PrintResource != nil && resource != nil {
			cc.PrintResource(s, cmd, resource)
		}
		return nil
	}

	if cc.ExperimentalF != nil {
		cmd = cc.ExperimentalF(s, cmd)
	}
	return cmd
}
