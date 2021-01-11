package network

import (
	"encoding/json"
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newDescribeCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] NETWORK",
		Short:                 "Describe a network",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.NetworkNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runNetworkDescribe),
	}
	output.AddFlag(cmd, output.OptionJSON(), output.OptionFormat())
	return cmd
}

func runNetworkDescribe(cli *state.State, cmd *cobra.Command, args []string) error {
	outputFlags := output.FlagsForCommand(cmd)

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
		return networkDescribeJSON(resp)
	case outputFlags.IsSet("format"):
		return util.DescribeFormat(network, outputFlags["format"][0])
	default:
		return networkDescribeText(cli, network)
	}
}

func networkDescribeText(cli *state.State, network *hcloud.Network) error {
	fmt.Printf("ID:\t\t%d\n", network.ID)
	fmt.Printf("Name:\t\t%s\n", network.Name)
	fmt.Printf("Created:\t%s (%s)\n", util.Datetime(network.Created), humanize.Time(network.Created))
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
			if subnet.Type == hcloud.NetworkSubnetTypeVSwitch {
				fmt.Printf("    vSwitch ID:\t\t%d\n", subnet.VSwitchID)
			}
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
	fmt.Printf("  Delete:\t%s\n", util.YesNo(network.Protection.Delete))

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

func networkDescribeJSON(resp *hcloud.Response) error {
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if network, ok := data["network"]; ok {
		return util.DescribeJSON(network)
	}
	if networks, ok := data["networks"].([]interface{}); ok {
		return util.DescribeJSON(networks[0])
	}
	return util.DescribeJSON(data)
}
