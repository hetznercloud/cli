package snapshot

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
)

func SuggestSnapshots(client hcapi2.Client) cobra.CompletionFunc {
	return cmpl.SuggestCandidatesCtx(func(cmd *cobra.Command, args []string) []string {
		if len(args) == 0 {
			return nil
		}

		storageBox, _, err := client.StorageBox().Get(cmd.Context(), args[0])
		if err != nil || storageBox == nil {
			return nil
		}

		snapshots, err := client.StorageBox().AllSnapshots(context.Background(), storageBox)
		if err != nil {
			return nil
		}

		snapshotNames := make([]string, 0, len(snapshots))
		for _, snapshot := range snapshots {
			snapshotNames = append(snapshotNames, snapshot.Name)
		}
		return snapshotNames
	})
}
