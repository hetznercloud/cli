package loadbalancer

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var CreateCmd = base.CreateCmd[*hcloud.LoadBalancer]{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "create [options] --name <name> --type <type>",
			Short:                 "Create a Load Balancer",
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().String("name", "", "Load Balancer name (required)")
		util.MarkFlagRequired(cmd, "name")

		cmd.Flags().String("type", "", "Load Balancer Type (ID or name) (required)")
		cmpl.RegisterFlagCompletion(cmd, "type", cmpl.SuggestCandidatesF(client.LoadBalancerType().Names))
		util.MarkFlagRequired(cmd, "type")

		cmd.Flags().String("algorithm-type", "", "Algorithm Type name (round_robin or least_connections)")
		cmpl.RegisterFlagCompletion(cmd, "algorithm-type", cmpl.SuggestCandidates(
			string(hcloud.LoadBalancerAlgorithmTypeLeastConnections),
			string(hcloud.LoadBalancerAlgorithmTypeRoundRobin),
		))
		cmd.Flags().String("location", "", "Location (ID or name)")
		cmpl.RegisterFlagCompletion(cmd, "location", cmpl.SuggestCandidatesF(client.Location().Names))

		cmd.Flags().String("network-zone", "", "Network Zone")
		cmpl.RegisterFlagCompletion(cmd, "network-zone", cmpl.SuggestCandidatesF(client.Location().NetworkZones))

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		cmd.Flags().StringSlice("enable-protection", []string{}, "Enable protection (delete) (default: none)")
		cmpl.RegisterFlagCompletion(cmd, "enable-protection", cmpl.SuggestCandidates("delete"))

		cmd.Flags().String("network", "", "Name or ID of the Network the Load Balancer should be attached to on creation")
		cmpl.RegisterFlagCompletion(cmd, "network", cmpl.SuggestCandidatesF(client.Network().Names))

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, _ []string) (*hcloud.LoadBalancer, any, error) {
		name, _ := cmd.Flags().GetString("name")
		loadBalancerTypeName, _ := cmd.Flags().GetString("type")
		algorithmType, _ := cmd.Flags().GetString("algorithm-type")
		location, _ := cmd.Flags().GetString("location")
		networkZone, _ := cmd.Flags().GetString("network-zone")
		labels, _ := cmd.Flags().GetStringToString("label")
		protection, _ := cmd.Flags().GetStringSlice("enable-protection")
		network, _ := cmd.Flags().GetString("network")

		protectionOpts, err := ChangeProtectionCmds.GetChangeProtectionOpts(true, protection)
		if err != nil {
			return nil, nil, err
		}

		loadBalancerType, _, err := s.Client().LoadBalancerType().Get(cmd.Context(), loadBalancerTypeName)
		if err != nil {
			return nil, nil, err
		}
		if loadBalancerType == nil {
			return nil, nil, fmt.Errorf("Load Balancer Type not found: %s", loadBalancerTypeName)
		}

		cmd.Print(deprecatedLoadBalancerTypeWarning(loadBalancerType))

		createOpts := hcloud.LoadBalancerCreateOpts{
			Name:             name,
			LoadBalancerType: loadBalancerType,
			Labels:           labels,
		}
		if algorithmType != "" {
			createOpts.Algorithm = &hcloud.LoadBalancerAlgorithm{Type: hcloud.LoadBalancerAlgorithmType(algorithmType)}
		}
		if networkZone != "" {
			createOpts.NetworkZone = hcloud.NetworkZone(networkZone)
		}
		if location != "" {
			createOpts.Location = &hcloud.Location{Name: location}
		}
		if network != "" {
			net, _, err := s.Client().Network().Get(cmd.Context(), network)
			if err != nil {
				return nil, nil, err
			}
			if net == nil {
				return nil, nil, fmt.Errorf("Network not found: %s", network)
			}
			createOpts.Network = net
		}
		result, _, err := s.Client().LoadBalancer().Create(cmd.Context(), createOpts)
		if err != nil {
			return nil, nil, err
		}

		if err := s.WaitForActions(cmd.Context(), cmd, result.Action); err != nil {
			return nil, nil, err
		}
		cmd.Printf("Load Balancer %d created\n", result.LoadBalancer.ID)

		if protectionOpts.Delete != nil {
			if err := ChangeProtectionCmds.ChangeProtection(s, cmd, result.LoadBalancer, true, protectionOpts); err != nil {
				return nil, nil, err
			}
		}

		loadBalancer, _, err := s.Client().LoadBalancer().GetByID(cmd.Context(), result.LoadBalancer.ID)
		if err != nil {
			return nil, nil, err
		}
		if loadBalancer == nil {
			return nil, nil, fmt.Errorf("Load Balancer not found: %d", result.LoadBalancer.ID)
		}

		return loadBalancer, util.Wrap("load_balancer", hcloud.SchemaFromLoadBalancer(loadBalancer)), nil
	},

	PrintResource: func(_ state.State, cmd *cobra.Command, loadBalancer *hcloud.LoadBalancer) error {
		_, err := fmt.Fprintf(cmd.OutOrStdout(), "IPv4: %s\nIPv6: %s\n", loadBalancer.PublicNet.IPv4.IP.String(), loadBalancer.PublicNet.IPv6.IP.String())
		return err
	},
}
