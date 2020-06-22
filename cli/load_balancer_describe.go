package cli

import (
	"fmt"
	humanize "github.com/dustin/go-humanize"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newLoadBalancerDescribeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] LOADBALANCER",
		Short:                 "Describe a Load Balancer",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runLoadBalancerDescribe),
	}
	addOutputFlag(cmd, outputOptionJSON(), outputOptionFormat())
	//cmd.Flags().Bool("expand-targets", false, "Expand all label_selector targets")
	return cmd
}

func runLoadBalancerDescribe(cli *CLI, cmd *cobra.Command, args []string) error {
	outputFlags := outputFlagsForCommand(cmd)
	idOrName := args[0]
	loadBalancer, resp, err := cli.Client().LoadBalancer.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if loadBalancer == nil {
		return fmt.Errorf("loadBalancer not found: %s", idOrName)
	}

	switch {
	case outputFlags.IsSet("json"):
		return serverDescribeJSON(resp)
	case outputFlags.IsSet("format"):
		return describeFormat(loadBalancer, outputFlags["format"][0])
	default:
		return loadBalancerDescribeText(cli, loadBalancer)
	}
}

func loadBalancerDescribeText(cli *CLI, loadBalancer *hcloud.LoadBalancer) error {
	fmt.Printf("ID:\t\t\t\t%d\n", loadBalancer.ID)
	fmt.Printf("Name:\t\t\t\t%s\n", loadBalancer.Name)
	fmt.Printf("Created:\t\t\t%s (%s)\n", datetime(loadBalancer.Created), humanize.Time(loadBalancer.Created))
	fmt.Printf("Public Net:\n")
	fmt.Printf("  Enabled:\t\t\t%s\n", yesno(loadBalancer.PublicNet.Enabled))
	fmt.Printf("  IPv4:\t\t\t\t%s\n", loadBalancer.PublicNet.IPv4.IP.String())
	fmt.Printf("  IPv6:\t\t\t\t%s\n", loadBalancer.PublicNet.IPv6.IP.String())

	fmt.Printf("Private Net:\n")
	if len(loadBalancer.PrivateNet) > 0 {
		for _, n := range loadBalancer.PrivateNet {
			network, _, err := cli.client.Network.GetByID(cli.Context, n.Network.ID)
			if err != nil {
				return fmt.Errorf("error fetching network: %v", err)
			}
			fmt.Printf("  - ID:\t\t\t%d\n", network.ID)
			fmt.Printf("    Name:\t\t%s\n", network.Name)
			fmt.Printf("    IP:\t\t\t%s\n", n.IP.String())
		}
	} else {
		fmt.Printf("    No Private Network\n")
	}
	fmt.Printf("Algorithm:\t\t\t%s\n", loadBalancer.Algorithm.Type)

	fmt.Printf("Load Balancer Type:\t\t%s (ID: %d)\n", loadBalancer.LoadBalancerType.Name, loadBalancer.LoadBalancerType.ID)
	fmt.Printf("  ID:\t\t\t\t%d\n", loadBalancer.LoadBalancerType.ID)
	fmt.Printf("  Name:\t\t\t\t%s\n", loadBalancer.LoadBalancerType.Name)
	fmt.Printf("  Description:\t\t\t%s\n", loadBalancer.LoadBalancerType.Description)
	fmt.Printf("  Max Services:\t\t\t%d\n", loadBalancer.LoadBalancerType.MaxServices)
	fmt.Printf("  Max Connections:\t\t%d\n", loadBalancer.LoadBalancerType.MaxConnections)
	fmt.Printf("  Max Targets:\t\t\t%d\n", loadBalancer.LoadBalancerType.MaxTargets)
	fmt.Printf("  Max assigned Certificates:\t%d\n", loadBalancer.LoadBalancerType.MaxAssignedCertificates)

	fmt.Printf("Services:\n")
	if len(loadBalancer.Services) == 0 {
		fmt.Print("  No services\n")
	} else {
		for _, service := range loadBalancer.Services {
			fmt.Printf("  - Protocol:\t\t\t%s\n", service.Protocol)
			fmt.Printf("    Listen Port:\t\t%d\n", service.ListenPort)
			fmt.Printf("    Destination Port:\t\t%d\n", service.DestinationPort)
			fmt.Printf("    Proxy Protocol:\t\t%s\n", yesno(service.Proxyprotocol))
			if service.Protocol != hcloud.LoadBalancerServiceProtocolTCP {
				fmt.Printf("    Sticky Sessions:\t\t%s\n",yesno(service.HTTP.StickySessions))
				if service.HTTP.StickySessions {
					fmt.Printf("    Sticky Cookie Name:\t\t%s\n", service.HTTP.CookieName)
					fmt.Printf("    Sticky Cookie Lifetime:\t%vs\n", service.HTTP.CookieLifetime.Seconds())
				}
				if service.Protocol == hcloud.LoadBalancerServiceProtocolHTTPS {
					fmt.Printf("    Certificates:\n")
					for _, cert := range service.HTTP.Certificates {
						fmt.Printf("      - ID: \t\t\t%v\n", cert.ID)
					}
				}
			}

			fmt.Printf("    Health Check:\n")
			fmt.Printf("      Protocol:\t\t\t%s\n", service.HealthCheck.Protocol)
			fmt.Printf("      Timeout:\t\t\t%vs\n", service.HealthCheck.Timeout.Seconds())
			fmt.Printf("      Interval:\t\t\tevery %vs\n", service.HealthCheck.Interval.Seconds())
			fmt.Printf("      Retries:\t\t\t%d\n", service.HealthCheck.Retries)
			if service.HealthCheck.Protocol != hcloud.LoadBalancerServiceProtocolTCP {
				fmt.Printf("      HTTP Domain:\t\t%s\n", service.HealthCheck.HTTP.Domain)
				fmt.Printf("      HTTP Path:\t\t%s\n", service.HealthCheck.HTTP.Path)
				fmt.Printf("      Response:\t\t%s\n", service.HealthCheck.HTTP.Response)
				fmt.Printf("      TLS:\t\t\t%s\n", yesno(service.HealthCheck.HTTP.TLS))
				fmt.Printf("      Status Codes:\t\t%v\n", service.HealthCheck.HTTP.StatusCodes)
			}
		}
	}

	fmt.Printf("Targets:\n")
	if len(loadBalancer.Targets) == 0 {
		fmt.Print("  No targets\n")
	} else {
		for _, target := range loadBalancer.Targets {
			fmt.Printf("  - Type:\t\t\t%s\n", target.Type)
			if target.Server != nil {
				fmt.Printf("    Server:\n")
				fmt.Printf("      ID:\t\t\t%d\n", target.Server.Server.ID)
				fmt.Printf("      Name:\t\t\t%s\n", cli.GetServerName(target.Server.Server.ID))
				fmt.Printf("    Use Private IP:\t\t%s\n", yesno(target.UsePrivateIP))
				fmt.Printf("    Status:\n")
				for _, healthStatus := range target.HealthStatus {
					fmt.Printf("    - Service:\t\t\t%d\n", healthStatus.ListenPort)
					fmt.Printf("      Status:\t\t\t%s\n", healthStatus.Status)
				}
			}
		}
	}

	fmt.Printf("Protection:\n")
	fmt.Printf("  Delete:\t%s\n", yesno(loadBalancer.Protection.Delete))

	fmt.Print("Labels:\n")
	if len(loadBalancer.Labels) == 0 {
		fmt.Print("  No labels\n")
	} else {
		for key, value := range loadBalancer.Labels {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	return nil
}
