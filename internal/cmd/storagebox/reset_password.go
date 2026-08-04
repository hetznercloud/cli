package storagebox

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ResetPasswordCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "reset-password --password <password> <storage-box>",
			Short:                 "Reset the password of a Storage Box",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.StorageBox().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().String("password", "", "New password for the Storage Box")
		util.MarkFlagRequired(cmd, "password")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		password, _ := cmd.Flags().GetString("password")

		storageBox, _, err := s.Client().StorageBox().Get(cmd.Context(), idOrName)
		if err != nil {
			return err
		}
		if storageBox == nil {
			return fmt.Errorf("Storage Box not found: %s", idOrName)
		}

		opts := hcloud.StorageBoxResetPasswordOpts{
			Password: password,
		}

		action, _, err := s.Client().StorageBox().ResetPassword(cmd.Context(), storageBox, opts)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(cmd.Context(), cmd, action); err != nil {
			return err
		}

		cmd.Printf("Password of Storage Box %d reset\n", storageBox.ID)
		return nil
	},
}
