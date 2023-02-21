package volume

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
			Use:   "disable-protection [FLAGS] VOLUME PROTECTIONLEVEL [PROTECTIONLEVEL...]",
			Short: "Disable resource protection for a volume",
			Args:  cobra.MinimumNArgs(2),
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.Volume().Names),
				cmpl.SuggestCandidates("delete"),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		volume, _, err := client.Volume().Get(ctx, args[0])
		if err != nil {
			return err
		}
		if volume == nil {
			return fmt.Errorf("volume not found: %s", args[0])
		}

		var unknown []string
		opts := hcloud.VolumeChangeProtectionOpts{}
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

		action, _, err := client.Volume().ChangeProtection(ctx, volume, opts)
		if err != nil {
			return err
		}

		if err := waiter.ActionProgress(ctx, action); err != nil {
			return err
		}

		fmt.Printf("Resource protection disabled for volume %d\n", volume.ID)
		return nil
	},
}
