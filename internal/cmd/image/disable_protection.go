package image

import (
	"context"
	"errors"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DisableProtectionCommand = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:   "disable-protection [FLAGS] IMAGE PROTECTIONLEVEL [PROTECTIONLEVEL...]",
			Short: "Disable resource protection for an image",
			Args:  cobra.MinimumNArgs(2),
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.Image().Names),
				cmpl.SuggestCandidates("delete"),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, command *cobra.Command, args []string) error {
		imageID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return errors.New("invalid image ID")
		}
		image := &hcloud.Image{ID: imageID}

		opts, err := getChangeProtectionOpts(false, args[1:])
		if err != nil {
			return err
		}

		return changeProtection(ctx, client, waiter, image, false, opts)
	},
}
