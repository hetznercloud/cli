package storagebox

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/storagebox/snapshot"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var RollbackSnapshotCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {

		cmd := &cobra.Command{
			Use:                   "rollback-snapshot --snapshot <snapshot> <storage-box>",
			Short:                 "Rolls back the Storage Box to the given Snapshot",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.StorageBox().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().String("snapshot", "", "The name or ID of the snapshot to roll back to")
		util.MarkFlagRequired(cmd, "snapshot")
		cmpl.RegisterFlagCompletion(cmd, "snapshot", snapshot.SuggestSnapshots(client))

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		snapshotIDOrName, _ := cmd.Flags().GetString("snapshot")

		storageBox, _, err := s.Client().StorageBox().Get(cmd.Context(), idOrName)
		if err != nil {
			return err
		}
		if storageBox == nil {
			return fmt.Errorf("Storage Box not found: %s", idOrName)
		}

		snapshot, _, err := s.Client().StorageBox().GetSnapshot(cmd.Context(), storageBox, snapshotIDOrName)
		if err != nil {
			return err
		}
		if snapshot == nil {
			return fmt.Errorf("Storage Box Snapshot not found: %s", snapshotIDOrName)
		}

		action, _, err := s.Client().StorageBox().RollbackSnapshot(cmd.Context(), storageBox, hcloud.StorageBoxRollbackSnapshotOpts{
			Snapshot: snapshot,
		})
		if err != nil {
			return err
		}

		if err := s.WaitForActions(cmd.Context(), cmd, action); err != nil {
			return err
		}

		cmd.Printf("Rolled back Storage Box %d to Snapshot %d\n", storageBox.ID, snapshot.ID)
		return nil
	},
}
