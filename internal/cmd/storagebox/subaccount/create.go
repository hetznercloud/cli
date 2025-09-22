package subaccount

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

var CreateCmd = base.CreateCmd[*hcloud.StorageBoxSubaccount]{
	BaseCobraCommand: func(c hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "create [options] --password <password> --home-directory <home-directory> <storage-box>",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(c.StorageBox().Names)),
			Short:                 "Create a Storage Box Subaccount",
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().String("password", "", "Password for the Subaccount (required)")
		_ = cmd.MarkFlagRequired("password")

		cmd.Flags().String("home-directory", "", "Home directory for the Subaccount (required)")
		_ = cmd.MarkFlagRequired("home-directory")

		cmd.Flags().String("description", "", "Description for the Subaccount")

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		cmd.Flags().Bool("enable-samba", false, "Whether the Samba subsystem should be enabled (true, false)")
		cmd.Flags().Bool("enable-ssh", false, "Whether the SSH subsystem should be enabled (true, false)")
		cmd.Flags().Bool("enable-webdav", false, "Whether the WebDAV subsystem should be enabled (true, false)")
		cmd.Flags().Bool("reachable-externally", false, "Whether the Storage Box should be accessible from outside the Hetzner network (true, false)")
		cmd.Flags().Bool("readonly", false, "Whether the Subaccount should be read-only (true, false)")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) (*hcloud.StorageBoxSubaccount, any, error) {
		storageBoxIDOrName := args[0]
		password, _ := cmd.Flags().GetString("password")
		homeDirectory, _ := cmd.Flags().GetString("home-directory")
		description, _ := cmd.Flags().GetString("description")
		labels, _ := cmd.Flags().GetStringToString("label")

		enableSamba, _ := cmd.Flags().GetBool("enable-samba")
		enableSSH, _ := cmd.Flags().GetBool("enable-ssh")
		enableWebDAV, _ := cmd.Flags().GetBool("enable-webdav")
		reachableExternally, _ := cmd.Flags().GetBool("reachable-externally")
		readonly, _ := cmd.Flags().GetBool("readonly")

		storageBox, _, err := s.Client().StorageBox().Get(s, storageBoxIDOrName)
		if err != nil {
			return nil, nil, err
		}
		if storageBox == nil {
			return nil, nil, fmt.Errorf("Storage Box not found: %s", storageBoxIDOrName)
		}

		opts := hcloud.StorageBoxSubaccountCreateOpts{
			Password:      password,
			HomeDirectory: &homeDirectory,
			AccessSettings: &hcloud.StorageBoxSubaccountAccessSettingsOpts{
				ReachableExternally: &reachableExternally,
				SambaEnabled:        &enableSamba,
				SSHEnabled:          &enableSSH,
				WebDAVEnabled:       &enableWebDAV,
				Readonly:            &readonly,
			},
			Labels: labels,
		}

		if cmd.Flags().Changed("description") {
			opts.Description = &description
		}

		result, _, err := s.Client().StorageBox().CreateSubaccount(s, storageBox, opts)
		if err != nil {
			return nil, nil, err
		}
		if err := s.WaitForActions(s, cmd, result.Action); err != nil {
			return nil, nil, err
		}

		subaccount, _, err := s.Client().StorageBox().GetSubaccountByID(s, storageBox, result.Subaccount.ID)
		if err != nil {
			return nil, nil, err
		}
		if subaccount == nil {
			return nil, nil, fmt.Errorf("Storage Box Subaccount not found: %d", result.Subaccount.ID)
		}

		cmd.Printf("Storage Box Subaccount %d created\n", subaccount.ID)

		return subaccount, util.Wrap("subaccount", hcloud.SchemaFromStorageBoxSubaccount(subaccount)), nil
	},
}
