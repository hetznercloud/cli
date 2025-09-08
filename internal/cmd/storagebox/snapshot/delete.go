package snapshot

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

// TODO: modify base delete command to support multiple positional arguments? -> delete multiple snapshots at once
var DeleteCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "delete <storage-box> <snapshot>",
			Short: "Delete a Storage Box Snapshot",
			Args:  util.Validate,
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.StorageBox().Names),
				SuggestSnapshots(client),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		storageBoxIDOrName := args[0]

		storageBox, _, err := s.Client().StorageBox().Get(s, storageBoxIDOrName)
		if err != nil {
			return err
		}
		if storageBox == nil {
			return fmt.Errorf("Storage Box not found: %s", storageBoxIDOrName)
		}

		snapshot, _, err := s.Client().StorageBox().GetSnapshot(s, storageBox, args[1])
		if err != nil {
			return err
		}

		action, _, err := s.Client().StorageBox().DeleteSnapshot(s, snapshot)
		if err != nil {
			return err
		}

		err = s.WaitForActions(s, cmd, action)
		if err != nil {
			return err
		}

		cmd.Printf("Storage Box Snapshot %d deleted\n", snapshot.ID)
		return nil
	},
}
