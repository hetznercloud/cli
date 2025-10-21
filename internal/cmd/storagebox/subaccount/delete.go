package subaccount

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
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
		return func(s state.State, _ *cobra.Command, idOrName string) (any, *hcloud.Response, error) {
			return s.Client().StorageBox().GetSubaccount(s, storageBox, idOrName)
		}, nil
	},

	Delete: func(s state.State, _ *cobra.Command, resource any) (*hcloud.Action, error) {
		subaccount := resource.(*hcloud.StorageBoxSubaccount)
		result, _, err := s.Client().StorageBox().DeleteSubaccount(s, subaccount)
		return result.Action, err
	},
	Experimental: experimental.StorageBoxes,
}
