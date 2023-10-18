package primaryip

import (
	"context"
	"fmt"
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
)

func getChangeProtectionOpts(enable bool, flags []string) (hcloud.PrimaryIPChangeProtectionOpts, error) {

	opts := hcloud.PrimaryIPChangeProtectionOpts{}

	var unknown []string
	for _, arg := range flags {
		switch strings.ToLower(arg) {
		case "delete":
			opts.Delete = enable
		default:
			unknown = append(unknown, arg)
		}
	}
	if len(unknown) > 0 {
		return opts, fmt.Errorf("unknown protection level: %s", strings.Join(unknown, ", "))
	}

	return opts, nil
}

func changeProtection(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, primaryIp *hcloud.PrimaryIP, enable bool, opts hcloud.PrimaryIPChangeProtectionOpts) error {

	opts.ID = primaryIp.ID

	action, _, err := client.PrimaryIP().ChangeProtection(ctx, opts)
	if err != nil {
		return err
	}

	if err := waiter.ActionProgress(ctx, action); err != nil {
		return err
	}

	if enable {
		fmt.Printf("Resource protection enabled for primary IP %d\n", opts.ID)
	} else {
		fmt.Printf("Resource protection disabled for primary IP %d\n", opts.ID)
	}
	return nil
}

var EnableProtectionCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "enable-protection PRIMARYIP [PROTECTIONLEVEL...]",
			Short: "Enable Protection for a Primary IP",
			Args:  cobra.MinimumNArgs(1),
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.PrimaryIP().Names),
				cmpl.SuggestCandidates("delete"),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, actionWaiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		primaryIP, _, err := client.PrimaryIP().Get(ctx, idOrName)
		if err != nil {
			return err
		}
		if primaryIP == nil {
			return fmt.Errorf("Primary IP not found: %v", idOrName)
		}

		// This command used to have the "delete" protection level as the default.
		// To avoid a breaking change, we now add it if no level is defined.
		if len(args) < 2 {
			args = append(args, "delete")
		}

		opts, err := getChangeProtectionOpts(true, args[1:])
		if err != nil {
			return err
		}

		return changeProtection(ctx, client, actionWaiter, primaryIP, true, opts)
	},
}
