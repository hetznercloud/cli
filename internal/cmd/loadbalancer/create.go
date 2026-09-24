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
		_ = cmd.MarkFlagRequired("name")

		cmd.Flags().String("type", "", "Load Balancer Type (ID or name) (required)")
		_ = cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidatesF(client.LoadBalancerType().Names))
		_ = cmd.MarkFlagRequired("type")

		cmd.Flags().String("algorithm-type", "", "Algorithm Type name (round_robin or least_connections)")
		_ = cmd.RegisterFlagCompletionFunc("algorithm-type", cmpl.SuggestCandidates(
			string(hcloud.LoadBalancerAlgorithmTypeLeastConnections),
			string(hcloud.LoadBalancerAlgorithmTypeRoundRobin),
		))
		cmd.Flags().String("location", "", "Location (ID or name)")
		_ = cmd.RegisterFlagCompletionFunc("location", cmpl.SuggestCandidatesF(client.Location().Names))

		cmd.Flags().String("network-zone", "", "Network Zone")
		_ = cmd.RegisterFlagCompletionFunc("network-zone", cmpl.SuggestCandidatesF(client.Location().NetworkZones))

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		cmd.Flags().StringSlice("enable-protection", []string{}, "Enable protection (delete) (default: none)")
		_ = cmd.RegisterFlagCompletionFunc("enable-protection", cmpl.SuggestCandidates("delete"))

		cmd.Flags().String("network", "", "Name or ID of the Network the Load Balancer should be attached to on creation")
		_ = cmd.RegisterFlagCompletionFunc("network", cmpl.SuggestCandidatesF(client.Network().Names))

		cmd.Flags().String("primary-ipv4", "", "Primary IPv4 (ID or name) to assign to the Load Balancer")
		_ = cmd.RegisterFlagCompletionFunc("primary-ipv4", cmpl.SuggestCandidatesF(client.PrimaryIP().Names(true, false, hcloud.Ptr(hcloud.PrimaryIPTypeIPv4))))

		cmd.Flags().String("primary-ipv6", "", "Primary IPv6 (ID or name) to assign to the Load Balancer")
		_ = cmd.RegisterFlagCompletionFunc("primary-ipv6", cmpl.SuggestCandidatesF(client.PrimaryIP().Names(true, false, hcloud.Ptr(hcloud.PrimaryIPTypeIPv6))))

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
		primaryIPv4IDOrName, _ := cmd.Flags().GetString("primary-ipv4")
		primaryIPv6IDOrName, _ := cmd.Flags().GetString("primary-ipv6")

		protectionOpts, err := ChangeProtectionCmds.GetChangeProtectionOpts(true, protection)
		if err != nil {
			return nil, nil, err
		}

		loadBalancerType, _, err := s.Client().LoadBalancerType().Get(s, loadBalancerTypeName)
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
			net, _, err := s.Client().Network().Get(s, network)
			if err != nil {
				return nil, nil, err
			}
			if net == nil {
				return nil, nil, fmt.Errorf("Network not found: %s", network)
			}
			createOpts.Network = net
		}
		publicNet := &hcloud.LoadBalancerCreateOptsPublicNet{}
		if primaryIPv4IDOrName != "" {
			primaryIPv4, _, err := s.Client().PrimaryIP().Get(s, primaryIPv4IDOrName)
			if err != nil {
				return nil, nil, err
			}
			if primaryIPv4 == nil {
				return nil, nil, fmt.Errorf("Primary IPv4 not found: %s", primaryIPv4IDOrName)
			}
			publicNet.IPv4 = primaryIPv4
		}
		if primaryIPv6IDOrName != "" {
			primaryIPv6, _, err := s.Client().PrimaryIP().Get(s, primaryIPv6IDOrName)
			if err != nil {
				return nil, nil, err
			}
			if primaryIPv6 == nil {
				return nil, nil, fmt.Errorf("Primary IPv6 not found: %s", primaryIPv6IDOrName)
			}
			publicNet.IPv6 = primaryIPv6
		}
		if primaryIPv4IDOrName != "" || primaryIPv6IDOrName != "" {
			createOpts.PublicNet = publicNet
		}
		result, _, err := s.Client().LoadBalancer().Create(s, createOpts)
		if err != nil {
			return nil, nil, err
		}

		if err := s.WaitForActions(s, cmd, result.Action); err != nil {
			return nil, nil, err
		}
		cmd.Printf("Load Balancer %d created\n", result.LoadBalancer.ID)

		if protectionOpts.Delete != nil {
			if err := ChangeProtectionCmds.ChangeProtection(s, cmd, result.LoadBalancer, true, protectionOpts); err != nil {
				return nil, nil, err
			}
		}

		loadBalancer, _, err := s.Client().LoadBalancer().GetByID(s, result.LoadBalancer.ID)
		if err != nil {
			return nil, nil, err
		}
		if loadBalancer == nil {
			return nil, nil, fmt.Errorf("Load Balancer not found: %d", result.LoadBalancer.ID)
		}

		return loadBalancer, util.Wrap("load_balancer", hcloud.SchemaFromLoadBalancer(loadBalancer)), nil
	},

	PrintResource: func(_ state.State, cmd *cobra.Command, loadBalancer *hcloud.LoadBalancer) {
		cmd.Printf("IPv4: %s\n", loadBalancer.PublicNet.IPv4.IP.String())
		cmd.Printf("IPv6: %s\n", loadBalancer.PublicNet.IPv6.IP.String())
	},
}
