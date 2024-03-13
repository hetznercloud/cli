package volume

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

var ResizeCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "resize --size <size> <volume>",
			Short:                 "Resize a volume",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Volume().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().Int("size", 0, "New size (GB) of the volume (required)")
		cmd.MarkFlagRequired("size")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		volume, _, err := s.Client().Volume().Get(s, args[0])
		if err != nil {
			return err
		}
		if volume == nil {
			return fmt.Errorf("volume not found: %s", args[0])
		}

		size, _ := cmd.Flags().GetInt("size")
		action, _, err := s.Client().Volume().Resize(s, volume, size)
		if err != nil {
			return err
		}

		if err := s.ActionProgress(cmd, s, action); err != nil {
			return err
		}

		cmd.Printf("Volume %d resized\n", volume.ID)
		cmd.Printf("You might need to adjust the filesystem size on the server too\n")
		return nil
	},
}
