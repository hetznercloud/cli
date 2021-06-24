package firewall

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func newApplyToResourceCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "apply-to-resource FIREWALL FLAGS",
		Short:                 "Applies a Firewall to a single resource",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.FirewallNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateApplyToResource, cli.EnsureToken),
		RunE:                  cli.Wrap(runApplyToResource),
	}
	cmd.Flags().String("type", "", "Resource Type (server) (required)")
	cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("server", "label_selector"))
	cmd.MarkFlagRequired("type")

	cmd.Flags().String("server", "", "Server name of ID (required when type is server)")
	cmd.RegisterFlagCompletionFunc("server", cmpl.SuggestCandidatesF(cli.ServerNames))

	cmd.Flags().StringP("label-selector", "l", "", "Label Selector")
	return cmd
}
func validateApplyToResource(cmd *cobra.Command, args []string) error {
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

	return nil
}
func runApplyToResource(cli *state.State, cmd *cobra.Command, args []string) error {
	resourceType, _ := cmd.Flags().GetString("type")
	serverIdOrName, _ := cmd.Flags().GetString("server")
	labelSelector, _ := cmd.Flags().GetString("label-selector")
	idOrName := args[0]
	firewall, _, err := cli.Client().Firewall.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if firewall == nil {
		return fmt.Errorf("Firewall not found: %v", idOrName)
	}
	opts := hcloud.FirewallResource{Type: hcloud.FirewallResourceType(resourceType)}

	switch opts.Type {
	case hcloud.FirewallResourceTypeServer:
		server, _, err := cli.Client().Server.Get(cli.Context, serverIdOrName)
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

	actions, _, err := cli.Client().Firewall.ApplyResources(cli.Context, firewall, []hcloud.FirewallResource{opts})
	if err != nil {
		return err
	}
	if err := cli.ActionsProgresses(cli.Context, actions); err != nil {
		return err
	}
	fmt.Printf("Firewall %d applied\n", firewall.ID)

	return nil
}
