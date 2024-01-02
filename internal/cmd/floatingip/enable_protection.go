package floatingip

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func getChangeProtectionOpts(enable bool, flags []string) (hcloud.FloatingIPChangeProtectionOpts, error) {

	opts := hcloud.FloatingIPChangeProtectionOpts{}

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

func changeProtection(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command,
	floatingIP *hcloud.FloatingIP, enable bool, opts hcloud.FloatingIPChangeProtectionOpts) error {

	if opts.Delete == nil {
		return nil
	}

	action, _, err := client.FloatingIP().ChangeProtection(ctx, floatingIP, opts)
	if err != nil {
		return err
	}

	if err := waiter.ActionProgress(cmd, ctx, action); err != nil {
		return err
	}

	if enable {
		cmd.Printf("Resource protection enabled for floating IP %d\n", floatingIP.ID)
	} else {
		cmd.Printf("Resource protection disabled for floating IP %d\n", floatingIP.ID)
	}
	return nil
}

var EnableProtectionCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:   "enable-protection [FLAGS] FLOATINGIP PROTECTIONLEVEL [PROTECTIONLEVEL...]",
			Short: "Enable resource protection for a Floating IP",
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

		opts, err := getChangeProtectionOpts(true, args[1:])
		if err != nil {
			return err
		}

		return changeProtection(ctx, client, waiter, cmd, floatingIP, true, opts)
	},
}
