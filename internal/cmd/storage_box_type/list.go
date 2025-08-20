package storageboxtype

import (
	"strconv"

	"github.com/dustin/go-humanize"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = base.ListCmd{
	ResourceNamePlural: "Storage Box Types",
	JSONKeyGetByName:   "storage_box_types",
	DefaultColumns:     []string{"id", "name", "description", "snapshot_limit", "automatic_snapshot_limit", "subaccounts_limit", "size"},

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, _ []string) ([]interface{}, error) {
		opts := hcloud.StorageBoxTypeListOpts{ListOpts: listOpts}
		storageBoxTypes, err := s.Client().StorageBoxType().AllWithOpts(s, opts)

		var resources []interface{}
		for _, n := range storageBoxTypes {
			resources = append(resources, n)
		}
		return resources, err
	},

	OutputTable: func(t *output.Table, _ hcapi2.Client) {
		t.
			AddAllowedFields(hcloud.StorageBoxType{}).
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
			}).
			AddFieldFn("size", func(obj interface{}) string {
				storageBoxType := obj.(*hcloud.StorageBoxType)
				return humanize.IBytes(uint64(storageBoxType.Size))
			})
	},

	Schema: func(resources []interface{}) interface{} {
		storageBoxTypeSchemas := make([]schema.StorageBoxType, 0, len(resources))
		for _, resource := range resources {
			storageBox := resource.(*hcloud.StorageBoxType)
			storageBoxTypeSchemas = append(storageBoxTypeSchemas, hcloud.SchemaFromStorageBoxType(storageBox))
		}
		return storageBoxTypeSchemas
	},
}
