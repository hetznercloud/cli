package snapshot

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
)

func SuggestSnapshots(client hcapi2.Client) cobra.CompletionFunc {
	return cmpl.SuggestCandidatesCtxE(func(cmd *cobra.Command, args []string) ([]string, error) {
		if len(args) == 0 {
			return nil, nil
		}

		storageBox, _, err := client.StorageBox().Get(cmd.Context(), args[0])
		if err != nil {
			return nil, err
		}
		if storageBox == nil {
			return nil, nil
		}

		snapshots, err := client.StorageBox().AllSnapshots(cmd.Context(), storageBox)
		if err != nil {
			return nil, err
		}

		snapshotNames := make([]string, 0, len(snapshots))
		for _, snapshot := range snapshots {
			snapshotNames = append(snapshotNames, snapshot.Name)
		}
		return snapshotNames, nil
	})
}
