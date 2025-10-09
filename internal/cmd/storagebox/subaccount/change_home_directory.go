package subaccount

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ChangeHomeDirectoryCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {

		cmd := &cobra.Command{
			Use:   "change-home-directory --home-directory <home-directory> <storage-box> <subaccount>",
			Short: "Update access settings of the Storage Box Subaccount",
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.StorageBox().Names),
				SuggestSubaccounts(client),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().String("home-directory", "", "Home directory of the Subaccount. Will be created if it doesn't exist yet")
		_ = cmd.MarkFlagRequired("home-directory")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		storageBoxIDOrName, subaccountIDOrName := args[0], args[1]
		homeDirectory, _ := cmd.Flags().GetString("home-directory")

		storageBox, _, err := s.Client().StorageBox().Get(s, storageBoxIDOrName)
		if err != nil {
			return err
		}
		if storageBox == nil {
			return fmt.Errorf("Storage Box not found: %s", storageBoxIDOrName)
		}

		subaccount, _, err := s.Client().StorageBox().GetSubaccount(s, storageBox, subaccountIDOrName)
		if err != nil {
			return err
		}
		if subaccount == nil {
			return fmt.Errorf("Storage Box Subaccount not found: %s", subaccountIDOrName)
		}

		action, _, err := s.Client().StorageBox().ChangeSubaccountHomeDirectory(s, subaccount, hcloud.StorageBoxSubaccountChangeHomeDirectoryOpts{
			HomeDirectory: &homeDirectory,
		})
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		cmd.Printf("Home directory updated for Storage Box Subaccount %d\n", subaccount.ID)
		return nil
	},
}
