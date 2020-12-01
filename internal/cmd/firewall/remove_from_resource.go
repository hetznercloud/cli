package firewall

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func newRemoveFromResourceCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "remove-from-resource FIREWALL FLAGS",
		Short:                 "Removes a Firewall from a single resource",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.FirewallNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateRemoveFromResource, cli.EnsureToken),
		RunE:                  cli.Wrap(runRemoveFromResource),
	}
	cmd.Flags().String("type", "", "Resource Type (server) (required)")
	cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("server"))
	cmd.MarkFlagRequired("type")

	cmd.Flags().String("server", "", "Server name of ID (required when type is server)")
	cmd.RegisterFlagCompletionFunc("server", cmpl.SuggestCandidatesF(cli.ServerNames))

	return cmd
}
func validateRemoveFromResource(cmd *cobra.Command, args []string) error {
	resourceType, _ := cmd.Flags().GetString("type")

	switch resourceType {
	case "server":
		server, _ := cmd.Flags().GetString("server")
		if server == "" {
			return fmt.Errorf("type %s need a --server specific", resourceType)
		}
	default:
		return fmt.Errorf("unknown type %s", resourceType)
	}

	return nil
}
func runRemoveFromResource(cli *state.State, cmd *cobra.Command, args []string) error {
	resourceType, _ := cmd.Flags().GetString("type")
	serverIdOrName, _ := cmd.Flags().GetString("server")

	idOrName := args[0]
	firewall, _, err := cli.Client().Firewall.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if firewall == nil {
		return fmt.Errorf("Firewall not found: %v", idOrName)
	}

	server, _, err := cli.Client().Server.Get(cli.Context, serverIdOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("Server not found: %v", serverIdOrName)
	}

	opts := hcloud.FirewallResource{Type: hcloud.FirewallResourceType(resourceType)}

	switch opts.Type {
	case hcloud.FirewallResourceTypeServer:
		opts.Server = &hcloud.FirewallResourceServer{ID: server.ID}
	default:
		return fmt.Errorf("unkown type %s", resourceType)
	}

	actions, _, err := cli.Client().Firewall.RemoveResources(cli.Context, firewall, []hcloud.FirewallResource{opts})
	if err != nil {
		return err
	}
	if err := cli.ActionsProgresses(cli.Context, actions); err != nil {
		return err
	}
	fmt.Printf("Firewall %d applied\n", firewall.ID)

	return nil
}
