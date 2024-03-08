package loadbalancer

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var CreateCmd = base.CreateCmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "create [options] --name <name> --type <type>",
			Short:                 "Create a Load Balancer",
			Args:                  util.ValidateExact,
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().String("name", "", "Load Balancer name (required)")
		cmd.MarkFlagRequired("name")

		cmd.Flags().String("type", "", "Load Balancer type (ID or name) (required)")
		cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidatesF(client.LoadBalancerType().Names))
		cmd.MarkFlagRequired("type")

		cmd.Flags().String("algorithm-type", "", "Algorithm Type name (round_robin or least_connections)")
		cmd.RegisterFlagCompletionFunc("algorithm-type", cmpl.SuggestCandidates(
			string(hcloud.LoadBalancerAlgorithmTypeLeastConnections),
			string(hcloud.LoadBalancerAlgorithmTypeRoundRobin),
		))
		cmd.Flags().String("location", "", "Location (ID or name)")
		cmd.RegisterFlagCompletionFunc("location", cmpl.SuggestCandidatesF(client.Location().Names))

		cmd.Flags().String("network-zone", "", "Network Zone")
		cmd.RegisterFlagCompletionFunc("network-zone", cmpl.SuggestCandidatesF(client.Location().NetworkZones))

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		cmd.Flags().StringSlice("enable-protection", []string{}, "Enable protection (delete) (default: none)")
		cmd.RegisterFlagCompletionFunc("enable-protection", cmpl.SuggestCandidates("delete"))

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) (any, any, error) {
		name, _ := cmd.Flags().GetString("name")
		serverType, _ := cmd.Flags().GetString("type")
		algorithmType, _ := cmd.Flags().GetString("algorithm-type")
		location, _ := cmd.Flags().GetString("location")
		networkZone, _ := cmd.Flags().GetString("network-zone")
		labels, _ := cmd.Flags().GetStringToString("label")
		protection, _ := cmd.Flags().GetStringSlice("enable-protection")

		protectionOpts, err := getChangeProtectionOpts(true, protection)
		if err != nil {
			return nil, nil, err
		}

		createOpts := hcloud.LoadBalancerCreateOpts{
			Name: name,
			LoadBalancerType: &hcloud.LoadBalancerType{
				Name: serverType,
			},
			Labels: labels,
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
		result, _, err := s.Client().LoadBalancer().Create(s, createOpts)
		if err != nil {
			return nil, nil, err
		}

		if err := s.ActionProgress(cmd, s, result.Action); err != nil {
			return nil, nil, err
		}
		loadBalancer, _, err := s.Client().LoadBalancer().GetByID(s, result.LoadBalancer.ID)
		if err != nil {
			return nil, nil, err
		}
		cmd.Printf("Load Balancer %d created\n", loadBalancer.ID)

		if err := changeProtection(s, cmd, loadBalancer, true, protectionOpts); err != nil {
			return nil, nil, err
		}

		return loadBalancer, util.Wrap("load_balancer", hcloud.SchemaFromLoadBalancer(loadBalancer)), nil
	},

	PrintResource: func(_ state.State, cmd *cobra.Command, resource any) {
		loadBalancer := resource.(*hcloud.LoadBalancer)
		cmd.Printf("IPv4: %s\n", loadBalancer.PublicNet.IPv4.IP.String())
		cmd.Printf("IPv6: %s\n", loadBalancer.PublicNet.IPv6.IP.String())
	},
}
