package snapshot

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var CreateCmd = base.CreateCmd[*hcloud.StorageBoxSnapshot]{
	BaseCobraCommand: func(c hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "create [--description <description>] <storage-box>",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(c.StorageBox().Names)),
			Short:                 "Create a Storage Box Snapshot",
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().String("description", "", "Description of the Storage Box Snapshot")
		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) (*hcloud.StorageBoxSnapshot, any, error) {
		storageBoxIDOrName := args[0]
		description, _ := cmd.Flags().GetString("description")
		labels, _ := cmd.Flags().GetStringToString("label")

		storageBox, _, err := s.Client().StorageBox().Get(s, storageBoxIDOrName)
		if err != nil {
			return nil, nil, err
		}
		if storageBox == nil {
			return nil, nil, fmt.Errorf("Storage Box not found: %s", storageBoxIDOrName)
		}

		opts := hcloud.StorageBoxSnapshotCreateOpts{
			Labels: labels,
		}
		if cmd.Flags().Changed("description") {
			opts.Description = &description
		}

		result, _, err := s.Client().StorageBox().CreateSnapshot(s, storageBox, opts)
		if err != nil {
			return nil, nil, err
		}
		if err := s.WaitForActions(s, cmd, result.Action); err != nil {
			return nil, nil, err
		}

		snapshot, _, err := s.Client().StorageBox().GetSnapshotByID(s, storageBox, result.Snapshot.ID)
		if err != nil {
			return nil, nil, err
		}
		if snapshot == nil {
			return nil, nil, fmt.Errorf("Storage Box Snapshot not found: %d", result.Snapshot.ID)
		}

		cmd.Printf("Storage Box Snapshot %d created\n", snapshot.ID)

		return snapshot, util.Wrap("snapshot", hcloud.SchemaFromStorageBoxSnapshot(snapshot)), nil
	},
	PrintResource: func(_ state.State, cmd *cobra.Command, snapshot *hcloud.StorageBoxSnapshot) {
		cmd.Printf("Name: %s\n", snapshot.Name)
		cmd.Printf("Size: %s\n", humanize.IBytes(snapshot.Stats.Size))
	},
}
