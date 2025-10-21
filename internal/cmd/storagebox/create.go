package storagebox

import (
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var CreateCmd = base.CreateCmd[*hcloud.StorageBox]{
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

		cmd.Flags().StringArray("ssh-key", []string{}, "SSH public keys in OpenSSH format or as the ID or name of an existing SSH key")
		_ = cmd.RegisterFlagCompletionFunc("ssh-key", cmpl.SuggestCandidatesF(client.SSHKey().Names))

		cmd.Flags().Bool("enable-samba", false, "Whether the Samba subsystem should be enabled (true, false)")
		cmd.Flags().Bool("enable-ssh", false, "Whether the SSH subsystem should be enabled (true, false)")
		cmd.Flags().Bool("enable-webdav", false, "Whether the WebDAV subsystem should be enabled (true, false)")
		cmd.Flags().Bool("enable-zfs", false, "Whether the ZFS Snapshot folder should be visible (true, false)")
		cmd.Flags().Bool("reachable-externally", false, "Whether the Storage Box should be accessible from outside the Hetzner network (true, false)")

		cmd.Flags().StringSlice("enable-protection", []string{}, "Enable protection (delete) (default: none)")
		_ = cmd.RegisterFlagCompletionFunc("enable-protection", cmpl.SuggestCandidates("delete"))

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, _ []string) (*hcloud.StorageBox, any, error) {
		name, _ := cmd.Flags().GetString("name")
		sbType, _ := cmd.Flags().GetString("type")
		location, _ := cmd.Flags().GetString("location")
		password, _ := cmd.Flags().GetString("password")
		sshKeys, _ := cmd.Flags().GetStringArray("ssh-key")
		labels, _ := cmd.Flags().GetStringToString("label")
		protection, _ := cmd.Flags().GetStringSlice("enable-protection")

		protectionOpts, err := getChangeProtectionOpts(true, protection)
		if err != nil {
			return nil, nil, err
		}

		enableSamba, _ := cmd.Flags().GetBool("enable-samba")
		enableSSH, _ := cmd.Flags().GetBool("enable-ssh")
		enableWebDAV, _ := cmd.Flags().GetBool("enable-webdav")
		enableZFS, _ := cmd.Flags().GetBool("enable-zfs")
		reachableExternally, _ := cmd.Flags().GetBool("reachable-externally")

		if !cmd.Flags().Changed("ssh-key") && config.OptionDefaultSSHKeys.Changed(s.Config()) {
			sshKeys, err = config.OptionDefaultSSHKeys.Get(s.Config())
			if err != nil {
				return nil, nil, err
			}
		}

		resolvedSSHKeys := make([]*hcloud.SSHKey, len(sshKeys))
		for i, sshKey := range sshKeys {
			resolvedSSHKeys[i], err = resolveSSHKey(s, sshKey)
			if err != nil {
				return nil, nil, err
			}
		}

		var accessSettings hcloud.StorageBoxCreateOptsAccessSettings
		if cmd.Flags().Changed("enable-samba") {
			accessSettings.SambaEnabled = &enableSamba
		}
		if cmd.Flags().Changed("enable-ssh") {
			accessSettings.SSHEnabled = &enableSSH
		}
		if cmd.Flags().Changed("enable-webdav") {
			accessSettings.WebDAVEnabled = &enableWebDAV
		}
		if cmd.Flags().Changed("enable-zfs") {
			accessSettings.ZFSEnabled = &enableZFS
		}
		if cmd.Flags().Changed("reachable-externally") {
			accessSettings.ReachableExternally = &reachableExternally
		}

		opts := hcloud.StorageBoxCreateOpts{
			Name:           name,
			StorageBoxType: &hcloud.StorageBoxType{Name: sbType},
			Location:       &hcloud.Location{Name: location},
			Labels:         labels,
			Password:       password,
			SSHKeys:        resolvedSSHKeys,
			AccessSettings: &accessSettings,
		}
		result, _, err := s.Client().StorageBox().Create(s, opts)
		if err != nil {
			return nil, nil, err
		}

		if err := s.WaitForActions(s, cmd, result.Action); err != nil {
			return nil, nil, err
		}
		cmd.Printf("Storage Box %d created\n", result.StorageBox.ID)

		storageBox, _, err := s.Client().StorageBox().GetByID(s, result.StorageBox.ID)
		if err != nil {
			return nil, nil, err
		}
		if storageBox == nil {
			return nil, nil, fmt.Errorf("Storage Box not found: %d", result.StorageBox.ID)
		}

		if err := changeProtection(s, cmd, storageBox, true, protectionOpts); err != nil {
			return nil, nil, err
		}

		return storageBox, util.Wrap("storage_box", hcloud.SchemaFromStorageBox(result.StorageBox)), nil
	},
	PrintResource: func(_ state.State, cmd *cobra.Command, storageBox *hcloud.StorageBox) {
		cmd.Printf("Server: %s\n", storageBox.Server)
		cmd.Printf("Username: %s\n", storageBox.Username)
	},
	Experimental: experimental.StorageBoxes,
}

// resolveSSHKey resolves the given pubKey by doing the following:
// - If pubKey is a valid public key in OpenSSH format, it is wrapped in a [hcloud.SSHKey].
// - Otherwise, it is treated as an ID or name of an existing SSH key in the project.
// - If an SSH key with the given ID or name exists, it is returned.
// - Otherwise, pubKey is returned wrapped in a [hcloud.SSHKey].
func resolveSSHKey(s state.State, pubKey string) (*hcloud.SSHKey, error) {
	_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(pubKey))
	if err == nil {
		return &hcloud.SSHKey{PublicKey: pubKey}, nil
	}

	sshKey, _, err := s.Client().SSHKey().Get(s, pubKey)
	if err != nil {
		return nil, err
	}
	if sshKey != nil {
		return sshKey, nil
	}

	return &hcloud.SSHKey{PublicKey: pubKey}, nil
}
