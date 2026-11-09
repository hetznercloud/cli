package volume

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var AttachCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "attach [--automount] --assignee-type <type> --assignee <assignee> <volume>",
			Short:                 "Attach a Volume to a Server",
			Long:                  "Attach a Volume to a Server.\n\n" + assigneeDeprecationNotice,
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Volume().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		addAssigneeFlags(cmd, client)
		cmd.MarkFlagsOneRequired("assignee", "server")
		cmd.Flags().Bool("automount", false, "Automount Volume after attach (true, false)")

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

		assignee, assigneeName, err := getAssignee(s, cmd)
		if err != nil {
			return err
		}
		automount, _ := cmd.Flags().GetBool("automount")
		action, _, err := s.Client().Volume().AttachWithOpts(s, volume, hcloud.VolumeAttachOpts{
			Assignee:  assignee,
			Automount: &automount,
		})

		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		cmd.Printf("Volume %d attached to %s\n", volume.ID, assigneeName)
		return nil
	},
}
