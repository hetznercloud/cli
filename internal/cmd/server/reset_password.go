package server

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

var ResetPasswordCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "reset-password [options] <server>",
			Short:                 "Reset the root password of a server",
			Args:                  cobra.ExactArgs(1),
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		output.AddFlag(cmd, output.OptionJSON(), output.OptionYAML())
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		outputFlags := output.FlagsForCommand(cmd)

		idOrName := args[0]
		server, _, err := s.Client().Server().Get(s, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", idOrName)
		}

		result, _, err := s.Client().Server().ResetPassword(s, server)
		if err != nil {
			return err
		}

		if err := s.ActionProgress(cmd, s, result.Action); err != nil {
			return err
		}

		if outputFlags.IsSet("json") || outputFlags.IsSet("yaml") {
			schema := make(map[string]interface{})
			schema["root_password"] = result.RootPassword
			if outputFlags.IsSet("json") {
				return util.DescribeJSON(schema)
			} else {
				return util.DescribeYAML(schema)
			}
		}

		cmd.Printf("Password of server %d reset to: %s\n", server.ID, result.RootPassword)
		return nil
	},
}
