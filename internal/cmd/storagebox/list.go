package storagebox

import (
	"time"

	"github.com/dustin/go-humanize"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = base.ListCmd[*hcloud.StorageBox, schema.StorageBox]{
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

	OutputTable: func(t *output.Table, _ hcapi2.Client) {
		t.
			AddAllowedFields(hcloud.StorageBox{}).
			AddFieldFn("type", func(obj any) string {
				storageBox := obj.(*hcloud.StorageBox)
				return storageBox.StorageBoxType.Name
			}).
			AddFieldFn("username", func(obj any) string {
				storageBox := obj.(*hcloud.StorageBox)
				return util.OptionalString(storageBox.Username, "-")
			}).
			AddFieldFn("server", func(obj any) string {
				storageBox := obj.(*hcloud.StorageBox)
				return util.OptionalString(storageBox.Server, "-")
			}).
			AddFieldFn("system", func(obj any) string {
				storageBox := obj.(*hcloud.StorageBox)
				return util.OptionalString(storageBox.System, "-")
			}).
			AddFieldFn("location", func(obj any) string {
				storageBox := obj.(*hcloud.StorageBox)
				return storageBox.Location.Name
			}).
			AddFieldFn("size", func(obj any) string {
				storageBox := obj.(*hcloud.StorageBox)
				return humanize.IBytes(storageBox.Stats.Size)
			}).
			AddFieldFn("labels", func(obj any) string {
				storageBox := obj.(*hcloud.StorageBox)
				return util.LabelsToString(storageBox.Labels)
			}).
			AddFieldFn("created", func(obj any) string {
				storageBox := obj.(*hcloud.StorageBox)
				return util.Datetime(storageBox.Created)
			}).
			AddFieldFn("age", func(obj any) string {
				storageBox := obj.(*hcloud.StorageBox)
				return util.Age(storageBox.Created, time.Now())
			})
	},
	Schema: hcloud.SchemaFromStorageBox,
}
