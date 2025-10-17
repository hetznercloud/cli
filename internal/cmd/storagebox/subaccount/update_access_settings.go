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

var UpdateAccessSettingsCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {

		cmd := &cobra.Command{
			Use:   "update-access-settings [options] <storage-box> <subaccount>",
			Short: "Update access settings of the Storage Box Subaccount",
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.StorageBox().Names),
				SuggestSubaccounts(client),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().Bool("enable-samba", false, "Whether the Samba subsystem should be enabled (true, false)")
		cmd.Flags().Bool("enable-ssh", false, "Whether the SSH subsystem should be enabled (true, false)")
		cmd.Flags().Bool("enable-webdav", false, "Whether the WebDAV subsystem should be enabled (true, false)")
		cmd.Flags().Bool("reachable-externally", false, "Whether the Storage Box should be accessible from outside the Hetzner network (true, false)")
		cmd.Flags().Bool("readonly", false, "Whether the Subaccount should be read-only (true, false)")

		cmd.MarkFlagsOneRequired("enable-samba", "enable-ssh", "enable-webdav", "reachable-externally", "readonly")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		storageBoxIDOrName, subaccountIDOrName := args[0], args[1]
		enableSamba, _ := cmd.Flags().GetBool("enable-samba")
		enableSSH, _ := cmd.Flags().GetBool("enable-ssh")
		enableWebDAV, _ := cmd.Flags().GetBool("enable-webdav")
		reachableExternally, _ := cmd.Flags().GetBool("reachable-externally")
		readonly, _ := cmd.Flags().GetBool("readonly")

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

		var opts hcloud.StorageBoxSubaccountAccessSettingsUpdateOpts
		if cmd.Flags().Changed("enable-samba") {
			opts.SambaEnabled = &enableSamba
		}
		if cmd.Flags().Changed("enable-ssh") {
			opts.SSHEnabled = &enableSSH
		}
		if cmd.Flags().Changed("enable-webdav") {
			opts.WebDAVEnabled = &enableWebDAV
		}
		if cmd.Flags().Changed("reachable-externally") {
			opts.ReachableExternally = &reachableExternally
		}
		if cmd.Flags().Changed("readonly") {
			opts.Readonly = &readonly
		}

		action, _, err := s.Client().StorageBox().UpdateSubaccountAccessSettings(s, subaccount, opts)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		cmd.Printf("Access settings updated for Storage Box Subaccount %d\n", subaccount.ID)
		return nil
	},
}
