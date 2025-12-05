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

var ChangeTypeCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "change-type <storage-box> <storage-box-type>",
			Short: "Change type of a Storage Box",
			Long: `Requests a Storage Box to be upgraded or downgraded to another Storage Box Type.
Please note that it is not possible to downgrade to a Storage Box Type that offers less disk space than you are currently using.`,
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.StorageBox().Names),
				cmpl.SuggestCandidatesF(client.StorageBoxType().Names),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		storageBox, _, err := s.Client().StorageBox().Get(s, idOrName)
		if err != nil {
			return err
		}
		if storageBox == nil {
			return fmt.Errorf("Storage Box not found: %s", idOrName)
		}

		storageBoxTypeIDOrName := args[1]
		storageBoxType, _, err := s.Client().StorageBoxType().Get(s, storageBoxTypeIDOrName)
		if err != nil {
			return err
		}
		if storageBoxType == nil {
			return fmt.Errorf("Storage Box Type not found: %s", storageBoxTypeIDOrName)
		}

		if storageBoxType.IsDeprecated() {
			cmd.Print(warningDeprecatedStorageBoxType(storageBoxType))
		}

		opts := hcloud.StorageBoxChangeTypeOpts{
			StorageBoxType: storageBoxType,
		}
		action, _, err := s.Client().StorageBox().ChangeType(s, storageBox, opts)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		cmd.Printf("Storage Box %d upgraded to type %s\n", storageBox.ID, storageBoxType.Name)
		return nil
	},
}
