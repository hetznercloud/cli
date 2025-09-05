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

var UpdateAccessSettingsCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {

		cmd := &cobra.Command{
			Use:                   "update-access-settings [options] <storage-box>",
			Short:                 "Update access settings of the primary Storage Box account",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.StorageBox().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().Bool("samba-enabled", false, "Whether the Samba subsystem should be enabled (true, false)")
		cmd.Flags().Bool("ssh-enabled", false, "Whether the SSH subsystem should be enabled (true, false)")
		cmd.Flags().Bool("webdav-enabled", false, "Whether the WebDAV subsystem should be enabled (true, false)")
		cmd.Flags().Bool("zfs-enabled", false, "Whether the ZFS Snapshot folder should be visible (true, false)")
		cmd.Flags().Bool("reachable-externally", false, "Whether the Storage Box should be accessible from outside the Hetzner network (true, false)")
		cmd.MarkFlagsOneRequired("samba-enabled", "ssh-enabled", "webdav-enabled", "zfs-enabled", "reachable-externally")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		enableSamba, _ := cmd.Flags().GetBool("samba-enabled")
		enableSSH, _ := cmd.Flags().GetBool("ssh-enabled")
		enableWebDAV, _ := cmd.Flags().GetBool("webdav-enabled")
		enableZFS, _ := cmd.Flags().GetBool("zfs-enabled")
		reachableExternally, _ := cmd.Flags().GetBool("reachable-externally")

		storageBox, _, err := s.Client().StorageBox().Get(s, idOrName)
		if err != nil {
			return err
		}
		if storageBox == nil {
			return fmt.Errorf("Storage Box not found: %s", idOrName)
		}

		var opts hcloud.StorageBoxUpdateAccessSettingsOpts
		if cmd.Flags().Changed("samba-enabled") {
			opts.SambaEnabled = &enableSamba
		}
		if cmd.Flags().Changed("ssh-enabled") {
			opts.SSHEnabled = &enableSSH
		}
		if cmd.Flags().Changed("webdav-enabled") {
			opts.WebDAVEnabled = &enableWebDAV
		}
		if cmd.Flags().Changed("zfs-enabled") {
			opts.ZFSEnabled = &enableZFS
		}
		if cmd.Flags().Changed("reachable-externally") {
			opts.ReachableExternally = &reachableExternally
		}

		action, _, err := s.Client().StorageBox().UpdateAccessSettings(s, storageBox, opts)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		cmd.Printf("Access settings updated for Storage Box %d\n", storageBox.ID)
		return nil
	},
}
