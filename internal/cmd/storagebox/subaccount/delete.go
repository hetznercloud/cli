package subaccount

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd[*hcloud.StorageBoxSubaccount]{
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

	FetchFunc: func(s state.State, cmd *cobra.Command, args []string) (base.FetchFunc[*hcloud.StorageBoxSubaccount], error) {
		storageBox, _, err := s.Client().StorageBox().Get(cmd.Context(), args[0])
		if err != nil {
			return nil, err
		}
		if storageBox == nil {
			return nil, fmt.Errorf("Storage Box not found: %s", args[0])
		}
		return func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.StorageBoxSubaccount, *hcloud.Response, error) {
			return s.Client().StorageBox().GetSubaccount(cmd.Context(), storageBox, idOrName)
		}, nil
	},

	Delete: func(s state.State, cmd *cobra.Command, subaccount *hcloud.StorageBoxSubaccount) ([]*hcloud.Action, error) {
		result, _, err := s.Client().StorageBox().DeleteSubaccount(cmd.Context(), subaccount)
		return []*hcloud.Action{result.Action}, err
	},
}
