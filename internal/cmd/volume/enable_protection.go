package volume

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func getChangeProtectionOpts(enable bool, flags []string) (hcloud.VolumeChangeProtectionOpts, error) {

	opts := hcloud.VolumeChangeProtectionOpts{}

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

func changeProtection(s state.State, cmd *cobra.Command,
	volume *hcloud.Volume, enable bool, opts hcloud.VolumeChangeProtectionOpts) error {

	if opts.Delete == nil {
		return nil
	}

	action, _, err := s.Client().Volume().ChangeProtection(s, volume, opts)
	if err != nil {
		return err
	}

	if err := s.ActionProgress(cmd, s, action); err != nil {
		return err
	}

	if enable {
		cmd.Printf("Resource protection enabled for volume %d\n", volume.ID)
	} else {
		cmd.Printf("Resource protection disabled for volume %d\n", volume.ID)
	}
	return nil
}

var EnableProtectionCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:   "enable-protection <volume> <protection-level>...",
			Short: "Enable resource protection for a volume",
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.Volume().Names),
				cmpl.SuggestCandidates("delete"),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		volume, _, err := s.Client().Volume().Get(s, args[0])
		if err != nil {
			return err
		}
		if volume == nil {
			return fmt.Errorf("volume not found: %s", args[0])
		}

		opts, err := getChangeProtectionOpts(true, args[1:])
		if err != nil {
			return err
		}

		return changeProtection(s, cmd, volume, true, opts)
	},
}
