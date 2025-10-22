package subaccount

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = base.ListCmd[*hcloud.StorageBoxSubaccount, schema.StorageBoxSubaccount]{
	ResourceNamePlural: "Storage Box Subaccounts",
	JSONKeyGetByName:   "subaccounts",

	DefaultColumns: []string{"id", "username", "home_directory", "description", "server", "age"},

	ValidArgsFunction: func(client hcapi2.Client) cobra.CompletionFunc {
		return cmpl.SuggestCandidatesF(client.StorageBox().Names)
	},

	PositionalArgumentOverride: []string{"storage-box"},
	SortOption:                 config.OptionSortStorageBoxSubaccount,

	FetchWithArgs: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string, args []string) ([]*hcloud.StorageBoxSubaccount, error) {
		storageBoxIDOrName := args[0]

		storageBox, _, err := s.Client().StorageBox().Get(s, storageBoxIDOrName)
		if err != nil {
			return nil, err
		}
		if storageBox == nil {
			return nil, fmt.Errorf("Storage Box not found: %s", storageBoxIDOrName)
		}

		opts := hcloud.StorageBoxSubaccountListOpts{LabelSelector: listOpts.LabelSelector}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		return s.Client().StorageBox().AllSubaccountsWithOpts(s, storageBox, opts)
	},

	OutputTable: func(t *output.Table[*hcloud.StorageBoxSubaccount], _ hcapi2.Client) {
		t.
			AddAllowedFields(hcloud.StorageBoxSubaccount{}).
			AddFieldFn("description", func(subaccount *hcloud.StorageBoxSubaccount) string {
				return util.NA(subaccount.Description)
			}).
			AddFieldFn("labels", func(subaccount *hcloud.StorageBoxSubaccount) string {
				return util.LabelsToString(subaccount.Labels)
			}).
			AddFieldFn("created", func(subaccount *hcloud.StorageBoxSubaccount) string {
				return util.Datetime(subaccount.Created)
			}).
			AddFieldFn("age", func(subaccount *hcloud.StorageBoxSubaccount) string {
				return util.Age(subaccount.Created, time.Now())
			})
	},

	Schema:       hcloud.SchemaFromStorageBoxSubaccount,
	Experimental: experimental.StorageBoxes,
}
