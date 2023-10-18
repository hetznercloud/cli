package network

import (
	"context"
	"fmt"
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
)

func getChangeProtectionOpts(enable bool, flags []string) (hcloud.NetworkChangeProtectionOpts, error) {

	opts := hcloud.NetworkChangeProtectionOpts{}

	var unknown []string
	for _, arg := range flags {
		switch strings.ToLower(arg) {
		case "delete":
			opts.Delete = hcloud.Ptr(enable)
		default:
			unknown = append(unknown, arg)
		}
	}
	if len(unknown) > 0 {
		return opts, fmt.Errorf("unknown protection level: %s", strings.Join(unknown, ", "))
	}

	return opts, nil
}

func changeProtection(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, network *hcloud.Network, enable bool, opts hcloud.NetworkChangeProtectionOpts) error {

	if opts.Delete == nil {
		return nil
	}

	action, _, err := client.Network().ChangeProtection(ctx, network, opts)
	if err != nil {
		return err
	}

	if err := waiter.ActionProgress(ctx, action); err != nil {
		return err
	}

	if enable {
		fmt.Printf("Resource protection enabled for network %d\n", network.ID)
	} else {
		fmt.Printf("Resource protection disabled for network %d\n", network.ID)
	}
	return nil
}

var EnableProtectionCommand = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:   "enable-protection [FLAGS] NETWORK PROTECTIONLEVEL [PROTECTIONLEVEL...]",
			Short: "Enable resource protection for a network",
			Args:  cobra.MinimumNArgs(2),
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.Network().Names),
				cmpl.SuggestCandidates("delete"),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		network, _, err := client.Network().Get(ctx, idOrName)
		if err != nil {
			return err
		}
		if network == nil {
			return fmt.Errorf("network not found: %s", idOrName)
		}

		opts, err := getChangeProtectionOpts(true, args[1:])
		if err != nil {
			return err
		}

		return changeProtection(ctx, client, waiter, network, true, opts)
	},
}
