package server

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

var IPCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "ip [--ipv6] <server>",
			Short:                 "Print a server's IP address",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().BoolP("ipv6", "6", false, "Print the first address of the Server's primary IPv6 network")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		ipv6, _ := cmd.Flags().GetBool("ipv6")
		idOrName := args[0]
		server, _, err := s.Client().Server().Get(s, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("Server not found: %s", idOrName)
		}
		if ipv6 {
			if server.PublicNet.IPv6.IsUnspecified() {
				return fmt.Errorf("Server %s has no Primary IPv6", idOrName)
			}
			cmd.Println(server.PublicNet.IPv6.IP.String() + "1")
		} else {
			if server.PublicNet.IPv4.IsUnspecified() {
				return fmt.Errorf("Server %s has no Primary IPv4", idOrName)
			}
			cmd.Println(server.PublicNet.IPv4.IP.String())
		}
		return nil
	},
}
