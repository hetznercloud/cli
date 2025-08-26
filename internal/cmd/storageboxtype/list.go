package storageboxtype

import (
	"strconv"

	"github.com/dustin/go-humanize"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = base.ListCmd[*hcloud.StorageBoxType, schema.StorageBoxType]{
	ResourceNamePlural: "Storage Box Types",
	JSONKeyGetByName:   "storage_box_types",
	DefaultColumns:     []string{"id", "name", "description", "size", "snapshot_limit", "automatic_snapshot_limit", "subaccounts_limit"},
	SortOption:         nil, // Storage Box Types do not support sorting

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, _ []string) ([]*hcloud.StorageBoxType, error) {
		opts := hcloud.StorageBoxTypeListOpts{ListOpts: listOpts}
		return s.Client().StorageBoxType().AllWithOpts(s, opts)
	},

	OutputTable: func(t *output.Table, _ hcapi2.Client) {
		t.
			AddAllowedFields(hcloud.StorageBoxType{}).
			AddFieldFn("size", func(obj interface{}) string {
				storageBoxType := obj.(*hcloud.StorageBoxType)
				return humanize.IBytes(uint64(storageBoxType.Size))
			}).
			AddFieldFn("deprecated", func(obj interface{}) string {
				storageBoxType := obj.(*hcloud.StorageBoxType)
				if !storageBoxType.IsDeprecated() {
					return "-"
				}
				return util.Datetime(storageBoxType.UnavailableAfter())
			}).
			AddFieldFn("snapshot_limit", func(obj interface{}) string {
				storageBoxType := obj.(*hcloud.StorageBoxType)
				if storageBoxType.SnapshotLimit == nil {
					return "-"
				}
				return strconv.Itoa(*storageBoxType.SnapshotLimit)
			}).
			AddFieldFn("automatic_snapshot_limit", func(obj interface{}) string {
				storageBoxType := obj.(*hcloud.StorageBoxType)
				if storageBoxType.AutomaticSnapshotLimit == nil {
					return "-"
				}
				return strconv.Itoa(*storageBoxType.AutomaticSnapshotLimit)
			})
	},

	Schema: hcloud.SchemaFromStorageBoxType,
}
