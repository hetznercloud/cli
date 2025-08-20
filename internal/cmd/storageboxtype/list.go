package storageboxtype

import (
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

var ListCmd = base.ListCmd{
	ResourceNamePlural: "Storage Box Types",
	JSONKeyGetByName:   "storage_box_types",
	DefaultColumns:     []string{"id", "name", "description", "size", "snapshot_limit", "automatic_snapshot_limit", "subaccounts_limit"},
	SortOption:         nil, // Storage Box Types do not support sorting

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, _ []string) ([]interface{}, error) {
		opts := hcloud.StorageBoxTypeListOpts{ListOpts: listOpts}
		storageBoxTypes, err := s.Client().StorageBoxType().AllWithOpts(s, opts)

		var resources []interface{}
		for _, r := range storageBoxTypes {
			resources = append(resources, r)
		}
		return resources, err
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
			})
	},

	Schema: func(resources []interface{}) interface{} {
		storageBoxTypeSchemas := make([]schema.StorageBoxType, 0, len(resources))
		for _, resource := range resources {
			storageBoxType := resource.(*hcloud.StorageBoxType)
			storageBoxTypeSchemas = append(storageBoxTypeSchemas, hcloud.SchemaFromStorageBoxType(storageBoxType))
		}
		return storageBoxTypeSchemas
	},
}
