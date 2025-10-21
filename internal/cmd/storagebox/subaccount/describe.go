package subaccount

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd[*hcloud.StorageBoxSubaccount]{
	ResourceNameSingular:       "Storage Box Subaccount",
	ShortDescription:           "Describe a Storage Box Subaccount",
	PositionalArgumentOverride: []string{"storage-box", "subaccount"},
	ValidArgsFunction: func(client hcapi2.Client) []cobra.CompletionFunc {
		return []cobra.CompletionFunc{
			cmpl.SuggestCandidatesF(client.StorageBox().Names),
			SuggestSubaccounts(client),
		}
	},
	FetchWithArgs: func(s state.State, _ *cobra.Command, args []string) (*hcloud.StorageBoxSubaccount, any, error) {
		storageBoxIDOrName, subaccountIDOrName := args[0], args[1]

		storageBox, _, err := s.Client().StorageBox().Get(s, storageBoxIDOrName)
		if err != nil {
			return nil, nil, err
		}
		if storageBox == nil {
			return nil, nil, fmt.Errorf("Storage Box not found: %s", storageBoxIDOrName)
		}

		subaccount, _, err := s.Client().StorageBox().GetSubaccount(s, storageBox, subaccountIDOrName)
		if err != nil {
			return nil, nil, err
		}
		return subaccount, hcloud.SchemaFromStorageBoxSubaccount(subaccount), nil
	},
	PrintText: func(_ state.State, cmd *cobra.Command, subaccount *hcloud.StorageBoxSubaccount) error {

		cmd.Printf("ID:\t\t\t%d\n", subaccount.ID)
		cmd.Printf("Description:\t\t%s\n", util.NA(subaccount.Description))
		cmd.Printf("Created:\t\t%s (%s)\n", util.Datetime(subaccount.Created), humanize.Time(subaccount.Created))
		cmd.Printf("Username:\t\t%s\n", subaccount.Username)
		cmd.Printf("Home Directory:\t\t%s\n", subaccount.HomeDirectory)
		cmd.Printf("Server:\t\t\t%s\n", subaccount.Server)

		accessSettings := subaccount.AccessSettings
		cmd.Println("Access Settings:")
		cmd.Printf("  Reachable Externally:\t%t\n", accessSettings.ReachableExternally)
		cmd.Printf("  Samba Enabled:\t%t\n", accessSettings.SambaEnabled)
		cmd.Printf("  SSH Enabled:\t\t%t\n", accessSettings.SSHEnabled)
		cmd.Printf("  WebDAV Enabled:\t%t\n", accessSettings.WebDAVEnabled)
		cmd.Printf("  Readonly:\t\t%t\n", accessSettings.Readonly)

		cmd.Println("Labels:")
		if len(subaccount.Labels) == 0 {
			cmd.Println("  No labels")
		} else {
			for key, value := range util.IterateInOrder(subaccount.Labels) {
				cmd.Printf("  %s: %s\n", key, value)
			}
		}

		cmd.Println("Storage Box:")
		cmd.Printf("  ID:\t\t\t%d\n", subaccount.StorageBox.ID)
		return nil
	},
	Experimental: experimental.StorageBoxes,
}
