package cli

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newLoadBalancerCreateCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create [FLAGS]",
		Short:                 "Create a Load Balancer",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runLoadBalancerCreate),
	}

	cmd.Flags().String("name", "", "Load Balancer name")
	cmd.MarkFlagRequired("name")

	cmd.Flags().String("type", "", "Load Balancer type (ID or name)")
	cmd.Flag("type").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_load_balancer_type_names"},
	}
	cmd.MarkFlagRequired("type")

	cmd.Flags().String("algorithm-type", "", "Algorithm Type name (round_robin or least_connections)")

	cmd.Flag("algorithm-type").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_load_balancer_algorithm_types"},
	}
	cmd.Flags().String("location", "", "Location (ID or name)")
	cmd.Flag("location").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_location_names"},
	}

	cmd.Flags().String("network-zone", "", "Network Zone")
	cmd.Flag("network-zone").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_network_zones"},
	}

	cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

	return cmd
}

func runLoadBalancerCreate(cli *CLI, cmd *cobra.Command, args []string) error {
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
