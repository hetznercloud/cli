package cmds

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func newServerRequestConsoleCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "request-console [FLAGS] SERVER",
		Short:                 "Request a WebSocket VNC console for a server",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.ServerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runServerRequestConsole),
	}
	addOutputFlag(cmd, outputOptionJSON())
	return cmd
}

func runServerRequestConsole(cli *state.State, cmd *cobra.Command, args []string) error {
	outOpts := outputFlagsForCommand(cmd)
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	result, _, err := cli.Client().Server.RequestConsole(cli.Context, server)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, result.Action); err != nil {
		return err
	}

	if outOpts.IsSet("json") {
		return util.DescribeJSON(struct {
			WSSURL   string
			Password string
		}{
			WSSURL:   result.WSSURL,
			Password: result.Password,
		})
	}

	fmt.Printf("Console for server %d:\n", server.ID)
	fmt.Printf("WebSocket URL: %s\n", result.WSSURL)
	fmt.Printf("VNC Password: %s\n", result.Password)
	return nil
}
