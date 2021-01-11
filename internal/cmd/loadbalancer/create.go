package loadbalancer

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newCreateCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create [FLAGS]",
		Short:                 "Create a Load Balancer",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runLoadBalancerCreate),
	}

	cmd.Flags().String("name", "", "Load Balancer name (required)")
	cmd.MarkFlagRequired("name")

	cmd.Flags().String("type", "", "Load Balancer type (ID or name) (required)")
	cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidatesF(cli.LoadBalancerTypeNames))
	cmd.MarkFlagRequired("type")

	cmd.Flags().String("algorithm-type", "", "Algorithm Type name (round_robin or least_connections)")
	cmd.RegisterFlagCompletionFunc("algorithm-type", cmpl.SuggestCandidates(
		string(hcloud.LoadBalancerAlgorithmTypeLeastConnections),
		string(hcloud.LoadBalancerAlgorithmTypeRoundRobin),
	))
	cmd.Flags().String("location", "", "Location (ID or name)")
	cmd.RegisterFlagCompletionFunc("location", cmpl.SuggestCandidatesF(cli.LocationNames))

	cmd.Flags().String("network-zone", "", "Network Zone")
	cmd.RegisterFlagCompletionFunc("network-zone", cmpl.SuggestCandidates("eu-central"))

	cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

	return cmd
}

func runLoadBalancerCreate(cli *state.State, cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	serverType, _ := cmd.Flags().GetString("type")
	algorithmType, _ := cmd.Flags().GetString("algorithm-type")
	location, _ := cmd.Flags().GetString("location")
	networkZone, _ := cmd.Flags().GetString("network-zone")
	labels, _ := cmd.Flags().GetStringToString("label")

	opts := hcloud.LoadBalancerCreateOpts{
		Name: name,
		LoadBalancerType: &hcloud.LoadBalancerType{
			Name: serverType,
		},
		Labels: labels,
	}
	if algorithmType != "" {
		opts.Algorithm = &hcloud.LoadBalancerAlgorithm{Type: hcloud.LoadBalancerAlgorithmType(algorithmType)}
	}
	if networkZone != "" {
		opts.NetworkZone = hcloud.NetworkZone(networkZone)
	}
	if location != "" {
		opts.Location = &hcloud.Location{Name: location}
	}
	result, _, err := cli.Client().LoadBalancer.Create(cli.Context, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, result.Action); err != nil {
		return err
	}
	loadBalancer, _, err := cli.Client().LoadBalancer.GetByID(cli.Context, result.LoadBalancer.ID)
	if err != nil {
		return err
	}
	fmt.Printf("LoadBalancer %d created\n", loadBalancer.ID)
	fmt.Printf("IPv4: %s\n", loadBalancer.PublicNet.IPv4.IP.String())
	fmt.Printf("IPv6: %s\n", loadBalancer.PublicNet.IPv6.IP.String())
	return nil
}
