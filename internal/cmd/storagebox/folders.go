package storagebox

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var FoldersCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "folders <storage-box>",
			Short:                 "List folders of a Storage Box",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.StorageBox().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String("path", "", "Relative path for which the listing is to be made")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		path, _ := cmd.Flags().GetString("path")

		storageBox, _, err := s.Client().StorageBox().Get(s, idOrName)
		if err != nil {
			return err
		}
		if storageBox == nil {
			return fmt.Errorf("Storage Box not found: %s", idOrName)
		}

		var opts hcloud.StorageBoxFoldersOpts
		if cmd.Flags().Changed("path") {
			opts.Path = path
		}

		result, _, err := s.Client().StorageBox().Folders(s, storageBox, opts)
		if err != nil {
			return err
		}

		if len(result.Folders) == 0 {
			cmd.Println("No folders found.")
		} else {
			cmd.Println("Folders:")
			for _, folder := range result.Folders {
				cmd.Printf("- %s\n", folder)
			}
		}
		return nil
	},
}
