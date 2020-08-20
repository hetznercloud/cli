package cli

import (
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newNetworkDescribeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] NETWORK",
		Short:                 "Describe a network",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.NetworkNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runNetworkDescribe),
	}
	addOutputFlag(cmd, outputOptionJSON(), outputOptionFormat())
	return cmd
}

func runNetworkDescribe(cli *CLI, cmd *cobra.Command, args []string) error {
	outputFlags := outputFlagsForCommand(cmd)

	idOrName := args[0]
	network, resp, err := cli.Client().Network.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if network == nil {
		return fmt.Errorf("network not found: %s", idOrName)
	}

	switch {
	case outputFlags.IsSet("json"):
		return serverDescribeJSON(resp)
	case outputFlags.IsSet("format"):
		return describeFormat(network, outputFlags["format"][0])
	default:
		return networkDescribeText(cli, network)
	}
}

func networkDescribeText(cli *CLI, network *hcloud.Network) error {
	fmt.Printf("ID:\t\t%d\n", network.ID)
	fmt.Printf("Name:\t\t%s\n", network.Name)
	fmt.Printf("Created:\t%s (%s)\n", datetime(network.Created), humanize.Time(network.Created))
	fmt.Printf("IP Range:\t%s\n", network.IPRange.String())

	fmt.Printf("Subnets:\n")
	if len(network.Subnets) == 0 {
		fmt.Print("  No subnets\n")
	} else {
		for _, subnet := range network.Subnets {
			fmt.Printf("  - Type:\t\t%s\n", subnet.Type)
			fmt.Printf("    Network Zone:\t%s\n", subnet.NetworkZone)
			fmt.Printf("    IP Range:\t\t%s\n", subnet.IPRange.String())
			fmt.Printf("    Gateway:\t\t%s\n", subnet.Gateway.String())
		}
	}

	fmt.Printf("Routes:\n")
	if len(network.Routes) == 0 {
		fmt.Print("  No routes\n")
	} else {
		for _, route := range network.Routes {
			fmt.Printf("  - Destination:\t%s\n", route.Destination.String())
			fmt.Printf("    Gateway:\t\t%s\n", route.Gateway.String())
		}
	}

	fmt.Printf("Protection:\n")
	fmt.Printf("  Delete:\t%s\n", yesno(network.Protection.Delete))

	fmt.Print("Labels:\n")
	if len(network.Labels) == 0 {
		fmt.Print("  No labels\n")
	} else {
		for key, value := range network.Labels {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	return nil
}
