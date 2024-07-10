package firewall

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ApplyToResourceCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "apply-to-resource (--type server --server <server> | --type label_selector --label-selector <label-selector>) <firewall>",
			Short:                 "Applies a Firewall to a single resource",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Firewall().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String("type", "", "Resource Type (server, label_selector) (required)")
		cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("server", "label_selector"))
		cmd.MarkFlagRequired("type")

		cmd.Flags().String("server", "", "Server name of ID (required when type is server)")
		cmd.RegisterFlagCompletionFunc("server", cmpl.SuggestCandidatesF(client.Server().Names))

		cmd.Flags().StringP("label-selector", "l", "", "Label Selector")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
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
		firewall, _, err := s.Client().Firewall().Get(s, idOrName)
		if err != nil {
			return err
		}
		if firewall == nil {
			return fmt.Errorf("Firewall not found: %v", idOrName)
		}
		opts := hcloud.FirewallResource{Type: hcloud.FirewallResourceType(resourceType)}

		switch opts.Type {
		case hcloud.FirewallResourceTypeServer:
			server, _, err := s.Client().Server().Get(s, serverIdOrName)
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

		actions, _, err := s.Client().Firewall().ApplyResources(s, firewall, []hcloud.FirewallResource{opts})
		if err != nil {
			return err
		}
		if err := s.WaitForActions(cmd, s, actions...); err != nil {
			return err
		}
		cmd.Printf("Firewall %d applied to resource\n", firewall.ID)

		return nil
	},
}
