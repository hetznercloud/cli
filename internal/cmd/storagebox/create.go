package storagebox

import (
	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
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
		_ = cmd.RegisterFlagCompletionFunc("location", cmpl.SuggestCandidatesF(client.Location().Names))

		cmd.Flags().String("password", "", "The password that will be set for this Storage Box (required)")
		_ = cmd.MarkFlagRequired("password")

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		// TODO: Can OpenSSH Public Keys have comma in them? If yes we should use StringArray instead
		// TODO: How to handle SSH Keys? --public-key-from-file?
		cmd.Flags().StringSlice("ssh-keys", []string{}, "SSH public keys in OpenSSH format which should be injected into the Storage Box")

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
	Run: func(state state.State, command *cobra.Command, strings []string) (any, any, error) {
		return nil, nil, nil
	},
	PrintResource: func(s state.State, command *cobra.Command, resource any) {
		storageBox := resource.(*hcloud.StorageBox)

	},
}
