package subaccount

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds[*hcloud.StorageBoxSubaccount]{
	ResourceNameSingular:   "Storage Box Subaccount",
	ShortDescriptionAdd:    "Add a label to a Storage Box Subaccount",
	ShortDescriptionRemove: "Remove a label from a Storage Box Subaccount",

	PositionalArgumentOverride: []string{"storage-box", "subaccount"},
	ValidArgsFunction: func(client hcapi2.Client) []cobra.CompletionFunc {
		return []cobra.CompletionFunc{
			cmpl.SuggestCandidatesF(client.StorageBox().Names),
			SuggestSubaccounts(client),
			cmpl.SuggestCandidatesCtx(func(_ *cobra.Command, args []string) []string {
				if len(args) < 2 {
					return nil
				}
				return client.StorageBox().SubaccountLabelKeys(args[0], args[1])
			}),
		}
	},

	FetchWithArgs: func(s state.State, args []string) (*hcloud.StorageBoxSubaccount, error) {
		storageBoxIDOrName, subaccountIDOrName := args[0], args[1]

		storageBox, _, err := s.Client().StorageBox().Get(s, storageBoxIDOrName)
		if err != nil {
			return nil, err
		}
		if storageBox == nil {
			return nil, fmt.Errorf("Storage Box not found: %s", storageBoxIDOrName)
		}

		subaccount, _, err := s.Client().StorageBox().GetSubaccount(s, storageBox, subaccountIDOrName)
		if err != nil {
			return nil, err
		}
		if subaccount == nil {
			return nil, fmt.Errorf("Storage Box Subaccount not found: %s", subaccountIDOrName)
		}
		return subaccount, nil
	},

	SetLabels: func(s state.State, subaccount *hcloud.StorageBoxSubaccount, labels map[string]string) error {
		opts := hcloud.StorageBoxSubaccountUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().StorageBox().UpdateSubaccount(s, subaccount, opts)
		return err
	},

	GetLabels: func(subaccount *hcloud.StorageBoxSubaccount) map[string]string {
		return subaccount.Labels
	},

	GetIDOrName: func(subaccount *hcloud.StorageBoxSubaccount) string {
		if subaccount.Name != "" {
			return subaccount.Name
		}
		return strconv.FormatInt(subaccount.ID, 10)
	},
}
