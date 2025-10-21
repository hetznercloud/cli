package subaccount

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ResetPasswordCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "reset-password --password <password> <storage-box> <subaccount>",
			Short: "Reset the password of a Storage Box Subaccount",
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.StorageBox().Names),
				SuggestSubaccounts(client),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().String("password", "", "New password for the Storage Box Subaccount")
		_ = cmd.MarkFlagRequired("password")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		storageBoxIDOrName, subaccountIDOrName := args[0], args[1]
		password, _ := cmd.Flags().GetString("password")

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

		opts := hcloud.StorageBoxSubaccountResetPasswordOpts{
			Password: password,
		}

		action, _, err := s.Client().StorageBox().ResetSubaccountPassword(s, subaccount, opts)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		cmd.Printf("Password of Storage Box Subaccount %d reset\n", subaccount.ID)
		return nil
	},
	Experimental: experimental.StorageBoxes,
}
