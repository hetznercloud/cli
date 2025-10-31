package storagebox

import (
	"time"

	"github.com/dustin/go-humanize"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = &base.ListCmd[*hcloud.StorageBox, schema.StorageBox]{
	ResourceNamePlural: "Storage Boxes",
	JSONKeyGetByName:   "storage_boxes",
	DefaultColumns:     []string{"id", "name", "username", "server", "type", "size", "location", "age"},
	SortOption:         config.OptionSortStorageBox,

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]*hcloud.StorageBox, error) {
		opts := hcloud.StorageBoxListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		return s.Client().StorageBox().AllWithOpts(s, opts)
	},

	OutputTable: func(t *output.Table[*hcloud.StorageBox], _ hcapi2.Client) {
		t.
			AddAllowedFields(&hcloud.StorageBox{}).
			AddFieldFn("type", func(storageBox *hcloud.StorageBox) string {
				return storageBox.StorageBoxType.Name
			}).
			AddFieldFn("location", func(storageBox *hcloud.StorageBox) string {
				return storageBox.Location.Name
			}).
			AddFieldFn("size", func(storageBox *hcloud.StorageBox) string {
				return humanize.IBytes(storageBox.Stats.Size)
			}).
			AddFieldFn("labels", func(storageBox *hcloud.StorageBox) string {
				return util.LabelsToString(storageBox.Labels)
			}).
			AddFieldFn("created", func(storageBox *hcloud.StorageBox) string {
				return util.Datetime(storageBox.Created)
			}).
			AddFieldFn("age", func(storageBox *hcloud.StorageBox) string {
				return util.Age(storageBox.Created, time.Now())
			})
	},
	Schema:       hcloud.SchemaFromStorageBox,
	Experimental: experimental.StorageBoxes,
}
