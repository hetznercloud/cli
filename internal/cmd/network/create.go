package network

import (
	"net"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var CreateCmd = base.CreateCmd{
	BaseCobraCommand: func(hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "create [options] --name <name> --ip-range <ip-range>",
			Short: "Create a network",
		}

		cmd.Flags().String("name", "", "Network name (required)")
		_ = cmd.MarkFlagRequired("name")

		cmd.Flags().IPNet("ip-range", net.IPNet{}, "Network IP range (required)")
		_ = cmd.MarkFlagRequired("ip-range")

		cmd.Flags().Bool("expose-routes-to-vswitch", false, "Expose routes from this network to the vSwitch connection. It only takes effect if a vSwitch connection is active.")

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		cmd.Flags().StringSlice("enable-protection", []string{}, "Enable protection (delete) (default: none)")
		_ = cmd.RegisterFlagCompletionFunc("enable-protection", cmpl.SuggestCandidates("delete"))
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, _ []string) (any, any, error) {
		name, _ := cmd.Flags().GetString("name")
		ipRange, _ := cmd.Flags().GetIPNet("ip-range")
		labels, _ := cmd.Flags().GetStringToString("label")
		exposeRoutesToVSwitch, _ := cmd.Flags().GetBool("expose-routes-to-vswitch")
		protection, _ := cmd.Flags().GetStringSlice("enable-protection")

		protectionOpts, err := getChangeProtectionOpts(true, protection)
		if err != nil {
			return nil, nil, err
		}

		createOpts := hcloud.NetworkCreateOpts{
			Name:                  name,
			IPRange:               &ipRange,
			Labels:                labels,
			ExposeRoutesToVSwitch: exposeRoutesToVSwitch,
		}

		network, _, err := s.Client().Network().Create(s, createOpts)
		if err != nil {
			return nil, nil, err
		}

		cmd.Printf("Network %d created\n", network.ID)

		if err := changeProtection(s, cmd, network, true, protectionOpts); err != nil {
			return nil, nil, err
		}

		return network, util.Wrap("network", hcloud.SchemaFromNetwork(network)), nil
	},
}
