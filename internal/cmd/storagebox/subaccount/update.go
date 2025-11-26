package subaccount

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var UpdateCmd = base.UpdateCmd[*hcloud.StorageBoxSubaccount]{
	ResourceNameSingular: "Storage Box Subaccount",
	ShortDescription:     "Update a Storage Box Subaccount",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.StorageBox().Names },
	ValidArgsFunction: func(client hcapi2.Client) []cobra.CompletionFunc {
		return []cobra.CompletionFunc{
			cmpl.SuggestCandidatesF(client.StorageBox().Names),
			SuggestSubaccounts(client),
		}
	},
	PositionalArgumentOverride: []string{"storage-box", "subaccount"},
	FetchWithArgs: func(s state.State, _ *cobra.Command, args []string) (*hcloud.StorageBoxSubaccount, *hcloud.Response, error) {
		storageBoxIDOrName, subaccountIDOrName := args[0], args[1]
		storageBox, _, err := s.Client().StorageBox().Get(s, storageBoxIDOrName)
		if err != nil {
			return nil, nil, err
		}
		if storageBox == nil {
			return nil, nil, fmt.Errorf("Storage Box not found: %s", storageBoxIDOrName)
		}
		return s.Client().StorageBox().GetSubaccount(s, storageBox, subaccountIDOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("description", "", "Description of the Storage Box Snapshot")
		cmd.MarkFlagsOneRequired("description")
	},
	Update: func(s state.State, cmd *cobra.Command, subaccount *hcloud.StorageBoxSubaccount, _ map[string]pflag.Value) error {
		var opts hcloud.StorageBoxSubaccountUpdateOpts
		if cmd.Flags().Changed("description") {
			description, _ := cmd.Flags().GetString("description")
			opts.Description = &description
		}
		_, _, err := s.Client().StorageBox().UpdateSubaccount(s, subaccount, opts)
		if err != nil {
			return err
		}
		return nil
	},
	Experimental: experimental.StorageBoxes,
}
