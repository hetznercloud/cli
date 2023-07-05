package firewall

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
)

var RemoveFromResourceCommand = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "remove-from-resource FIREWALL FLAGS",
			Short:                 "Removes a Firewall from a single resource",
			Args:                  cobra.ExactArgs(1),
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Firewall().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String("type", "", "Resource Type (server) (required)")
		cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("server", "label_selector"))
		cmd.MarkFlagRequired("type")

		cmd.Flags().String("server", "", "Server name of ID (required when type is server)")
		cmd.RegisterFlagCompletionFunc("server", cmpl.SuggestCandidatesF(client.Server().Names))

		cmd.Flags().StringP("label-selector", "l", "", "Label Selector")
		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		resourceType, _ := cmd.Flags().GetString("type")

		switch resourceType {
		case string(hcloud.FirewallResourceTypeServer):
			server, _ := cmd.Flags().GetString("server")
			if server == "" {
				return fmt.Errorf("type %s need a --server specific", resourceType)
			}
		case string(hcloud.FirewallResourceTypeLabelSelector):
			labelSelector, _ := cmd.Flags().GetString("label-selector")
			if labelSelector == "" {
				return fmt.Errorf("type %s need a --label-selector specific", resourceType)
			}
		default:
			return fmt.Errorf("unknown type %s", resourceType)
		}
		serverIdOrName, _ := cmd.Flags().GetString("server")
		labelSelector, _ := cmd.Flags().GetString("label-selector")

		idOrName := args[0]
		firewall, _, err := client.Firewall().Get(ctx, idOrName)
		if err != nil {
			return err
		}
		if firewall == nil {
			return fmt.Errorf("Firewall not found: %v", idOrName)
		}
		opts := hcloud.FirewallResource{Type: hcloud.FirewallResourceType(resourceType)}

		switch opts.Type {
		case hcloud.FirewallResourceTypeServer:
			server, _, err := client.Server().Get(ctx, serverIdOrName)
			if err != nil {
				return err
			}
			if server == nil {
				return fmt.Errorf("Server not found: %v", serverIdOrName)
			}
			opts.Server = &hcloud.FirewallResourceServer{ID: server.ID}
		case hcloud.FirewallResourceTypeLabelSelector:
			opts.LabelSelector = &hcloud.FirewallResourceLabelSelector{Selector: labelSelector}
		default:
			return fmt.Errorf("unknown type %s", opts.Type)
		}
		actions, _, err := client.Firewall().RemoveResources(ctx, firewall, []hcloud.FirewallResource{opts})
		if err != nil {
			return err
		}
		if err := waiter.WaitForActions(ctx, actions); err != nil {
			return err
		}
		fmt.Printf("Firewall %d applied\n", firewall.ID)

		return nil
	},
}
