package subaccount

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var UpdateCmd = base.UpdateCmd{
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
	FetchWithArgs: func(s state.State, _ *cobra.Command, args []string) (any, *hcloud.Response, error) {
		storageBox, _, err := s.Client().StorageBox().Get(s, args[0])
		if err != nil {
			return nil, nil, err
		}
		if storageBox == nil {
			return nil, nil, fmt.Errorf("Storage Box not found: %s", args[0])
		}
		id, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid Storage Box Subaccount ID: %s", args[1])
		}
		return s.Client().StorageBox().GetSubaccountByID(s, storageBox, id)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("description", "", "Description of the Storage Box Snapshot")
		cmd.MarkFlagsOneRequired("description")
	},
	Update: func(s state.State, cmd *cobra.Command, resource interface{}, _ map[string]pflag.Value) error {
		subaccount := resource.(*hcloud.StorageBoxSubaccount)
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
}
