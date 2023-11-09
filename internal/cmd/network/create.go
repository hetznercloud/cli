package network

import (
	"context"
	"net"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var CreateCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "create [FLAGS]",
			Short: "Create a network",
			Args:  cobra.NoArgs,
		}

		cmd.Flags().String("name", "", "Network name (required)")
		cmd.MarkFlagRequired("name")

		cmd.Flags().IPNet("ip-range", net.IPNet{}, "Network IP range (required)")
		cmd.MarkFlagRequired("ip-range")

		cmd.Flags().Bool("expose-routes-to-vswitch", false, "Expose routes from this network to the vSwitch connection. It only takes effect if a vSwitch connection is active.")

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		cmd.Flags().StringSlice("enable-protection", []string{}, "Enable protection (delete) (default: none)")
		cmd.RegisterFlagCompletionFunc("enable-protection", cmpl.SuggestCandidates("delete"))
		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		ipRange, _ := cmd.Flags().GetIPNet("ip-range")
		labels, _ := cmd.Flags().GetStringToString("label")
		exposeRoutesToVSwitch, _ := cmd.Flags().GetBool("expose-routes-to-vswitch")
		protection, _ := cmd.Flags().GetStringSlice("enable-protection")

		protectionOpts, err := getChangeProtectionOpts(true, protection)
		if err != nil {
			return err
		}

		createOpts := hcloud.NetworkCreateOpts{
			Name:                  name,
			IPRange:               &ipRange,
			Labels:                labels,
			ExposeRoutesToVSwitch: exposeRoutesToVSwitch,
		}

		network, _, err := client.Network().Create(ctx, createOpts)
		if err != nil {
			return err
		}

		cmd.Printf("Network %d created\n", network.ID)

		return changeProtection(ctx, client, waiter, cmd, network, true, protectionOpts)
	},
}
