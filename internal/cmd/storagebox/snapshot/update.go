package snapshot

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var UpdateCmd = base.UpdateCmd{
	ResourceNameSingular:       "Storage Box Snapshot",
	ShortDescription:           "Update a Storage Box Snapshot",
	NameSuggestions:            func(c hcapi2.Client) func() []string { return c.StorageBox().Names },
	PositionalArgumentOverride: []string{"storage-box", "snapshot"},
	FetchWithArgs: func(s state.State, _ *cobra.Command, args []string) (any, *hcloud.Response, error) {
		storageBox, _, err := s.Client().StorageBox().Get(s, args[0])
		if err != nil {
			return nil, nil, err
		}
		if storageBox == nil {
			return nil, nil, fmt.Errorf("Storage Box not found: %s", args[0])
		}
		return s.Client().StorageBox().GetSnapshot(s, storageBox, args[1])
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("description", "", "Description of the Storage Box Snapshot")
		cmd.MarkFlagsOneRequired("description")
	},
	Update: func(s state.State, cmd *cobra.Command, resource interface{}, _ map[string]pflag.Value) error {
		snapshot := resource.(*hcloud.StorageBoxSnapshot)
		var opts hcloud.StorageBoxSnapshotUpdateOpts
		if cmd.Flags().Changed("description") {
			description, _ := cmd.Flags().GetString("description")
			opts.Description = &description
		}
		_, _, err := s.Client().StorageBox().UpdateSnapshot(s, snapshot, opts)
		if err != nil {
			return err
		}
		return nil
	},
	Experimental: experimental.StorageBoxes,
}
