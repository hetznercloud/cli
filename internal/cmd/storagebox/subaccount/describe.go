package subaccount

import (
	"fmt"
	"io"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
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
	PrintText: func(_ state.State, _ *cobra.Command, out io.Writer, subaccount *hcloud.StorageBoxSubaccount) error {
		fmt.Fprintf(out, "ID:\t%d\n", subaccount.ID)
		fmt.Fprintf(out, "Name:\t%s\n", subaccount.Name)
		fmt.Fprintf(out, "Description:\t%s\n", util.NA(subaccount.Description))
		fmt.Fprintf(out, "Created:\t%s (%s)\n", util.Datetime(subaccount.Created), humanize.Time(subaccount.Created))
		fmt.Fprintf(out, "Username:\t%s\n", subaccount.Username)
		fmt.Fprintf(out, "Home Directory:\t%s\n", subaccount.HomeDirectory)
		fmt.Fprintf(out, "Server:\t%s\n", subaccount.Server)

		accessSettings := subaccount.AccessSettings
		fmt.Fprintln(out)
		fmt.Fprintf(out, "Access Settings:\n")
		fmt.Fprintf(out, "  Reachable Externally:\t%t\n", accessSettings.ReachableExternally)
		fmt.Fprintf(out, "  Samba Enabled:\t%t\n", accessSettings.SambaEnabled)
		fmt.Fprintf(out, "  SSH Enabled:\t%t\n", accessSettings.SSHEnabled)
		fmt.Fprintf(out, "  WebDAV Enabled:\t%t\n", accessSettings.WebDAVEnabled)
		fmt.Fprintf(out, "  Readonly:\t%t\n", accessSettings.Readonly)

		fmt.Fprintln(out)
		util.DescribeLabels(out, subaccount.Labels, "")

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Storage Box:\n")
		fmt.Fprintf(out, "  ID:\t%d\n", subaccount.StorageBox.ID)

		return nil
	},
}
