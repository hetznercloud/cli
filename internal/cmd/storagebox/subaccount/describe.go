package subaccount

import (
	"fmt"
	"io"
	"strings"

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
	PrintText: func(_ state.State, _ *cobra.Command, out io.Writer, subaccount *hcloud.StorageBoxSubaccount) error {
		_, _ = fmt.Fprint(out, DescribeSubaccount(subaccount))
		return nil
	},
	Experimental: experimental.StorageBoxes,
}

func DescribeSubaccount(subaccount *hcloud.StorageBoxSubaccount) string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, "ID:\t%d\n", subaccount.ID)
	_, _ = fmt.Fprintf(&sb, "Description:\t%s\n", util.NA(subaccount.Description))
	_, _ = fmt.Fprintf(&sb, "Created:\t%s (%s)\n", util.Datetime(subaccount.Created), humanize.Time(subaccount.Created))
	_, _ = fmt.Fprintf(&sb, "Username:\t%s\n", subaccount.Username)
	_, _ = fmt.Fprintf(&sb, "Home Directory:\t%s\n", subaccount.HomeDirectory)
	_, _ = fmt.Fprintf(&sb, "Server:\t%s\n", subaccount.Server)

	accessSettings := subaccount.AccessSettings
	_, _ = fmt.Fprintf(&sb, "Access Settings:\n")
	_, _ = fmt.Fprintf(&sb, "  Reachable Externally:\t%t\n", accessSettings.ReachableExternally)
	_, _ = fmt.Fprintf(&sb, "  Samba Enabled:\t%t\n", accessSettings.SambaEnabled)
	_, _ = fmt.Fprintf(&sb, "  SSH Enabled:\t%t\n", accessSettings.SSHEnabled)
	_, _ = fmt.Fprintf(&sb, "  WebDAV Enabled:\t%t\n", accessSettings.WebDAVEnabled)
	_, _ = fmt.Fprintf(&sb, "  Readonly:\t%t\n", accessSettings.Readonly)

	util.DescribeLabels(&sb, subaccount.Labels, "")

	_, _ = fmt.Fprintf(&sb, "Storage Box:\n")
	_, _ = fmt.Fprintf(&sb, "  ID:\t%d\n", subaccount.StorageBox.ID)
	return sb.String()
}
