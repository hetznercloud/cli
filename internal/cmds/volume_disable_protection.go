package cmds

import (
	"fmt"
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newVolumeDisableProtectionCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable-protection [FLAGS] VOLUME PROTECTIONLEVEL [PROTECTIONLEVEL...]",
		Short: "Disable resource protection for a volume",
		Args:  cobra.MinimumNArgs(2),
		ValidArgsFunction: cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(cli.VolumeNames),
			cmpl.SuggestCandidates("delete"),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runVolumeDisableProtection),
	}
	return cmd
}

func runVolumeDisableProtection(cli *state.State, cmd *cobra.Command, args []string) error {
	volume, _, err := cli.Client().Volume.Get(cli.Context, args[0])
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

	action, _, err := cli.Client().Volume.ChangeProtection(cli.Context, volume, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Resource protection disabled for volume %d\n", volume.ID)
	return nil
}
