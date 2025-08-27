package storagebox

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var CreateCmd = base.CreateCmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "create [options] --name <name> --type <type> --location <location> --password <password>",
			Short: `Create a new Storage Box`,
		}

		cmd.Flags().String("name", "", "Storage Box name (required)")
		_ = cmd.MarkFlagRequired("name")

		cmd.Flags().String("type", "", "Storage Box Type (ID or name) (required)")
		_ = cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidatesF(client.StorageBoxType().Names))
		_ = cmd.MarkFlagRequired("type")

		cmd.Flags().String("location", "", "Location (ID or name) (required)")
		_ = cmd.MarkFlagRequired("location")
		_ = cmd.RegisterFlagCompletionFunc("location", cmpl.SuggestCandidatesF(client.Location().Names))

		cmd.Flags().String("password", "", "The password that will be set for this Storage Box (required)")
		_ = cmd.MarkFlagRequired("password")

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		// TODO: How to handle SSH Keys? --public-key-from-file?
		// TODO: Fetch SSH Keys from the project
		cmd.Flags().StringArray("ssh-key", []string{}, "SSH public keys in OpenSSH format which should be injected into the Storage Box")

		// TODO: Are we fine with dropping the nested object key ("access_settings") from the flag names?
		cmd.Flags().Bool("enable-samba", false, "Whether the Samba subsystem should be enabled")
		cmd.Flags().Bool("enable-ssh", false, "Whether the SSH subsystem should be enabled")
		cmd.Flags().Bool("enable-webdav", false, "Whether the WebDAV subsystem should be enabled")
		cmd.Flags().Bool("enable-zfs", false, "Whether the ZFS Snapshot folder should be visible")
		cmd.Flags().Bool("reachable-externally", false, "Whether the Storage Box should be accessible from outside the Hetzner network")

		cmd.Flags().StringSlice("enable-protection", []string{}, "Enable protection (delete) (default: none)")
		_ = cmd.RegisterFlagCompletionFunc("enable-protection", cmpl.SuggestCandidates("delete"))

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, strings []string) (any, any, error) {
		name, _ := cmd.Flags().GetString("name")
		sbType, _ := cmd.Flags().GetString("type")
		location, _ := cmd.Flags().GetString("location")
		password, _ := cmd.Flags().GetString("password")
		sshKeys, _ := cmd.Flags().GetStringArray("ssh-key")
		labels, _ := cmd.Flags().GetStringToString("label")

		enableSamba, _ := cmd.Flags().GetBool("enable-samba")
		enableSSH, _ := cmd.Flags().GetBool("enable-ssh")
		enableWebDAV, _ := cmd.Flags().GetBool("enable-webdav")
		enableZFS, _ := cmd.Flags().GetBool("enable-zfs")
		reachableExternally, _ := cmd.Flags().GetBool("reachable-externally")

		opts := hcloud.StorageBoxCreateOpts{
			Name:           name,
			StorageBoxType: &hcloud.StorageBoxType{Name: sbType},
			Location:       &hcloud.Location{Name: location},
			Labels:         labels,
			Password:       password,
			SSHKeys:        sshKeys,
			AccessSettings: &hcloud.StorageBoxCreateOptsAccessSettings{
				ReachableExternally: &reachableExternally,
				SambaEnabled:        &enableSamba,
				SSHEnabled:          &enableSSH,
				WebDAVEnabled:       &enableWebDAV,
				ZFSEnabled:          &enableZFS,
			},
		}
		result, _, err := s.Client().StorageBox().Create(s, opts)
		if err != nil {
			return nil, nil, err
		}

		if err := s.WaitForActions(s, cmd, result.Action); err != nil {
			return nil, nil, err
		}
		cmd.Printf("Storage Box %d created\n", result.StorageBox.ID)

		// TODO change protection here once change-protection is implemented

		return result.StorageBox, util.Wrap("storage_box", hcloud.SchemaFromStorageBox(result.StorageBox)), nil
	},
	PrintResource: func(s state.State, command *cobra.Command, resource any) {
		// storageBox := resource.(*hcloud.StorageBox)

		// TODO should we wait until the storage box is done initializing to display username/server?
	},
}
