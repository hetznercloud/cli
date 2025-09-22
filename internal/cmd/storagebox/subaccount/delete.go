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

var DeleteCmd = base.DeleteCmd{
	ResourceNameSingular:       "Storage Box Subaccount",
	ResourceNamePlural:         "Storage Box Subaccounts",
	ShortDescription:           "Delete a Storage Box Subaccount",
	PositionalArgumentOverride: []string{"storage-box", "subaccount"},
	ValidArgsFunction: func(client hcapi2.Client) []cobra.CompletionFunc {
		return []cobra.CompletionFunc{
			cmpl.SuggestCandidatesF(client.StorageBox().Names),
			SuggestSubaccounts(client),
		}
	},

	FetchFunc: func(s state.State, _ *cobra.Command, args []string) (base.FetchFunc, error) {
		storageBox, _, err := s.Client().StorageBox().Get(s, args[0])
		if err != nil {
			return nil, err
		}
		if storageBox == nil {
			return nil, fmt.Errorf("Storage Box not found: %s", args[0])
		}
		return func(s state.State, _ *cobra.Command, idStr string) (any, *hcloud.Response, error) {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid Storage Box Subaccount ID: %s", idStr)
			}
			return s.Client().StorageBox().GetSubaccountByID(s, storageBox, id)
		}, nil
	},

	Delete: func(s state.State, _ *cobra.Command, resource any) (*hcloud.Action, error) {
		subaccount := resource.(*hcloud.StorageBoxSubaccount)
		action, _, err := s.Client().StorageBox().DeleteSubaccount(s, subaccount)
		return action, err
	},
}
