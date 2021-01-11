package server

import (
	"fmt"
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newEnableProtectionCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable-protection [FLAGS] SERVER PROTECTIONLEVEL [PROTECTIONLEVEL...]",
		Short: "Enable resource protection for a server",
		Args:  cobra.MinimumNArgs(2),
		ValidArgsFunction: cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(cli.ServerNames),
			cmpl.SuggestCandidates("delete", "rebuild"),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runServerEnableProtection),
	}
	return cmd
}

func runServerEnableProtection(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	var unknown []string
	opts := hcloud.ServerChangeProtectionOpts{}
	for _, arg := range args[1:] {
		switch strings.ToLower(arg) {
		case "delete":
			opts.Delete = hcloud.Bool(true)
		case "rebuild":
			opts.Rebuild = hcloud.Bool(true)
		default:
			unknown = append(unknown, arg)
		}
	}
	if len(unknown) > 0 {
		return fmt.Errorf("unknown protection level: %s", strings.Join(unknown, ", "))
	}

	action, _, err := cli.Client().Server.ChangeProtection(cli.Context, server, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Resource protection enabled for server %d\n", server.ID)
	return nil
}
