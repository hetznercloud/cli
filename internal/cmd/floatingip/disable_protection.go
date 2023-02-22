package floatingip

import (
	"context"
	"fmt"
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var DisableProtectionCommand = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:   "disable-protection [FLAGS] FLOATINGIP PROTECTIONLEVEL [PROTECTIONLEVEL...]",
			Short: "Disable resource protection for a Floating IP",
			Args:  cobra.MinimumNArgs(2),
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.FloatingIP().Names),
				cmpl.SuggestCandidates("delete"),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		floatingIP, _, err := client.FloatingIP().Get(ctx, idOrName)
		if err != nil {
			return err
		}
		if floatingIP == nil {
			return fmt.Errorf("Floating IP not found: %v", idOrName)
		}

		var unknown []string
		opts := hcloud.FloatingIPChangeProtectionOpts{}
		for _, arg := range args[1:] {
			switch strings.ToLower(arg) {
			case "delete":
				opts.Delete = hcloud.Bool(false)
			default:
				unknown = append(unknown, arg)
			}
		}
		if len(unknown) > 0 {
			return fmt.Errorf("unknown protection level: %s", strings.Join(unknown, ", "))
		}

		action, _, err := client.FloatingIP().ChangeProtection(ctx, floatingIP, opts)
		if err != nil {
			return err
		}

		if err := waiter.ActionProgress(ctx, action); err != nil {
			return err
		}

		fmt.Printf("Resource protection disabled for Floating IP %d\n", floatingIP.ID)
		return nil
	},
}
