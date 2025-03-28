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

var RequestConsoleCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "request-console [options] <server>",
			Short:                 "Request a WebSocket VNC console for a Server",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		output.AddFlag(cmd, output.OptionJSON(), output.OptionYAML())
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		outOpts := output.FlagsForCommand(cmd)
		idOrName := args[0]
		server, _, err := s.Client().Server().Get(s, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("Server not found: %s", idOrName)
		}

		result, _, err := s.Client().Server().RequestConsole(s, server)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, result.Action); err != nil {
			return err
		}

		if outOpts.IsSet("json") || outOpts.IsSet("yaml") {
			schema := struct {
				WSSURL   string `json:"wss_url"`
				Password string `json:"password"`
			}{
				WSSURL:   result.WSSURL,
				Password: result.Password,
			}

			if outOpts.IsSet("json") {
				return util.DescribeJSON(cmd.OutOrStdout(), schema)
			}
			return util.DescribeYAML(cmd.OutOrStdout(), schema)
		}

		cmd.Printf("Console for Server %d:\n", server.ID)
		cmd.Printf("WebSocket URL: %s\n", result.WSSURL)
		cmd.Printf("VNC Password: %s\n", result.Password)
		return nil
	},
}
