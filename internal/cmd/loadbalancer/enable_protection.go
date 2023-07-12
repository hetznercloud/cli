package loadbalancer

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

var EnableProtectionCommand = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:   "enable-protection [FLAGS] LOADBALANCER PROTECTIONLEVEL [PROTECTIONLEVEL...]",
			Short: "Enable resource protection for a Load Balancer",
			Args:  cobra.MinimumNArgs(2),
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.LoadBalancer().Names),
				cmpl.SuggestCandidates("delete"),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		LoadBalancer, _, err := client.LoadBalancer().Get(ctx, idOrName)
		if err != nil {
			return err
		}
		if LoadBalancer == nil {
			return fmt.Errorf("Load Balancer not found: %s", idOrName)
		}

		var unknown []string
		opts := hcloud.LoadBalancerChangeProtectionOpts{}
		for _, arg := range args[1:] {
			switch strings.ToLower(arg) {
			case "delete":
				opts.Delete = hcloud.Bool(true)
			default:
				unknown = append(unknown, arg)
			}
		}
		if len(unknown) > 0 {
			return fmt.Errorf("unknown protection level: %s", strings.Join(unknown, ", "))
		}

		action, _, err := client.LoadBalancer().ChangeProtection(ctx, LoadBalancer, opts)
		if err != nil {
			return err
		}

		if err := waiter.ActionProgress(ctx, action); err != nil {
			return err
		}

		fmt.Printf("Resource protection enabled Load Balancer %d\n", LoadBalancer.ID)
		return nil
	},
}
