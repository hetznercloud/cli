package subaccount

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
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
		idOrName := args[0]
		password, _ := cmd.Flags().GetString("password")

		storageBox, _, err := s.Client().StorageBox().Get(s, idOrName)
		if err != nil {
			return err
		}
		if storageBox == nil {
			return fmt.Errorf("Storage Box not found: %s", idOrName)
		}

		id, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid Storage Box Subaccount ID: %s", args[1])
		}
		subaccount, _, err := s.Client().StorageBox().GetSubaccountByID(s, storageBox, id)
		if err != nil {
			return err
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
}
